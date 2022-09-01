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

package ebpf

import (
	"encoding/binary"
	"os"
	"time"

	"github.com/cilium/ebpf/perf"
	"github.com/go-logr/logr"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/types/known/timestamppb"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events/api"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

//go:generate go run -mod=mod github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS renamelink renamelink_filter.c -- -I./headers
//go:generate ../../hack/add_license.sh ./renamelink_bpfeb.go
//go:generate ../../hack/add_license.sh ./renamelink_bpfel.go

type RenameLinkFilter struct {
	log    logr.Logger
	objs   renamelinkObjects
	reader *perfReader

	tracepoints tracepointCollection
}

var _ api.EventsSource = &RenameLinkFilter{}

type RenameEvent struct {
	NetNs   uint32
	PID     uint32
	TS      uint64
	Syscall uint32
	Comm    [16]byte
	Oldname [200]byte
	Newname [200]byte
	Ret     int32
}

func (o *RenameLinkFilter) Load(log logr.Logger) (err error) {
	o.log = log
	if err := loadRenamelinkObjects(&o.objs, GetManagerInstance().GetDefaultCollectionOptions()); err != nil {
		return err
	}
	defer closeOnError(&o.objs, err)

	o.tracepoints = tracepointCollection{
		Tracepoints: []tracepoint{
			{
				Group:   "syscalls",
				Name:    "sys_exit_rename",
				Program: o.objs.SysExitRename,
			},
			{
				Group:   "syscalls",
				Name:    "sys_exit_renameat",
				Program: o.objs.SysExitRename,
			},
			{
				Group:   "syscalls",
				Name:    "sys_exit_link",
				Program: o.objs.SysExitRename,
			},
			{
				Group:   "syscalls",
				Name:    "sys_exit_linkat",
				Program: o.objs.SysExitRename,
			},
			{
				Group:   "syscalls",
				Name:    "sys_enter_rename",
				Program: o.objs.SysEnterRename,
			},
			{
				Group:   "syscalls",
				Name:    "sys_enter_renameat",
				Program: o.objs.SysEnterRenameat,
			},
			{
				Group:   "syscalls",
				Name:    "sys_enter_link",
				Program: o.objs.SysEnterLink,
			},
			{
				Group:   "syscalls",
				Name:    "sys_enter_linkat",
				Program: o.objs.SysEnterLinkat,
			},
		},
	}

	err = o.tracepoints.Load()
	if err != nil {
		return err
	}

	rd, err := perf.NewReader(o.objs.Pb, os.Getpagesize()*2048)
	o.reader = NewPerfReader(o.log, rd, func(sample []byte) (*proto.Event, error) {

		event := RenameEvent{
			NetNs:   binary.LittleEndian.Uint32(sample[:4]),
			PID:     binary.LittleEndian.Uint32(sample[4:8]),
			TS:      binary.LittleEndian.Uint64(sample[8:16]),
			Syscall: binary.LittleEndian.Uint32(sample[16:20]),
			Ret:     int32(binary.LittleEndian.Uint32(sample[436:440])),
		}
		copy(event.Comm[:], sample[20:36])
		copy(event.Oldname[:], sample[36:236])
		copy(event.Newname[:], sample[236:436])

		comm := cleanAfterNull(byteSlice2String(event.Comm[:]))
		oldname := cleanAfterNull(byteSlice2String(event.Oldname[:]))
		newname := cleanAfterNull(byteSlice2String(event.Newname[:]))

		outputEvent := proto.Event{}
		outputEvent.Source = GetManagerInstance().GetId(event.NetNs)
		outputEvent.Timestamp = timestamppb.New(time.Unix(0, int64(event.TS)))
		outputEvent.EventType = &proto.Event_Syscall{
			Syscall: &proto.Event_SyscallEvent{
				Comm:    comm,
				RetCode: int64(event.Ret),
			},
		}
		if event.Syscall == unix.SYS_RENAME || event.Syscall == unix.SYS_RENAMEAT || event.Syscall == unix.SYS_RENAMEAT2 {
			outputEvent.GetSyscall().Syscall = &proto.Event_SyscallEvent_Rename{
				Rename: &proto.Event_SyscallEvent_RenameSyscall{
					Oldname: oldname,
					Newname: newname,
				},
			}
		} else {
			outputEvent.GetSyscall().Syscall = &proto.Event_SyscallEvent_Link{
				Link: &proto.Event_SyscallEvent_LinkSyscall{
					Oldname: oldname,
					Newname: newname,
				},
			}
		}
		return &outputEvent, nil
	})
	if err != nil {
		return
	}

	o.reader.Start()

	return nil
}

func (o *RenameLinkFilter) Events() <-chan *proto.Event {
	return o.reader.Events()
}

func (o *RenameLinkFilter) Close() error {
	if o.reader != nil {
		o.reader.Stop()
	}
	o.tracepoints.Unload()
	return o.objs.Close()
}
