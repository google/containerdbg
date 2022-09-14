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

package dump

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"github.com/google/containerdbg/pkg/events"
)

func NewDumpCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "dump",
		Hidden: true,
		Short:  "dumps collected information file in textual format",
		Long:   "converts the containerdbg protobuf format to textual format",
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()
			reader := events.NewEventReader(f)

			for event, err := reader.Read(); err == nil; event, err = reader.Read() {
				fmt.Fprintln(streams.Out, prototext.MarshalOptions{Multiline: false}.Format(event))
			}

			return nil
		},
	}

	return cmd
}
