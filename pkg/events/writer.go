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
	pb "github.com/google/containerdbg/proto"
)

type EventWriter struct {
	writer io.WriteCloser
}

func NewEventWriter(writer io.WriteCloser) *EventWriter {
	return &EventWriter{
		writer: writer,
	}
}

func (writer *EventWriter) Write(event *pb.Event) error {
	data, err := proto.Marshal((*pb.Event)(event))
	if err != nil {
		return err
	}
	sb := make([]byte, 4)
	binary.LittleEndian.PutUint32(sb, uint32(len(data)))
	_, err = writer.writer.Write(sb)
	if err != nil {
		return err
	}
	_, err = writer.writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (writer *EventWriter) Close() error {
	return writer.writer.Close()
}
