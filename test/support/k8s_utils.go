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
	"testing"

	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	"k8s.io/client-go/kubernetes"
)


func NewK8sClient(ctx context.Context, t *testing.T, cfg *envconf.Config) (klient.Client, *kubernetes.Clientset) {
	t.Helper()
	client, err := cfg.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(client.RESTConfig())
	if err != nil {
		t.Fatal(err)
	}
	return client, clientset
}
