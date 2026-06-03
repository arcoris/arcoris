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
	"arcoris.dev/admissioncatalog"
)

func TestComponentDescriptors(t *testing.T) {
	t.Parallel()

	kinds := NewKindRegistry()
	descriptors := ComponentDescriptors()
	wantComponents := map[admission.ComponentID]bool{
		"resilience.bulkhead":    false,
		"resilience.deadline":    false,
		"resilience.retrybudget": false,
	}

	for _, descriptor := range descriptors {
		if !descriptor.IsValid() {
			t.Fatalf("standard descriptor should be valid: %+v", descriptor)
		}
		if !kinds.Contains(descriptor.Kind) {
			t.Fatalf("standard descriptor references unknown kind %q", descriptor.Kind)
		}
		if found, known := wantComponents[descriptor.ID]; !known {
			t.Fatalf("unexpected standard component %q", descriptor.ID)
		} else if found {
			t.Fatalf("duplicate standard component %q", descriptor.ID)
		}
		wantComponents[descriptor.ID] = true
	}
	for id, found := range wantComponents {
		if !found {
			t.Fatalf("missing standard component %q", id)
		}
	}
}

func TestComponentDescriptorsReturnsCopy(t *testing.T) {
	t.Parallel()

	descriptors := ComponentDescriptors()
	descriptors[0].ID = "resilience.mutated"

	fresh := ComponentDescriptors()
	if fresh[0].ID == "resilience.mutated" {
		t.Fatal("mutating returned descriptors should not mutate standard catalog")
	}
}

func TestComponentDescriptorCapabilities(t *testing.T) {
	t.Parallel()

	descriptors := ComponentDescriptors()
	tests := []struct {
		id           admission.ComponentID
		kind         admission.ComponentKind
		capabilities []admissioncatalog.Capability
	}{
		{
			id:   "resilience.bulkhead",
			kind: KindBulkhead,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectOwned,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			id:   "resilience.retrybudget",
			kind: KindRetryBudget,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectCommitted,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			id:   "resilience.deadline",
			kind: KindDeadline,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.id.String(), func(t *testing.T) {
			t.Parallel()

			descriptor := requireComponentDescriptor(
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
			if test.id == "resilience.deadline" &&
				descriptor.Capabilities.Has(admissioncatalog.CapabilityDefer) {
				t.Fatal("resilience.deadline should not advertise defer")
			}
		})
	}
}

func TestNewComponentRegistry(t *testing.T) {
	t.Parallel()

	registry := NewComponentRegistry(NewKindRegistry())
	for _, descriptor := range ComponentDescriptors() {
		if got, ok := registry.Lookup(descriptor.ID); !ok || got != descriptor {
			t.Fatalf("Lookup(%q) = (%+v, %v), want standard descriptor", descriptor.ID, got, ok)
		}
	}
}

func requireComponentDescriptor(
	t *testing.T,
	descriptors []admissioncatalog.ComponentDescriptor,
	id admission.ComponentID,
) admissioncatalog.ComponentDescriptor {
	t.Helper()

	for _, descriptor := range descriptors {
		if descriptor.ID == id {
			return descriptor
		}
	}
	t.Fatalf("missing standard component descriptor %q", id)
	return admissioncatalog.ComponentDescriptor{}
}
