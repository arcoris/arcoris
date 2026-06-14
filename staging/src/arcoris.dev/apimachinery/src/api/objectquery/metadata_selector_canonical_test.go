// Copyright 2026 The ARCORIS Authors.
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

package objectquery

import "testing"

func TestCanonicalMetadataRequirementsSortsDeduplicatesAndDetaches(t *testing.T) {
	raw := []metadataRequirement{
		{key: "tier", op: OperatorEquals, values: []string{"backend"}},
		{key: "env", op: OperatorIn, values: []string{"qa", "prod"}},
		{key: "env", op: OperatorIn, values: []string{"prod", "qa"}},
	}

	got, err := canonicalMetadataRequirements("query.labels.requirements", raw, validateLabelKey, validateLabelValue)
	requireNoError(t, err)
	raw[1].values[0] = "mutated"

	if len(got) != 2 {
		t.Fatalf("len = %d; want 2", len(got))
	}
	requireRequirement(t, got[0], "env", OperatorIn, "prod", "qa")
	requireRequirement(t, got[1], "tier", OperatorEquals, "backend")
}

func TestCompactMetadataRequirementsEmpty(t *testing.T) {
	if got := compactMetadataRequirements(nil); got != nil {
		t.Fatalf("compactMetadataRequirements(nil) = %#v; want nil", got)
	}
}
