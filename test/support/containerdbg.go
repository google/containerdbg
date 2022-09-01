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

package support

import (
	"context"
	"os"
	"testing"
	"time"

	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/consts"
)

func DumpContainerDbgLogs(ctx context.Context, t *testing.T, clientset *kubernetes.Clientset) {
	daemon := appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{Name: "containerdbg-daemonset", Namespace: consts.ContainerdbgNamespace},
	}
	DumpLogs(ctx, t, clientset.CoreV1(), &daemon)
}

func RunContainerDebug(t *testing.T, ctx context.Context, cfg *envconf.Config, outpath string, namespace string, debugParams ...string) context.Context {
	client, clientset := NewK8sClient(ctx, t, cfg)
	if err := client.Resources().Create(ctx, &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}); err != nil && !k8serrors.IsAlreadyExists(err) {
		t.Fatal(err)
	}
	intCh := make(chan os.Signal, 1)
	cmdErrCh := make(chan error, 1)
	cmdCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd, _, _, _ := NewTestRootCmd(cfg, intCh)
	cmd.SetArgs(append([]string{"debug", "-o", outpath}, debugParams...))
	go func() {
		cmdErrCh <- cmd.ExecuteContext(cmdCtx)
	}()

	list := v1.PodList{}

	for {
		if err := client.Resources(namespace).List(ctx, &list); err != nil {
			cancel()
			t.Fatal(err)
		}

		if len(list.Items) > 0 {
			break
		}
		select {
		case <-time.After(time.Second * 2):
			continue
		case err := <-cmdErrCh:
			t.Fatalf("command exited with error: %s", err)
		}
	}
	pod := list.Items[0]

	err := wait.For(conditions.New(client.Resources(namespace)).PodReady(&pod), wait.WithTimeout(time.Minute*1))
	if testing.Verbose() {
		DumpEvents(ctx, t, clientset.CoreV1(), &pod, 200)
		DumpLogs(ctx, t, clientset.CoreV1(), &pod)
		DumpContainerDbgLogs(ctx, t, clientset)
	}

	if err != nil {
		DumpEvents(ctx, t, clientset.CoreV1(), &pod, 200)
		DumpLogs(ctx, t, clientset.CoreV1(), &pod)
		t.Fatal(err)
	}

	time.Sleep(time.Second * 30)

	intCh <- os.Interrupt

	if err := <-cmdErrCh; err != nil {
		t.Fatal(err)
	}

	return ctx
}
