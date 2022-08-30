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
	"path/filepath"
	"strings"
)

type libraryAnalyzer struct {
	foundLibs map[string]string
	filters   *Filters
}

func newLibraryAnalyzer(filters *Filters) *libraryAnalyzer {
	return &libraryAnalyzer{
		foundLibs: map[string]string{},
		filters:   filters,
	}
}

func (analyzer *libraryAnalyzer) doesHaveEquivalentLibrary(library string) bool {
	// if there is a library with the same name but different equivalent extension
	// for example compiled version of a library with so extension we will also remove it from the missing files
	parts := strings.SplitN(library, ".", 2)
	if len(parts) < 2 {
		return false
	}
	noExt := parts[0]
	if equivalents, ok := analyzer.filters.LibEquivalent[filepath.Ext(library)]; ok {
		for _, equiv := range equivalents {
			if _, ok := analyzer.foundLibs[noExt+equiv]; ok {
				return true
			}
		}
	}

	return false

}

func (analyzer *libraryAnalyzer) filterOutFoundLibraries(missingFiles map[string]any) {
	for fname := range missingFiles {
		if !analyzer.filters.IsLibrary(fname) {
			// if the file is not a library we skip it
			continue
		}

		library := filepath.Base(fname)
		if _, ok := analyzer.foundLibs[library]; ok {
			// if a missing library appears in the found files
			// we remove it from the missing files
			delete(missingFiles, fname)
			continue
		}

		if analyzer.doesHaveEquivalentLibrary(library) {
			delete(missingFiles, fname)
		}
	}
}

func (analyzer *libraryAnalyzer) handleFoundFile(path string) bool {
	if analyzer.filters.IsLibrary(path) {
		analyzer.foundLibs[filepath.Base(path)] = path
		return true
	}
	return false
}
