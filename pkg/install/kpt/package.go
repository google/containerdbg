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

package kpt

import (
	"context"
	"fmt"
	"time"

	"github.com/GoogleContainerTools/kpt/pkg/live"
	"github.com/GoogleContainerTools/kpt/pkg/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/cli-utils/pkg/apply"
	"sigs.k8s.io/cli-utils/pkg/common"
	"sigs.k8s.io/cli-utils/pkg/inventory"
	"sigs.k8s.io/cli-utils/pkg/kstatus/polling"
	"sigs.k8s.io/cli-utils/pkg/printers"
)

func InstallRG(ctx context.Context, f cmdutil.Factory) error {
	return (&live.ResourceGroupInstaller{
		Factory: f,
	}).InstallRG(ctx)
}

type KptPackage struct {
	objs         []*unstructured.Unstructured
	invInfo      inventory.Info
	invClient    inventory.Client
	statusPoller *polling.StatusPoller
}

func LoadPackage(ctx context.Context, f cmdutil.Factory, path string) (*KptPackage, error) {
	result := &KptPackage{}
	if !live.ResourceGroupCRDApplied(f) {
		if err := InstallRG(ctx, f); err != nil {
			return nil, err
		}
	}
	objs, inv, err := live.Load(f, path, "", nil)
	if err != nil {
		return nil, err
	}
	result.objs = objs

	invInfo, err := live.ToInventoryInfo(inv)
	if err != nil {
		return nil, err
	}
	result.invInfo = invInfo

	invClient, err := inventory.NewClient(f, live.WrapInventoryObj, live.InvToUnstructuredFunc, inventory.StatusPolicyAll)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %v", err)
	}
	result.invClient = invClient

	statusPoller, err := status.NewStatusPoller(f)
	if err != nil {
		return nil, err
	}
	result.statusPoller = statusPoller

	return result, nil
}

func (pkg *KptPackage) Install(ctx context.Context, f cmdutil.Factory, streams genericclioptions.IOStreams, fieldmanager string) error {
	fmt.Fprintf(streams.Out, "Installing containerdbg node daemon\n")
	applier, err := apply.NewApplierBuilder().
		WithFactory(f).
		WithInventoryClient(pkg.invClient).
		WithStatusPoller(pkg.statusPoller).
		Build()

	if err != nil {
		return err
	}

	ch := applier.Run(ctx, pkg.invInfo, pkg.objs, apply.ApplierOptions{
		ServerSideOptions: common.ServerSideOptions{
			ServerSideApply: false,
			FieldManager:    fieldmanager,
		},
		PollInterval:           2 * time.Second,
		ReconcileTimeout:       1 * time.Minute,
		EmitStatusEvents:       true,
		DryRunStrategy:         common.DryRunNone,
		PrunePropagationPolicy: metav1.DeletePropagationBackground,
		PruneTimeout:           time.Duration(0),
		InventoryPolicy:        inventory.PolicyAdoptIfNoInventory,
	})

	printer := printers.GetPrinter(printers.TablePrinter, streams)
	return printer.Print(ch, common.DryRunNone, true)
}

func (pkg *KptPackage) Uninstall(ctx context.Context, f cmdutil.Factory, streams genericclioptions.IOStreams) error {
	fmt.Fprintf(streams.Out, "Uninstalling containerdbg node daemon\n")
	destroyer, err := apply.NewDestroyer(f, pkg.invClient)
	if err != nil {
		return err
	}
	destroyer.StatusPoller = pkg.statusPoller

	options := apply.DestroyerOptions{
		InventoryPolicy:  inventory.PolicyMustMatch,
		DryRunStrategy:   common.DryRunNone,
		EmitStatusEvents: true,
	}

	ch := destroyer.Run(ctx, pkg.invInfo, options)

	printer := printers.GetPrinter(printers.TablePrinter, streams)
	return printer.Print(ch, common.DryRunNone, true)
}
