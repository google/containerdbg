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
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/google/containerdbg/pkg/consts"
	"github.com/google/containerdbg/pkg/daemon"
	"github.com/google/containerdbg/pkg/dnsproxy"
	"github.com/google/containerdbg/proto"
)

func main() {
	if err := xmain(); err != nil {
		panic(err)
	}
}

func waitForResolvModification() {
	for {
		data, err := os.ReadFile(dnsproxy.ResolveConfPath)
		if err != nil {
			// unknown error, we will ignore in this case
			fmt.Printf("got error while trying to read %s: %s", dnsproxy.ResolveConfPath, err)
			break
		}

		if bytes.Contains(data, []byte(dnsproxy.ContainerdbgComment)) {
			return
		}
		time.Sleep(time.Second)
	}
}

func xmain() error {

	sharedDir, err := daemon.GetAndPrepareSharedDir()
	if err != nil {
		panic(err)
	}

	client, err := daemon.CreateNodeDaemonClient(sharedDir)
	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Printf("trying to execute %+v", os.Args[1:])
	path, err := exec.LookPath(os.Args[1])
	// if err != nil && !errors.Is(err, exec.ErrDot) { // restore in go 1.19
	// 	return err
	// } else {
	// 	fmt.Printf("found binary in %s", path)
	// }

	waitForResolvModification()

	_, err = client.Monitor(context.Background(), &proto.MonitorPodRequest{
		Id: &proto.SourceId{
			Type: "container",
			Id:   hostname + "-" + os.Getenv(consts.ContainerNameEnv),
		},
	})
	if err != nil {
		panic(err)
	}
	if err := syscall.Exec(path, os.Args[1:], os.Environ()); err != nil {
		panic(err)
	}

	return nil
}
