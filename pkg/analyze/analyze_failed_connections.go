// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analyze

import (
	"strings"

	"github.com/google/containerdbg/proto"
)

type connectionTuple struct {
	ip      string
	srcIp   string
	dstPort int32
}

type connectionAnalyzer struct {
	connections         map[connectionTuple]*proto.ContainerAnalysisSummary_Connection
	initiatedConnection map[connectionTuple]*proto.ContainerAnalysisSummary_Connection
	failedConnections   map[connectionTuple]*proto.ContainerAnalysisSummary_Connection
	ipToDns             map[string]string
	failedDns           map[string]*proto.DnsQueryError
	successfulDns       map[string]string
	searchPaths         []string
	sourceId            *proto.SourceId
}

func newConnectionAnalyzer(searchPaths []string, sourceId *proto.SourceId) *connectionAnalyzer {
	return &connectionAnalyzer{
		connections:         make(map[connectionTuple]*proto.ContainerAnalysisSummary_Connection),
		failedConnections:   make(map[connectionTuple]*proto.ContainerAnalysisSummary_Connection),
		initiatedConnection: make(map[connectionTuple]*proto.ContainerAnalysisSummary_Connection),
		ipToDns:             make(map[string]string),
		failedDns:           make(map[string]*proto.DnsQueryError),
		successfulDns:       make(map[string]string),
		searchPaths:         searchPaths,
		sourceId:            sourceId,
	}
}

var _ analyzer = &connectionAnalyzer{}

func (an *connectionAnalyzer) handleNetEvent(source *proto.SourceId, net *proto.Event_NetworkEvent) bool {

	connection := &proto.ContainerAnalysisSummary_Connection{
		TargetIp: net.GetDstAddr(),
		Port:     net.GetDstPort(),
	}

	index := connectionTuple{
		ip:      net.GetDstAddr(),
		srcIp:   net.GetSrcAddr(),
		dstPort: net.GetDstPort(),
	}

	if source.GetType() != "host" {
		an.connections[index] = connection
	}

	if net.GetEventType() == proto.Event_NetworkEvent_FAILED_CONNECTION {
		if source.GetType() == "host" {
			// if we got it on the host we need to check if there was a corresponding connection starting from the container
			if _, ok := an.initiatedConnection[index]; ok {
				an.failedConnections[index] = connection
			}
		} else {
			an.failedConnections[index] = connection
		}
		return true
	} else if source.GetType() != "host" && net.GetEventType() == proto.Event_NetworkEvent_INITIATE_CONNECTION {
		an.initiatedConnection[index] = connection
		return true
	}

	return false
}

func (an *connectionAnalyzer) cleanDnsRequest(query string) string {
	query = strings.TrimSuffix(query, ".")
	if an.searchPaths == nil {
		return query
	}
	result := query
	for _, p := range an.searchPaths {
		trimmed := strings.TrimSuffix(query, p)
		if trimmed != query {
			if len(trimmed) < len(result) {
				result = trimmed
			}
		}
	}

	result = strings.TrimSuffix(result, ".")

	return result
}

func (an *connectionAnalyzer) handleDnsEvent(dns *proto.Event_DnsQueryEvent) bool {
	cleanReq := an.cleanDnsRequest(dns.GetQuery())
	if dns.GetError() != nil {
		if _, ok := an.successfulDns[cleanReq]; !ok {
			an.failedDns[cleanReq] = dns.GetError()
			return true
		}
	}
	an.ipToDns[dns.GetIp()] = cleanReq
	an.successfulDns[cleanReq] = dns.GetIp()
	delete(an.failedDns, cleanReq)
	return true
}

func (an *connectionAnalyzer) isEventRelevant(event *proto.Event) bool {
	// TODO: compare the host id
	if event.GetSource().GetType() == "host" {
		return true
	}
	if event.GetSource().GetType() == "container" {
		if event.GetSource().GetId() == an.sourceId.GetId() {
			return true
		}
	} else if event.GetSource().GetType() == "pod" {
		if event.GetSource().GetId() == an.sourceId.GetParent() {
			return true
		}
	}

	return false

}

func (an *connectionAnalyzer) handleEvent(event *proto.Event) bool {
	if !an.isEventRelevant(event) {
		return true
	}
	net := event.GetNetwork()
	if net != nil {
		return an.handleNetEvent(event.GetSource(), net)
	}
	dns := event.GetDnsQuery()
	if dns != nil {
		return an.handleDnsEvent(dns)
	}

	return false
}

func (an *connectionAnalyzer) updateSummary(summary *proto.ContainerAnalysisSummary) {
	failedConn := []*proto.ContainerAnalysisSummary_Connection{}
	for _, conn := range an.failedConnections {
		if dns, ok := an.ipToDns[conn.GetTargetIp()]; ok {
			conn.TargetFqdn = dns
		}
		failedConn = append(failedConn, conn)
	}

	staticIPs := []*proto.ContainerAnalysisSummary_Connection{}
	for _, conn := range an.connections {
		if conn.GetTargetIp() == "0.0.0.0" {
			continue
		}
		if _, ok := an.ipToDns[conn.GetTargetIp()]; !ok {
			staticIPs = append(staticIPs, conn)
		}
	}

	failedDns := []*proto.ContainerAnalysisSummary_DnsFailure{}
	for query, dnsError := range an.failedDns {
		failedDns = append(failedDns, &proto.ContainerAnalysisSummary_DnsFailure{
			Query: query,
			Error: dnsError,
		})
	}

	summary.DnsFailures = failedDns

	summary.ConnectionFailures = failedConn

	summary.StaticIps = staticIPs
}
