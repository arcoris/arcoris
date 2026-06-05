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

// EffectCapability declares one side-effect class that metadata says a
// component may produce.
//
// Effect capabilities are documentation and validation metadata. They do not
// implement rollback, ownership, leases, queue state, or admission.Result shape
// checks.
type EffectCapability uint8

const (
	// EffectCapabilityNone declares that no-side-effect decisions are part of
	// the documented component surface.
	EffectCapabilityNone EffectCapability = 1 << iota

	// EffectCapabilityCommitted declares that spend-only committed side effects
	// are part of the documented component surface.
	EffectCapabilityCommitted

	// EffectCapabilityOwned declares that caller-owned grants are part of the
	// documented component surface.
	EffectCapabilityOwned

	// EffectCapabilityQueued declares that system-owned queued work is part of
	// the documented component surface.
	EffectCapabilityQueued
)

const knownEffectCapabilityMask = EffectCapabilityNone |
	EffectCapabilityCommitted |
	EffectCapabilityOwned |
	EffectCapabilityQueued

// EffectSet is a compact set of declared side-effect capabilities.
//
// The zero value is valid and means the descriptor does not specify an effect
// surface. Unknown bits are invalid because they would make catalog metadata
// impossible to interpret safely.
type EffectSet uint8

// NewEffectSet returns a set containing capabilities.
//
// Unknown capability bits are preserved so IsValid can report malformed
// metadata instead of silently normalizing it.
func NewEffectSet(capabilities ...EffectCapability) EffectSet {
	var set EffectSet
	for _, capability := range capabilities {
		set = set.With(capability)
	}
	return set
}

// With returns s plus capability.
func (s EffectSet) With(capability EffectCapability) EffectSet {
	return s | EffectSet(capability)
}

// Has reports whether s contains capability.
func (s EffectSet) Has(capability EffectCapability) bool {
	return capability != 0 && EffectCapability(s)&capability == capability
}

// IsValid reports whether s contains only known effect capability bits.
func (s EffectSet) IsValid() bool {
	return EffectCapability(s)&^knownEffectCapabilityMask == 0
}

// IsZero reports whether s leaves effect capabilities unspecified.
func (s EffectSet) IsZero() bool {
	return s == 0
}
