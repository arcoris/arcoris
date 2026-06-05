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

import (
	"testing"

	"arcoris.dev/admission"
)

func TestComponentDescriptors(t *testing.T) {
	descriptors := ComponentDescriptors()
	if len(descriptors) == 0 {
		t.Fatal("ComponentDescriptors returned empty slice")
	}
	requireComponent(t, descriptors, ComponentResilienceBulkhead)
	requireComponent(t, descriptors, ComponentResilienceRetryBudget)
	requireComponent(t, descriptors, ComponentResilienceDeadline)
}

func TestComponentDescriptorsAreFresh(t *testing.T) {
	first := ComponentDescriptors()
	second := ComponentDescriptors()
	first[0].ID = "mutated.component"

	if second[0].ID == "mutated.component" {
		t.Fatal("ComponentDescriptors shared mutable slice storage")
	}
}

func TestComponentDescriptorsAreValid(t *testing.T) {
	for _, descriptor := range ComponentDescriptors() {
		if !descriptor.IsValid() {
			t.Fatalf("descriptor is invalid: %+v", descriptor)
		}
		if descriptor.Summary == "" {
			t.Fatalf("descriptor has empty summary: %+v", descriptor)
		}
	}
}

func TestRetryBudgetComponentIDUsesSnakeCase(t *testing.T) {
	if got, want := ComponentResilienceRetryBudget, admission.ComponentID("resilience.retry_budget"); got != want {
		t.Fatalf("ComponentResilienceRetryBudget = %s, want %s", got, want)
	}
	if ComponentResilienceRetryBudget == admission.ComponentID("resilience.retrybudget") {
		t.Fatal("retry budget component ID still uses package-name spelling")
	}
}
