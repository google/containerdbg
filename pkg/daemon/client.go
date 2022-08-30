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

package daemon

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/consts"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

func dialer(addr string, timeout time.Duration) (net.Conn, error) {
	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	return net.DialTimeout(url.Scheme, url.Path, timeout)
}

func CreateClient(serverAddr string) (proto.NodeDaemonServiceClient, error) {

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithDialer(dialer))
	if err != nil {
		return nil, err
	}

	return proto.NewNodeDaemonServiceClient(conn), nil
}

func GetAndPrepareSharedDir() (string, error) {

	sharedDir := os.Getenv(consts.SharedDirectoryEnv)

	if sharedDir == "" {
		return "", fmt.Errorf("SHARED_DIRECTORY env was not provided")
	}

	if err := os.MkdirAll(sharedDir, 0770); err != nil {
		return "", err
	}

	return sharedDir, nil
}

func GetServerAddr(sharedDir string) string {
	return "passthrough:///unix://" + filepath.Join(sharedDir, consts.NodeDaemonSocketName)
}

func CreateNodeDaemonClient(sharedDir string) (proto.NodeDaemonServiceClient, error) {

	// passthrough prefix - see: https://github.com/grpc/grpc-go/issues/1911
	// and https://github.com/grpc/grpc-go/issues/1846#issuecomment-362634790
	url := GetServerAddr(sharedDir)

	// TODO: Send grpc to node monitoring daemon to notify about container creation
	client, err := CreateClient(url)
	if err != nil {
		return nil, err
	}

	return client, nil
}
