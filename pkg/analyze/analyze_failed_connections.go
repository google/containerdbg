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

func (analyzer *connectionAnalyzer) handleNetEvent(source *proto.SourceId, net *proto.Event_NetworkEvent) bool {

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
		analyzer.connections[index] = connection
	}

	if net.GetEventType() == proto.Event_NetworkEvent_FAILED_CONNECTION {
		if source.GetType() == "host" {
			// if we got it on the host we need to check if there was a corresponding connection starting from the container
			if _, ok := analyzer.initiatedConnection[index]; ok {
				analyzer.failedConnections[index] = connection
			}
		} else {
			analyzer.failedConnections[index] = connection
		}
		return true
	} else if source.GetType() != "host" && net.GetEventType() == proto.Event_NetworkEvent_INITIATE_CONNECTION {
		analyzer.initiatedConnection[index] = connection
		return true
	}

	return false
}

func (analyzer *connectionAnalyzer) cleanDnsRequest(query string) string {
	query = strings.TrimSuffix(query, ".")
	if analyzer.searchPaths == nil {
		return query
	}
	result := query
	for _, p := range analyzer.searchPaths {
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

func (analyzer *connectionAnalyzer) handleDnsEvent(dns *proto.Event_DnsQueryEvent) bool {
	cleanReq := analyzer.cleanDnsRequest(dns.GetQuery())
	if dns.GetError() != nil {
		if _, ok := analyzer.successfulDns[cleanReq]; !ok {
			analyzer.failedDns[cleanReq] = dns.GetError()
			return true
		}
	}
	analyzer.ipToDns[dns.GetIp()] = cleanReq
	analyzer.successfulDns[cleanReq] = dns.GetIp()
	delete(analyzer.failedDns, cleanReq)
	return true
}

func (analyzer *connectionAnalyzer) isEventRelevant(event *proto.Event) bool {
	// TODO: compare the host id
	if event.GetSource().GetType() == "host" {
		return true
	}
	if event.GetSource().GetType() == "container" {
		if event.GetSource().GetId() == analyzer.sourceId.GetId() {
			return true
		}
	} else if event.GetSource().GetType() == "pod" {
		if event.GetSource().GetId() == analyzer.sourceId.GetParent() {
			return true
		}
	}

	return false

}

func (analyzer *connectionAnalyzer) handleEvent(event *proto.Event) bool {
	if !analyzer.isEventRelevant(event) {
		return true
	}
	net := event.GetNetwork()
	if net != nil {
		return analyzer.handleNetEvent(event.GetSource(), net)
	}
	dns := event.GetDnsQuery()
	if dns != nil {
		return analyzer.handleDnsEvent(dns)
	}

	return false
}

func (analyzer *connectionAnalyzer) updateSummary(summary *proto.ContainerAnalysisSummary) {
	failedConn := []*proto.ContainerAnalysisSummary_Connection{}
	for _, conn := range analyzer.failedConnections {
		if dns, ok := analyzer.ipToDns[conn.GetTargetIp()]; ok {
			conn.TargetFqdn = dns
		}
		failedConn = append(failedConn, conn)
	}

	staticIPs := []*proto.ContainerAnalysisSummary_Connection{}
	for _, conn := range analyzer.connections {
		if conn.GetTargetIp() == "0.0.0.0" {
			continue
		}
		if _, ok := analyzer.ipToDns[conn.GetTargetIp()]; !ok {
			staticIPs = append(staticIPs, conn)
		}
	}

	failedDns := []*proto.ContainerAnalysisSummary_DnsFailure{}
	for query, dnsError := range analyzer.failedDns {
		failedDns = append(failedDns, &proto.ContainerAnalysisSummary_DnsFailure{
			Query: query,
			Error: dnsError,
		})
	}

	summary.DnsFailures = failedDns

	summary.ConnectionFailures = failedConn

	summary.StaticIps = staticIPs
}
