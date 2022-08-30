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

package test

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	e2ek8s "sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"velostrata-internal.googlesource.com/containerdbg.git/pkg/events"
	"velostrata-internal.googlesource.com/containerdbg.git/pkg/rand"
	pb "velostrata-internal.googlesource.com/containerdbg.git/proto"
	"velostrata-internal.googlesource.com/containerdbg.git/test/support"
)

func helperTestNetEventIsRecorded(t *testing.T, ctx context.Context, cfg *envconf.Config, expectedEvent *pb.Event_Network, dnsEvent *pb.Event_DnsQueryEvent, namespace string, debugParams ...string) context.Context {
	tmpFileName := path.Join(t.TempDir(), "events.json")

	support.RunContainerDebug(t, ctx, cfg, tmpFileName, namespace, debugParams...)

	if !helperFindNetListenEvent(t, tmpFileName, expectedEvent) {
		t.Fatal("did not find listen event")
	}

	if !helperFindDnsEvent(t, tmpFileName, dnsEvent) {
		t.Fatal("did not find dns event")
	}

	return ctx
}

func helperFindNetListenEvent(t *testing.T, filename string, expectedEvent *pb.Event_Network) bool {
	t.Helper()
	recordsFile, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	reader := events.NewEventReader(recordsFile)
	for event, err := reader.Read(); err == nil; event, err = reader.Read() {
		if listenEqual(t, expectedEvent, event) {
			return true
		}
	}

	return false
}

func helperFindDnsEvent(t *testing.T, filename string, expectedEvent *pb.Event_DnsQueryEvent) bool {
	t.Helper()
	if expectedEvent == nil {
		return true
	}
	recordsFile, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	reader := events.NewEventReader(recordsFile)
	for event, err := reader.Read(); err == nil; event, err = reader.Read() {
		if dnsEventEqual(t, expectedEvent, event) {
			return true
		}
	}

	return false
}

func listenEqual(t *testing.T, expected *pb.Event_Network, actualEvent *pb.Event) bool {
	actual, ok := actualEvent.EventType.(*pb.Event_Network)
	if !ok {
		return false
	}

	t.Logf("network event line %+v", actual)

	if expected.Network.GetComm() != "" && actual.Network.GetComm() != expected.Network.GetComm() {
		return false
	}

	if actual.Network.GetEventType() != expected.Network.GetEventType() {
		return false
	}

	if expected.Network.GetSrcPort() != 0 && actual.Network.GetSrcPort() != expected.Network.GetSrcPort() {
		return false
	}

	if expected.Network.GetSrcAddr() != "" && actual.Network.GetSrcAddr() != expected.Network.GetSrcAddr() {
		return false
	}

	return true
}

func dnsEventEqual(t *testing.T, expected *pb.Event_DnsQueryEvent, actualEvent *pb.Event) bool {
	actual, ok := actualEvent.EventType.(*pb.Event_DnsQuery)
	if !ok {
		return false
	}

	t.Logf("dns event line %+v", actual)

	if expected.GetQuery() != "" && actual.DnsQuery.GetQuery() != expected.GetQuery() {
		return false
	}

	if expected.GetIp() != "" && actual.DnsQuery.GetIp() != expected.GetIp() {
		return false
	}

	return true
}

func TestNetworkE2EFlow(t *testing.T) {
	networkTests := features.New("network tests").
		WithLabel("type", "e2e").
		Assess("test listen functionality", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			expectedEvent := &pb.Event_Network{
				Network: &pb.Event_NetworkEvent{
					EventType: pb.Event_NetworkEvent_LISTEN,
					Comm:      "nginx",
					SrcAddr:   "0.0.0.0",
					SrcPort:   80,
				},
			}
			namespace := "debug-" + rand.RandStringRunes(10)
			return helperTestNetEventIsRecorded(t, ctx, cfg, expectedEvent, nil, namespace, "-n", namespace, "-f", "../../examples/nginx.yaml")
		}).
		Assess("test connection functionality", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			namespace := "debug-" + rand.RandStringRunes(10)
			client, _ := support.NewK8sClient(ctx, t, cfg)
			if err := client.Resources().Create(ctx, &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}); err != nil {
				t.Fatal(err)
			}
			petclinicPostgresService := v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "petclinic-postgres",
					Namespace: namespace,
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:       "petclinic-postgres-port",
							Port:       5432,
							Protocol:   v1.ProtocolTCP,
							TargetPort: intstr.FromInt(5432),
						},
					},
					Selector: map[string]string{
						"app": "petclinic-postgres",
					},
					Type: v1.ServiceTypeClusterIP,
				},
			}
			if err := client.Resources(namespace).Create(ctx, &petclinicPostgresService); err != nil {
				t.Fatalf("failed to create postgres service: %v", err)
			}
			err := wait.For(conditions.New(client.Resources(namespace)).ResourceMatch(&petclinicPostgresService, func(object e2ek8s.Object) bool {
				service := object.(*v1.Service)
				return service.Spec.ClusterIP != ""
			}), wait.WithTimeout(time.Minute*1))
			if err != nil {
				t.Fatalf("failed to wait for cluster IP: %s", err)
			}
			expectedEvent := &pb.Event_Network{
				Network: &pb.Event_NetworkEvent{
					EventType: pb.Event_NetworkEvent_FAILED_CONNECTION,
					DstPort:   5432,
				},
			}
			expectedDnsEvent := &pb.Event_DnsQueryEvent{
				Query: fmt.Sprintf("petclinic-postgres.%s.svc.cluster.local.", namespace),
			}
			return helperTestNetEventIsRecorded(t, ctx, cfg, expectedEvent, expectedDnsEvent, namespace, "-n", namespace, "-f", "../../examples/tomcat-external-db.yaml")
		}).
		Feature()
	testenv.Test(t, networkTests)
}
