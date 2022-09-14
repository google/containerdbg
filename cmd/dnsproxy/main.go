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

package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"google.golang.org/grpc"

	"github.com/google/containerdbg/pkg/consts"
	"github.com/google/containerdbg/pkg/daemon"
	"github.com/google/containerdbg/pkg/dnsproxy"
	"github.com/google/containerdbg/pkg/linux"
	"github.com/google/containerdbg/proto"
)

const dnsHost = ":53"

var searchRegexp = regexp.MustCompile(`search\s+(.*)`)

func retrieveSearchList() ([]string, error) {
	resolveConfContent, err := os.ReadFile(dnsproxy.ResolveConfPath)
	if err != nil {
		return nil, err
	}

	matches := searchRegexp.FindStringSubmatch(string(resolveConfContent))
	if len(matches) < 2 {
		// there was no search expression
		fmt.Printf("failed to find any submatches: %+v in\n%s\n", matches, resolveConfContent)
		return nil, nil
	}

	return strings.Fields(matches[1]), nil
}

func CreateDaemonServiceClient(sourceId *proto.SourceId) proto.NodeDaemonServiceClient {
	var client proto.NodeDaemonServiceClient
	daemonProxy := os.Getenv(consts.DaemonProxyEnv)
	if daemonProxy == "" {
		sharedDir, err := daemon.GetAndPrepareSharedDir()
		if err != nil {
			panic(err)
		}

		client, err = daemon.CreateNodeDaemonClient(sharedDir)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		fmt.Printf("connecting to grpc at %s\n", daemonProxy)
		nsId, err := linux.GetNetNsId(int32(os.Getpid()))
		if err != nil {
			panic(err)
		}
		conn, err := grpc.Dial(daemonProxy+":"+strconv.Itoa(consts.DaemonProxyPort), grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		client = proto.NewNodeDaemonServiceClient(conn)
		_, err = client.Monitor(context.Background(), &proto.MonitorPodRequest{
			Id:    sourceId,
			Netns: nsId,
		})
		if err != nil {
			panic(err)
		}
	}

	return client
}

func main() {
	fmt.Printf("Starting server on %s\n", dnsHost)
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	sourceId := &proto.SourceId{
		Type: "container",
		Id:   hostname,
	}
	client := CreateDaemonServiceClient(sourceId)
	search, err := retrieveSearchList()
	if err != nil {
		panic(err)

	}
	if _, err := client.ReportDnsSearchValues(context.Background(), &proto.ReportDnsSearchValuesRequest{
		Search: search,
		Id:     sourceId,
	}); err != nil {
		panic(err)
	}

	proxy, err := dnsproxy.NewDnsProxy(client, sourceId)
	if err != nil {
		fmt.Printf("Failed to install proxy: %+v\n", err)
		os.Exit(1)
	}
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m, err := proxy.GetResponse(r)
		if err != nil {
			fmt.Printf("Failed getting response: %+v\n", err)
		}
		m.SetReply(r)
		w.WriteMsg(m)
	})

	server := dns.Server{Addr: dnsHost, Net: "udp"}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed starting server: %+v\n", err)
	}
	os.Exit(0)
}
