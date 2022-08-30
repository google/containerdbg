// Copyright 2021 Google LLC All Rights Reserved.
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
	"bytes"
	_ "embed"
	"io"
	"path/filepath"
	"regexp"

	"k8s.io/apimachinery/pkg/util/yaml"
)

type filtersFileData struct {
	CommFilter        []string            `json:"commFilter"`
	FileRegexFilter   []string            `json:"fileRegexFilter"`
	LibraryExtensions []string            `json:"libraryExtensions"`
	LibEquivalent     map[string][]string `json:"libEquivalents"`
}

type Filters struct {
	CommFilter        map[string]any
	FileRegexFilter   []*regexp.Regexp
	LibraryExtensions map[string]any
	LibEquivalent     map[string][]string
}

//go:embed filters.yaml
var defaultFiltersFile []byte

func (filters *Filters) IsLibrary(path string) bool {
	_, ok := filters.LibraryExtensions[filepath.Ext(path)]
	return ok
}

func LoadFilters(reader io.Reader) (*Filters, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	filters := filtersFileData{}
	if err := yaml.Unmarshal(data, &filters); err != nil {
		return nil, err
	}

	result := &Filters{
		CommFilter:        map[string]any{},
		FileRegexFilter:   []*regexp.Regexp{},
		LibraryExtensions: map[string]any{},
		LibEquivalent:     filters.LibEquivalent,
	}

	for _, comm := range filters.CommFilter {
		result.CommFilter[comm] = nil
	}

	for _, extension := range filters.LibraryExtensions {
		result.LibraryExtensions[extension] = nil
	}

	for _, fileFilter := range filters.FileRegexFilter {
		r, err := regexp.Compile(fileFilter)
		if err != nil {
			return nil, err
		}
		result.FileRegexFilter = append(result.FileRegexFilter, r)
	}

	return result, nil
}

var defaultFilters *Filters

func init() {
	var err error

	defaultFilters, err = LoadFilters(bytes.NewBuffer(defaultFiltersFile))

	if err != nil {
		panic(err)
	}
}
