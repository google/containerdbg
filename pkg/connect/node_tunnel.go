// Copyright 2021 Google LLC All Rights Reserved.
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

package connect

import (
	"context"
	"fmt"
	"net"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"github.com/google/containerdbg/pkg/polymorphichelpers"
)

type tunnelHandler func(port int, podName string, namespace string) error

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func TunnelForPod(ctx context.Context, config *rest.Config, pod *corev1.Pod, port int, handle tunnelHandler) error {
	nextPort, err := GetFreePort()
	if err != nil {
		return err
	}
	ports := []string{fmt.Sprintf("%d:%d", nextPort, port)}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	readyChannel := make(chan struct{}, 1)
	errChannel := make(chan error, 1)
	go func() {
		err := StartTunnel(config, pod.Name, pod.Namespace, ports, WithStopChannel(ctx.Done()),
			WithReadyChannel(readyChannel))
		if err != nil {
			errChannel <- err
		}
	}()

	select {
	case <-readyChannel:
	case err := <-errChannel:
		return err
	}

	return handle(nextPort, pod.Name, pod.Namespace)
}

func TunnelForEachPod(ctx context.Context, config *rest.Config, clientset *kubernetes.Clientset, object runtime.Object, port int, handle tunnelHandler) error {

	pods, err := polymorphichelpers.PodsForObject(ctx, clientset.CoreV1(), object)
	if err != nil {
		return err
	}

	for _, pod := range pods {
		TunnelForPod(ctx, config, pod, port, handle)
	}

	return nil
}
