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

import "testing"

// TestUint64GaugeZeroValueIsUsable verifies a gauge can be used without explicit initialization.
func TestUint64GaugeZeroValueIsUsable(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge

	if got := gauge.Load(); got != 0 {
		t.Fatalf("zero-value Uint64Gauge.Load() = %d, want 0", got)
	}
}

// TestUint64GaugeStoreAndLoad verifies owner-controlled publication through Store.
func TestUint64GaugeStoreAndLoad(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(42)

	if got := gauge.Load(); got != 42 {
		t.Fatalf("Uint64Gauge.Load() after Store(42) = %d, want 42", got)
	}
}

// TestUint64GaugeAddSubAndZeroDeltas verifies invariant-preserving arithmetic.
func TestUint64GaugeAddSubAndZeroDeltas(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge

	if got := gauge.Add(10); got != 10 {
		t.Fatalf("Uint64Gauge.Add(10) = %d, want 10", got)
	}
	if got := gauge.Add(0); got != 10 {
		t.Fatalf("Uint64Gauge.Add(0) = %d, want 10", got)
	}
	if got := gauge.Add(5); got != 15 {
		t.Fatalf("Uint64Gauge.Add(5) = %d, want 15", got)
	}
	if got := gauge.Sub(4); got != 11 {
		t.Fatalf("Uint64Gauge.Sub(4) = %d, want 11", got)
	}
	if got := gauge.Sub(0); got != 11 {
		t.Fatalf("Uint64Gauge.Sub(0) = %d, want 11", got)
	}
	if got := gauge.Load(); got != 11 {
		t.Fatalf("Uint64Gauge.Load() after Add/Sub sequence = %d, want 11", got)
	}
}

// TestUint64GaugeTryAddSuccess verifies successful checked addition updates state.
func TestUint64GaugeTryAddSuccess(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(10)

	got, ok := gauge.TryAdd(5)
	if !ok {
		t.Fatal("Uint64Gauge.TryAdd(5) ok = false, want true")
	}
	if got != 15 {
		t.Fatalf("Uint64Gauge.TryAdd(5) value = %d, want 15", got)
	}
	if loaded := gauge.Load(); loaded != 15 {
		t.Fatalf("Uint64Gauge.Load() after TryAdd(5) = %d, want 15", loaded)
	}
}

// TestUint64GaugeTryAddOverflowLeavesStateUnchanged verifies checked admission
// failure does not corrupt the current state. A failed TryAdd must be safe to
// use as a capacity rejection path.
func TestUint64GaugeTryAddOverflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(maxUint64)

	got, ok := gauge.TryAdd(1)
	if ok {
		t.Fatal("Uint64Gauge.TryAdd(1) at max ok = true, want false")
	}
	if got != maxUint64 {
		t.Fatalf("Uint64Gauge.TryAdd(1) failure value = %d, want %d", got, maxUint64)
	}
	if loaded := gauge.Load(); loaded != maxUint64 {
		t.Fatalf("Uint64Gauge.Load() after failed TryAdd = %d, want %d", loaded, maxUint64)
	}
}

// TestUint64GaugeTrySubSuccess verifies successful checked subtraction updates state.
func TestUint64GaugeTrySubSuccess(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(10)

	got, ok := gauge.TrySub(4)
	if !ok {
		t.Fatal("Uint64Gauge.TrySub(4) ok = false, want true")
	}
	if got != 6 {
		t.Fatalf("Uint64Gauge.TrySub(4) value = %d, want 6", got)
	}
	if loaded := gauge.Load(); loaded != 6 {
		t.Fatalf("Uint64Gauge.Load() after TrySub(4) = %d, want 6", loaded)
	}
}

// TestUint64GaugeTrySubUnderflowLeavesStateUnchanged verifies checked release
// failure does not hide an accounting imbalance by changing the gauge.
func TestUint64GaugeTrySubUnderflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(3)

	got, ok := gauge.TrySub(4)
	if ok {
		t.Fatal("Uint64Gauge.TrySub(4) from 3 ok = true, want false")
	}
	if got != 3 {
		t.Fatalf("Uint64Gauge.TrySub(4) failure value = %d, want 3", got)
	}
	if loaded := gauge.Load(); loaded != 3 {
		t.Fatalf("Uint64Gauge.Load() after failed TrySub = %d, want 3", loaded)
	}
}

