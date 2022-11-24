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

package linux

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

func GetMntNsId(pid int32) (uint64, error) {
	info, err := os.Stat(fmt.Sprintf("/proc/%d/ns/mnt", pid))
	if err != nil {
		return 0, err
	}

	unixStat := info.Sys().(*syscall.Stat_t)

	return unixStat.Ino, nil
}

func IsFileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return false // might or might not exist, assume it doesn't
	}
}
