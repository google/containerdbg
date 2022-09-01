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

package table

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/cli-utils/pkg/apply/event"
	pe "sigs.k8s.io/cli-utils/pkg/kstatus/polling/event"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
	"sigs.k8s.io/cli-utils/pkg/object"
	"sigs.k8s.io/cli-utils/pkg/object/validation"
	"sigs.k8s.io/cli-utils/pkg/print/stats"
	"sigs.k8s.io/cli-utils/pkg/print/table"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/decoder"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/polymorphichelpers"
)

// mostly copied from https://github.com/kubernetes-sigs/cli-utils/blob/master/pkg/printers/table/collector.go
// with changes to the processEvents to look up pods information

const InvalidStatus status.Status = "Invalid"

type resourceStateCollector struct {
	mux sync.RWMutex

	resourceInfos map[object.ObjMetadata]*resourceInfo

	stats stats.Stats

	err error

	podsclient corev1client.PodsGetter

	wasInstalled bool
}

type resourceInfo struct {
	// identifier contains the information that identifies a
	// single resource.
	identifier object.ObjMetadata

	// resourceStatus contains the latest status information
	// about the resource.
	resourceStatus *pe.ResourceStatus

	// ResourceAction defines the action we are performing
	// on this particular resource. This can be either Apply
	// or Prune.
	ResourceAction event.ResourceAction

	// Error is set if an error occurred trying to perform
	// the desired action on the resource.
	Error error

	// ApplyStatus contains the result after
	// a resource has been applied to the cluster.
	ApplyStatus event.ApplyEventOperation

	// PruneStatus contains the result after
	// a prune operation on a resource
	PruneStatus event.PruneEventOperation

	// DeleteStatus contains the result after
	// a delete operation on a resource
	DeleteStatus event.DeleteEventOperation

	// WaitStatus contains the result after
	// a wait operation on a resource
	WaitStatus event.WaitEventOperation
}

func (r *resourceInfo) Identifier() object.ObjMetadata {
	return r.identifier
}

func (r *resourceInfo) ResourceStatus() *pe.ResourceStatus {
	return r.resourceStatus
}

func (r *resourceInfo) SubResources() []table.Resource {
	var resources []table.Resource

	for _, res := range r.resourceStatus.GeneratedResources {
		resources = append(resources, &subResourceInfo{
			resourceStatus: res,
		})
	}

	return resources
}

type subResourceInfo struct {
	resourceStatus *pe.ResourceStatus
}

func (r *subResourceInfo) Identifier() object.ObjMetadata {
	return r.resourceStatus.Identifier
}

func (r *subResourceInfo) ResourceStatus() *pe.ResourceStatus {
	return r.resourceStatus
}

func (r *subResourceInfo) SubResources() []table.Resource {
	var resources []table.Resource
	for _, res := range r.resourceStatus.GeneratedResources {
		resources = append(resources, &subResourceInfo{
			resourceStatus: res,
		})
	}

	return resources
}

func newResourceStateCollector(resourceGroups []event.ActionGroup, podsClient corev1client.PodsGetter) *resourceStateCollector {

	resourceInfos := make(map[object.ObjMetadata]*resourceInfo)
	for _, group := range resourceGroups {
		action := group.Action

		if action == event.WaitAction {
			continue
		}

		for _, identifier := range group.Identifiers {
			resourceInfos[identifier] = &resourceInfo{
				identifier: identifier,
				resourceStatus: &pe.ResourceStatus{
					Identifier: identifier,
					Status:     status.UnknownStatus,
				},
				ResourceAction: action,
			}
		}
	}

	return &resourceStateCollector{
		resourceInfos: resourceInfos,
		podsclient:    podsClient,
		wasInstalled:  true,
	}
}

// Listen starts a new goroutine that will listen for events on the provided eventChannel and keep track of the latest state for the resources
func (r *resourceStateCollector) Listen(eventChannel <-chan event.Event) <-chan listenerResult {

	completed := make(chan listenerResult)
	go func() {
		defer close(completed)
		for ev := range eventChannel {
			if err := r.processEvent(ev); err != nil {
				completed <- listenerResult{err: err}
				return
			}
		}
	}()

	return completed
}

type listenerResult struct {
	err error
}

func (r *resourceStateCollector) processEvent(ev event.Event) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	switch ev.Type {
	case event.ValidationType:
		return r.processValidationEvent(ev.ValidationEvent)
	case event.StatusType:
		return r.processStatusEvent(ev.StatusEvent)
	case event.ApplyType:
		r.processApplyEvent(ev.ApplyEvent)
	case event.PruneType:
		r.processPruneEvent(ev.PruneEvent)
	case event.WaitType:
		r.processWaitEvent(ev.WaitEvent)
	case event.ErrorType:
		return ev.ErrorEvent.Err
	}

	return nil
}

func (r *resourceStateCollector) processValidationEvent(e event.ValidationEvent) error {
	err := e.Error
	if vErr, ok := err.(*validation.Error); ok {
		err = vErr.Unwrap()
	}

	if len(e.Identifiers) == 0 {
		// no objects, invalid event

		return fmt.Errorf("invalid validation event: no identifiers: %w", err)
	}

	for _, id := range e.Identifiers {
		previous, found := r.resourceInfos[id]

		if !found {
			continue
		}

		previous.resourceStatus = &pe.ResourceStatus{
			Identifier: id,
			Status:     InvalidStatus,
			Message:    e.Error.Error(),
		}
	}

	return nil
}

