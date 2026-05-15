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

// PaddedInt64 is an atomic int64 isolated from neighboring fields by explicit
// leading padding and trailing line-completion padding.
//
// PaddedInt64 is the signed counterpart of PaddedUint64. It is the lowest-level
// signed padded atomic primitive in this package and intentionally exposes raw
// atomic int64 operations. It does not enforce overflow or underflow checks.
//
// Use PaddedInt64 when negative values are part of the valid state model, for
// example:
//
//   - signed correction deltas;
//   - budget movement;
//   - controller drift;
//   - relative capacity adjustments;
//   - signed runtime offsets;
//   - temporary reconciliation error values.
//
// For signed current gauges that must reject int64 overflow/underflow, use
// Int64Gauge instead.
//
// Most lifetime event counters should not use PaddedInt64. They should use
// Uint64Counter because event counters normally move forward and are sampled via
// monotonic unsigned deltas.
//
// Padding reduces false sharing in the same way as PaddedUint64:
//
//	[noCopy marker][leading pad][atomic int64][trailing line-completion pad]
//
// The layout follows the default 64-byte ARCORIS policy, not a hardware
// guarantee. The leading pad separates the atomic value from the previous field.
// The trailing line-completion pad fills the rest of the atomic value's 64-byte
// slot, which matters when padded values appear in arrays, slices, or before
// following struct fields.
//
// The noCopy marker is intentional even though the embedded atomic value already
// follows the same rule. It makes the copy boundary explicit at the atomicx type
// declaration and allows static analysis tools such as go vet -copylocks to
// report accidental copies of the wrapper.
//
// PaddedInt64 is zero-value usable.
//
// PaddedInt64 must not be copied after first use. Copying a live atomic value
// can split one logical state cell into independent copies and produce incorrect
// runtime accounting. Construct it in place, pass it by pointer when sharing,
// and do not copy containing structs after use.
type PaddedInt64 struct {
	noCopy noCopy
	_      CacheLinePad
	value  atomic.Int64
	_      [CacheLinePadSize - atomicInt64Size]byte
}

// Load atomically returns the current int64 value.
//
// Load observes a single atomic value. It does not make a larger multi-field
// state object globally consistent.
func (p *PaddedInt64) Load() int64 {
	return p.value.Load()
}

// Store atomically replaces the current int64 value.
//
// Store is appropriate for initialization, owner-controlled publication, tests,
// or explicit handoff semantics. For signed gauge accounting with invariant
// checks, prefer Int64Gauge.
func (p *PaddedInt64) Store(val int64) {
	p.value.Store(val)
}

// Add atomically adds delta to the current int64 value and returns the new value.
//
// Add follows the semantics of atomic.Int64.Add. It does not guard against
// signed overflow or underflow. Use Int64Gauge when crossing int64 boundaries
// must be treated as an accounting invariant violation.
func (p *PaddedInt64) Add(delta int64) int64 {
	return p.value.Add(delta)
}

// Inc atomically adds one to the current int64 value and returns the new value.
//
// Inc is a convenience wrapper around Add(1).
func (p *PaddedInt64) Inc() int64 {
	return p.value.Add(1)
}

// Dec atomically subtracts one from the current int64 value and returns the new
// value.
//
// Dec is a convenience wrapper around Add(-1). It uses raw signed atomic
// arithmetic and does not check for int64 underflow.
func (p *PaddedInt64) Dec() int64 {
	return p.value.Add(-1)
}

// Swap atomically stores newValue and returns the previous value.
//
// Swap is useful for explicit owner-controlled state handoff, reset-style
// transitions, or tests. It should not be used to bypass higher-level accounting
// invariants when Int64Gauge would be more appropriate.
func (p *PaddedInt64) Swap(newValue int64) int64 {
	return p.value.Swap(newValue)
}

// CompareAndSwap atomically replaces oldValue with newValue when the current
// value still equals oldValue.
//
// CompareAndSwap is intended for explicit expected-value transitions. It returns
// true when the replacement was performed and false when the current value no
// longer matched oldValue.
func (p *PaddedInt64) CompareAndSwap(oldValue, newValue int64) bool {
	return p.value.CompareAndSwap(oldValue, newValue)
}
