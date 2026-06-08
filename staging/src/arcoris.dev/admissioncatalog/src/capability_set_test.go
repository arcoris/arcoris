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

package admissioncatalog

import "testing"

func TestCapabilitySet(t *testing.T) {
	set := NewCapabilitySet(
		NewOutcomeSet(OutcomeCapabilityAdmit, OutcomeCapabilityQueue),
		NewEffectSet(EffectCapabilityNone, EffectCapabilityQueued),
	)

	if !set.HasOutcome(OutcomeCapabilityAdmit) {
		t.Fatal("set does not contain admit outcome")
	}
	if !set.HasEffect(EffectCapabilityQueued) {
		t.Fatal("set does not contain queued effect")
	}
	if set.HasOutcome(OutcomeCapabilityDeny) {
		t.Fatal("set unexpectedly contains deny outcome")
	}
	if set.HasEffect(EffectCapabilityOwned) {
		t.Fatal("set unexpectedly contains owned effect")
	}
	if !set.IsValid() {
		t.Fatal("set is invalid")
	}
}

func TestCapabilitySetZeroIsValidAndUnspecified(t *testing.T) {
	var set CapabilitySet
	if !set.IsValid() {
		t.Fatal("zero capability set is invalid")
	}
	if !set.IsZero() {
		t.Fatal("zero capability set is not zero")
	}
}

func TestCapabilitySetRejectsUnknownBits(t *testing.T) {
	tests := []struct {
		name string
		set  CapabilitySet
	}{
		{name: "outcome", set: NewCapabilitySet(OutcomeSet(1<<7), NewEffectSet())},
		{name: "effect", set: NewCapabilitySet(NewOutcomeSet(), EffectSet(1<<7))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set.IsValid() {
				t.Fatal("unknown capability bits were accepted")
			}
		})
	}
}

func TestCapabilitySetDoesNotEnforceAdmissionResultValidity(t *testing.T) {
	set := NewCapabilitySet(
		NewOutcomeSet(OutcomeCapabilityDeny),
		NewEffectSet(EffectCapabilityOwned),
	)
	if !set.IsValid() {
		t.Fatal("descriptive capability combination should remain valid metadata")
	}
}

func TestCapabilitySetDimensionsAreIndependent(t *testing.T) {
	outcomeOnly := NewCapabilitySet(NewOutcomeSet(OutcomeCapabilityQueue), NewEffectSet())
	if !outcomeOnly.HasOutcome(OutcomeCapabilityQueue) {
		t.Fatal("outcome capability is missing")
	}
	if outcomeOnly.HasEffect(EffectCapabilityQueued) {
		t.Fatal("outcome capability leaked into effect dimension")
	}

	effectOnly := NewCapabilitySet(NewOutcomeSet(), NewEffectSet(EffectCapabilityQueued))
	if !effectOnly.HasEffect(EffectCapabilityQueued) {
		t.Fatal("effect capability is missing")
	}
	if effectOnly.HasOutcome(OutcomeCapabilityQueue) {
		t.Fatal("effect capability leaked into outcome dimension")
	}
}
