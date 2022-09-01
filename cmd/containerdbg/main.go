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

package main

import (
	"context"
	"os"
	"os/signal"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/cmd"
)

func main() {

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()

	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	factory := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	streams := genericclioptions.IOStreams{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
		In:     os.Stdin,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	root := cmd.NewRootCmd(factory, streams, done)
	kubeConfigFlags.AddFlags(root.PersistentFlags())

	if err := root.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}
