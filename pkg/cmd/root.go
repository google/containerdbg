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

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"github.com/google/containerdbg/pkg/cmd/analyze"
	"github.com/google/containerdbg/pkg/cmd/attach"
	"github.com/google/containerdbg/pkg/cmd/collect"
	"github.com/google/containerdbg/pkg/cmd/convert"
	"github.com/google/containerdbg/pkg/cmd/debug"
	"github.com/google/containerdbg/pkg/cmd/dump"
	"github.com/google/containerdbg/pkg/cmd/install"
	"github.com/google/containerdbg/pkg/cmd/uninstall"
	"github.com/google/containerdbg/pkg/cmd/version"
)

func NewRootCmd(f cmdutil.Factory, streams genericclioptions.IOStreams, interruptCh <-chan os.Signal) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "containerdbg SUBCOMMAND",
		SilenceUsage: true,
	}

	cmd.AddCommand(collect.NewCollectCmd(f, streams),
		debug.NewDebugCmd(f, streams, interruptCh),
		analyze.NewAnalyzeCmd(streams),
		dump.NewDumpCmd(streams),
		convert.NewConvertCmd(streams),
		attach.NewAttachCmd(f, streams, interruptCh),
		version.NewVersionCmd(),
		uninstall.NewUninstallCmd(f, streams),
		install.NewInstallCmd(f, streams, interruptCh))

	return cmd
}
