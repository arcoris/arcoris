/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package atomicx

import "sync/atomic"

// PaddedUint64 is an atomic uint64 isolated from neighboring fields by explicit
// leading padding and trailing line-completion padding.
//
// PaddedUint64 is the lowest-level unsigned padded atomic primitive in this
// package. It intentionally exposes raw unsigned atomic operations and does not
// impose higher-level accounting semantics. In particular:
//
//   - Add follows ordinary uint64 arithmetic and may wrap;
//   - Store and Swap are allowed;
//   - no overflow invariant is enforced;
//   - no non-negative gauge invariant is enforced;
//   - no lifetime-counter reset policy is enforced.
//
// This separation is intentional. PaddedUint64 is the raw storage cell used by
// higher-level atomicx types. Use semantic wrappers when the value represents a
// specific kind of runtime accounting state:
//
//   - use Uint64Counter for monotonic lifetime event counters;
//   - use Uint64Gauge for non-negative current quantities;
//   - use Int64Gauge for signed current quantities or correction values.
//
// PaddedUint64 is useful for hot component/runtime fields such as queue-depth
// state, in-flight operation counts, cache statistics, admission counters,
// dispatch counters, controller-local accounting values, and other frequently
// written uint64 fields shared across goroutines.
//
// Padding reduces false sharing. False sharing can occur when independent hot
// variables occupy the same CPU cache line and different goroutines update them
// from different cores. Even though the logical variables are unrelated, writes
// to one value can invalidate the cache line containing another value.
//
// Layout:
//
//	[noCopy marker][leading pad][atomic uint64][trailing line-completion pad]
//
// The layout follows the default 64-byte ARCORIS policy, not a hardware
// guarantee:
//
//   - the leading pad separates the atomic value from the previous field;
//   - the trailing line-completion pad fills the rest of the atomic value's
//     64-byte slot;
//   - the trailing completion matters when padded values appear in arrays,
//     slices, or before following struct fields.
//
// The noCopy marker is a zero-size static-analysis marker. It does not
// participate in synchronization and does not provide runtime protection. Its
// purpose is to make accidental value copies visible to tools such as:
//
//	go vet -copylocks
//
// PaddedUint64 is zero-value usable.
//
// PaddedUint64 must not be copied after first use. Copying a live atomic value
// can split one logical state cell into independent copies and produce incorrect
// accounting. Treat PaddedUint64 like sync/atomic typed values: construct it in
// place, pass it by pointer when sharing, and do not copy containing structs
// after the value becomes active.
type PaddedUint64 struct {
	noCopy noCopy
	_      CacheLinePad
	value  atomic.Uint64
	_      [CacheLinePadSize - atomicUint64Size]byte
}

// Load atomically returns the current uint64 value.
//
// Load provides the same synchronization semantics as atomic.Uint64.Load. It
// observes exactly one atomic value. It does not make a larger multi-field
// structure globally consistent.
//
// If a caller needs a consistent snapshot across multiple counters, gauges, or
// state fields, the caller must provide additional synchronization at the owner
// level.
func (p *PaddedUint64) Load() uint64 {
	return p.value.Load()
}

// Store atomically replaces the current uint64 value.
//
// Store is a raw state publication operation. It is appropriate for
// initialization, tests, owner-controlled publication, or explicit handoff
// semantics.
//
// Store must be used carefully for values with semantic accounting meaning. For
// example, ordinary runtime code should not reset lifetime counters through a
// raw padded atomic cell. Use the higher-level wrappers when accounting
// semantics matter:
//
//   - Uint64Counter intentionally does not expose Store;
//   - Uint64Gauge exposes Store only for controlled current-state publication.
//
// Store does not validate the value and cannot fail.
func (p *PaddedUint64) Store(val uint64) {
	p.value.Store(val)
}

// Add atomically adds delta to the current uint64 value and returns the new
// value.
//
// Add follows ordinary unsigned uint64 arithmetic. It does not check for
// overflow and may wrap from math.MaxUint64 to zero. That behavior is correct
// for this raw primitive and for monotonic lifetime counters that are sampled
// with wrap-aware delta logic.
//
// Do not use PaddedUint64.Add directly for non-negative gauges that must never
// wrap. Use Uint64Gauge for current quantities where overflow must be treated as
// an accounting invariant violation.
//
// Do not use PaddedUint64 to model values that can legitimately become negative.
// Use PaddedInt64 or Int64Gauge for signed values.
//
// Add does not return an error. Overflow is not an error at this layer; it is a
// semantic concern handled by higher-level wrappers such as Uint64Gauge.
func (p *PaddedUint64) Add(delta uint64) uint64 {
	return p.value.Add(delta)
}

// Inc atomically adds one to the current uint64 value and returns the new value.
//
// Inc is a convenience wrapper around Add(1). It inherits the same raw unsigned
// arithmetic semantics as Add, including possible uint64 wraparound.
func (p *PaddedUint64) Inc() uint64 {
	return p.value.Add(1)
}

// Swap atomically stores newValue and returns the previous value.
//
// Swap is intended for explicit owner-controlled state handoff. Examples include
// test setup, state publication, snapshot-and-reset algorithms, or runtime
// transitions where the caller deliberately owns reset semantics.
//
// Do not use Swap to hide accounting errors in counters or gauges. If the value
// has semantic meaning, prefer the appropriate higher-level wrapper and its
// invariant-preserving operations.
//
// Swap does not validate newValue and cannot fail.
func (p *PaddedUint64) Swap(newValue uint64) uint64 {
	return p.value.Swap(newValue)
}

// CompareAndSwap atomically replaces oldValue with newValue when the current
// value still equals oldValue.
//
// CompareAndSwap is a low-level coordination primitive for callers that own an
// explicit expected-value transition. It should be used when the transition
// itself is conditional on the current atomic value.
//
// The method returns true when the replacement was performed and false when the
// current value no longer matched oldValue.
//
// CompareAndSwap does not validate newValue and cannot fail. Any semantic
// constraints, such as gauge overflow rules or admission-capacity rules, must be
// enforced by higher-level code.
func (p *PaddedUint64) CompareAndSwap(oldValue, newValue uint64) bool {
	return p.value.CompareAndSwap(oldValue, newValue)
}
