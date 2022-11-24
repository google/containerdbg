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

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/google/containerdbg/pkg/events"
	"github.com/google/containerdbg/pkg/rand"
	"github.com/google/containerdbg/test/support"
)

func helperFindEvent(t *testing.T, filename string, containers int) bool {

	t.Helper()
	// Very damp scanning for the failed open file
	recordsFile, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	reader := events.NewEventReader(recordsFile)
	containersFound := map[string]any{}
	for event, err := reader.Read(); err == nil; event, err = reader.Read() {
		t.Logf("event line %+v", event)
		syscall := event.GetSyscall()
		if syscall == nil {
			continue
		}
		if syscall.GetOpen().GetPath() == "/doesnotexists" {
			containersFound[event.Source.Id] = nil
			if len(containersFound) >= containers {
				return true
			}
		}
	}

	return false
}

func helperTestOpenIsRecorded(t *testing.T, ctx context.Context, cfg *envconf.Config, namespace string, containers int, debugParams ...string) context.Context {
	tmpFileName := path.Join(t.TempDir(), "events.json")

	support.RunContainerDebug(t, ctx, cfg, tmpFileName, namespace, debugParams...)

	if !helperFindEvent(t, tmpFileName, containers) {
		t.Fatal("did not find open file event")
	}

	return ctx
}

func TestFullE2EFlow(t *testing.T) {
	systemInstallation := features.New("open is recorded").
		WithLabel("type", "e2e").
		Assess("containerdbg debug captures the open", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			namespace := "debug-" + rand.RandStringRunes(10)
			return helperTestOpenIsRecorded(t, ctx, cfg, namespace, 1, "-n", namespace, "-i", "ko.local/test-openfile")
		}).Assess("containerdbg debug captures the open for yaml file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		namespace := "debug-" + rand.RandStringRunes(10)
		return helperTestOpenIsRecorded(t, ctx, cfg, namespace, 2, "-n", namespace, "-f", "../../examples/normal_deployment.yaml")
	}).Feature()

	testenv.Test(t, systemInstallation)
}
