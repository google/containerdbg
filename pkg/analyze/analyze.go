// Package analyze is responsible for failure analysis logic from debugged container
//
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
	"errors"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/google/containerdbg/pkg/events"
	"github.com/google/containerdbg/proto"
)

func analyzeContainer(inputFilename string, sourceId *proto.SourceId, searchPath []string, filters *Filters) (*proto.ContainerAnalysisSummary, error) {
	f, err := os.Open(inputFilename)
	if err != nil {
		return nil, err
	}

	reader := events.NewEventReader(f)

	analyzers := []analyzer{
		newOpenAnalyzer(filters, sourceId),
		newExdevAnalyzer(sourceId),
		newConnectionAnalyzer(searchPath, sourceId),
	}

readLoop:
	for event, err := reader.Read(); err == nil; event, err = reader.Read() {
		for _, analyzer := range analyzers {
			if !analyzer.handleEvent(event) {
				continue readLoop
			}
		}
	}

	containerSummary := &proto.ContainerAnalysisSummary{}

	for _, analyzer := range analyzers {
		analyzer.updateSummary(containerSummary)
	}

	return containerSummary, nil
}

type sourceTuple struct {
	Type string
	Id   string
}

type sources []*proto.SourceId

func (s sources) Len() int {
	return len([]*proto.SourceId(s))
}

func (s sources) Less(i, j int) bool {

	if s[i].Type != s[j].Type {
		return strings.Compare(s[i].Type, s[j].Type) < 0
	}

	return strings.Compare(s[i].Id, s[j].Id) < 0
}

func (s sources) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func findAllSources(inputFilename string) (result []*proto.SourceId, searchResult [][]string, err error) {
	sourceToId := map[sourceTuple]*proto.SourceId{}
	searchToQueries := map[sourceTuple][]string{}

	f, err := os.Open(inputFilename)
	if err != nil {
		return nil, nil, err
	}

	reader := events.NewEventReader(f)

	var event *proto.Event

	for event, err = reader.Read(); true; event, err = reader.Read() {
		if err != nil {
			break
		}
		if event.GetSource() == nil {
			continue
		}
		k := sourceTuple{
			Type: event.GetSource().GetType(),
			Id:   event.GetSource().GetId(),
		}

		sourceToId[k] = event.Source
		if event.GetDnsSearch() != nil {
			searchToQueries[k] = event.GetDnsSearch().GetSearch()
		}
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return nil, nil, err
	}

	for _, s := range sourceToId {
		result = append(result, s)
	}

	sort.Sort(sources(result))

	for _, source := range result {
		searchResult = append(searchResult, searchToQueries[sourceTuple{
			Type: source.GetType(),
			Id:   source.GetId(),
		}])
	}

	return
}

func Analyze(inputFilename string, filters *Filters) (*proto.AnalysisSummary, error) {
	if filters == nil {
		filters = defaultFilters
	}

	sources, searchPaths, err := findAllSources(inputFilename)
	if err != nil {
		return nil, err
	}

	summary := &proto.AnalysisSummary{
		ContainerSummaries: []*proto.AnalysisSummary_ContainerSummaryTuple{},
	}

	for i, s := range sources {
		if s.GetType() == "host" || s.GetType() == "pod" {
			continue
		}
		sum, err := analyzeContainer(inputFilename, s, searchPaths[i], filters)
		if err != nil {
			return nil, err
		}
		summary.ContainerSummaries = append(summary.ContainerSummaries, &proto.AnalysisSummary_ContainerSummaryTuple{
			Source:  s,
			Summary: sum,
		})

	}

	return summary, nil
}
