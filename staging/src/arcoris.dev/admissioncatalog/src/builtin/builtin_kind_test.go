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

func TestKindDescriptors(t *testing.T) {
	t.Parallel()

	descriptors := KindDescriptors()
	wantKinds := map[admission.ComponentKind]bool{
		KindBulkhead:        false,
		KindRetryBudget:     false,
		KindDeadline:        false,
		KindRateLimiter:     false,
		KindQueue:           false,
		KindScheduler:       false,
		KindWorkerPool:      false,
		KindOverloadGate:    false,
		KindTenantIsolation: false,
	}

	for _, descriptor := range descriptors {
		if !descriptor.IsValid() {
			t.Fatalf("standard descriptor should be valid: %+v", descriptor)
		}
		if descriptor.Capabilities.IsZero() {
			t.Fatalf("standard descriptor should declare capabilities: %+v", descriptor)
		}
		if _, known := wantKinds[descriptor.Kind]; !known {
			t.Fatalf("unexpected standard kind %q", descriptor.Kind)
		}
		wantKinds[descriptor.Kind] = true
	}
	for kind, found := range wantKinds {
		if !found {
			t.Fatalf("missing standard kind %q", kind)
		}
	}
}

func TestKindDescriptorsReturnsCopy(t *testing.T) {
	t.Parallel()

	descriptors := KindDescriptors()
	descriptors[0].Kind = "mutated_kind"

	fresh := KindDescriptors()
	if fresh[0].Kind == "mutated_kind" {
		t.Fatal("mutating returned descriptors should not mutate standard catalog")
	}
}

func TestKindDescriptorCapabilities(t *testing.T) {
	t.Parallel()

	descriptors := KindDescriptors()
	tests := []struct {
		kind         admission.ComponentKind
		capabilities []admissioncatalog.Capability
	}{
		{
			kind: KindBulkhead,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectOwned,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			kind: KindRetryBudget,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectCommitted,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			kind: KindDeadline,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			kind: KindQueue,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityQueue,
				admissioncatalog.CapabilityEffectQueued,
			},
		},
		{
			kind: KindWorkerPool,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityQueue,
				admissioncatalog.CapabilityEffectOwned,
				admissioncatalog.CapabilityEffectQueued,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.kind.String(), func(t *testing.T) {
			t.Parallel()

			descriptor := requireKindDescriptor(
				t,
				descriptors,
				test.kind,
			)
			for _, capability := range test.capabilities {
				requireCapability(t, descriptor.Capabilities, capability)
			}
		})
	}
}

func TestNewKindRegistry(t *testing.T) {
	t.Parallel()

	registry := NewKindRegistry()
	for _, descriptor := range KindDescriptors() {
		if got, ok := registry.Lookup(descriptor.Kind); !ok || got != descriptor {
			t.Fatalf("Lookup(%q) = (%+v, %v), want standard descriptor", descriptor.Kind, got, ok)
		}
	}
}

func requireKindDescriptor(
	t *testing.T,
	descriptors []admissioncatalog.ComponentKindDescriptor,
	kind admission.ComponentKind,
) admissioncatalog.ComponentKindDescriptor {
	t.Helper()

	for _, descriptor := range descriptors {
		if descriptor.Kind == kind {
			return descriptor
		}
	}
	t.Fatalf("missing standard kind descriptor %q", kind)
	return admissioncatalog.ComponentKindDescriptor{}
}
