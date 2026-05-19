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

func TestBuiltinReasonDescriptors(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinReasonDescriptors()
	wantReasons := map[Reason]bool{
		ReasonAdmitted:          false,
		ReasonDenied:            false,
		ReasonQueued:            false,
		ReasonDeferred:          false,
		ReasonCapacityExhausted: false,
		ReasonBudgetExhausted:   false,
		ReasonRateLimited:       false,
		ReasonOverloaded:        false,
		ReasonBackpressured:     false,
		ReasonClosed:            false,
		ReasonDraining:          false,
		ReasonDeadlineExceeded:  false,
		ReasonCanceled:          false,
		ReasonPolicyDenied:      false,
	}

	for _, descriptor := range descriptors {
		if !descriptor.IsValid() {
			t.Fatalf("built-in descriptor should be valid: %+v", descriptor)
		}
		if descriptor.Capabilities.IsZero() {
			t.Fatalf("built-in descriptor should declare capabilities: %+v", descriptor)
		}
		if found, known := wantReasons[descriptor.Reason]; !known {
			t.Fatalf("unexpected built-in reason %q", descriptor.Reason)
		} else if found {
			t.Fatalf("duplicate built-in reason %q", descriptor.Reason)
		}
		wantReasons[descriptor.Reason] = true
	}
	for reason, found := range wantReasons {
		if !found {
			t.Fatalf("missing built-in reason %q", reason)
		}
	}
}

func TestBuiltinReasonDescriptorsReturnsCopy(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinReasonDescriptors()
	descriptors[0].Reason = "mutated_reason"

	fresh := BuiltinReasonDescriptors()
	if fresh[0].Reason == "mutated_reason" {
		t.Fatal("mutating returned descriptors should not mutate built-in catalog")
	}
}

func TestBuiltinReasonDescriptorCapabilities(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinReasonDescriptors()
	tests := []struct {
		reason       Reason
		capabilities []Capability
	}{
		{
			reason: ReasonAdmitted,
			capabilities: []Capability{
				CapabilityAdmit,
				CapabilityEffectNone,
				CapabilityEffectCommitted,
				CapabilityEffectOwned,
			},
		},
		{
			reason: ReasonDenied,
			capabilities: []Capability{
				CapabilityDeny,
				CapabilityEffectNone,
			},
		},
		{
			reason: ReasonQueued,
			capabilities: []Capability{
				CapabilityQueue,
				CapabilityEffectQueued,
			},
		},
		{
			reason: ReasonDeferred,
			capabilities: []Capability{
				CapabilityDefer,
				CapabilityEffectNone,
			},
		},
		{
			reason: ReasonCapacityExhausted,
			capabilities: []Capability{
				CapabilityDeny,
				CapabilityEffectNone,
			},
		},
		{
			reason: ReasonBudgetExhausted,
			capabilities: []Capability{
				CapabilityDeny,
				CapabilityEffectNone,
			},
		},
		{
			reason: ReasonRateLimited,
			capabilities: []Capability{
				CapabilityDeny,
				CapabilityEffectNone,
			},
		},
		{
			reason: ReasonDeadlineExceeded,
			capabilities: []Capability{
				CapabilityDeny,
				CapabilityEffectNone,
			},
		},
		{
			reason: ReasonPolicyDenied,
			capabilities: []Capability{
				CapabilityDeny,
				CapabilityEffectNone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.reason.String(), func(t *testing.T) {
			t.Parallel()

			descriptor := requireBuiltinReasonDescriptor(
				t,
				descriptors,
				test.reason,
			)
			for _, capability := range test.capabilities {
				requireCapability(t, descriptor.Capabilities, capability)
			}
		})
	}
}

func TestNewBuiltinReasonRegistry(t *testing.T) {
	t.Parallel()

	registry := NewBuiltinReasonRegistry()
	for _, descriptor := range BuiltinReasonDescriptors() {
		if got, ok := registry.Lookup(descriptor.Reason); !ok || got != descriptor {
			t.Fatalf("Lookup(%q) = (%+v, %v), want built-in descriptor", descriptor.Reason, got, ok)
		}
	}
}

func requireBuiltinReasonDescriptor(
	t *testing.T,
	descriptors []ReasonDescriptor,
	reason Reason,
) ReasonDescriptor {
	t.Helper()

	for _, descriptor := range descriptors {
		if descriptor.Reason == reason {
			return descriptor
		}
	}
	t.Fatalf("missing built-in reason descriptor %q", reason)
	return ReasonDescriptor{}
}
