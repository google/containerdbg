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

package component

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cilium/ebpf/rlimit"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"github.com/google/containerdbg/pkg/ebpf"
)

var test env.Environment

func TestMain(m *testing.M) {
	test = env.New()
	test.Setup(func(ctx context.Context, _ *envconf.Config) (context.Context, error) {

		if err := rlimit.RemoveMemlock(); err != nil {
			return ctx, err
		}

		if err := ebpf.GetManagerInstance().Init(); err != nil {
			return ctx, fmt.Errorf("failed to load maps: %s", err)
		}

		return ctx, nil
	})

	os.Exit(test.Run(m))
}
