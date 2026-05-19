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

package admission

import "testing"

func TestBuiltinComponentDescriptors(t *testing.T) {
	t.Parallel()

	kinds := NewBuiltinKindRegistry()
	descriptors := BuiltinComponentDescriptors()
	wantComponents := map[ComponentID]bool{
		"resilience.bulkhead":    false,
		"resilience.deadline":    false,
		"resilience.retrybudget": false,
	}

	for _, descriptor := range descriptors {
		if !descriptor.IsValid() {
			t.Fatalf("built-in descriptor should be valid: %+v", descriptor)
		}
		if !kinds.Contains(descriptor.Kind) {
			t.Fatalf("built-in descriptor references unknown kind %q", descriptor.Kind)
		}
		if found, known := wantComponents[descriptor.ID]; !known {
			t.Fatalf("unexpected built-in component %q", descriptor.ID)
		} else if found {
			t.Fatalf("duplicate built-in component %q", descriptor.ID)
		}
		wantComponents[descriptor.ID] = true
	}
	for id, found := range wantComponents {
		if !found {
			t.Fatalf("missing built-in component %q", id)
		}
	}
}

func TestBuiltinComponentDescriptorsReturnsCopy(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinComponentDescriptors()
	descriptors[0].ID = "resilience.mutated"

	fresh := BuiltinComponentDescriptors()
	if fresh[0].ID == "resilience.mutated" {
		t.Fatal("mutating returned descriptors should not mutate built-in catalog")
	}
}

func TestBuiltinComponentDescriptorCapabilities(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinComponentDescriptors()
	tests := []struct {
		id           ComponentID
		kind         ComponentKind
		capabilities []Capability
	}{
		{
			id:   "resilience.bulkhead",
			kind: KindBulkhead,
			capabilities: []Capability{
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityEffectOwned,
				CapabilityEffectNone,
			},
		},
		{
			id:   "resilience.retrybudget",
			kind: KindRetryBudget,
			capabilities: []Capability{
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityEffectCommitted,
				CapabilityEffectNone,
			},
		},
		{
			id:   "resilience.deadline",
			kind: KindDeadline,
			capabilities: []Capability{
				CapabilityAdmit,
				CapabilityDeny,
				CapabilityEffectNone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.id.String(), func(t *testing.T) {
			t.Parallel()

			descriptor := requireBuiltinComponentDescriptor(
				t,
				descriptors,
				test.id,
			)
			if descriptor.Kind != test.kind {
				t.Fatalf("kind = %q, want %q", descriptor.Kind, test.kind)
			}
			for _, capability := range test.capabilities {
				requireCapability(t, descriptor.Capabilities, capability)
			}
		})
	}
}

func TestNewBuiltinComponentRegistry(t *testing.T) {
	t.Parallel()

	registry := NewBuiltinComponentRegistry(NewBuiltinKindRegistry())
	for _, descriptor := range BuiltinComponentDescriptors() {
		if got, ok := registry.Lookup(descriptor.ID); !ok || got != descriptor {
			t.Fatalf("Lookup(%q) = (%+v, %v), want built-in descriptor", descriptor.ID, got, ok)
		}
	}
}

func requireBuiltinComponentDescriptor(
	t *testing.T,
	descriptors []ComponentDescriptor,
	id ComponentID,
) ComponentDescriptor {
	t.Helper()

	for _, descriptor := range descriptors {
		if descriptor.ID == id {
			return descriptor
		}
	}
	t.Fatalf("missing built-in component descriptor %q", id)
	return ComponentDescriptor{}
}
