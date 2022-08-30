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

package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/go-logr/logr"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

type NodeDaemonServiceServer struct {
	proto.UnimplementedNodeDaemonServiceServer
	Manager *events.EventsSourceManager
}

var _ proto.NodeDaemonServiceServer = &NodeDaemonServiceServer{}

func UnixAuthFromContext(ctx context.Context) *UnixAuthInfo {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil
	}
	info, ok := p.AuthInfo.(*UnixAuthInfo)
	if !ok {
		return nil
	}

	return info
}

func getNsId(pid int32) (uint64, error) {
	info, err := os.Stat(fmt.Sprintf("/proc/%d/ns/net", pid))
	if err != nil {
		return 0, err
	}

	unixStat := info.Sys().(*syscall.Stat_t)

	return unixStat.Ino, nil
}

func (srv *NodeDaemonServiceServer) Monitor(ctx context.Context, request *proto.MonitorPodRequest) (*proto.MonitorPodResponse, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	authInfo := UnixAuthFromContext(ctx)
	log.Info("Got request from", "authInfo", authInfo)
	pid := authInfo.creds.Pid

	nsId, err := getNsId(pid)
	if err != nil {
		return nil, err
	}

	// TODO change RegisterContainer to accept uint64
	srv.Manager.RegisterContainer(uint32(nsId), request.GetId())

	return &proto.MonitorPodResponse{}, nil
}
