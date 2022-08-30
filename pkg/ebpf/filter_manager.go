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
	"os"
	"path"
	"sync"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS pinnedmaps pinned_maps.c -- -I./headers
//go:generate ../../hack/add_license.sh ./pinnedmaps_bpfeb.go
//go:generate ../../hack/add_license.sh ./pinnedmaps_bpfel.go

type filtersManager struct {
	pinnedMaps pinnedmapsObjects
	idMap      sync.Map
}

var filtersMgr *filtersManager = nil

func GetManagerInstance() *filtersManager {
	if filtersMgr == nil {
		filtersMgr = &filtersManager{}
	}
	return filtersMgr
}

func (*filtersManager) PinnedMapsPath() string {
	return path.Join(mapRoot, "containerdbg")
}

func (mgr *filtersManager) GetDefaultCollectionOptions() *ebpf.CollectionOptions {
	return &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: mgr.PinnedMapsPath(),
		},
	}
}

func (mgr *filtersManager) Init() error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return err
	}
	if err := checkOrMountFS(); err != nil {
		return err
	}
	if err := os.MkdirAll(mgr.PinnedMapsPath(), os.ModePerm); err != nil {
		return err
	}
	err := loadPinnedmapsObjects(&mgr.pinnedMaps, &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: mgr.PinnedMapsPath(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (mgr *filtersManager) RegisterContainer(nsId uint32, id *proto.SourceId) error {
	mgr.idMap.Store(nsId, id)
	if err := mgr.pinnedMaps.NetNs.Put(nsId, uint8(1)); err != nil {
		mgr.idMap.Delete(nsId)
		return err
	}
	return nil
}

func (mgr *filtersManager) GetId(nsId uint32) *proto.SourceId {
	v, ok := mgr.idMap.Load(nsId)
	if !ok {
		return nil
	}
	return v.(*proto.SourceId)
}
