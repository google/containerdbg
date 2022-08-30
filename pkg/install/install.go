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

package install

import (
	"context"
	"fmt"
	"os"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"velostrata-internal.googlesource.com/containerdbg.git/deploy"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/consts"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/install/kpt"
)

type Object interface {
	runtime.Object
	metav1.Object
}

func GetContainerDbgPackage() (string, error) {
	path := path.Join(os.Getenv("HOME"), ".cache", "containerdbg")
	if err := deploy.ExpandToDir(path); err != nil {
		return "", err
	}

	return path, nil
}

func EnsureInstallation(ctx context.Context, f cmdutil.Factory, streams genericclioptions.IOStreams) (bool, error) {
	f = NewNoNamespaceFactory(f)
	path, err := GetContainerDbgPackage()
	if err != nil {
		return false, err
	}

	pkg, err := kpt.LoadPackage(ctx, f, path)
	if err != nil {
		return false, fmt.Errorf("failed to load package: %v", err)
	}
	wasInstalled, err := pkg.Install(ctx, f, streams, consts.ContainerdbgFieldManagerName)
	if err != nil {
		pkg.Uninstall(ctx, f, streams) // rollback and ignore error
		return false, fmt.Errorf("installation failed: %v", err)
	}
	return wasInstalled, nil
}

func Uninstall(ctx context.Context, f cmdutil.Factory, streams genericclioptions.IOStreams) error {
	f = NewNoNamespaceFactory(f)
	path, err := GetContainerDbgPackage()
	if err != nil {
		return err
	}

	pkg, err := kpt.LoadPackage(ctx, f, path)
	if err != nil {
		return err
	}

	return pkg.Uninstall(ctx, f, streams)
}
