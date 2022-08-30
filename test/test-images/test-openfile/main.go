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
	"os"
	"os/signal"
)

func main() {
	int := make(chan os.Signal, 1)
	signal.Notify(int, os.Interrupt)
	// Just trying to open non existent file and check wether it is recorded by our eBPF filter
	if len(os.Args) >= 2 {
		os.Open(os.Args[1])
	} else {
		os.Open("/doesnotexists")
	}
	<-int
}
