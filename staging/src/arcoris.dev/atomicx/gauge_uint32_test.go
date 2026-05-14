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

// TestUint32GaugeZeroValueIsUsable verifies a bounded gauge works without initialization.
func TestUint32GaugeZeroValueIsUsable(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge

	if got := gauge.Load(); got != 0 {
		t.Fatalf("zero-value Uint32Gauge.Load() = %d, want 0", got)
	}
}

// TestUint32GaugeStoreAndLoad verifies owner-controlled publication through Store.
func TestUint32GaugeStoreAndLoad(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(42)

	if got := gauge.Load(); got != 42 {
		t.Fatalf("Uint32Gauge.Load() after Store(42) = %d, want 42", got)
	}
}

// TestUint32GaugeAddSubAndZeroDeltas verifies checked arithmetic and no-op deltas.
func TestUint32GaugeAddSubAndZeroDeltas(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge

	if got := gauge.Add(10); got != 10 {
		t.Fatalf("Uint32Gauge.Add(10) = %d, want 10", got)
	}
	if got := gauge.Add(0); got != 10 {
		t.Fatalf("Uint32Gauge.Add(0) = %d, want 10", got)
	}
	if got := gauge.Add(5); got != 15 {
		t.Fatalf("Uint32Gauge.Add(5) = %d, want 15", got)
	}
	if got := gauge.Sub(4); got != 11 {
		t.Fatalf("Uint32Gauge.Sub(4) = %d, want 11", got)
	}
	if got := gauge.Sub(0); got != 11 {
		t.Fatalf("Uint32Gauge.Sub(0) = %d, want 11", got)
	}
	if got := gauge.Load(); got != 11 {
		t.Fatalf("Uint32Gauge.Load() after Add/Sub sequence = %d, want 11", got)
	}
}

// TestUint32GaugeTryAddSuccess verifies successful checked addition updates state.
func TestUint32GaugeTryAddSuccess(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(10)

	got, ok := gauge.TryAdd(5)
	if !ok {
		t.Fatal("Uint32Gauge.TryAdd(5) ok = false, want true")
	}
	if got != 15 {
		t.Fatalf("Uint32Gauge.TryAdd(5) value = %d, want 15", got)
	}
	if loaded := gauge.Load(); loaded != 15 {
		t.Fatalf("Uint32Gauge.Load() after TryAdd(5) = %d, want 15", loaded)
	}
}

// TestUint32GaugeTryAddOverflowLeavesStateUnchanged verifies checked admission
// failure does not corrupt the current state. A failed TryAdd must be safe to
// use as a bounded-capacity rejection path.
func TestUint32GaugeTryAddOverflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(maxUint32)

	got, ok := gauge.TryAdd(1)
	if ok {
		t.Fatal("Uint32Gauge.TryAdd(1) at max ok = true, want false")
	}
	if got != maxUint32 {
		t.Fatalf("Uint32Gauge.TryAdd(1) failure value = %d, want %d", got, maxUint32)
	}
	if loaded := gauge.Load(); loaded != maxUint32 {
		t.Fatalf("Uint32Gauge.Load() after failed TryAdd = %d, want %d", loaded, maxUint32)
	}
}

// TestUint32GaugeTrySubSuccess verifies successful checked subtraction updates state.
func TestUint32GaugeTrySubSuccess(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(10)

	got, ok := gauge.TrySub(4)
	if !ok {
		t.Fatal("Uint32Gauge.TrySub(4) ok = false, want true")
	}
	if got != 6 {
		t.Fatalf("Uint32Gauge.TrySub(4) value = %d, want 6", got)
	}
	if loaded := gauge.Load(); loaded != 6 {
		t.Fatalf("Uint32Gauge.Load() after TrySub(4) = %d, want 6", loaded)
	}
}

// TestUint32GaugeTrySubUnderflowLeavesStateUnchanged verifies checked release
// failure does not hide an accounting imbalance by changing the gauge.
func TestUint32GaugeTrySubUnderflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(3)

	got, ok := gauge.TrySub(4)
	if ok {
		t.Fatal("Uint32Gauge.TrySub(4) from 3 ok = true, want false")
	}
	if got != 3 {
		t.Fatalf("Uint32Gauge.TrySub(4) failure value = %d, want 3", got)
	}
	if loaded := gauge.Load(); loaded != 3 {
		t.Fatalf("Uint32Gauge.Load() after failed TrySub = %d, want 3", loaded)
	}
}

