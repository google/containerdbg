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
	"os"
	"path"
	"testing"
	"time"

	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/consts"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/rand"
	"velostrata-internal.googlesource.com/containerdbg.git/test/support"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func helperFindEvent(t *testing.T, filename string) bool {

	t.Helper()
	// Very damp scanning for the failed open file
	recordsFile, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	reader := events.NewEventReader(recordsFile)
	for event, err := reader.Read(); err == nil; event, err = reader.Read() {
		t.Logf("event line %+v", event)
		syscall := event.GetSyscall()
		if syscall == nil {
			continue
		}
		if syscall.GetOpen().GetPath() == "/doesnotexists" {
			return true
		}
	}

	return false
}

func helperTestOpenIsRecorded(t *testing.T, ctx context.Context, cfg *envconf.Config, namespace string, debugParams ...string) context.Context {
	client, err := cfg.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Resources().Create(ctx, &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}); err != nil {
		t.Fatal(err)
	}
	tmpFileName := path.Join(t.TempDir(), "events.json")
	intCh := make(chan os.Signal, 1)
	cmdErrCh := make(chan error, 1)
	cmdCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd, _, _, _ := support.NewTestRootCmd(cfg, intCh)
	cmd.SetArgs(append([]string{"debug", "-o", tmpFileName}, debugParams...))
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

	err = wait.For(conditions.New(client.Resources(namespace)).PodReady(&pod), wait.WithTimeout(time.Minute*1))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 30)

	intCh <- os.Interrupt

	if err := <-cmdErrCh; err != nil {
		t.Fatal(err)
	}

	if !helperFindEvent(t, tmpFileName) {
		t.Fatal("did not find open file event")
	}

	return ctx
}

func TestFullE2EFlow(t *testing.T) {
	systemInstallation := features.New("open is recorded").
		WithLabel("type", "e2e").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			clientset, err := kubernetes.NewForConfig(client.RESTConfig())
			if err != nil {
				t.Fatal(err)
			}

			daemon := appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Name: "containerdbg-daemonset", Namespace: consts.ContainerdbgNamespace},
			}

			nodes := v1.NodeList{}

			if err := client.Resources().List(ctx, &nodes); err != nil {
				t.Fatal(err)
			}
			t.Logf("there are %d nodes in the cluster", len(nodes.Items))

			err = wait.For(conditions.New(client.Resources()).ResourceMatch(&daemon, func(object k8s.Object) bool {
				d := object.(*appsv1.DaemonSet)
				return d.Status.NumberAvailable == int32(len(nodes.Items))
			}), wait.WithTimeout(time.Minute*1))
			if err != nil {
				support.DumpLogs(ctx, t, clientset.CoreV1(), &daemon)
				t.Fatal(err)
			}
			t.Logf("daemon set is available: %d", daemon.Status.NumberReady)
			return ctx
		}).Assess("open command shows up with collect", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		client, err := cfg.NewClient()
		if err != nil {
			t.Fatal(err)
		}
		clientset, err := kubernetes.NewForConfig(client.RESTConfig())
		if err != nil {
			t.Fatal(err)
		}

		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "modified-pod", Namespace: "default"},
		}

		support.ApplyYaml(ctx, client.Resources(), "../../examples/modified_pod.yaml")
		err = wait.For(conditions.New(client.Resources()).PodReady(&pod), wait.WithTimeout(time.Minute*1))
		if err != nil {
			support.DumpLogs(ctx, t, clientset.CoreV1(), &pod)
			t.Fatal(err)
		}

		time.Sleep(time.Second * 30)

		tmpFileName := path.Join(t.TempDir(), "events.json")

		cmd, _, _, _ := support.NewTestRootCmd(cfg, nil)
		cmd.SetArgs([]string{"collect", "-o", tmpFileName})

		if err := cmd.ExecuteContext(ctx); err != nil {
			t.Fatal(err)
		}

		if !helperFindEvent(t, tmpFileName) {
			t.Fatal("did not find open file event")
		}

		return ctx
	}).Assess("containerdbg debug captures the open", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		namespace := "debug-" + rand.RandStringRunes(10)
		return helperTestOpenIsRecorded(t, ctx, cfg, namespace, "-n", namespace, "ko.local/test-openfile")
	}).Assess("containerdbg debug captures the open for yaml file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		namespace := "debug-" + rand.RandStringRunes(10)
		return helperTestOpenIsRecorded(t, ctx, cfg, namespace, "-n", namespace, "-f", "../../examples/normal_deployment.yaml")
	}).Feature()

	testenv.Test(t, systemInstallation)
}
