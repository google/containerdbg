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

package version

import (
	"github.com/spf13/cobra"
	"github.com/google/containerdbg/pkg/build"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the version of containerdbg",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Printf("Version: %s-%s\n", build.Version, build.GitSha)
			return nil
		},
	}

	return cmd
}
