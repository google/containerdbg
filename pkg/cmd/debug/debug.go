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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	dbgcmdutil "github.com/google/containerdbg/pkg/cmdutil"
	"github.com/google/containerdbg/pkg/collect"
	"github.com/google/containerdbg/pkg/debug"
	"github.com/google/containerdbg/pkg/decoder"
	"github.com/google/containerdbg/pkg/install"
	"github.com/google/containerdbg/pkg/k8s"
)

type debugOptions struct {
	keepInstall    bool
	outputFilename string
	yamlFilename   string
	imageName      string
}

func (o *debugOptions) validate() error {
	if o.imageName == "" && o.yamlFilename == "" {
		return fmt.Errorf("either image name of filename must be provided")
	}

	if o.imageName != "" && o.yamlFilename != "" {
		return fmt.Errorf("only one of a filename or an image name may be provided")
	}

	if o.outputFilename == "" {
		return fmt.Errorf("output filename must be supplied")
	}

	return nil
}

func NewDebugCmd(f cmdutil.Factory, streams genericclioptions.IOStreams, interruptCh <-chan os.Signal) *cobra.Command {
	o := debugOptions{}
	cmd := &cobra.Command{
		Use:   "debug {-f yaml | -i imagespec}",
		Short: "debug either using a yaml file or an image name",
		Long: "debug will deploy either an image as a deployment or an existing yaml file - " +
			"adding all of the neccassary options for debuging with containerdbg",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.validate(); err != nil {
				return err
			}
			return o.debugImage(cmd.Context(), f, streams, interruptCh)
		},
	}

	cmd.Flags().BoolVarP(&o.keepInstall, "keep", "k", false, "if specified will not uninstall containerdbg from cluster")
	cmd.Flags().StringVarP(&o.outputFilename, "output", "o", "", "the output filename to which we will save the events.pb file")
	cmd.Flags().StringVarP(&o.yamlFilename, "filename", "f", "", "A yaml describing a deployment containerdbg should debug. Cannot be supplied if specifying an image")
	cmd.Flags().StringVarP(&o.imageName, "image", "i", "", "An image path to deploy and debug. Cannot be supplied with a yaml file")

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
	}, decoder.WithMutation(mutatePodSpec))
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

func mutatePodSpec(obj k8s.Object) error {
	switch o := obj.(type) {
	case *appsv1.Deployment:
		return debug.ModifyPodSpec(&o.Spec.Template.Spec)
	case *appsv1.ReplicaSet:
		return debug.ModifyPodSpec(&o.Spec.Template.Spec)
	case *appsv1.StatefulSet:
		return debug.ModifyPodSpec(&o.Spec.Template.Spec)
	case *appsv1.DaemonSet:
		return debug.ModifyPodSpec(&o.Spec.Template.Spec)
	case *v1.Pod:
		return debug.ModifyPodSpec(&o.Spec)
	default:
		return nil
	}
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

	interrupttedCtx, cancel := context.WithCancel(ctx)

	go func() {
		<-interruptCh
		cancel()
	}()
	wasInstalled, err := install.EnsureInstallation(interrupttedCtx, f, streams)
	if err != nil {
		return fmt.Errorf("failed to ensure installation of containerdbg system: %v", err)
	}

	// avoid uninstalling if it was installed independently of this command
	if !o.keepInstall && !wasInstalled {
		defer install.Uninstall(ctx, f, streams)
	}

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

	<-interrupttedCtx.Done()

	fmt.Fprintf(streams.ErrOut, "Debug session ended succefully, cleaning up...")

	if err := collect.CollectRecordedData(ctx, f, outFile); err != nil {
		return fmt.Errorf("failed to collect data from system %v", err)
	}

	return nil
}
