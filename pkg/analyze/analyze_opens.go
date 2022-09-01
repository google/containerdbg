// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analyze

import (
	"path/filepath"
	"sort"

	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

type analyzeOpen struct {
	missingFiles map[string]any
	filters      *Filters
	*libraryAnalyzer
	sourceId *proto.SourceId
}

var _ analyzer = &analyzeOpen{}

func newOpenAnalyzer(filters *Filters, sourceId *proto.SourceId) *analyzeOpen {
	return &analyzeOpen{
		missingFiles:    map[string]any{},
		filters:         filters,
		libraryAnalyzer: newLibraryAnalyzer(filters),
		sourceId:        sourceId,
	}
}

func (analyzer *analyzeOpen) handleEvent(event *proto.Event) bool {
	if event.GetSource().GetId() != analyzer.sourceId.GetId() || event.GetSource().GetType() != analyzer.sourceId.GetType() {
		return true
	}
	syscall := event.GetSyscall()
	if syscall == nil {
		return true
	}

	openSyscall := syscall.GetOpen()
	if openSyscall == nil {
		return true
	}

	if syscall.GetRetCode() >= 0 {
		analyzer.libraryAnalyzer.handleFoundFile(openSyscall.GetPath())
		return false
	}

	if _, ok := analyzer.filters.CommFilter[syscall.GetComm()]; ok {
		return false
	}

	for _, re := range analyzer.filters.FileRegexFilter {
		if re.MatchString(openSyscall.GetPath()) {
			return false
		}
	}

	analyzer.missingFiles[openSyscall.GetPath()] = nil
	return true
}

func (analyzer *analyzeOpen) updateSummary(summary *proto.ContainerAnalysisSummary) {
	analyzer.libraryAnalyzer.filterOutFoundLibraries(analyzer.missingFiles)

	missingFilesSlice := []string{}
	missingLibs := map[string]any{}
	for fname := range analyzer.missingFiles {
		if analyzer.filters.IsLibrary(fname) {
			missingLibs[filepath.Base(fname)] = nil
		} else {
			missingFilesSlice = append(missingFilesSlice, fname)
		}
	}

	missingLibsSlice := []string{}
	for fname := range missingLibs {
		missingLibsSlice = append(missingLibsSlice, fname)
	}

	sort.Strings(missingLibsSlice)
	sort.Strings(missingFilesSlice)

	summary.MissingFiles = missingFilesSlice
	summary.MissingLibraries = missingLibsSlice
}
