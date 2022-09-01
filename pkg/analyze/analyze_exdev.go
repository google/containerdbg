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

	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

type exdevAnalyzer struct {
	failedRenames map[string]*proto.ContainerAnalysisSummary_MoveFailure
	sourceId      *proto.SourceId
}

func newExdevAnalyzer(sourceId *proto.SourceId) *exdevAnalyzer {
	return &exdevAnalyzer{
		failedRenames: map[string]*proto.ContainerAnalysisSummary_MoveFailure{},
		sourceId:      sourceId,
	}
}

var _ analyzer = &exdevAnalyzer{}

func (analyzer *exdevAnalyzer) handleEvent(event *proto.Event) bool {
	if event.GetSource().GetId() != analyzer.sourceId.GetId() || event.GetSource().GetType() != analyzer.sourceId.GetType() {
		return true
	}
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
		analyzer.failedRenames[renameSyscall.Oldname] = &proto.ContainerAnalysisSummary_MoveFailure{
			Source:      renameSyscall.GetOldname(),
			Destination: renameSyscall.GetNewname(),
		}
		return false
	}

	return true
}

func (analyzer *exdevAnalyzer) updateSummary(summary *proto.ContainerAnalysisSummary) {
	failedRenamesSlice := []*proto.ContainerAnalysisSummary_MoveFailure{}
	for _, rename := range analyzer.failedRenames {
		failedRenamesSlice = append(failedRenamesSlice, rename)
	}

	summary.MoveFailures = failedRenamesSlice
}
