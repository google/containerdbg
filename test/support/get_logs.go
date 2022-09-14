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
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeresource "k8s.io/cli-runtime/pkg/resource"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/reference"
	"k8s.io/kubectl/pkg/scheme"
	"github.com/google/containerdbg/pkg/polymorphichelpers"
)

func GetLogsForObject(ctx context.Context, podsGetter corev1client.PodsGetter, object runtime.Object) (string, error) {
	pods, err := polymorphichelpers.PodsForObject(ctx, podsGetter, object)
	if err != nil {
		return "", err
	}
	firstPod := pods[0]
	result := ""
	for _, container := range firstPod.Spec.Containers {
		resp := podsGetter.Pods(firstPod.Namespace).GetLogs(firstPod.Name, &v1.PodLogOptions{
			Container: container.Name,
		})
		data, err := resp.Do(ctx).Raw()
		if err != nil {
			return "", err
		}

		result += string(data)
	}

	return result, nil
}

func DumpEvents(ctx context.Context, t *testing.T, client corev1client.EventsGetter, objOrRef runtime.Object, limit int64) {

	ref, err := reference.GetReference(scheme.Scheme, objOrRef)
	if err != nil {
		t.Logf("failed to get ref for obj %v: %s", objOrRef, err)
		return
	}

	stringRefKind := string(ref.Kind)
	var refKind *string

	if len(stringRefKind) > 0 {
		refKind = &stringRefKind
	}

	stringRefUID := string(ref.UID)
	var refUID *string
	if len(stringRefUID) > 0 {
		refUID = &stringRefUID
	}

	e := client.Events(ref.Namespace)

	fieldSelector := e.GetFieldSelector(&ref.Name, &ref.Namespace, refKind, refUID)
	initialOpts := metav1.ListOptions{FieldSelector: fieldSelector.String(), Limit: limit}
	eventList := &v1.EventList{}
	err = runtimeresource.FollowContinue(&initialOpts,
		func(options metav1.ListOptions) (runtime.Object, error) {
			newEvents, err := e.List(ctx, options)
			if err != nil {
				return nil, runtimeresource.EnhanceListError(err, options, "events")
			}
			eventList.Items = append(eventList.Items, newEvents.Items...)
			return newEvents, nil
		})

	if err != nil {
		t.Logf("failed to get events for %v: %s", objOrRef, err)
	}

	for _, e := range eventList.Items {
		t.Logf("%v\t%v\t%v\t%v", e.Type, e.Reason, e.Source.Component, e.Message)
	}
}

func DumpLogs(ctx context.Context, t *testing.T, podsGetter corev1client.PodsGetter, object runtime.Object) {
	logs, logsErr := GetLogsForObject(ctx, podsGetter, object)
	if logsErr == nil {
		t.Log(logs)
	} else {
		t.Logf("failed to get logs: %s", logsErr)
	}
}
