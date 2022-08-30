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

package collect

import (
	"context"
	"fmt"
	"net/http"

	"io"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/connect"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/consts"
)

func CollectRecordedData(ctx context.Context, f cmdutil.Factory, outputStream io.Writer) error {
	config, err := f.ToRESTConfig()
	if err != nil {
		return err
	}

	clientset, err := f.KubernetesClientSet()
	if err != nil {
		return err
	}

	daemonset, err := clientset.AppsV1().DaemonSets(consts.ContainerdbgNamespace).Get(ctx, consts.ContainerdbgDaemonsetName, v1.GetOptions{})
	if err != nil {
		return err
	}

	err = connect.TunnelForEachPod(ctx, config, clientset, daemonset, 8080, func(port int, podName, namespace string) error {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/collect", port))
		if err != nil {
			return err
		}
		_, err = io.Copy(outputStream, resp.Body)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
