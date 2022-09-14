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
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"github.com/google/containerdbg/pkg/ebpf"
	"github.com/google/containerdbg/pkg/events/api"
	"github.com/google/containerdbg/pkg/linux"
	pb "github.com/google/containerdbg/proto"
	"github.com/google/containerdbg/test/support"
)

func runBinaryWithNewNSAndAttach(t *testing.T, path string, args []string, exitChan chan<- interface{}) {

	cmd := exec.Command(path, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNET,
	}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("failed to open pipe", err)
	}
	cmd.ExtraFiles = append(cmd.ExtraFiles, w)

	if err := cmd.Start(); err != nil {
		t.Fatal("failed to run test binary", err)
	}

	netId, err := linux.GetNetNsId(int32(cmd.Process.Pid))
	if err != nil {
		t.Fatal("failed to get namespace of command", err)
	}

	t.Logf("namespace is %d\n", netId)

	ebpf.GetManagerInstance().RegisterContainer(uint32(netId), &pb.SourceId{Type: "container", Id: "self"})

	t.Logf("Registered container\n")

	data := make([]byte, 5)
	_, err = r.Read(data)
	if err != nil {
		t.Fatal("failed to read from file", err)
	}

	cmd.Process.Signal(os.Interrupt)

	go func() {
		cmd.Wait()
		close(exitChan)
	}()

}

func helperWaitForEvent(t *testing.T, events <-chan *pb.Event, expectedEvent *pb.Event, timeout time.Duration, times int) int {
	t.Helper()
	count := 0
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return count
			}
			event.Timestamp = timestamppb.New(time.Unix(0, 0))
			if support.EqualProto(t, event, expectedEvent) {
				count++
				if count >= times {
					return count
				}
			}
		case <-time.After(timeout):
			t.Logf("timeout waiting for event")
			// after 10 seconds without a new event we have failed
			return count
		}
	}

}

func helperBinaryProduceEvent(t *testing.T, filter api.EventsSource, expectedEvent *pb.Event, timeout time.Duration, times int, binary string, args []string) {
	t.Helper()

	exitChan := make(chan interface{})
	runBinaryWithNewNSAndAttach(t, binary, args, exitChan)

	doneSignal := make(chan int, 1)

	go func() {
		doneSignal <- helperWaitForEvent(t, filter.Events(), expectedEvent, timeout, times)
	}()

	select {
	case <-exitChan:
	case <-time.After(timeout):
		t.Fatalf("timeout reached while trying to wait for executable to exit")
	}

	var count int
	select {
	case <-time.After(time.Second * 200):
	case count = <-doneSignal:
	}
	if count < times {
		t.Fatalf("expected event didn't occur enough times %f%%", 100*float64(count)/float64(times))
	}
}

func TestOpenIsRecorded(t *testing.T) {
	fileFilterFeature := features.New("ebpf filters").
		WithLabel("type", "component").
		Assess("filtered open is recorded", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {

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
						Comm:    "test-binary",
						RetCode: -2,
						Syscall: &pb.Event_SyscallEvent_Open{
							Open: &pb.Event_SyscallEvent_OpenSyscall{
								Path: "/filethatdoesnotexist",
							},
						},
					},
				},
			}

			helperBinaryProduceEvent(t, &openFilesFilter, expectedObj, time.Second*20, 1, "../../out/linux_amd64/release/test-binary", []string{"/filethatdoesnotexist"})

			return ctx
		}).
		Assess("filtered rename is recorded", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {

			renameFilter := ebpf.RenameLinkFilter{}
			if err := renameFilter.Load(testr.NewWithOptions(t, testr.Options{
				Verbosity: 5,
			})); err != nil {
				t.Fatal("failed to load ebpf filter", err)
			}
			defer func() {
				t.Logf("cleaning filter")
				if err := renameFilter.Close(); err != nil {
					t.Error("failed to remove ebpf filter", err)
				}
			}()

			expectedObj := &pb.Event{
				Timestamp: timestamppb.New(time.Unix(0, 0)),
				Source:    &pb.SourceId{Type: "container", Id: "self"},
				EventType: &pb.Event_Syscall{
					Syscall: &pb.Event_SyscallEvent{
						Comm:    "test-binary",
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

			helperBinaryProduceEvent(t, &renameFilter, expectedObj, time.Second*20, 1, "../../out/linux_amd64/release/test-binary", []string{"/filethatdoesnotexist"})

			return ctx
		}).
		Assess("filtered link is recorded", func(ctx context.Context, t *testing.T, _ *envconf.Config) context.Context {

			renameFilter := ebpf.RenameLinkFilter{}
			if err := renameFilter.Load(testr.NewWithOptions(t, testr.Options{
				Verbosity: 5,
			})); err != nil {
				t.Fatal("failed to load ebpf filter", err)
			}
			defer func() {
				t.Logf("cleaning filter")
				if err := renameFilter.Close(); err != nil {
					t.Error("failed to remove ebpf filter", err)
				}
			}()

			expectedObj := &pb.Event{
				Timestamp: timestamppb.New(time.Unix(0, 0)),
				Source:    &pb.SourceId{Type: "container", Id: "self"},
				EventType: &pb.Event_Syscall{
					Syscall: &pb.Event_SyscallEvent{
						Comm:    "test-binary",
						RetCode: -2,
						Syscall: &pb.Event_SyscallEvent_Link{
							Link: &pb.Event_SyscallEvent_LinkSyscall{
								Oldname: "/filethatdoesnotexist",
								Newname: "/www",
							},
						},
					},
				},
			}

			helperBinaryProduceEvent(t, &renameFilter, expectedObj, time.Second*20, 1, "../../out/linux_amd64/release/test-binary", []string{"/filethatdoesnotexist"})

			return ctx
		}).
		Feature()

	test.Test(t, fileFilterFeature)
}