func (r *resourceStateCollector) processStatusEvent(e event.StatusEvent) error {
	previous, found := r.resourceInfos[e.Identifier]
	if !found {
		return nil
	}

	previous.resourceStatus = e.PollResourceInfo

	err := r.processResource(previous.resourceStatus, e.Resource)
	if err != nil {
		previous.resourceStatus.Error = err
		previous.resourceStatus.Status = status.FailedStatus
	}

	return err
}

func (r *resourceStateCollector) processResource(previous *pe.ResourceStatus, res *unstructured.Unstructured) error {
	if res == nil {
		return nil
	}
	obj, err := decoder.DecodeUnstructured(res)
	if err != nil {
		return nil
	}

	pods, err := polymorphichelpers.PodsForObject(context.TODO(), r.podsclient, obj)
	if err != nil {
		return nil
	}
	var firstError error
	for _, pod := range pods {
		var resourceError error
		currentStatus := status.InProgressStatus
		for _, container := range pod.Status.ContainerStatuses {
			if container.LastTerminationState.Terminated != nil {
				currentStatus = status.FailedStatus
				resourceError = errors.New(container.LastTerminationState.Terminated.Message)
				if firstError == nil {
					firstError = resourceError
				}
			} else if container.Ready {
				currentStatus = status.CurrentStatus
			}

		}
		previous.GeneratedResources = append(previous.GeneratedResources, &pe.ResourceStatus{
			Identifier: object.ObjMetadata{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				GroupKind: schema.GroupKind{Kind: "Pod"},
			},
			Status: currentStatus,
			Error:  resourceError,
		})
	}

	if firstError != nil {
		return firstError
	}

	return nil
}

func (r *resourceStateCollector) processApplyEvent(e event.ApplyEvent) {
	previous, found := r.resourceInfos[e.Identifier]
	if !found {
		return
	}
	if e.Error != nil {
		previous.Error = e.Error
	}
	if e.Operation == event.Created {
		r.wasInstalled = false
	}
	previous.ApplyStatus = e.Operation
}

func (r *resourceStateCollector) WasInstalledBefore() bool {
	return r.wasInstalled
}

func (r *resourceStateCollector) processPruneEvent(e event.PruneEvent) {
	previous, found := r.resourceInfos[e.Identifier]
	if !found {
		return
	}
	if e.Error != nil {
		previous.Error = e.Error
	}
	previous.PruneStatus = e.Operation
}

func (r *resourceStateCollector) processWaitEvent(e event.WaitEvent) {
	previous, found := r.resourceInfos[e.Identifier]
	if !found {
		return
	}
	previous.WaitStatus = e.Operation
}

type ResourceInfos []*resourceInfo

func (ri ResourceInfos) Len() int {
	return len(ri)
}

func (ri ResourceInfos) Less(i, j int) bool {
	idI := ri[i].identifier
	idJ := ri[j].identifier

	if idI.Namespace != idJ.Namespace {
		return idI.Namespace < idJ.Namespace
	}

	if idI.GroupKind.Group != idJ.GroupKind.Group {
		return idI.GroupKind.Group < idJ.GroupKind.Group
	}

	if idI.GroupKind.Kind != idJ.GroupKind.Kind {
		return idI.GroupKind.Kind < idI.GroupKind.Kind
	}

	return idI.Name < idJ.Name
}

func (ri ResourceInfos) Swap(i, j int) {
	ri[i], ri[j] = ri[j], ri[i]
}

type ResourceState struct {
	resourceInfos ResourceInfos
	err           error
}

func (r *ResourceState) Resources() []table.Resource {
	var resources []table.Resource

	for _, res := range r.resourceInfos {
		resources = append(resources, res)
	}

	return resources
}

func (r *ResourceState) Error() error {
	return r.err
}

func (r *resourceStateCollector) LatestState() *ResourceState {
	r.mux.Lock()
	defer r.mux.Unlock()

	var resourceInfos ResourceInfos
	for _, ri := range r.resourceInfos {
		resourceInfos = append(resourceInfos, &resourceInfo{
			identifier:     ri.identifier,
			resourceStatus: ri.resourceStatus,
			ResourceAction: ri.ResourceAction,
			ApplyStatus:    ri.ApplyStatus,
			PruneStatus:    ri.PruneStatus,
			DeleteStatus:   ri.DeleteStatus,
			WaitStatus:     ri.WaitStatus,
		})
	}

	sort.Sort(resourceInfos)

	return &ResourceState{
		resourceInfos: resourceInfos,
		err:           r.err,
	}
}

func (r *resourceStateCollector) Stats() stats.Stats {
	var s stats.Stats

	for _, res := range r.resourceInfos {
		switch res.ResourceAction {
		case event.ApplyAction:
			if res.Error != nil {
				s.ApplyStats.IncFailed()
			}
			s.ApplyStats.Inc(res.ApplyStatus)
		case event.PruneAction:
			if res.Error != nil {
				s.PruneStats.IncFailed()
			}
			s.PruneStats.Inc(res.PruneStatus)
		case event.DeleteAction:
			if res.Error != nil {
				s.DeleteStats.IncFailed()
			}
			s.DeleteStats.Inc(res.DeleteStatus)
		}
		s.WaitStats.Inc(res.WaitStatus)
	}

	return s
}