// TestUint64GaugeIncAndDec verifies single-unit gauge arithmetic.
func TestUint64GaugeIncAndDec(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge

	if got := gauge.Inc(); got != 1 {
		t.Fatalf("Uint64Gauge.Inc() = %d, want 1", got)
	}
	if got := gauge.Inc(); got != 2 {
		t.Fatalf("second Uint64Gauge.Inc() = %d, want 2", got)
	}
	if got := gauge.Dec(); got != 1 {
		t.Fatalf("Uint64Gauge.Dec() = %d, want 1", got)
	}
	if got := gauge.Load(); got != 1 {
		t.Fatalf("Uint64Gauge.Load() after Inc/Dec sequence = %d, want 1", got)
	}
}

// TestUint64GaugeSwap verifies explicit owner-controlled replacement semantics.
func TestUint64GaugeSwap(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(10)

	if old := gauge.Swap(25); old != 10 {
		t.Fatalf("Uint64Gauge.Swap(25) old value = %d, want 10", old)
	}
	if got := gauge.Load(); got != 25 {
		t.Fatalf("Uint64Gauge.Load() after Swap(25) = %d, want 25", got)
	}
}

// TestUint64GaugeCompareAndSwap verifies conditional owner-controlled transitions.
func TestUint64GaugeCompareAndSwap(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(10)

	if swapped := gauge.CompareAndSwap(9, 20); swapped {
		t.Fatal("Uint64Gauge.CompareAndSwap(9, 20) = true, want false")
	}
	if got := gauge.Load(); got != 10 {
		t.Fatalf("Uint64Gauge.Load() after failed CAS = %d, want 10", got)
	}
	if swapped := gauge.CompareAndSwap(10, 20); !swapped {
		t.Fatal("Uint64Gauge.CompareAndSwap(10, 20) = false, want true")
	}
	if got := gauge.Load(); got != 20 {
		t.Fatalf("Uint64Gauge.Load() after successful CAS = %d, want 20", got)
	}
}

// TestUint64GaugeExactBoundaryOperations verifies that a gauge may legally
// reach the numeric boundary. The invariant violation starts only when an
// operation attempts to cross the boundary.
func TestUint64GaugeExactBoundaryOperations(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge

	if got := gauge.Add(maxUint64); got != maxUint64 {
		t.Fatalf("Uint64Gauge.Add(maxUint64) = %d, want %d", got, maxUint64)
	}
	if got := gauge.Sub(maxUint64); got != 0 {
		t.Fatalf("Uint64Gauge.Sub(maxUint64) = %d, want 0", got)
	}
}

// TestUint64GaugePanicsOnOverflow verifies that current-state gauges do not
// silently wrap to a smaller value. Silent wraparound would corrupt admission,
// queue-depth, retained-byte, and in-flight accounting.
func TestUint64GaugePanicsOnOverflow(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(maxUint64)

	mustPanicWithValue(t, errUint64GaugeOverflow, func() {
		_ = gauge.Add(1)
	})
}

// TestUint64GaugePanicsOnUnderflow verifies gauges reject negative current
// state instead of wrapping subtraction to a very large unsigned value.
func TestUint64GaugePanicsOnUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(10)

	mustPanicWithValue(t, errUint64GaugeUnderflow, func() {
		_ = gauge.Sub(11)
	})
}

// TestUint64GaugeIncPanicsOnOverflow verifies Inc preserves the same overflow
// invariant as Add, so the convenience method cannot bypass gauge safety.
func TestUint64GaugeIncPanicsOnOverflow(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge
	gauge.Store(maxUint64)

	mustPanicWithValue(t, errUint64GaugeOverflow, func() {
		_ = gauge.Inc()
	})
}

// TestUint64GaugeDecPanicsOnUnderflow verifies Dec preserves the same underflow
// invariant as Sub, so single-unit release cannot bypass gauge safety.
func TestUint64GaugeDecPanicsOnUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Uint64Gauge

	mustPanicWithValue(t, errUint64GaugeUnderflow, func() {
		_ = gauge.Dec()
	})
}

// TestUint64GaugeConcurrentBalancedAccounting verifies deterministic current
// state accounting under contention without sleeps or timing assumptions.
func TestUint64GaugeConcurrentBalancedAccounting(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const iterations = 10_000

	var gauge Uint64Gauge
	runConcurrent(t, goroutines, func() {
		for range iterations {
			gauge.Add(2)
			gauge.Sub(1)
		}
	})

	want := uint64(goroutines * iterations)
	if got := gauge.Load(); got != want {
		t.Fatalf("Uint64Gauge.Load() after concurrent balanced accounting = %d, want %d", got, want)
	}
}
