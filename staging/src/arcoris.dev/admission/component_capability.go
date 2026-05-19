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

// Capability describes one supported outcome or effect class of a component.
//
// Capabilities are metadata, not enforcement. A descriptor can advertise that a
// component may admit, deny, queue, defer, or produce a particular effect class,
// while the actual decision still lives in Result.
type Capability uint16

const (
	// CapabilityAdmit declares that a component can immediately admit work.
	CapabilityAdmit Capability = 1 << iota

	// CapabilityDeny declares that a component can reject an admission attempt.
	CapabilityDeny

	// CapabilityQueue declares that a component can take ownership of waiting
	// work instead of returning control to the caller immediately.
	CapabilityQueue

	// CapabilityDefer declares that a component can ask the caller to retry or
	// reconsider later without taking ownership of waiting work.
	CapabilityDefer

	// CapabilityEffectNone declares that a component can produce decisions with
	// no ownership, accounting, queueing, or other side effects.
	CapabilityEffectNone

	// CapabilityEffectCommitted declares that a component can commit a
	// spend-only side effect, such as consuming a retry-budget attempt.
	CapabilityEffectCommitted

	// CapabilityEffectOwned declares that a component can return a caller-owned
	// grant that has a later domain lifecycle.
	CapabilityEffectOwned

	// CapabilityEffectQueued declares that a component can accept
	// system-owned waiting state.
	CapabilityEffectQueued
)

// knownCapabilityMask contains every currently defined Capability bit.
//
// It is kept as a single mask so CapabilitySet.IsValid can reject future or
// corrupted bits without allocating or iterating.
const knownCapabilityMask = CapabilityAdmit |
	CapabilityDeny |
	CapabilityQueue |
	CapabilityDefer |
	CapabilityEffectNone |
	CapabilityEffectCommitted |
	CapabilityEffectOwned |
	CapabilityEffectQueued

// CapabilitySet is compact catalog metadata for outcome and effect capabilities.
//
// The set intentionally keeps outcome-like bits, such as admit or deny, and
// effect-like bits, such as owned or queued, in one small value. It is
// descriptive metadata for catalogs, docs, and config validation; it is not
// enforcement. Result validity is still governed by Decision, Effect, and
// grant-shape invariants. If registry users later need stronger dimensional
// validation, the outcome and effect dimensions can be split without changing
// Result semantics. The zero set is valid and means that capabilities are
// unspecified.
type CapabilitySet uint16

// NewCapabilitySet returns a set containing capabilities.
//
// Unknown bits are intentionally preserved by With so IsValid can detect
// malformed descriptors instead of silently normalizing them.
func NewCapabilitySet(capabilities ...Capability) CapabilitySet {
	var set CapabilitySet
	for _, capability := range capabilities {
		set = set.With(capability)
	}
	return set
}

// With returns s plus capability.
//
// Invalid capability bits are preserved. This makes descriptor validation
// strict: callers can discover malformed metadata through IsValid instead of
// losing the original error by truncation.
func (s CapabilitySet) With(capability Capability) CapabilitySet {
	return s | CapabilitySet(capability)
}

// Has reports whether s contains capability.
//
// A zero capability is never reported as present because zero is not a real bit
// in the advertised capability vocabulary.
func (s CapabilitySet) Has(capability Capability) bool {
	return capability != 0 && Capability(s)&capability == capability
}

// IsValid reports whether s contains only known capability bits.
//
// The zero value is valid and means unspecified capabilities. Non-zero unknown
// bits are rejected because they would make descriptor meaning ambiguous.
func (s CapabilitySet) IsValid() bool {
	return Capability(s)&^knownCapabilityMask == 0
}

// IsZero reports whether s leaves capabilities unspecified.
func (s CapabilitySet) IsZero() bool {
	return s == 0
}
