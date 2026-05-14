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
	// errUint32GaugeOverflow is the panic value used when a non-negative uint32
	// gauge addition would wrap past the largest representable uint32 value.
	errUint32GaugeOverflow = "atomicx.Uint32Gauge: addition overflows uint32"

	// errUint32GaugeUnderflow is the panic value used when a non-negative uint32
	// gauge subtraction would move the value below zero.
	errUint32GaugeUnderflow = "atomicx.Uint32Gauge: subtraction underflows uint32"
)

const (
	// maxUint32 is kept local to Uint32Gauge because this file owns uint32 gauge
	// boundary checks.
	maxUint32 = ^uint32(0)
)

// Uint32Gauge is a padded non-negative uint32 gauge.
//
// Uint32Gauge represents current runtime state that can move both up and down,
// but must never become negative, and whose 32-bit range is an intentional part
// of the model.
//
// Prefer Uint64Gauge for general counts, bytes, queue depths, in-flight counts,
// and long-lived component accounting. Once explicit cache-line padding is
// present, the memory difference between uint32 and uint64 is usually not the
// deciding factor. Uint32Gauge should be used only when a bounded 32-bit range,
// compact state model, or external protocol boundary is deliberate.
//
// Valid use cases include:
//
//   - small bounded current-state values;
//   - compact state counters with known upper bounds;
//   - protocol-shaped current quantities stored as uint32;
//   - test or simulation gauges where the smaller boundary is intentional.
//
// Uint32Gauge is not a lifetime event counter. For monotonically increasing
// event accounting, use Uint64Counter or, when explicitly bounded, Uint32Counter.
// For signed current values where negative states are meaningful, use Int32Gauge
// or Int64Gauge.
//
// Uint32Gauge uses PaddedUint32 internally to reduce false sharing when the
// gauge is stored near other hot fields in component/runtime state.
//
// Add and Sub enforce non-negative gauge invariants:
//
//   - Add panics on uint32 overflow;
//   - Sub panics on uint32 underflow.
//
// These panics indicate internal accounting bugs. Use TryAdd and TrySub when
// failure is an expected control-flow outcome.
//
// Uint32Gauge is zero-value usable.
//
// Uint32Gauge must not be copied after first use. Copying a live gauge can split
// one logical current-state value into independent copies and corrupt runtime
// accounting. Construct it in place, pass it by pointer when sharing, and do not
// copy containing structs after the gauge becomes active.
type Uint32Gauge struct {
	noCopy noCopy
	value  PaddedUint32
}

// Load atomically returns the current gauge value.
//
// Load observes exactly one atomic value. It does not make a multi-field
// accounting snapshot globally consistent.
func (g *Uint32Gauge) Load() uint32 {
	return g.value.Load()
}

// Store atomically replaces the current gauge value.
//
// Store is appropriate for initialization, tests, owner-controlled publication,
// or explicit state handoff. Ordinary runtime accounting should prefer Add and
// Sub so overflow and underflow are detected at the update point.
func (g *Uint32Gauge) Store(value uint32) {
	g.value.Store(value)
}

// Add atomically adds delta to the gauge and returns the new value.
//
// Add panics if the operation would overflow uint32. Use TryAdd when overflow
// should be handled as a normal rejection path.
func (g *Uint32Gauge) Add(delta uint32) uint32 {
	next, ok := g.TryAdd(delta)
	if !ok {
		panic(errUint32GaugeOverflow)
	}
	return next
}

// TryAdd atomically adds delta to the gauge when the operation would not
// overflow uint32.
//
// On success, TryAdd returns the new value and true.
//
// On failure, TryAdd returns the current value observed at the failing attempt
// and false. The gauge is not modified.
func (g *Uint32Gauge) TryAdd(delta uint32) (uint32, bool) {
	if delta == 0 {
		return g.value.Load(), true
	}

	for {
		current := g.value.Load()
		if current > maxUint32-delta {
			return current, false
		}

		next := current + delta
		if g.value.CompareAndSwap(current, next) {
			return next, true
		}
	}
}

// Sub atomically subtracts delta from the gauge and returns the new value.
//
// Sub panics if the operation would make the gauge negative. Use TrySub when
// underflow should be handled as a normal rejection path.
func (g *Uint32Gauge) Sub(delta uint32) uint32 {
	next, ok := g.TrySub(delta)
	if !ok {
		panic(errUint32GaugeUnderflow)
	}
	return next
}

// TrySub atomically subtracts delta from the gauge when the operation would not
// move the value below zero.
//
// On success, TrySub returns the new value and true.
//
// On failure, TrySub returns the current value observed at the failing attempt
// and false. The gauge is not modified.
func (g *Uint32Gauge) TrySub(delta uint32) (uint32, bool) {
	if delta == 0 {
		return g.value.Load(), true
	}

	for {
		current := g.value.Load()
		if current < delta {
			return current, false
		}

		next := current - delta
		if g.value.CompareAndSwap(current, next) {
			return next, true
		}
	}
}

// Inc atomically adds one to the gauge and returns the new value.
//
// Inc panics if the increment would overflow uint32. Use TryAdd(1) when overflow
// should be handled as a normal control-flow result.
func (g *Uint32Gauge) Inc() uint32 {
	return g.Add(1)
}

// Dec atomically subtracts one from the gauge and returns the new value.
//
// Dec panics if the decrement would move the gauge below zero. Use TrySub(1)
// when insufficient current value should be handled as a normal control-flow
// result.
func (g *Uint32Gauge) Dec() uint32 {
	return g.Sub(1)
}

// Swap atomically stores newValue and returns the previous value.
//
// Swap is useful for explicit owner-controlled handoff, reset-style transitions,
// test setup, or state publication. It should not be used to hide accounting
// bugs that should be expressed through Add or Sub.
func (g *Uint32Gauge) Swap(newValue uint32) uint32 {
	return g.value.Swap(newValue)
}

// CompareAndSwap atomically replaces oldValue with newValue when the current
// value still equals oldValue.
//
// CompareAndSwap is exposed for advanced internal state transitions where the
// caller owns the expected-value protocol. Callers that need
// invariant-preserving arithmetic should use Add, TryAdd, Sub, or TrySub.
func (g *Uint32Gauge) CompareAndSwap(oldValue, newValue uint32) bool {
	return g.value.CompareAndSwap(oldValue, newValue)
}
