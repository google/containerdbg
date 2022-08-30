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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/analyze"
)

type AnalyzeOptions struct {
	collectionFilename string
	filtersFile        string
}

func NewAnalyzeCmd(streams genericclioptions.IOStreams) *cobra.Command {
	options := AnalyzeOptions{}
	cmd := &cobra.Command{
		Use:   "analyze -f <file.pb> [-t filters]",
		Short: "analyze the collected events and print a list of probable errors with deployment",
		Long: "analyze takes as input a collected event file and lists probable errors, " +
			"optionally can recieve a custom list of files to ignore",
		RunE: func(cmd *cobra.Command, args []string) error {
			var filters *analyze.Filters
			if options.filtersFile != "" {
				var err error
				f, err := os.Open(options.filtersFile)
				if err != nil {
					return err
				}
				filters, err = analyze.LoadFilters(f)
				if err != nil {
					return err
				}
			}
			sum, err := analyze.Analyze(options.collectionFilename, filters)

			for _, contsum := range sum.GetContainerSummaries() {
				fmt.Fprintf(cmd.OutOrStdout(), "Findings for container %v\n", contsum.Source)
				sum := contsum.GetSummary()
				fmt.Fprintln(cmd.OutOrStdout(), "================================")
				fmt.Fprintln(cmd.Parent().OutOrStdout())

				if len(sum.MissingFiles) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "While executing the container the following files were missing:")
					fmt.Fprintln(cmd.OutOrStdout(), "--------------------------------")
					for _, fname := range sum.MissingFiles {
						fmt.Fprintf(cmd.OutOrStdout(), "%s is missing\n", fname)
					}
					fmt.Fprintln(cmd.Parent().OutOrStdout())
				}

				if len(sum.MissingLibraries) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "While executing the container the library type files were missing:")
					fmt.Fprintln(cmd.OutOrStdout(), "---------------------------------")
					for _, fname := range sum.MissingLibraries {
						fmt.Fprintf(cmd.OutOrStdout(), "%s is missing\n", fname)
					}
					fmt.Fprintln(cmd.Parent().OutOrStdout())
				}

				if len(sum.MoveFailures) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "While executing the container the following files where attempted to be moved but failed to docker limitation:")
					fmt.Fprintln(cmd.OutOrStdout(), "-------------------------------------------------------")
					for _, fname := range sum.MoveFailures {
						fmt.Fprintf(cmd.OutOrStdout(), "%s is attempted to be renamed atomically and failes do to docker filesystem limitation\n", fname.Source)
					}
					fmt.Fprintln(cmd.Parent().OutOrStdout())
				}

				if len(sum.ConnectionFailures) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "While executing the container the following connections failed:")
					fmt.Fprintln(cmd.OutOrStdout(), "-------------------------------------------------------")
					for _, conn := range sum.ConnectionFailures {
						if conn.TargetFqdn != "" {
							fmt.Fprintf(cmd.OutOrStdout(), "%s(%s):%d\n", conn.TargetFqdn, conn.TargetIp, conn.Port)
						} else {
							fmt.Fprintf(cmd.OutOrStdout(), "%s:%d\n", conn.TargetIp, conn.Port)
						}
					}
					fmt.Fprintln(cmd.Parent().OutOrStdout())
				}

				if len(sum.DnsFailures) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "While executing the container the following dns queries failed:")
					fmt.Fprintln(cmd.OutOrStdout(), "-------------------------------------------------------")
					for _, conn := range sum.DnsFailures {
						fmt.Fprintf(cmd.OutOrStdout(), "%s\n", conn.Query)
					}
					fmt.Fprintln(cmd.Parent().OutOrStdout())
				}

				if len(sum.StaticIps) > 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "While executing the container the following static IPs were detected, consider replacing those with DNS names:")
					fmt.Fprintln(cmd.OutOrStdout(), "-------------------------------------------------------")
					for _, conn := range sum.StaticIps {
						fmt.Fprintf(cmd.OutOrStdout(), "%s:%d\n", conn.TargetIp, conn.Port)
					}
					fmt.Fprintln(cmd.Parent().OutOrStdout())
				}

			}
			return err
		},
	}

	cmd.Flags().StringVarP(&options.collectionFilename, "filename", "f", "", "the file containing the input for the analysis")
	cmd.MarkFlagRequired("filename")
	cmd.MarkFlagFilename("filename")

	cmd.Flags().StringVarP(&options.filtersFile, "filters", "t", "", "the file containing the filters to use for the analysis")
	cmd.MarkFlagFilename("filters")
	cmd.Flags().MarkHidden("filters")

	return cmd
}
