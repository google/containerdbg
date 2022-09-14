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
	"testing"

	pb "github.com/google/containerdbg/proto"
)

func openEqual(actual, expected *pb.Event_SyscallEvent_Open) bool {
	if actual.Open.GetPath() != expected.Open.GetPath() {
		return false
	}

	return true
}

func renameEqual(actual, expected *pb.Event_SyscallEvent_Rename) bool {
	if actual.Rename.GetOldname() != expected.Rename.GetOldname() {
		return false
	}
	if actual.Rename.GetNewname() != expected.Rename.GetNewname() {
		return false
	}

	return true
}

func linkEqual(actual, expected *pb.Event_SyscallEvent_Link) bool {
	if actual.Link.GetOldname() != expected.Link.GetOldname() {
		return false
	}
	if actual.Link.GetNewname() != expected.Link.GetNewname() {
		return false
	}

	return true
}

func syscallEquality(actual, expected *pb.Event_Syscall) bool {
	if actual.Syscall.GetComm() != expected.Syscall.GetComm() {
		return false
	}
	if actual.Syscall.GetRetCode() != expected.Syscall.GetRetCode() {
		return false
	}

	switch actual := actual.Syscall.GetSyscall().(type) {
	case *pb.Event_SyscallEvent_Open:
		expected, ok := expected.Syscall.GetSyscall().(*pb.Event_SyscallEvent_Open)
		if !ok {
			return false
		}
		return openEqual(actual, expected)
	case *pb.Event_SyscallEvent_Rename:
		expected, ok := expected.Syscall.GetSyscall().(*pb.Event_SyscallEvent_Rename)
		if !ok {
			return false
		}
		return renameEqual(actual, expected)
	case *pb.Event_SyscallEvent_Link:
		expected, ok := expected.Syscall.GetSyscall().(*pb.Event_SyscallEvent_Link)
		if !ok {
			return false
		}
		return linkEqual(actual, expected)
	default:
		return false
	}
}

func networkEquality(actual, expected *pb.Event_Network) bool {
	if actual.Network.GetComm() != expected.Network.GetComm() {
		return false
	}

	if actual.Network.GetEventType() != expected.Network.GetEventType() {
		return false
	}

	if actual.Network.GetAddrFam() != expected.Network.GetAddrFam() {
		return false
	}

	if actual.Network.GetSrcAddr() != expected.Network.GetSrcAddr() {
		return false
	}

	if actual.Network.GetSrcPort() != expected.Network.GetSrcPort() {
		return false
	}

	if actual.Network.GetDstAddr() != expected.Network.GetDstAddr() {
		return false
	}

	if actual.Network.GetDstPort() != expected.Network.GetDstPort() {
		return false
	}
	return true;
}

// Although we could use proto.Equal to check for protobuf equality
// that function uses reflection and slows down the test validation hard enough that we are not reading from the channel fast enough to clean the ring buffer, thus we have to implement the equality check directly.
func EqualProto(t *testing.T, actual, expected *pb.Event) bool {
	t.Helper()
	if actual.Source.GetType() != expected.Source.GetType() {
		return false
	}
	if actual.Source.GetId() != expected.Source.GetId() {
		return false
	}
	if !actual.Timestamp.AsTime().Equal(expected.Timestamp.AsTime()) {
		return false
	}
	switch actual := actual.EventType.(type) {
	case *pb.Event_Syscall:
		expected, ok := expected.EventType.(*pb.Event_Syscall)
		if !ok {
			return false
		}
		return syscallEquality(actual, expected)
	case *pb.Event_Network :
		expected, ok := expected.EventType.(*pb.Event_Network)
		if !ok {
			return false
		}
		return networkEquality(actual, expected)
	default:
		return false
	}
}
