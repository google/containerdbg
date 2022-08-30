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

package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/consts"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

func main() {
	if err := xmain(); err != nil {
		panic(err)
	}
}

func Dialer(addr string, timeout time.Duration) (net.Conn, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	return net.DialTimeout(url.Scheme, url.Path, timeout)
}

func createClient(serverAddr string) (proto.NodeDaemonServiceClient, error) {

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithDialer(Dialer))
	if err != nil {
		return nil, err
	}

	return proto.NewNodeDaemonServiceClient(conn), nil
}

func xmain() error {
	sharedDir := os.Getenv(consts.SharedDirectoryEnv)

	if sharedDir == "" {
		return fmt.Errorf("SHARED_DIRECTORY env was not provided")
	}

	if err := os.MkdirAll(sharedDir, 0770); err != nil {
		return err
	}

	// passthrough prefix - see: https://github.com/grpc/grpc-go/issues/1911
	// and https://github.com/grpc/grpc-go/issues/1846#issuecomment-362634790
	url := "passthrough:///unix://" + filepath.Join(sharedDir, consts.NodeDaemonSocketName)

	// TODO: Send grpc to node monitoring daemon to notify about container creation
	client, err := createClient(url)
	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	_, err = client.Monitor(context.Background(), &proto.MonitorPodRequest{
		Id: &proto.SourceId{
			Type: "container",
			Id:   hostname,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("trying to execute %+v", os.Args[1:])
	if err := syscall.Exec(os.Args[1], os.Args[1:], os.Environ()); err != nil {
		panic(err)
	}

	return nil
}
