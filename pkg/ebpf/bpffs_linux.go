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
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

const defaultMapFallback = "/run/containerdbg/bpffs"

var mapRoot = "/sys/fs/bpf"

// Code mostly copied from cilium

func checkOrMountFS() error {
	if err := checkOrMountDefaultLocations(); err != nil {
		return err
	}

	return nil
}

func setMapRoot(bpfRoot string) {
	mapRoot = bpfRoot
}

func checkOrMountDefaultLocations() error {
	mounted, bpffsInstance, err := isMountFS(unix.BPF_FS_MAGIC, mapRoot)
	if err != nil {
		return err
	}

	if !mounted {
		if err := mountFS(); err != nil {
			return err
		}
	}

	if !bpffsInstance {
		setMapRoot(defaultMapFallback)

		cMounted, cBpffsInstance, err := isMountFS(unix.BPF_FS_MAGIC, mapRoot)
		if err != nil {
			return err
		}
		if !cMounted {
			if err := mountFS(); err != nil {
				return err
			}
		} else if !cBpffsInstance {
			// TODO fatal error
		}
	}

	// Detect mounted BPF filesystem

	return nil
}

func mountFS() error {
	mapRootStat, err := os.Stat(mapRoot)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(mapRoot, 0755); err != nil {
				return fmt.Errorf("unable to create bpf mount directory: %s", err)
			}
		} else {
			return fmt.Errorf("failed to stat the mount path %s: %s", mapRoot, err)
		}
	} else if !mapRootStat.IsDir() {
		return fmt.Errorf("%s is a file which is not a directory", mapRoot)
	}

	if err := unix.Mount(mapRoot, mapRoot, "bpf", 0, ""); err != nil {
		return fmt.Errorf("failed to mount %s: %s", mapRoot, err)
	}
	return nil
}

// isMountFS returns two boolean values, checking
// - whether the path is a mount point
// - if yes whether its filesystem type is mntType
func isMountFS(mntType int64, path string) (bool, bool, error) {
	var st, pst unix.Stat_t

	err := unix.Lstat(path, &st)
	if err != nil {
		if errors.Is(err, unix.ENOENT) {
			return false, false, nil
		}
		return false, false, &os.PathError{Op: "lstat", Path: path, Err: err}
	}

	parent := filepath.Dir(path)
	err = unix.Lstat(parent, &pst)
	if err != nil {
		return false, false, &os.PathError{Op: "lstat", Path: parent, Err: err}
	}
	if st.Dev == pst.Dev {
		// parent has the same dev -- not a mount point
		return false, false, nil
	}

	fst := unix.Statfs_t{}
	err = unix.Statfs(path, &fst)
	if err != nil {
		return true, false, &os.PathError{Op: "statfs", Path: path, Err: err}
	}

	return true, fst.Type == mntType, nil
}
