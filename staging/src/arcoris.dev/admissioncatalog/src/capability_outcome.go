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

// OutcomeCapability declares one outcome class that metadata says a component
// may produce.
//
// Outcome capabilities describe catalog surface area only. They do not make a
// runtime admission decision and do not change admission.Result validity.
type OutcomeCapability uint8

const (
	// OutcomeCapabilityAdmit declares that immediate admission is part of the
	// documented component surface.
	OutcomeCapabilityAdmit OutcomeCapability = 1 << iota

	// OutcomeCapabilityDeny declares that rejection is part of the documented
	// component surface.
	OutcomeCapabilityDeny

	// OutcomeCapabilityQueue declares that system-owned queued waiting is part
	// of the documented component surface.
	OutcomeCapabilityQueue

	// OutcomeCapabilityDefer declares that caller-owned retry or reconsideration
	// is part of the documented component surface.
	OutcomeCapabilityDefer
)

const knownOutcomeCapabilityMask = OutcomeCapabilityAdmit |
	OutcomeCapabilityDeny |
	OutcomeCapabilityQueue |
	OutcomeCapabilityDefer

// OutcomeSet is a compact set of declared outcome capabilities.
//
// The zero value is valid and means the descriptor does not specify an outcome
// surface. Unknown bits are invalid because their meaning would be ambiguous in
// documentation and configuration checks.
type OutcomeSet uint8

// NewOutcomeSet returns a set containing capabilities.
//
// Unknown capability bits are preserved so IsValid can report malformed
// metadata instead of silently normalizing it.
func NewOutcomeSet(capabilities ...OutcomeCapability) OutcomeSet {
	var set OutcomeSet
	for _, capability := range capabilities {
		set = set.With(capability)
	}
	return set
}

// With returns s plus capability.
func (s OutcomeSet) With(capability OutcomeCapability) OutcomeSet {
	return s | OutcomeSet(capability)
}

// Has reports whether s contains capability.
func (s OutcomeSet) Has(capability OutcomeCapability) bool {
	return capability != 0 && OutcomeCapability(s)&capability == capability
}

// IsValid reports whether s contains only known outcome capability bits.
func (s OutcomeSet) IsValid() bool {
	return OutcomeCapability(s)&^knownOutcomeCapabilityMask == 0
}

// IsZero reports whether s leaves outcome capabilities unspecified.
func (s OutcomeSet) IsZero() bool {
	return s == 0
}
