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
	"fmt"
	"path/filepath"
	"sort"
)

type EventFiles []string

func (a EventFiles) Len() int      { return len(a) }
func (a EventFiles) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a EventFiles) Less(i, j int) bool {
	x, _ := getEventFileIndex(a[i]) // on error x will be 0
	y, _ := getEventFileIndex(a[j])
	return x < y
}

func getEventFileIndex(filename string) (int, error) {
	var x int
	_, err := fmt.Sscanf(filename, EventsFilePath+".%d", &x)
	return x, err
}

func sortEventFilesList(files []string) {
	sort.Sort(EventFiles(files))
}

func getEventFilesList() ([]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s.*", EventsFilePath))
	if err != nil {
		return nil, err
	}
	sortEventFilesList(files)

	return files, nil
}

func findLastFileIndex() (int, error) {
	files, err := getEventFilesList()
	if err != nil {
		return 0, err
	}

	if len(files) == 0 {
		return 0, nil
	}

	for i := 1; i <= len(files); i++ {
		var val int
		val, err = getEventFileIndex(files[len(files)-i])
		if err == nil {
			return val, nil
		}
	}

	return 0, err
}
