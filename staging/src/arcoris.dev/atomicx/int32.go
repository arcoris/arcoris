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

// PaddedInt32 is an atomic int32 isolated from neighboring fields by explicit
// leading and trailing padding.
//
// PaddedInt32 is a low-level primitive for hot signed 32-bit state. It is
// intended for compact runtime state where negative values are meaningful and
// the state is naturally bounded to int32, for example:
//
//   - compact signed state-machine markers;
//   - signed small-range correction values;
//   - signed bit-pattern-compatible state;
//   - compact owner-controlled runtime state;
//   - low-cardinality signed state transitions.
//
// PaddedInt32 intentionally exposes raw signed atomic operations. It does not
// impose higher-level accounting semantics. In particular:
//
//   - Add follows atomic.Int32.Add semantics;
//   - Add does not check for signed overflow or underflow;
//   - Store and Swap are allowed;
//   - negative values are allowed;
//   - no gauge invariant is enforced here.
//
// Do not use PaddedInt32 merely to save memory for counters or gauges. Once
// explicit cache-line padding is present, the memory difference between int32
// and int64 is usually irrelevant. For signed current quantities, controller
// corrections, budget movement, and other accounting values that can reasonably
// exceed small ranges, prefer PaddedInt64 or Int64Gauge.
//
// Use PaddedInt32 only when the value is semantically a signed 32-bit state
// word, not when it is just a smaller counter.
//
// For signed current gauges that must reject overflow or underflow, use a
// higher-level gauge type. For monotonic lifetime event counters, use
// Uint64Counter. For non-negative current quantities, use Uint64Gauge.
//
// Padding reduces false sharing. False sharing can occur when independent hot
// variables occupy the same CPU cache line and different goroutines update them
// from different cores. Even though the logical variables are unrelated, writes
// to one value can invalidate the cache line containing another value.
//
// Layout:
//
//	[noCopy marker][leading pad][atomic int32][trailing pad]
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
// PaddedInt32 is zero-value usable.
//
// PaddedInt32 must not be copied after first use. Copying a live atomic value
// can split one logical state cell into independent copies and produce incorrect
// runtime state. Construct it in place, pass it by pointer when sharing, and do
// not copy containing structs after the value becomes active.
type PaddedInt32 struct {
	noCopy noCopy
	_      CacheLinePad
	value  atomic.Int32
	_      CacheLinePad
}

// Load atomically returns the current int32 value.
//
// Load provides the same synchronization semantics as atomic.Int32.Load. It
// observes exactly one atomic value. It does not make a larger multi-field state
// object globally consistent.
//
// If a caller needs a consistent snapshot across multiple fields, the caller
// must provide additional synchronization at the owner level.
func (p *PaddedInt32) Load() int32 {
	return p.value.Load()
}

// Store atomically replaces the current int32 value.
//
// Store is a raw state publication operation. It is appropriate for
// initialization, tests, owner-controlled publication, or explicit state
// handoff.
//
// Store does not validate the value. If the int32 represents a constrained
// state machine or bounded signed state, the caller is responsible for
// validating allowed states and transitions.
func (p *PaddedInt32) Store(value int32) {
	p.value.Store(value)
}

// Add atomically adds delta to the current int32 value and returns the new value.
//
// Add follows atomic.Int32.Add semantics. It does not check for signed overflow
// or underflow. Crossing int32 boundaries is raw atomic behavior at this layer.
//
// This method is intentionally raw. Do not use PaddedInt32.Add for accounting
// values that require non-wrapping semantics. Use a dedicated higher-level type
// when overflow or underflow must be treated as an invariant violation.
func (p *PaddedInt32) Add(delta int32) int32 {
	return p.value.Add(delta)
}

// Inc atomically adds one to the current int32 value and returns the new value.
//
// Inc is a convenience wrapper around Add(1). It inherits the same raw signed
// arithmetic semantics as Add.
func (p *PaddedInt32) Inc() int32 {
	return p.value.Add(1)
}

// Dec atomically subtracts one from the current int32 value and returns the new
// value.
//
// Dec is a convenience wrapper around Add(-1). It inherits the same raw signed
// arithmetic semantics as Add and does not check for int32 underflow.
func (p *PaddedInt32) Dec() int32 {
	return p.value.Add(-1)
}

// Swap atomically stores newValue and returns the previous value.
//
// Swap is intended for explicit owner-controlled state handoff. It is useful for
// tests, state publication, state-machine replacement, reset-style transitions,
// or runtime transitions where the caller deliberately owns the semantics.
//
// Swap does not validate newValue. If the value represents a constrained state
// machine or bounded signed state, the caller must validate the transition
// before calling Swap.
func (p *PaddedInt32) Swap(newValue int32) int32 {
	return p.value.Swap(newValue)
}

// CompareAndSwap atomically replaces oldValue with newValue when the current
// value still equals oldValue.
//
// CompareAndSwap is the preferred operation for lock-free state-machine
// transitions backed by PaddedInt32. The caller owns the state model and must
// validate that newValue is a legal state before attempting the transition.
//
// The method returns true when the replacement was performed and false when the
// current value no longer matched oldValue.
//
// CompareAndSwap does not validate newValue and does not enforce signed bounds.
// Any semantic constraints must be enforced by higher-level code.
func (p *PaddedInt32) CompareAndSwap(oldValue, newValue int32) bool {
	return p.value.CompareAndSwap(oldValue, newValue)
}
