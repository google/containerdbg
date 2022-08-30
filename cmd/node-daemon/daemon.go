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
	"github.com/go-logr/logr"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/ebpf"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/linux"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

type NodeDaemonServiceServer struct {
	proto.UnimplementedNodeDaemonServiceServer
	Manager *events.EventsSourceManager
	events.DynamicSource
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

func (srv *NodeDaemonServiceServer) Monitor(ctx context.Context, request *proto.MonitorPodRequest) (*proto.MonitorPodResponse, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	nsId := request.Netns
	if nsId == 0 {
		authInfo := UnixAuthFromContext(ctx)
		log.Info("Got request from", "authInfo", authInfo)
		pid := authInfo.creds.Pid
		nsId, err = linux.GetNetNsId(pid)
		if err != nil {
			return nil, err
		}
	}

	// TODO change RegisterContainer to accept uint64
	err = srv.Manager.RegisterContainer(uint32(nsId), request.GetId())
	if err != nil {
		return nil, err
	}

	return &proto.MonitorPodResponse{}, nil
}

func (srv *NodeDaemonServiceServer) ReportDnsQuery(ctx context.Context, request *proto.ReportDnsQueryResultRequest) (*proto.ReportDnsQueryResultResponse, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	source := request.Id
	if source == nil {
		authInfo := UnixAuthFromContext(ctx)
		pid := authInfo.creds.Pid

		netns, err := linux.GetNetNsId(pid)
		if err != nil {
			return nil, err
		}
		source = ebpf.GetManagerInstance().GetId(uint32(netns))
	}

	log.Info("got dns event", "event", request)

	queryEvent := &proto.Event_DnsQueryEvent{
		Query: request.GetDnsQuery(),
	}

	event := proto.Event{
		EventType: &proto.Event_DnsQuery{
			DnsQuery: queryEvent,
		},
	}
	if request.GetError() != nil {
		queryEvent.Answer = &proto.Event_DnsQueryEvent_Error{
			Error: request.GetError(),
		}
	} else {
		queryEvent.Answer = &proto.Event_DnsQueryEvent_Ip{
			Ip: request.GetReturnedIp(),
		}
	}
	event.Source = source
	srv.DynamicSource.SendEvent(&event)
	return &proto.ReportDnsQueryResultResponse{}, nil
}

func (srv *NodeDaemonServiceServer) ReportDnsSearchValues(ctx context.Context, request *proto.ReportDnsSearchValuesRequest) (*proto.ReportDnsSearchValuesResponse, error) {
	log, err := logr.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Info("got search value post", "search", request.Search)

	searchEvent := &proto.Event_DnsSearchParametersProbe{
		Search: request.GetSearch(),
	}
	event := proto.Event{
		Source: request.Id,
		EventType: &proto.Event_DnsSearch{
			DnsSearch: searchEvent,
		},
	}
	srv.DynamicSource.SendEvent(&event)
	return &proto.ReportDnsSearchValuesResponse{}, nil
}
