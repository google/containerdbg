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

package events

import (
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"
	pb "velostrata-internal.googlesource.com/containerdbg.git/proto"
)

type EventReader struct {
	reader io.Reader
}

func NewEventReader(reader io.Reader) *EventReader {
	return &EventReader{
		reader: reader,
	}
}

func (reader *EventReader) Read() (*pb.Event, error) {
	sb := make([]byte, 4)
	_, err := reader.reader.Read(sb)
	if err != nil {
		return nil, err
	}
	size := binary.LittleEndian.Uint32(sb)

	eventData := make([]byte, size)
	_, err = reader.reader.Read(eventData)
	if err != nil {
		return nil, err
	}

	event := &pb.Event{}

	err = proto.Unmarshal(eventData, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
