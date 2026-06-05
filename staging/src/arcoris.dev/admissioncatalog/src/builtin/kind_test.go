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

package builtin

import "testing"

func TestKindDescriptors(t *testing.T) {
	descriptors := KindDescriptors()
	if len(descriptors) == 0 {
		t.Fatal("KindDescriptors returned empty slice")
	}
	requireKind(t, descriptors, KindBulkhead)
	requireKind(t, descriptors, KindRetryBudget)
	requireKind(t, descriptors, KindDeadline)
}

func TestKindDescriptorsAreFresh(t *testing.T) {
	first := KindDescriptors()
	second := KindDescriptors()
	first[0].Kind = "mutated_kind"

	if second[0].Kind == "mutated_kind" {
		t.Fatal("KindDescriptors shared mutable slice storage")
	}
}

func TestKindDescriptorsAreValid(t *testing.T) {
	for _, descriptor := range KindDescriptors() {
		if !descriptor.IsValid() {
			t.Fatalf("descriptor is invalid: %+v", descriptor)
		}
		if descriptor.Summary == "" {
			t.Fatalf("descriptor has empty summary: %+v", descriptor)
		}
	}
}
