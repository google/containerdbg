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
)

func TestAnalyze(t *testing.T) {

	cases := []struct {
		name     string
		dataFile string
		expected AnalyzeSummary
	}{
		{
			name:     "missing files",
			dataFile: "./testdata/collected.pb",
			expected: AnalyzeSummary{
				MissingFiles: []string{
					"/etc/crypto-policies/back-ends/java.config",
					"/opt/jboss/wildfly/modules/com/sun/jsf-impl",
					"/opt/jboss/wildfly/modules/javax/faces/api",
					"/opt/jboss/wildfly/modules/org/jboss/as/jsf-injection",
					"/opt/jboss/wildfly/modules/system/add-ons",
					"Config/Config.properties",
					"Config/Properties/application-system.properties",
				},
				ExdevFailures: []string{},
				MissingLibs:   []string{},
			},
		},
		{
			name:     "EXDEV error",
			dataFile: "./testdata/rename-fail.pb",
			expected: AnalyzeSummary{
				MissingFiles: []string{},
				ExdevFailures: []string{
					"/opt/jboss-7.1.1/standalone/configuration/nickelodeonarabia.com/standalone_xml_history/current",
				},
				MissingLibs: []string{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			summary, err := Analyze(tc.dataFile, nil)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(&tc.expected, summary); diff != "" {
				t.Fatalf("analysis summary is different from expected result (-got, +want):\n %s", diff)
			}

		})
	}
}
