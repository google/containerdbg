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

package events

import (
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/protobuf/types/known/timestamppb"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events/api"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

type DynamicSource struct {
	eventCh chan *proto.Event
}

var _ api.EventsSource = &DynamicSource{}

func NewDynamicSource() *DynamicSource {
	return &DynamicSource{
		eventCh: make(chan *proto.Event),
	}
}

func (o *DynamicSource) Load(log logr.Logger) error {
	return nil
}

func (o *DynamicSource) SendEvent(event *proto.Event) {
	event.Timestamp = timestamppb.New(time.Now())
	o.eventCh <- event
}

func (o *DynamicSource) Events() <-chan *proto.Event {
	return o.eventCh
}

func (o *DynamicSource) Close() error {
	close(o.eventCh)
	return nil
}
