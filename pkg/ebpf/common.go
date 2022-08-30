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
	"io"
	"strings"
	"unsafe"
)

func closeOnError(closable io.Closer, err error) {
	if err != nil {
		closable.Close()
	}
}

func cleanAfterNull(data string) string {
	i := strings.Index(data, "\x00")
	if i >= 0 {
		return data[:i]
	}
	return data
}

// Copied from go source code https://go.dev/src/strings/builder.go#L45
// this is a faster method to convert byte slice to string without copying the data
func byteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
