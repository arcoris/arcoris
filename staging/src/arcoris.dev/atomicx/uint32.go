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

// PaddedUint32 is an atomic uint32 isolated from neighboring fields by explicit
// leading and trailing padding.
//
// PaddedUint32 is a low-level primitive for hot uint32 state. It is intended for
// compact runtime state that is naturally represented as a 32-bit unsigned
// value, for example:
//
//   - small state-machine markers;
//   - bitmask state;
//   - compact enum-like values;
//   - hot flags stored as a single numeric state word;
//   - owner-controlled component state transitions.
//
// PaddedUint32 intentionally exposes raw unsigned atomic operations. It does not
// impose counter or gauge semantics. In particular:
//
//   - Add follows ordinary uint32 arithmetic and may wrap;
//   - Store and Swap are allowed;
//   - no overflow invariant is enforced;
//   - no non-negative gauge invariant is enforced;
//   - no lifecycle-state validation is enforced.
//
// Do not use PaddedUint32 merely to save memory for counters or byte accounting.
// Once explicit cache-line padding is present, the memory difference between
// uint32 and uint64 is usually irrelevant. For event counts, byte totals, queue
// depths, in-flight counts, and other accounting values, prefer the semantic
// uint64 wrappers:
//
//   - Uint64Counter for monotonic lifetime event counters;
//   - Uint64Gauge for non-negative current quantities.
//
// Use PaddedUint32 only when the value is semantically a uint32 state word, not
// when it is just a smaller counter.
//
// Padding reduces false sharing. False sharing can occur when independent hot
// variables occupy the same CPU cache line and different goroutines update them
// from different cores. Even though the logical variables are unrelated, writes
// to one value can invalidate the cache line containing another value.
//
// Layout:
//
//	[noCopy marker][leading pad][atomic uint32][trailing pad]
//
// Both pads are intentional:
//
//   - the leading pad separates the atomic value from the previous field;
//   - the trailing pad separates it from the next field;
//   - the trailing pad also matters when padded values appear in arrays/slices.
//
// The noCopy marker is a static-analysis marker. It does not participate in
// synchronization and does not provide runtime protection. Its purpose is to make
// accidental value copies visible to tools such as:
//
//	go vet -copylocks
//
// PaddedUint32 is zero-value usable.
//
// PaddedUint32 must not be copied after first use. Copying a live atomic value
// can split one logical state cell into independent copies and produce incorrect
// runtime state. Construct it in place, pass it by pointer when sharing, and do
// not copy containing structs after the value becomes active.
type PaddedUint32 struct {
	noCopy noCopy
	_      CacheLinePad
	value  atomic.Uint32
	_      CacheLinePad
}

// Load atomically returns the current uint32 value.
//
// Load provides the same synchronization semantics as atomic.Uint32.Load. It
// observes exactly one atomic value. It does not make a larger multi-field state
// object globally consistent.
//
// If a caller needs a consistent snapshot across multiple fields, the caller
// must provide additional synchronization at the owner level.
func (p *PaddedUint32) Load() uint32 {
	return p.value.Load()
}

// Store atomically replaces the current uint32 value.
//
// Store is a raw state publication operation. It is appropriate for
// initialization, tests, owner-controlled publication, or explicit state
// handoff.
//
// Store does not validate the value. If the uint32 represents an enum-like state
// machine, the caller is responsible for validating allowed states and
// transitions.
func (p *PaddedUint32) Store(value uint32) {
	p.value.Store(value)
}

// Add atomically adds delta to the current uint32 value and returns the new
// value.
//
// Add follows ordinary unsigned uint32 arithmetic. It does not check for
// overflow and may wrap from math.MaxUint32 to zero.
//
// This method is intentionally raw. Do not use PaddedUint32.Add for accounting
// values that require non-wrapping semantics. Use Uint64Gauge or a dedicated
// higher-level type when overflow must be treated as an invariant violation.
func (p *PaddedUint32) Add(delta uint32) uint32 {
	return p.value.Add(delta)
}

// Inc atomically adds one to the current uint32 value and returns the new value.
//
// Inc is a convenience wrapper around Add(1). It inherits the same raw unsigned
// arithmetic semantics as Add, including possible uint32 wraparound.
func (p *PaddedUint32) Inc() uint32 {
	return p.value.Add(1)
}

// Swap atomically stores newValue and returns the previous value.
//
// Swap is intended for explicit owner-controlled state handoff. It is useful for
// tests, state publication, state-machine replacement, or reset-style
// transitions where the caller owns the semantics.
//
// Swap does not validate newValue. If the value represents a constrained state
// machine, the caller must validate the transition before calling Swap.
func (p *PaddedUint32) Swap(newValue uint32) uint32 {
	return p.value.Swap(newValue)
}

// CompareAndSwap atomically replaces oldValue with newValue when the current
// value still equals oldValue.
//
// CompareAndSwap is the preferred operation for lock-free state-machine
// transitions backed by PaddedUint32. The caller owns the state model and must
// validate that newValue is a legal state before attempting the transition.
//
// The method returns true when the replacement was performed and false when the
// current value no longer matched oldValue.
func (p *PaddedUint32) CompareAndSwap(oldValue, newValue uint32) bool {
	return p.value.CompareAndSwap(oldValue, newValue)
}
