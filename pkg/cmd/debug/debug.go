// Copyright 2022 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package debug

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	dbgcmdutil "velostrata-internal.googlesource.com/containerdbg.git/pkg/cmdutil"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/collect"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/debug"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/decoder"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/install"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/k8s"
)

type debugOptions struct {
	outputFilename string
	yamlFilename   string
	imageName      string
}

func (o *debugOptions) validate() error {
	if o.imageName == "" && o.yamlFilename == "" {
		return fmt.Errorf("either image name of filename must be provided")
	}

	if o.outputFilename == "" {
		return fmt.Errorf("output filename must be supplied")
	}

	return nil
}

func NewDebugCmd(f cmdutil.Factory, streams genericclioptions.IOStreams, interruptCh <-chan os.Signal) *cobra.Command {
	o := debugOptions{}
	cmd := &cobra.Command{
		Use:  "debug [image name]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				o.imageName = args[0]
			}
			if err := o.validate(); err != nil {
				return err
			}
			return o.debugImage(cmd.Context(), f, streams, interruptCh)
		},
	}

	cmd.Flags().StringVarP(&o.outputFilename, "output", "o", "", "the output filename to which we will save the events.json file")
	cmd.Flags().StringVarP(&o.yamlFilename, "filename", "f", "", "A yaml describing a deployment containerdbg should debug. If provided image name is ignored.")

	return cmd
}

func deployFromImageName(ctx context.Context, f cmdutil.Factory, imagename, namespace string) (func() error, error) {
	clientset, err := f.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	deployment, err := debug.CreateDebugDeploymentForImage(imagename, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed creating deployment for image %v", err)
	}
	created, err := clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment %v", err)
	}
	return func() error {
		return clientset.AppsV1().Deployments(namespace).Delete(ctx, created.Name, metav1.DeleteOptions{})
	}, nil
}

func deployFromYaml(ctx context.Context, f cmdutil.Factory, yamlFilename string, namespace string) (func() error, error) {
	appliedObjs := []k8s.Object{}
	cl, err := dbgcmdutil.CreateControllerClient(f)
	if err != nil {
		return nil, err
	}
	yamlFile, err := os.Open(yamlFilename)
	if err != nil {
		return nil, err
	}
	defer yamlFile.Close()
	err = decoder.DecodeEach(ctx, yamlFile, func(ctx context.Context, obj k8s.Object) error {
		obj.SetNamespace(namespace)
		if err := cl.Create(ctx, obj); err != nil {
			return err
		}
		appliedObjs = append(appliedObjs, obj)
		return nil
	}, decoder.WithMutation(func(obj k8s.Object) error {
		dep, ok := obj.(*appsv1.Deployment)
		if !ok {
			return nil
		}
		return debug.MutateDeployment(dep)
	}))
	if err != nil {
		return nil, err
	}

	// cleanup
	return func() error {
		for _, o := range appliedObjs {
			cl.Delete(ctx, o)
		}
		return nil
	}, nil
}

func (o *debugOptions) debugImage(ctx context.Context, f cmdutil.Factory, streams genericclioptions.IOStreams, interruptCh <-chan os.Signal) error {

	outFile, err := os.Create(o.outputFilename)
	if err != nil {
		return err
	}

	namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	if err := install.EnsureInstallation(ctx, f, streams); err != nil {
		return fmt.Errorf("failed to ensure installation of containerdbg system: %v", err)
	}
	defer install.Uninstall(ctx, f, streams)

	if o.yamlFilename == "" {
		cleanup, err := deployFromImageName(ctx, f, o.imageName, namespace)
		if err != nil {
			return err
		}
		defer cleanup()
	} else {
		cleanup, err := deployFromYaml(ctx, f, o.yamlFilename, namespace)
		if err != nil {
			return err
		}
		defer cleanup()
	}

	fmt.Fprintf(streams.ErrOut, "Press Ctrl-C to finish the debugging session and download the collected report\n\n")

	<-interruptCh

	if err := collect.CollectRecordedData(ctx, f, outFile); err != nil {
		return fmt.Errorf("failed to collect data from system %v", err)
	}

	return nil
}
