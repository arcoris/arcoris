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

func TestReasonDescriptors(t *testing.T) {
	t.Parallel()

	descriptors := ReasonDescriptors()
	wantReasons := map[admission.Reason]bool{
		admission.ReasonAdmitted: false,
		admission.ReasonDenied:   false,
		admission.ReasonQueued:   false,
		admission.ReasonDeferred: false,
		ReasonCapacityExhausted:  false,
		ReasonBudgetExhausted:    false,
		ReasonRateLimited:        false,
		ReasonOverloaded:         false,
		ReasonBackpressured:      false,
		ReasonClosed:             false,
		ReasonDraining:           false,
		ReasonDeadlineExceeded:   false,
		ReasonCanceled:           false,
		ReasonPolicyDenied:       false,
	}

	for _, descriptor := range descriptors {
		if !descriptor.IsValid() {
			t.Fatalf("standard descriptor should be valid: %+v", descriptor)
		}
		if descriptor.Capabilities.IsZero() {
			t.Fatalf("standard descriptor should declare capabilities: %+v", descriptor)
		}
		if found, known := wantReasons[descriptor.Reason]; !known {
			t.Fatalf("unexpected standard reason %q", descriptor.Reason)
		} else if found {
			t.Fatalf("duplicate standard reason %q", descriptor.Reason)
		}
		wantReasons[descriptor.Reason] = true
	}
	for reason, found := range wantReasons {
		if !found {
			t.Fatalf("missing standard reason %q", reason)
		}
	}
}

func TestReasonDescriptorsReturnsCopy(t *testing.T) {
	t.Parallel()

	descriptors := ReasonDescriptors()
	descriptors[0].Reason = "mutated_reason"

	fresh := ReasonDescriptors()
	if fresh[0].Reason == "mutated_reason" {
		t.Fatal("mutating returned descriptors should not mutate standard catalog")
	}
}

func TestReasonDescriptorCapabilities(t *testing.T) {
	t.Parallel()

	descriptors := ReasonDescriptors()
	tests := []struct {
		reason       admission.Reason
		capabilities []admissioncatalog.Capability
	}{
		{
			reason: admission.ReasonAdmitted,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityAdmit,
				admissioncatalog.CapabilityEffectNone,
				admissioncatalog.CapabilityEffectCommitted,
				admissioncatalog.CapabilityEffectOwned,
			},
		},
		{
			reason: admission.ReasonDenied,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			reason: admission.ReasonQueued,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityQueue,
				admissioncatalog.CapabilityEffectQueued,
			},
		},
		{
			reason: admission.ReasonDeferred,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityDefer,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			reason: ReasonCapacityExhausted,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			reason: ReasonBudgetExhausted,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			reason: ReasonRateLimited,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			reason: ReasonDeadlineExceeded,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			},
		},
		{
			reason: ReasonPolicyDenied,
			capabilities: []admissioncatalog.Capability{
				admissioncatalog.CapabilityDeny,
				admissioncatalog.CapabilityEffectNone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.reason.String(), func(t *testing.T) {
			t.Parallel()

			descriptor := requireReasonDescriptor(
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

func TestNewReasonRegistry(t *testing.T) {
	t.Parallel()

	registry := NewReasonRegistry()
	for _, descriptor := range ReasonDescriptors() {
		if got, ok := registry.Lookup(descriptor.Reason); !ok || got != descriptor {
			t.Fatalf("Lookup(%q) = (%+v, %v), want standard descriptor", descriptor.Reason, got, ok)
		}
	}
}

func requireReasonDescriptor(
	t *testing.T,
	descriptors []admissioncatalog.ReasonDescriptor,
	reason admission.Reason,
) admissioncatalog.ReasonDescriptor {
	t.Helper()

	for _, descriptor := range descriptors {
		if descriptor.Reason == reason {
			return descriptor
		}
	}
	t.Fatalf("missing standard reason descriptor %q", reason)
	return admissioncatalog.ReasonDescriptor{}
}
