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
	"runtime"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/google/containerdbg/pkg/ebpf"
	pb "github.com/google/containerdbg/proto"
)

func TestStress(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	fileFilterFeature := features.New("ebpf filters stress").
		WithLabel("type", "stability").
		Assess("open no event drops", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			openFilesFilter := ebpf.OpenFilesFilter{}
			if err := openFilesFilter.Load(testr.NewWithOptions(t, testr.Options{
				Verbosity: 5,
			})); err != nil {
				t.Fatal("failed to load ebpf filter", err)
			}
			defer func() {
				t.Logf("cleaning filter")
				if err := openFilesFilter.Close(); err != nil {
					t.Error("failed to remove ebpf filter", err)
				}
			}()

			expectedObj := &pb.Event{
				Timestamp: timestamppb.New(time.Unix(0, 0)),
				Source:    &pb.SourceId{Type: "container", Id: "self"},
				EventType: &pb.Event_Syscall{
					Syscall: &pb.Event_SyscallEvent{
						Comm:    "scale-binary",
						RetCode: -2,
						Syscall: &pb.Event_SyscallEvent_Open{
							Open: &pb.Event_SyscallEvent_OpenSyscall{
								Path: "/filethatdoesnotexist",
							},
						},
					},
				},
			}

			runtime.GC()

			helperBinaryProduceEvent(t, &openFilesFilter, expectedObj, time.Minute*2, 1000000, "../../out/linux_amd64/release/scale-binary", []string{"open", "1000000", "/filethatdoesnotexist"})

			return ctx
		}).
		Assess("rename no event drops", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {
			renameFilesFilter := ebpf.RenameLinkFilter{}
			if err := renameFilesFilter.Load(testr.NewWithOptions(t, testr.Options{
				Verbosity: 5,
			})); err != nil {
				t.Fatal("failed to load ebpf filter", err)
			}
			defer func() {
				t.Logf("cleaning filter")
				if err := renameFilesFilter.Close(); err != nil {
					t.Error("failed to remove ebpf filter", err)
				}
			}()

			expectedObj := &pb.Event{
				Timestamp: timestamppb.New(time.Unix(0, 0)),
				Source:    &pb.SourceId{Type: "container", Id: "self"},
				EventType: &pb.Event_Syscall{
					Syscall: &pb.Event_SyscallEvent{
						Comm:    "scale-binary",
						RetCode: -2,
						Syscall: &pb.Event_SyscallEvent_Rename{
							Rename: &pb.Event_SyscallEvent_RenameSyscall{
								Oldname: "/filethatdoesnotexist",
								Newname: "/www",
							},
						},
					},
				},
			}

			runtime.GC()

			helperBinaryProduceEvent(t, &renameFilesFilter, expectedObj, time.Minute*2, 1000000, "../../out/linux_amd64/release/scale-binary", []string{"rename", "1000000", "/filethatdoesnotexist"})

			return ctx
		}).
		Feature()

	test.Test(t, fileFilterFeature)
}
