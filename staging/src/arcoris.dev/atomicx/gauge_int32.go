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

const (
	// errInt32GaugeOverflow is the panic value used when a signed int32 gauge
	// operation would exceed the largest representable int32 value.
	errInt32GaugeOverflow = "atomicx.Int32Gauge: operation overflows int32"

	// errInt32GaugeUnderflow is the panic value used when a signed int32 gauge
	// operation would move below the smallest representable int32 value.
	errInt32GaugeUnderflow = "atomicx.Int32Gauge: operation underflows int32"
)

const (
	// maxInt32 is kept local to Int32Gauge because this file owns int32 gauge
	// boundary checks.
	maxInt32 = int32(1<<31 - 1)

	// minInt32 is the smallest signed 32-bit value.
	minInt32 = -maxInt32 - 1
)

// Int32Gauge is a padded signed int32 gauge.
//
// Int32Gauge represents current runtime state where negative values are
// meaningful and the 32-bit range is an intentional part of the model.
//
// Prefer Int64Gauge for general signed accounting, correction deltas, budget
// movement, drift, and long-lived component state. Once explicit cache-line
// padding is present, the memory difference between int32 and int64 is usually
// not the deciding factor. Int32Gauge should be used only when a bounded 32-bit
// range, compact state model, or external protocol boundary is deliberate.
//
// Valid use cases include:
//
//   - small signed correction values with known bounds;
//   - compact signed state transitions;
//   - protocol-shaped signed current quantities stored as int32;
//   - test or simulation gauges where the smaller boundary is intentional.
//
// Int32Gauge is not a lifetime event counter. For values that only move forward,
// use Uint64Counter or Uint32Counter. For non-negative current quantities, use
// Uint32Gauge or Uint64Gauge.
//
// Int32Gauge uses PaddedInt32 internally to reduce false sharing when the gauge
// is stored near other hot fields in component/runtime state.
//
// Add and Sub enforce int32 boundaries:
//
//   - Add panics when the operation would overflow or underflow int32;
//   - Sub panics when the operation would overflow or underflow int32.
//
// These panics indicate internal accounting bugs. Use TryAdd and TrySub when
// failure is an expected control-flow outcome.
//
// Int32Gauge is zero-value usable.
//
// Int32Gauge must not be copied after first use. Copying a live gauge can split
// one logical current-state value into independent copies and corrupt runtime
// accounting. Construct it in place, pass it by pointer when sharing, and do not
// copy containing structs after the gauge becomes active.
type Int32Gauge struct {
	noCopy noCopy
	value  PaddedInt32
}

// Load atomically returns the current signed gauge value.
//
// Load observes exactly one atomic value. It does not make a multi-field
// accounting snapshot globally consistent.
func (g *Int32Gauge) Load() int32 {
	return g.value.Load()
}

// Store atomically replaces the current signed gauge value.
//
// Store is appropriate for initialization, tests, owner-controlled publication,
// or explicit state handoff. Ordinary runtime accounting should prefer Add and
// Sub so overflow and underflow are detected at the update point.
func (g *Int32Gauge) Store(val int32) {
	g.value.Store(val)
}

// Add atomically adds delta to the signed gauge and returns the new value.
//
// Add panics if the operation would overflow or underflow int32. Silent signed
// wraparound would corrupt current-state accounting.
//
// Add should be used when crossing int32 bounds is an internal invariant
// violation. If failure should instead be handled as a normal rejection path, use
// TryAdd.
func (g *Int32Gauge) Add(delta int32) int32 {
	next, ok := g.TryAdd(delta)
	if !ok {
		if delta < 0 {
			panic(errInt32GaugeUnderflow)
		}
		panic(errInt32GaugeOverflow)
	}
	return next
}

// TryAdd atomically adds delta to the signed gauge when the operation would not
// overflow or underflow int32.
//
// On success, TryAdd returns the new value and true.
//
// On failure, TryAdd returns the current value observed at the failing attempt
// and false. The gauge is not modified.
func (g *Int32Gauge) TryAdd(delta int32) (int32, bool) {
	if delta == 0 {
		return g.value.Load(), true
	}

	for {
		cur := g.value.Load()

		if delta > 0 && cur > maxInt32-delta {
			return cur, false
		}
		if delta < 0 && cur < minInt32-delta {
			return cur, false
		}

		next := cur + delta
		if g.value.CompareAndSwap(cur, next) {
			return next, true
		}
	}
}

// Sub atomically subtracts delta from the signed gauge and returns the new value.
//
// Sub panics if the operation would overflow or underflow int32. For example,
// subtracting a negative delta can overflow upward, and subtracting a positive
// delta can underflow downward.
//
// Sub should be used when crossing int32 bounds is an internal invariant
// violation. If failure should instead be handled as a normal rejection path, use
// TrySub.
func (g *Int32Gauge) Sub(delta int32) int32 {
	next, ok := g.TrySub(delta)
	if !ok {
		if delta < 0 {
			panic(errInt32GaugeOverflow)
		}
		panic(errInt32GaugeUnderflow)
	}
	return next
}

// TrySub atomically subtracts delta from the signed gauge when the operation
// would not overflow or underflow int32.
//
// On success, TrySub returns the new value and true.
//
// On failure, TrySub returns the current value observed at the failing attempt
// and false. The gauge is not modified.
//
// TrySub is not implemented as TryAdd(-delta), because -minInt32 is not
// representable as int32.
func (g *Int32Gauge) TrySub(delta int32) (int32, bool) {
	if delta == 0 {
		return g.value.Load(), true
	}

	for {
		cur := g.value.Load()

		if delta > 0 && cur < minInt32+delta {
			return cur, false
		}
		if delta < 0 && cur > maxInt32+delta {
			return cur, false
		}

		next := cur - delta
		if g.value.CompareAndSwap(cur, next) {
			return next, true
		}
	}
}

// Inc atomically adds one to the signed gauge and returns the new value.
//
// Inc panics if the increment would overflow int32. Use TryAdd(1) when overflow
// should be handled as a normal control-flow result.
func (g *Int32Gauge) Inc() int32 {
	return g.Add(1)
}

// Dec atomically subtracts one from the signed gauge and returns the new value.
//
// Dec panics if the decrement would underflow int32. Use TrySub(1) when underflow
// should be handled as a normal control-flow result.
func (g *Int32Gauge) Dec() int32 {
	return g.Sub(1)
}

// Swap atomically stores newValue and returns the previous value.
//
// Swap is useful for explicit owner-controlled handoff, reset-style transitions,
// test setup, or state publication. It should not be used to hide accounting
// bugs that should be expressed through Add or Sub.
func (g *Int32Gauge) Swap(newValue int32) int32 {
	return g.value.Swap(newValue)
}

// CompareAndSwap atomically replaces oldValue with newValue when the current
// value still equals oldValue.
//
// CompareAndSwap is exposed for advanced internal state transitions where the
// caller owns the expected-value protocol. Callers that need
// invariant-preserving arithmetic should use Add, TryAdd, Sub, or TrySub.
func (g *Int32Gauge) CompareAndSwap(oldValue, newValue int32) bool {
	return g.value.CompareAndSwap(oldValue, newValue)
}
