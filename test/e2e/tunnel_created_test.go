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

package test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/connect"
	"velostrata-internal.googlesource.com/containerdbg.git/test/support"
)

func TestTunnelWorks(t *testing.T) {
	tunnelWorks := features.New("daemonset tunnel").
		WithLabel("type", "e2e").
		Assess("http server is reachable through tunnel", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			nodes := v1.NodeList{}

			if err := client.Resources().List(ctx, &nodes); err != nil {
				t.Fatal(err)
			}

			daemonset := appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "echo-server", Namespace: "default"},
			}

			support.ApplyYaml(ctx, client.Resources(), "./resources/echo-server.yaml")

			err = wait.For(conditions.New(client.Resources()).ResourceMatch(&daemonset, func(object k8s.Object) bool {
				d := object.(*appsv1.DaemonSet)
				return d.Status.NumberAvailable == int32(len(nodes.Items))
			}), wait.WithTimeout(time.Minute*1))
			t.Logf("done waiting for daemonset")

			if err := client.Resources().Get(ctx, "echo-server", "default", &daemonset); err != nil {
				t.Fatal(err)
			}

			ctxTimeout, cancel := context.WithTimeout(ctx, time.Minute*1)
			defer cancel()

			err = connect.TunnelForEachPod(ctxTimeout, client.RESTConfig(), kubernetes.NewForConfigOrDie(client.RESTConfig()), &daemonset, 8080, func(port int, podName, namespace string) error {
				cookie := "key=1234key1234"
				resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/?%s", port, cookie))
				if err != nil {
					t.Fatal(err)
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}

				if !strings.Contains(string(body), cookie) {
					t.Fatalf("returned body from http echo server does not contain the sent request: %s", body)
				}

				return nil

			})

			if err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()

	testenv.Test(t, tunnelWorks)
}
