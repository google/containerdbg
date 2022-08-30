// Copyright 2022 Google LLC All Rights Reserved.
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

package ebpf

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/cilium/ebpf/perf"
	"github.com/go-logr/logr"
	"google.golang.org/protobuf/types/known/timestamppb"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events/api"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

//go:generate go run -mod=mod github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS openfiles file_opens_filter.c -- -I./headers
//go:generate ../../hack/add_license.sh ./openfiles_bpfeb.go
//go:generate ../../hack/add_license.sh ./openfiles_bpfel.go

type OpenFilesFilter struct {
	log    logr.Logger
	objs   openfilesObjects
	reader *perfReader

	tracepoints tracepointCollection
}

var _ api.EventsSource = &OpenFilesFilter{}

type OpenFileEvent struct {
	NetNs uint32
	PID   uint32
	TS    uint64
	Comm  [16]byte
	Path  [200]byte
	Ret   int32
}

var openFileEventSize = binary.Size(&OpenFileEvent{})

func (o *OpenFilesFilter) Load(log logr.Logger) (err error) {
	o.log = log
	if err := loadOpenfilesObjects(&o.objs, GetManagerInstance().GetDefaultCollectionOptions()); err != nil {
		return err
	}
	defer closeOnError(&o.objs, err)

	o.tracepoints = tracepointCollection{
		Tracepoints: []tracepoint{
			{
				Group:   "syscalls",
				Name:    "sys_exit_open",
				Program: o.objs.SysExitOpen,
			},
			{
				Group:   "syscalls",
				Name:    "sys_exit_openat",
				Program: o.objs.SysExitOpen,
			},
			{
				Group:           "syscalls",
				Name:            "sys_exit_openat2",
				Program:         o.objs.SysExitOpen,
				IgnoreNotExists: true,
			},
			{
				Group:   "syscalls",
				Name:    "sys_enter_open",
				Program: o.objs.SysEnterOpen,
			},
			{
				Group:   "syscalls",
				Name:    "sys_enter_openat",
				Program: o.objs.SysEnterOpenat,
			},
			{
				Group:           "syscalls",
				Name:            "sys_enter_openat2",
				Program:         o.objs.SysEnterOpenat2,
				IgnoreNotExists: true,
			},
		},
	}

	err = o.tracepoints.Load()
	if err != nil {
		return err
	}

	rd, err := perf.NewReader(o.objs.Pb, os.Getpagesize()*1024)
	o.reader = NewPerfReader(o.log, rd, func(sample []byte) (*proto.Event, error) {
		if len(sample) < openFileEventSize {
			return nil, fmt.Errorf("sample is to small: %s", sample)
		}
		// using binary.Read is too slow due to the usage of reflection.
		// Directly using the correct endianity parsers is faster although not as clean
		event := OpenFileEvent{
			NetNs: binary.LittleEndian.Uint32(sample[:4]),
			PID:   binary.LittleEndian.Uint32(sample[4:8]),
			TS:    binary.LittleEndian.Uint64(sample[8:16]),
			Ret:   int32(binary.LittleEndian.Uint32(sample[232:236])),
		}
		copy(event.Comm[:], sample[16:32])
		copy(event.Path[:], sample[32:232])
		path := byteSlice2String(event.Path[:])
		path = cleanAfterNull(path)
		comm := byteSlice2String(event.Comm[:])
		comm = cleanAfterNull(comm)
		outputEvent := proto.Event{}
		outputEvent.Timestamp = timestamppb.New(time.Unix(0, int64(event.TS)))
		outputEvent.Source = GetManagerInstance().GetId(event.NetNs)
		outputEvent.EventType = &proto.Event_Syscall{
			Syscall: &proto.Event_SyscallEvent{
				Comm:    comm,
				RetCode: int64(event.Ret),
				Syscall: &proto.Event_SyscallEvent_Open{
					Open: &proto.Event_SyscallEvent_OpenSyscall{
						Path: path,
					},
				},
			},
		}

		return &outputEvent, nil

	})

	o.reader.Start()

	return nil
}

func (o *OpenFilesFilter) Events() <-chan *proto.Event {
	return o.reader.Events()
}

func (o *OpenFilesFilter) Close() error {
	o.reader.Stop()
	o.tracepoints.Unload()
	return o.objs.Close()
}
