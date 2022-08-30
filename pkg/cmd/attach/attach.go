// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attach

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/consts"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/debug"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/install"
)

type attachOptions struct {
	podname string
}

func NewAttachCmd(f cmdutil.Factory, streams genericclioptions.IOStreams, interruptCh <-chan os.Signal) *cobra.Command {
	o := attachOptions{}
	cmd := &cobra.Command{
		Use:  "attach <pod name>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.podname = args[0]
			return o.attach(cmd.Context(), f, streams, interruptCh)
		},
	}

	return cmd
}

func createProxyPod(namespace string, nodename string) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "proxy-pod-",
			Namespace:    namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "socat",
					Image: "alpine/socat",
					Args: []string{
						fmt.Sprintf("tcp-listen:%d,fork,reuseaddr", consts.DaemonProxyPort),
						fmt.Sprintf("unix-connect:%s/%s", "/var/run/containerdbg/daemon", consts.NodeDaemonSocketName),
					},
					VolumeMounts: []v1.VolumeMount{
						{
							MountPath: "/var/run/containerdbg/daemon",
							Name:      "shareddir",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "shareddir",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/var/run/containerdbg/daemon",
						},
					},
				},
			},
			NodeName: nodename,
		},
	}
}

func (o *attachOptions) attach(ctx context.Context, f cmdutil.Factory, streams genericclioptions.IOStreams, interruptCh <-chan os.Signal) error {
	clientset, err := f.KubernetesClientSet()
	if err != nil {
		return err
	}

	namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, o.podname, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if pod.Spec.NodeName == "" {
		return fmt.Errorf("debugged pod is not yet scheduled, please try again in a few minutes")
	}

	interrupttedCtx, cancel := context.WithCancel(ctx)

	go func() {
		<-interruptCh
		cancel()
	}()

	if _, err := install.EnsureInstallation(interrupttedCtx, f, streams); err != nil {
		return fmt.Errorf("failed to ensure installation of containerdbg system: %v", err)
	}

	proxypod := createProxyPod(namespace, pod.Spec.NodeName)

	proxypod, err = clientset.CoreV1().Pods(namespace).Create(ctx, proxypod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	w, err := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.SingleObject(proxypod.ObjectMeta))
	for e := range w.ResultChan() {
		if e.Type != watch.Modified {
			continue
		}
		proxypod = e.Object.(*v1.Pod)
		if proxypod.Status.PodIP != "" {
			break
		}
	}
	w.Stop()

	podIP := proxypod.Status.PodIP

	container := v1.EphemeralContainer{
		EphemeralContainerCommon: debug.GetDnsProxyContainer(),
	}
	container.Env = append(container.Env, v1.EnvVar{
		Name:  consts.DaemonProxyEnv,
		Value: podIP,
	})

	patch := make([]map[string]interface{}, 1)
	patch[0] = map[string]interface{}{}
	patch[0]["op"] = "add"
	if pod.Spec.EphemeralContainers == nil {
		patch[0]["path"] = "/spec/ephemeralContainers"
		patch[0]["value"] = []v1.EphemeralContainer{container}
	} else {
		patch[0]["path"] = "/spec/ephemeralContainers/-"
		patch[0]["value"] = container
	}
	patchJS, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	result := clientset.CoreV1().RESTClient().Patch(types.JSONPatchType).
		Namespace(namespace).
		Resource("pods").
		Name(o.podname).
		SubResource("ephemeralcontainers").
		Body(patchJS).
		Do(ctx)

	err = result.Error()

	if err != nil {
		if serr, ok := err.(*errors.StatusError); ok && serr.Status().Reason == metav1.StatusReasonNotFound && serr.ErrStatus.Details.Name == "" {
			return fmt.Errorf("ephemeral containers are disabled for this cluster (error from server: %q).", err)
		}
		return err
	}

	return nil
}
