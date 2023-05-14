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

package dnsproxy

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/google/containerdbg/proto"
	"github.com/miekg/dns"
)

const ResolveConfPath = "/etc/resolv.conf"
const nameserverRegexp = "nameserver\\s+.*\\n"

var nameserverRedirect = fmt.Sprintf("\nnameserver 127.0.0.1\n%s\n", ContainerdbgComment)

const kubeDns = "kube-dns.kube-system"

type DNSProxy struct {
	dnsServer    string
	client       *dns.Client
	daemonClient proto.NodeDaemonServiceClient
	sourceId     *proto.SourceId
}

func NewDnsProxy(client proto.NodeDaemonServiceClient, sourceId *proto.SourceId) (*DNSProxy, error) {
	realDnsServer, err := installResolveHook()
	if err != nil {
		return nil, err
	}
	result := DNSProxy{dnsServer: realDnsServer + ":53", daemonClient: client}
	result.client = new(dns.Client)
	result.client.Net = "udp"
	result.sourceId = sourceId
	return &result, nil
}

func installResolveHook() (string, error) {
	realDnsServerIP, err := net.LookupIP(kubeDns)
	if err != nil {
		return "", err
	}
	realDnsServer := realDnsServerIP[0].String()

	err = sedResolveConf()
	if err != nil {
		return "", nil
	}

	return realDnsServer, nil
}

func sedResolveConf() error {
	resolveConfContent, err := os.ReadFile(ResolveConfPath)
	if err != nil {
		return fmt.Errorf("Failed reading resolv.conf %w", err)
	}
	replacer := regexp.MustCompile(nameserverRegexp)
	fmt.Printf("resolv.conf before\n%s", resolveConfContent)
	resolveConfContent = replacer.ReplaceAll(resolveConfContent, []byte(nameserverRedirect))
	fmt.Printf("resolv.conf after\n%s", resolveConfContent)
	err = os.WriteFile(ResolveConfPath, resolveConfContent, 0644)
	if err != nil {
		return fmt.Errorf("Failed writing resolv conf: %w", err)
	}
	return nil
}

func (proxy *DNSProxy) GetResponse(requestMsg *dns.Msg) (*dns.Msg, error) {
	responseMsg := new(dns.Msg)
	if len(requestMsg.Question) > 0 {
		question := requestMsg.Question[0]

		answer, err := proxy.forwardRequest(&question, requestMsg)
		if err != nil {
			proxy.reportError(question, err)
			return responseMsg, err
		}
		proxy.reportAnswer(question, *answer)
		responseMsg.Answer = append(responseMsg.Answer, *answer)
	}
	return responseMsg, nil
}

func (proxy *DNSProxy) reportAnswer(question dns.Question, answer dns.RR) {
	if arec, ok := answer.(*dns.A); ok {
		proxy.daemonClient.ReportDnsQuery(context.TODO(), &proto.ReportDnsQueryResultRequest{
			DnsQuery: question.Name,
			Id:       proxy.sourceId,
			Result: &proto.ReportDnsQueryResultRequest_ReturnedIp{
				ReturnedIp: arec.A.String(),
			},
		})
		fmt.Println("    [+] Type is A")
		// Need to parse the below RR to get the IP address returned
		fmt.Println("    [+] Answer: ", arec)
	}
}

func (proxy *DNSProxy) reportError(question dns.Question, err error) {
	if _, err := proxy.daemonClient.ReportDnsQuery(context.TODO(), &proto.ReportDnsQueryResultRequest{
		DnsQuery: question.Name,
		Id:       proxy.sourceId,
		Result: &proto.ReportDnsQueryResultRequest_Error{
			Error: &proto.DnsQueryError{
				Code: proto.DnsQueryError_UNKNOWN,
			},
		},
	}); err != nil {
		fmt.Printf("failed to send query: %s", err)
	}
}

func (proxy *DNSProxy) forwardRequest(q *dns.Question, requestMsg *dns.Msg) (*dns.RR, error) {
	queryMsg := new(dns.Msg)
	requestMsg.CopyTo(queryMsg)
	queryMsg.Question = []dns.Question{*q}

	msg, err := proxy.lookup(proxy.dnsServer, queryMsg)
	if err != nil {
		return nil, err
	}

	if len(msg.Answer) > 0 {
		return &msg.Answer[0], nil
	}

	fmt.Printf("Failed to find: %s\n", q.Name)
	return nil, fmt.Errorf("not found")
}

func (proxy *DNSProxy) lookup(server string, m *dns.Msg) (*dns.Msg, error) {
	response, _, err := proxy.client.Exchange(m, proxy.dnsServer)
	if err != nil {
		return nil, err
	}

	return response, nil
}
