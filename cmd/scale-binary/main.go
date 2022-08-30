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
	"os"
	"os/signal"
	"strconv"
)

func main() {
	// Just trying to open non existent file and check wether it is recorded by our eBPF filter
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	f := os.NewFile(3, "sync file")
	f.Write([]byte("registered"))

	<-c

	numberOfOpens, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	if os.Args[1] == "open" {
		for i := 0; i < numberOfOpens; i++ {
			os.Open(os.Args[3])
		}
	} else if os.Args[1] == "rename" {
		for i := 0; i < numberOfOpens; i++ {
			os.Rename(os.Args[3], "/www")
		}
	}
}
