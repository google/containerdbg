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

package analyze

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"velostrata-internal.googlesource.com/containerdbg.git/proto"
)

func TestAnalyze(t *testing.T) {

	cases := []struct {
		name     string
		dataFile string
		expected *proto.AnalysisSummary
	}{
		{
			name:     "missing files",
			dataFile: "./testdata/collected.pb",
			expected: &proto.AnalysisSummary{
				ContainerSummaries: []*proto.AnalysisSummary_ContainerSummaryTuple{
					{
						Source: &proto.SourceId{},
						Summary: &proto.ContainerAnalysisSummary{
							MissingFiles: []string{
								"/etc/crypto-policies/back-ends/java.config",
								"/opt/jboss/wildfly/modules/com/sun/jsf-impl",
								"/opt/jboss/wildfly/modules/javax/faces/api",
								"/opt/jboss/wildfly/modules/org/jboss/as/jsf-injection",
								"/opt/jboss/wildfly/modules/system/add-ons",
								"Config/Config.properties",
								"Config/Properties/application-system.properties",
							},
							MoveFailures:       []*proto.ContainerAnalysisSummary_MoveFailure{},
							MissingLibraries:   []string{},
							ConnectionFailures: []*proto.ContainerAnalysisSummary_Connection{},
						},
					},
				},
			},
		},
		{
			name:     "EXDEV error",
			dataFile: "./testdata/rename-fail.pb",
			expected: &proto.AnalysisSummary{
				ContainerSummaries: []*proto.AnalysisSummary_ContainerSummaryTuple{
					{
						Source: &proto.SourceId{},
						Summary: &proto.ContainerAnalysisSummary{
							MissingFiles: []string{},
							MoveFailures: []*proto.ContainerAnalysisSummary_MoveFailure{
								{
									Source:      "f1",
									Destination: "f2",
								},
							},
							MissingLibraries:   []string{},
							ConnectionFailures: []*proto.ContainerAnalysisSummary_Connection{},
						},
					},
				},
			},
		},
		{
			name:     "connection failed to external service error",
			dataFile: "./testdata/petclinic-tomcat.pb",
			expected: &proto.AnalysisSummary{
				ContainerSummaries: []*proto.AnalysisSummary_ContainerSummaryTuple{
					{
						Source: &proto.SourceId{Type: "container", Id: "barp-tomcat-petclinic-tomcat-petclinic-5b2b62c3-bc9cb7f68-bxphj"},
						Summary: &proto.ContainerAnalysisSummary{
							MissingFiles: []string{
								"/usr/local/openjdk-11/conf/jndi.properties",
								"/usr/local/tomcat/work/Catalina/localhost/petclinic/SESSIONS.ser",
							},
							MoveFailures:     []*proto.ContainerAnalysisSummary_MoveFailure{},
							MissingLibraries: []string{},
							ConnectionFailures: []*proto.ContainerAnalysisSummary_Connection{
								{
									TargetIp: "10.108.0.105",
									Port:     5432,
								},
							},
							StaticIps: []*proto.ContainerAnalysisSummary_Connection{
								{
									TargetIp: "10.108.0.105",
									Port:     5432,
								},
							},
						},
					},
					{
						Source: &proto.SourceId{Type: "container", Id: "barp-tomcat-petclinic-tomcat-petclinic-5b2b62c3-bc9cb7f68-lk64p"},
						Summary: &proto.ContainerAnalysisSummary{
							MissingFiles: []string{
								"/usr/local/openjdk-11/conf/jndi.properties",
								"/usr/local/tomcat/work/Catalina/localhost/petclinic/SESSIONS.ser",
							},
							MoveFailures:     []*proto.ContainerAnalysisSummary_MoveFailure{},
							MissingLibraries: []string{},
							ConnectionFailures: []*proto.ContainerAnalysisSummary_Connection{
								{
									TargetIp: "10.108.0.105",
									Port:     5432,
								},
							},
							StaticIps: []*proto.ContainerAnalysisSummary_Connection{
								{
									TargetIp: "10.108.0.105",
									Port:     5432,
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "connection failed to external service error",
			dataFile: "./testdata/tomcat-ubuntu-host.pb",
			expected: &proto.AnalysisSummary{
				ContainerSummaries: []*proto.AnalysisSummary_ContainerSummaryTuple{
					{
						Source: &proto.SourceId{Type: "container", Id: "petclinic-tomcat-64fd676cd6-vf5x6"},
						Summary: &proto.ContainerAnalysisSummary{
							MissingFiles: []string{
								"/usr/local/openjdk-11/conf/jndi.properties",
								"/usr/local/tomcat/work/Catalina/localhost/petclinic/SESSIONS.ser",
							},
							MoveFailures:     []*proto.ContainerAnalysisSummary_MoveFailure{},
							MissingLibraries: []string{},
							ConnectionFailures: []*proto.ContainerAnalysisSummary_Connection{
								{
									TargetFqdn: "petclinic-postgres",
									TargetIp:   "10.55.240.47",
									Port:       5432,
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "connection failed to external service with dns",
			dataFile: "./testdata/tomcat-dns.pb",
			expected: &proto.AnalysisSummary{
				ContainerSummaries: []*proto.AnalysisSummary_ContainerSummaryTuple{
					{
						Source: &proto.SourceId{Type: "container", Id: "petclinic-tomcat-6cbcb9998-kgbk7"},
						Summary: &proto.ContainerAnalysisSummary{
							MissingFiles: []string{
								"/usr/local/openjdk-11/conf/jndi.properties",
								"/usr/local/tomcat/work/Catalina/localhost/petclinic/SESSIONS.ser",
							},
							MoveFailures:     []*proto.ContainerAnalysisSummary_MoveFailure{},
							MissingLibraries: []string{},
							ConnectionFailures: []*proto.ContainerAnalysisSummary_Connection{
								{
									TargetFqdn: "petclinic-postgres.default.svc.cluster.local",
									TargetIp:   "10.108.12.9",
									Port:       5432,
								},
							},
						},
					},
					{
						Source: &proto.SourceId{Type: "container", Id: "petclinic-tomcat-6cbcb9998-vxr4b"},
						Summary: &proto.ContainerAnalysisSummary{
							MissingFiles: []string{
								"/usr/local/openjdk-11/conf/jndi.properties",
								"/usr/local/tomcat/work/Catalina/localhost/petclinic/SESSIONS.ser",
							},
							MoveFailures:     []*proto.ContainerAnalysisSummary_MoveFailure{},
							MissingLibraries: []string{},
						},
					},
				},
			},
		},
		{
			name:     "dns query failure detection while ignoring successful queries",
			dataFile: "./testdata/failed_dns.pb",
			expected: &proto.AnalysisSummary{
				ContainerSummaries: []*proto.AnalysisSummary_ContainerSummaryTuple{
					{
						Source: &proto.SourceId{Type: "container", Id: "ubuntu-deployment-b59db5469-gvwtx"},
						Summary: &proto.ContainerAnalysisSummary{
							MissingFiles: []string{
								"/usr/bin/runc",
							},
							MoveFailures:       []*proto.ContainerAnalysisSummary_MoveFailure{},
							MissingLibraries:   []string{},
							ConnectionFailures: []*proto.ContainerAnalysisSummary_Connection{},
							DnsFailures: []*proto.ContainerAnalysisSummary_DnsFailure{
								{
									Query: "asdasdqwe12321312.com",
									Error: &proto.DnsQueryError{
										Code: proto.DnsQueryError_UNKNOWN,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			summary, err := Analyze(tc.dataFile, nil)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(summary, tc.expected, protocmp.Transform()); diff != "" {
				t.Fatalf("analysis summary is different from expected result (-got, +want):\n %s", diff)
			}

		})
	}
}
