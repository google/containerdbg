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
)

func main() {
	// Just trying to open non existent file and check wether it is recorded by our eBPF filter
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	f := os.NewFile(3, "sync file")
	f.Write([]byte("registered"))

	<-c

	os.Open(os.Args[1])
	os.Rename(os.Args[1], "/www")
	os.Link(os.Args[1], "/www")
}
