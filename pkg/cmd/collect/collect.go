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

package collect

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"github.com/google/containerdbg/pkg/collect"
)

type CollectOptions struct {
	outputFilename string
}

func NewCollectCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	options := CollectOptions{}
	cmd := &cobra.Command{
		Use:     "collect -o <filename.pb>",
		Short:   "collect all recorded events from the current containerdbg deployment",
		Long:    "collect can be used while running a containerdbg deployment to dump the events before finishing the debug session",
		Example: "collect -o events.pb",
		RunE: func(cmd *cobra.Command, args []string) error {
			return options.collect(cmd.Context(), f, streams)

		},
	}

	cmd.Flags().StringVarP(&options.outputFilename, "output", "o", "", "The path to the output file")
	cmd.MarkFlagRequired("output")
	cmd.MarkFlagFilename("output")

	return cmd
}

func (o *CollectOptions) collect(ctx context.Context, f cmdutil.Factory, streams genericclioptions.IOStreams) error {
	outFile, err := os.Create(o.outputFilename)
	if err != nil {
		return err
	}

	return collect.CollectRecordedData(ctx, f, outFile)
}
