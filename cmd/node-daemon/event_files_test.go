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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSortFileNames(t *testing.T) {
	list := []string{
		EventsFilePath + ".1",
		EventsFilePath + ".12",
		EventsFilePath + ".2",
		EventsFilePath + ".21",
		EventsFilePath + ".3",
	}

	sortEventFilesList(list)

	diff := cmp.Diff(list, []string{
		EventsFilePath + ".1",
		EventsFilePath + ".2",
		EventsFilePath + ".3",
		EventsFilePath + ".12",
		EventsFilePath + ".21",
	})
	if diff != "" {
		t.Fatalf("expected value is different from actual (+got, -actual)\n%s\n", diff)
	}
}
