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

package convert

import (
	"bufio"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

func NewConvertCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "convert <input> <output>",
		Short:  "convert converts a binary event file to a text file",
		Long:   "convert converts a binary event file to a text file - for custom parsing",
		Hidden: true,
		Args:   cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			out, err := os.Create(args[1])
			if err != nil {
				return err
			}

			writer := events.NewEventWriter(out)
			defer writer.Close()

			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				line := scanner.Bytes()
				m := proto.Event{}
				if err := prototext.Unmarshal(line, &m); err != nil {
					return err
				}
				if err := writer.Write(&m); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}
