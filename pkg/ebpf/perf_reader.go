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
	"errors"
	"sync"

	"github.com/cilium/ebpf/perf"
	"github.com/go-logr/logr"
	"github.com/google/containerdbg/proto"
)

type eventParserFunc func(sample []byte) (*proto.Event, error)

type perfReader struct {
	log            logr.Logger
	wg             sync.WaitGroup
	readThreadDone chan bool
	rd             *perf.Reader
	events         chan *proto.Event
	samplesChan    chan []byte
	parser         eventParserFunc
}

func (reader *perfReader) Start() {
	reader.wg.Add(1)
	go reader.readEventsThread()
}

func NewPerfReader(log logr.Logger, rd *perf.Reader, parser eventParserFunc) *perfReader {
	return &perfReader{
		log:            log,
		rd:             rd,
		events:         make(chan *proto.Event, 100),
		samplesChan:    make(chan []byte, 100),
		readThreadDone: make(chan bool),
		parser:         parser,
	}
}

func (reader *perfReader) readEventsThread() {
	defer reader.wg.Done()
	for {
		record, err := reader.rd.Read()
		if err != nil {
			if errors.Is(err, perf.ErrClosed) {
				break
			}
		}

		if record.LostSamples != 0 {
			reader.log.V(1).Info("perf event ring buffer full", "dropped samples", record.LostSamples)
			continue
		}

		event, err := reader.parser(record.RawSample)
		if err != nil {
			reader.log.V(1).Info("event parsing failed: %s", err)
			continue
		}
		reader.events <- event

	}

	close(reader.readThreadDone)
}

func (reader *perfReader) Events() chan *proto.Event {
	return reader.events
}

func (reader *perfReader) cleanEventChannel() {
loop:
	for {
		select {
		case <-reader.readThreadDone:
			break loop
		case <-reader.events:
		}
	}
}

func (reader *perfReader) Stop() error {
	if reader.rd != nil {
		if err := reader.rd.Close(); err != nil {
			return err
		}
		reader.cleanEventChannel()
		close(reader.events)
	}
	reader.wg.Wait()
	return nil
}
