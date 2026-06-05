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

// CapabilitySet groups declared outcome and side-effect capabilities.
//
// Capabilities are metadata declarations, not runtime enforcement. The zero
// value is valid and means the descriptor leaves both dimensions unspecified.
type CapabilitySet struct {
	// Outcomes declares the outcome classes described by a descriptor.
	Outcomes OutcomeSet

	// Effects declares the side-effect classes described by a descriptor.
	Effects EffectSet
}

// NewCapabilitySet returns a capability set from explicit outcome and effect
// dimensions.
func NewCapabilitySet(outcomes OutcomeSet, effects EffectSet) CapabilitySet {
	return CapabilitySet{
		Outcomes: outcomes,
		Effects:  effects,
	}
}

// HasOutcome reports whether s declares capability in the outcome dimension.
func (s CapabilitySet) HasOutcome(capability OutcomeCapability) bool {
	return s.Outcomes.Has(capability)
}

// HasEffect reports whether s declares capability in the side-effect dimension.
func (s CapabilitySet) HasEffect(capability EffectCapability) bool {
	return s.Effects.Has(capability)
}

// IsValid reports whether both capability dimensions contain only known bits.
func (s CapabilitySet) IsValid() bool {
	return s.Outcomes.IsValid() && s.Effects.IsValid()
}

// IsZero reports whether both capability dimensions are unspecified.
func (s CapabilitySet) IsZero() bool {
	return s.Outcomes.IsZero() && s.Effects.IsZero()
}
