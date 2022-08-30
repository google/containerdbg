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

package analyze

import (
	"os"
	"path/filepath"
	"sort"

	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

type AnalyzeSummary struct {
	MissingFiles  []string
	ExdevFailures []string
	MissingLibs   []string
}

func filterOpen(filters *Filters, event *proto.Event, missingFiles map[string]any, libAnalyzer *libraryAnalyzer) bool {
	syscall := event.GetSyscall()
	if syscall == nil {
		return true
	}

	openSyscall := syscall.GetOpen()
	if openSyscall == nil {
		return true
	}

	if syscall.GetRetCode() >= 0 {
		libAnalyzer.handleFoundFile(openSyscall.GetPath())
		return false
	}

	if _, ok := filters.CommFilter[syscall.GetComm()]; ok {
		return false
	}

	for _, re := range filters.FileRegexFilter {
		if re.MatchString(openSyscall.GetPath()) {
			return false
		}
	}

	missingFiles[openSyscall.GetPath()] = nil

	return true
}

func filterExdev(event *proto.Event, failedRenames map[string]any) bool {
	syscall := event.GetSyscall()
	if syscall == nil {
		return true
	}

	renameSyscall := syscall.GetRename()
	if renameSyscall == nil {
		return true
	}

	if syscall.GetRetCode() != -18 {
		return false
	}

	if filepath.Dir(renameSyscall.Oldname) == filepath.Dir(renameSyscall.Newname) {
		failedRenames[renameSyscall.Oldname] = nil
		return false
	}

	return true
}

func Analyze(inputFilename string, filters *Filters) (*AnalyzeSummary, error) {
	if filters == nil {
		filters = defaultFilters
	}
	f, err := os.Open(inputFilename)
	if err != nil {
		return nil, err
	}

	reader := events.NewEventReader(f)

	missingFiles := map[string]any{}
	failedRenames := map[string]any{}
	libAnalyzer := newLibraryAnalyzer(filters)

	for event, err := reader.Read(); err == nil; event, err = reader.Read() {
		if !filterOpen(filters, event, missingFiles, libAnalyzer) {
			continue
		}
		if !filterExdev(event, failedRenames) {
			continue
		}
	}

	libAnalyzer.filterOutFoundLibraries(missingFiles)

	missingFilesSlice := []string{}
	missingLibs := map[string]any{}
	for fname := range missingFiles {
		if filters.IsLibrary(fname) {
			missingLibs[filepath.Base(fname)] = nil
		} else {
			missingFilesSlice = append(missingFilesSlice, fname)
		}
	}

	missingLibsSlice := []string{}
	for fname := range missingLibs {
		missingLibsSlice = append(missingLibsSlice, fname)
	}

	failedRenamesSlice := []string{}
	for fname := range failedRenames {
		failedRenamesSlice = append(failedRenamesSlice, fname)
	}

	sort.Strings(missingLibsSlice)
	sort.Strings(missingFilesSlice)

	return &AnalyzeSummary{MissingFiles: missingFilesSlice, ExdevFailures: failedRenamesSlice, MissingLibs: missingLibsSlice}, nil
}
