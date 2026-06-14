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

func TestCompareMetadataRequirementOrdersByKeyOperatorValues(t *testing.T) {
	first := metadataRequirement{key: "env", op: OperatorIn, values: []string{"prod"}}
	second := metadataRequirement{key: "tier", op: OperatorEquals, values: []string{"backend"}}
	if compareMetadataRequirement(first, second) >= 0 {
		t.Fatal("env requirement did not sort before tier requirement")
	}

	left := metadataRequirement{key: "env", op: OperatorEquals, values: []string{"prod"}}
	right := metadataRequirement{key: "env", op: OperatorIn, values: []string{"prod"}}
	if compareMetadataRequirement(left, right) >= 0 {
		t.Fatal("equals requirement did not sort before in requirement")
	}
}

func TestSameMetadataRequirement(t *testing.T) {
	left := metadataRequirement{key: "env", op: OperatorIn, values: []string{"prod", "qa"}}
	right := metadataRequirement{key: "env", op: OperatorIn, values: []string{"prod", "qa"}}
	if !sameMetadataRequirement(left, right) {
		t.Fatal("sameMetadataRequirement = false; want true")
	}

	right.values[1] = "staging"
	if sameMetadataRequirement(left, right) {
		t.Fatal("sameMetadataRequirement = true for different values")
	}
}

func TestCompareStringSlices(t *testing.T) {
	if compareStringSlices([]string{"prod"}, []string{"qa"}) >= 0 {
		t.Fatal("prod did not sort before qa")
	}
	if compareStringSlices([]string{"prod"}, []string{"prod", "qa"}) >= 0 {
		t.Fatal("shorter equal-prefix slice did not sort first")
	}
	if compareStringSlices([]string{"prod"}, []string{"prod"}) != 0 {
		t.Fatal("equal slices did not compare equal")
	}
}