// TestUint32GaugeIncAndDec verifies single-unit gauge arithmetic.
func TestUint32GaugeIncAndDec(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge

	if got := gauge.Inc(); got != 1 {
		t.Fatalf("Uint32Gauge.Inc() = %d, want 1", got)
	}
	if got := gauge.Inc(); got != 2 {
		t.Fatalf("second Uint32Gauge.Inc() = %d, want 2", got)
	}
	if got := gauge.Dec(); got != 1 {
		t.Fatalf("Uint32Gauge.Dec() = %d, want 1", got)
	}
	if got := gauge.Load(); got != 1 {
		t.Fatalf("Uint32Gauge.Load() after Inc/Dec sequence = %d, want 1", got)
	}
}

// TestUint32GaugeSwap verifies explicit owner-controlled replacement semantics.
func TestUint32GaugeSwap(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(10)

	if old := gauge.Swap(25); old != 10 {
		t.Fatalf("Uint32Gauge.Swap(25) old value = %d, want 10", old)
	}
	if got := gauge.Load(); got != 25 {
		t.Fatalf("Uint32Gauge.Load() after Swap(25) = %d, want 25", got)
	}
}

// TestUint32GaugeCompareAndSwap verifies conditional owner-controlled transitions.
func TestUint32GaugeCompareAndSwap(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(25)

	if swapped := gauge.CompareAndSwap(10, 40); swapped {
		t.Fatal("Uint32Gauge.CompareAndSwap(10, 40) = true, want false")
	}
	if got := gauge.Load(); got != 25 {
		t.Fatalf("Uint32Gauge.Load() after failed CAS = %d, want 25", got)
	}
	if swapped := gauge.CompareAndSwap(25, 40); !swapped {
		t.Fatal("Uint32Gauge.CompareAndSwap(25, 40) = false, want true")
	}
	if got := gauge.Load(); got != 40 {
		t.Fatalf("Uint32Gauge.Load() after successful CAS = %d, want 40", got)
	}
}

// TestUint32GaugeExactBoundaryOperations verifies that a bounded gauge may
// legally reach the numeric boundary. The invariant violation starts only when
// an operation attempts to cross the boundary.
func TestUint32GaugeExactBoundaryOperations(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge

	if got := gauge.Add(maxUint32); got != maxUint32 {
		t.Fatalf("Uint32Gauge.Add(maxUint32) = %d, want %d", got, maxUint32)
	}
	if got := gauge.Sub(maxUint32); got != 0 {
		t.Fatalf("Uint32Gauge.Sub(maxUint32) = %d, want 0", got)
	}
}

// TestUint32GaugePanicsOnOverflow verifies bounded gauges do not silently wrap.
func TestUint32GaugePanicsOnOverflow(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(maxUint32)

	mustPanicWithValue(t, errUint32GaugeOverflow, func() {
		_ = gauge.Add(1)
	})
}

// TestUint32GaugePanicsOnUnderflow verifies bounded gauges reject negative
// current state instead of wrapping subtraction to a very large unsigned value.
func TestUint32GaugePanicsOnUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(10)

	mustPanicWithValue(t, errUint32GaugeUnderflow, func() {
		_ = gauge.Sub(11)
	})
}

// TestUint32GaugeIncPanicsOnOverflow verifies Inc preserves the same overflow
// invariant as Add, so the convenience method cannot bypass gauge safety.
func TestUint32GaugeIncPanicsOnOverflow(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge
	gauge.Store(maxUint32)

	mustPanicWithValue(t, errUint32GaugeOverflow, func() {
		_ = gauge.Inc()
	})
}

// TestUint32GaugeDecPanicsOnUnderflow verifies Dec preserves the same underflow
// invariant as Sub, so single-unit release cannot bypass gauge safety.
func TestUint32GaugeDecPanicsOnUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Uint32Gauge

	mustPanicWithValue(t, errUint32GaugeUnderflow, func() {
		_ = gauge.Dec()
	})
}

// TestUint32GaugeConcurrentBalancedAccounting verifies the bounded gauge keeps
// exact current-state accounting under deterministic concurrent updates.
func TestUint32GaugeConcurrentBalancedAccounting(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const iterations = 10_000

	var gauge Uint32Gauge
	runConcurrent(t, goroutines, func() {
		for range iterations {
			gauge.Add(2)
			gauge.Sub(1)
		}
	})

	want := uint32(goroutines * iterations)
	if got := gauge.Load(); got != want {
		t.Fatalf("Uint32Gauge.Load() after concurrent balanced accounting = %d, want %d", got, want)
	}
}
