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
	// errInt64GaugeOverflow is the panic value used when a signed int64 gauge
	// operation would exceed the largest representable int64 value.
	errInt64GaugeOverflow = "atomicx.Int64Gauge: operation overflows int64"

	// errInt64GaugeUnderflow is the panic value used when a signed int64 gauge
	// operation would move below the smallest representable int64 value.
	errInt64GaugeUnderflow = "atomicx.Int64Gauge: operation underflows int64"
)

const (
	// maxInt64 is kept local to Int64Gauge because this file owns int64 gauge
	// boundary checks.
	maxInt64 = int64(1<<63 - 1)

	// minInt64 is the smallest signed 64-bit value.
	minInt64 = -maxInt64 - 1
)

// Int64Gauge is a padded signed int64 gauge.
//
// Int64Gauge represents current runtime state where negative values are
// meaningful. Typical examples include:
//
//   - signed correction deltas;
//   - signed budget movement;
//   - controller drift;
//   - relative capacity adjustments;
//   - signed reconciliation offsets;
//   - temporary signed runtime adjustments;
//   - capacity debt or surplus.
//
// Int64Gauge is not a lifetime event counter. For values that only move forward,
// use Uint64Counter. For non-negative current quantities, use Uint64Gauge. For
// deliberately small signed state, use Int32Gauge.
//
// Int64Gauge uses PaddedInt64 internally to reduce false sharing when the gauge
// is stored near other hot fields in component/runtime state.
//
// Add and Sub enforce int64 boundaries:
//
//   - Add panics when the operation would overflow or underflow int64;
//   - Sub panics when the operation would overflow or underflow int64.
//
// These panics indicate internal accounting bugs. Use TryAdd and TrySub when
// failure is an expected control-flow outcome.
//
// Int64Gauge is zero-value usable.
//
// Int64Gauge must not be copied after first use. Copying a live gauge can split
// one logical current-state value into independent copies and corrupt runtime
// accounting. Construct it in place, pass it by pointer when sharing, and do not
// copy containing structs after the gauge becomes active.
type Int64Gauge struct {
	noCopy noCopy
	value  PaddedInt64
}

// Load atomically returns the current signed gauge value.
//
// Load observes exactly one atomic value. It does not make a multi-field
// accounting snapshot globally consistent.
func (g *Int64Gauge) Load() int64 {
	return g.value.Load()
}

// Store atomically replaces the current signed gauge value.
//
// Store is appropriate for initialization, tests, owner-controlled publication,
// or explicit state handoff. Ordinary runtime accounting should prefer Add and
// Sub so overflow and underflow are detected at the update point.
func (g *Int64Gauge) Store(value int64) {
	g.value.Store(value)
}

// Add atomically adds delta to the signed gauge and returns the new value.
//
// Add panics if the operation would overflow or underflow int64. Silent signed
// wraparound would corrupt correction, budget, drift, or adjustment accounting.
//
// Add should be used when crossing int64 bounds is an internal invariant
// violation. If failure should instead be handled as a normal rejection path, use
// TryAdd.
func (g *Int64Gauge) Add(delta int64) int64 {
	next, ok := g.TryAdd(delta)
	if !ok {
		if delta < 0 {
			panic(errInt64GaugeUnderflow)
		}
		panic(errInt64GaugeOverflow)
	}
	return next
}

// TryAdd atomically adds delta to the signed gauge when the operation would not
// overflow or underflow int64.
//
// On success, TryAdd returns the new value and true.
//
// On failure, TryAdd returns the current value observed at the failing attempt
// and false. The gauge is not modified.
func (g *Int64Gauge) TryAdd(delta int64) (int64, bool) {
	if delta == 0 {
		return g.value.Load(), true
	}

	for {
		current := g.value.Load()

		if delta > 0 && current > maxInt64-delta {
			return current, false
		}
		if delta < 0 && current < minInt64-delta {
			return current, false
		}

		next := current + delta
		if g.value.CompareAndSwap(current, next) {
			return next, true
		}
	}
}

// Sub atomically subtracts delta from the signed gauge and returns the new value.
//
// Sub panics if the operation would overflow or underflow int64. For example,
// subtracting a negative delta can overflow upward, and subtracting a positive
// delta can underflow downward.
//
// Sub should be used when crossing int64 bounds is an internal invariant
// violation. If failure should instead be handled as a normal rejection path, use
// TrySub.
func (g *Int64Gauge) Sub(delta int64) int64 {
	next, ok := g.TrySub(delta)
	if !ok {
		if delta < 0 {
			panic(errInt64GaugeOverflow)
		}
		panic(errInt64GaugeUnderflow)
	}
	return next
}

// TrySub atomically subtracts delta from the signed gauge when the operation
// would not overflow or underflow int64.
//
// On success, TrySub returns the new value and true.
//
// On failure, TrySub returns the current value observed at the failing attempt
// and false. The gauge is not modified.
//
// TrySub is not implemented as TryAdd(-delta), because -minInt64 is not
// representable as int64.
func (g *Int64Gauge) TrySub(delta int64) (int64, bool) {
	if delta == 0 {
		return g.value.Load(), true
	}

	for {
		current := g.value.Load()

		if delta > 0 && current < minInt64+delta {
			return current, false
		}
		if delta < 0 && current > maxInt64+delta {
			return current, false
		}

		next := current - delta
		if g.value.CompareAndSwap(current, next) {
			return next, true
		}
	}
}

// Inc atomically adds one to the signed gauge and returns the new value.
//
// Inc panics if the increment would overflow int64. Use TryAdd(1) when overflow
// should be handled as a normal control-flow result.
func (g *Int64Gauge) Inc() int64 {
	return g.Add(1)
}

// Dec atomically subtracts one from the signed gauge and returns the new value.
//
// Dec panics if the decrement would underflow int64. Use TrySub(1) when underflow
// should be handled as a normal control-flow result.
func (g *Int64Gauge) Dec() int64 {
	return g.Sub(1)
}

// Swap atomically stores newValue and returns the previous value.
//
// Swap is useful for explicit owner-controlled handoff, reset-style transitions,
// test setup, or state publication. It should not be used to hide accounting
// bugs that should be expressed through Add or Sub.
func (g *Int64Gauge) Swap(newValue int64) int64 {
	return g.value.Swap(newValue)
}

// CompareAndSwap atomically replaces oldValue with newValue when the current
// value still equals oldValue.
//
// CompareAndSwap is exposed for advanced internal state transitions where the
// caller owns the expected-value protocol. Callers that need
// invariant-preserving arithmetic should use Add, TryAdd, Sub, or TrySub.
func (g *Int64Gauge) CompareAndSwap(oldValue, newValue int64) bool {
	return g.value.CompareAndSwap(oldValue, newValue)
}
