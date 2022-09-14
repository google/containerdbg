// Copyright 2021 Google LLC All Rights Reserved.
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

package support

import (
	"bytes"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	utilpointer "k8s.io/utils/pointer"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"github.com/google/containerdbg/pkg/cmd"
)

func NewTestRootCmd(cfg *envconf.Config, interruptCh <-chan os.Signal) (*cobra.Command, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {

	streams, in, out, errOut := genericclioptions.NewTestIOStreams()

	config := genericclioptions.NewConfigFlags(false)
	config.KubeConfig = utilpointer.String(cfg.KubeconfigFile())

	factory := cmdutil.NewFactory(config)

	rootcmd := cmd.NewRootCmd(factory, streams, interruptCh)
	config.AddFlags(rootcmd.PersistentFlags())

	return rootcmd, in, out, errOut
}
