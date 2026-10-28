// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package istiovalues

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/istio-ecosystem/sail-operator/api/v1"

	"k8s.io/utils/ptr"
)

func TestApplyOverrides(t *testing.T) {
	tests := []struct {
		name           string
		revision       string
		namespace      string
		values         v1.Values
		expectedValues v1.Values
	}{
		{
			name:      "default-revision",
			revision:  v1.DefaultRevision,
			namespace: "ns1",
			values:    v1.Values{},
			expectedValues: v1.Values{
				Revision: ptr.To(""),
				Global: &v1.GlobalConfig{
					IstioNamespace: ptr.To("ns1"),
				},
				DefaultRevision: ptr.To(""),
			},
		},
		{
			name:      "non-default-revision",
			revision:  "my-revision",
			namespace: "ns1",
			values:    v1.Values{},
			expectedValues: v1.Values{
				Revision: ptr.To("my-revision"),
				Global: &v1.GlobalConfig{
					IstioNamespace: ptr.To("ns1"),
				},
				DefaultRevision: ptr.To(""),
			},
		},
		{
			name:      "revision-already-set",
			revision:  "my-revision",
			namespace: "ns1",
			values: v1.Values{
				Revision: ptr.To("this-should-be-overridden"),
			},
			expectedValues: v1.Values{
				Revision: ptr.To("my-revision"),
				Global: &v1.GlobalConfig{
					IstioNamespace: ptr.To("ns1"),
				},
				DefaultRevision: ptr.To(""),
			},
		},
		{
			name:      "namespace-already-set",
			revision:  "my-revision",
			namespace: "ns1",
			values: v1.Values{
				Global: &v1.GlobalConfig{
					IstioNamespace: ptr.To("this-should-be-overridden"),
				},
			},
			expectedValues: v1.Values{
				Revision: ptr.To("my-revision"),
				Global: &v1.GlobalConfig{
					IstioNamespace: ptr.To("ns1"),
				},
				DefaultRevision: ptr.To(""),
			},
		},
		{
			name:      "defaultRevision-is-ignored-when-set-by-user",
			revision:  "my-revision",
			namespace: "ns1",
			values: v1.Values{
				DefaultRevision: ptr.To("my-revision"),
			},
			expectedValues: v1.Values{
				Revision: ptr.To("my-revision"),
				Global: &v1.GlobalConfig{
					IstioNamespace: ptr.To("ns1"),
				},
				DefaultRevision: ptr.To(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyOverrides(tt.revision, tt.namespace, &tt.values)
			if diff := cmp.Diff(tt.expectedValues, tt.values); diff != "" {
				t.Fatalf("unexpected output values (-expected +actual):\n%s", diff)
			}
		})
	}
}
