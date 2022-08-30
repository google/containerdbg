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

package ebpf

import (
	"errors"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

type tracepointCollection struct {
	Tracepoints []tracepoint
	links       []link.Link
}

type tracepoint struct {
	Group           string
	Name            string
	Program         *ebpf.Program
	IgnoreNotExists bool
}

func (o *tracepointCollection) Load() error {

	for _, tp := range o.Tracepoints {
		linked, err := link.Tracepoint(tp.Group, tp.Name, tp.Program)
		if err != nil {
			if !tp.IgnoreNotExists || !errors.Is(err, os.ErrNotExist) {
				o.Unload()
				return err
			} else {
				continue
			}
		}
		o.links = append(o.links, linked)
	}

	return nil
}

func (o *tracepointCollection) Unload() {
	for i := len(o.links) - 1; i >= 0; i-- {
		o.links[i].Close()
	}
}
