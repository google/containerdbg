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

package install

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/install"
)

func NewInstallCmd(f cmdutil.Factory, streams genericclioptions.IOStreams, interruptCh <-chan os.Signal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install containerdbg daemonset",
		RunE: func(cmd *cobra.Command, args []string) error {
			interrupttedCtx, cancel := context.WithCancel(cmd.Context())

			go func() {
				<-interruptCh
				cancel()
			}()
			_, err := install.EnsureInstallation(interrupttedCtx, f, streams)
			return err
		},
	}

	return cmd
}
