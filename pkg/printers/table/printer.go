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

package table

import (
	"fmt"
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/cli-utils/pkg/apply/event"
	"sigs.k8s.io/cli-utils/pkg/common"
	printcommon "sigs.k8s.io/cli-utils/pkg/print/common"
	"sigs.k8s.io/cli-utils/pkg/print/table"
)

// mostly copied from https://github.com/kubernetes-sigs/cli-utils/blob/v0.29.4/pkg/printers/table/printer.go
// modified to collect pods information and termination messages

type Printer struct {
	IOStreams  genericclioptions.IOStreams
	PodsClient corev1client.PodsGetter
}

func (t *Printer) Print(ch <-chan event.Event, _ common.DryRunStrategy, _ bool) (bool, error) {
	var initEvent event.InitEvent
	for e := range ch {
		if e.Type == event.InitType {
			initEvent = e.InitEvent
			break
		}

		// If we get an error event, we just return it
		if e.Type == event.ErrorType {
			return false, e.ErrorEvent.Err
		}
	}

	coll := newResourceStateCollector(initEvent.ActionGroups, t.PodsClient)

	stop := make(chan struct{})

	printCompleted := t.runPrintLoop(coll, stop)

	done := coll.Listen(ch)

	var err error

	for msg := range done {
		err = msg.err
	}

	close(stop)

	<-printCompleted

	if err != nil {
		return false, err
	}

	return coll.WasInstalledBefore(), printcommon.ResultErrorFromStats(coll.Stats())
}

var columns = []table.ColumnDefinition{
	table.MustColumn("namespace"),
	table.MustColumn("resource"),
	table.MustColumn("status"),
	table.MustColumn("conditions"),
	table.MustColumn("age"),
	table.MustColumn("message"),
}

func (t *Printer) runPrintLoop(coll *resourceStateCollector, stop chan struct{}) chan struct{} {
	finished := make(chan struct{})

	baseTablePrinter := table.BaseTablePrinter{
		IOStreams: t.IOStreams,
		Columns:   columns,
	}

	linesPrinted := baseTablePrinter.PrintTable(coll.LatestState(), 0)

	go func() {
		defer close(finished)

		ticker := time.NewTicker(500 * time.Millisecond)

		for {
			select {
			case <-stop:
				ticker.Stop()
				latestState := coll.LatestState()
				linesPrinted = baseTablePrinter.PrintTable(latestState, linesPrinted)
				_, _ = fmt.Fprint(t.IOStreams.Out, "\n")
				return
			case <-ticker.C:
				latestState := coll.LatestState()
				linesPrinted = baseTablePrinter.PrintTable(latestState, linesPrinted)
			}
		}
	}()

	return finished
}
