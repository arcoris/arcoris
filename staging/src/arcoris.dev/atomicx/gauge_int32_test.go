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

// TestInt32GaugeZeroValueIsUsable verifies bounded signed gauges work without initialization.
func TestInt32GaugeZeroValueIsUsable(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge

	if got := gauge.Load(); got != 0 {
		t.Fatalf("zero-value Int32Gauge.Load() = %d, want 0", got)
	}
}

// TestInt32GaugeStoreAndLoad verifies owner-controlled signed state publication.
func TestInt32GaugeStoreAndLoad(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(-42)

	if got := gauge.Load(); got != -42 {
		t.Fatalf("Int32Gauge.Load() after Store(-42) = %d, want -42", got)
	}
}

// TestInt32GaugeAddSubAndZeroDeltas verifies signed arithmetic in both directions.
func TestInt32GaugeAddSubAndZeroDeltas(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge

	if got := gauge.Add(10); got != 10 {
		t.Fatalf("Int32Gauge.Add(10) = %d, want 10", got)
	}
	if got := gauge.Add(0); got != 10 {
		t.Fatalf("Int32Gauge.Add(0) = %d, want 10", got)
	}
	if got := gauge.Add(-4); got != 6 {
		t.Fatalf("Int32Gauge.Add(-4) = %d, want 6", got)
	}
	if got := gauge.Sub(10); got != -4 {
		t.Fatalf("Int32Gauge.Sub(10) = %d, want -4", got)
	}
	if got := gauge.Sub(0); got != -4 {
		t.Fatalf("Int32Gauge.Sub(0) = %d, want -4", got)
	}
	if got := gauge.Sub(-6); got != 2 {
		t.Fatalf("Int32Gauge.Sub(-6) = %d, want 2", got)
	}
}

// TestInt32GaugeTryAddSuccess verifies checked signed addition updates state.
func TestInt32GaugeTryAddSuccess(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(-10)

	got, ok := gauge.TryAdd(15)
	if !ok {
		t.Fatal("Int32Gauge.TryAdd(15) ok = false, want true")
	}
	if got != 5 {
		t.Fatalf("Int32Gauge.TryAdd(15) value = %d, want 5", got)
	}
	if loaded := gauge.Load(); loaded != 5 {
		t.Fatalf("Int32Gauge.Load() after TryAdd(15) = %d, want 5", loaded)
	}
}

// TestInt32GaugeTryAddOverflowLeavesStateUnchanged verifies positive overflow
// is reported without mutating state. Failed checked arithmetic must be safe for
// callers that treat refusal as normal control flow.
func TestInt32GaugeTryAddOverflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(maxInt32)

	got, ok := gauge.TryAdd(1)
	if ok {
		t.Fatal("Int32Gauge.TryAdd(1) at max ok = true, want false")
	}
	if got != maxInt32 {
		t.Fatalf("Int32Gauge.TryAdd(1) failure value = %d, want %d", got, maxInt32)
	}
	if loaded := gauge.Load(); loaded != maxInt32 {
		t.Fatalf("Int32Gauge.Load() after failed TryAdd overflow = %d, want %d", loaded, maxInt32)
	}
}

// TestInt32GaugeTryAddUnderflowLeavesStateUnchanged verifies negative underflow
// is reported without mutating state. Signed gauges must not wrap past the lower
// boundary.
func TestInt32GaugeTryAddUnderflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(minInt32)

	got, ok := gauge.TryAdd(-1)
	if ok {
		t.Fatal("Int32Gauge.TryAdd(-1) at min ok = true, want false")
	}
	if got != minInt32 {
		t.Fatalf("Int32Gauge.TryAdd(-1) failure value = %d, want %d", got, minInt32)
	}
	if loaded := gauge.Load(); loaded != minInt32 {
		t.Fatalf("Int32Gauge.Load() after failed TryAdd underflow = %d, want %d", loaded, minInt32)
	}
}

// TestInt32GaugeTrySubSuccess verifies checked signed subtraction updates state.
func TestInt32GaugeTrySubSuccess(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(10)

	got, ok := gauge.TrySub(15)
	if !ok {
		t.Fatal("Int32Gauge.TrySub(15) ok = false, want true")
	}
	if got != -5 {
		t.Fatalf("Int32Gauge.TrySub(15) value = %d, want -5", got)
	}
	if loaded := gauge.Load(); loaded != -5 {
		t.Fatalf("Int32Gauge.Load() after TrySub(15) = %d, want -5", loaded)
	}
}

// TestInt32GaugeTrySubOverflowLeavesStateUnchanged verifies subtracting negative
// deltas checks the upper bound without computing through a wrapped value.
func TestInt32GaugeTrySubOverflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(maxInt32)

	got, ok := gauge.TrySub(-1)
	if ok {
		t.Fatal("Int32Gauge.TrySub(-1) at max ok = true, want false")
	}
	if got != maxInt32 {
		t.Fatalf("Int32Gauge.TrySub(-1) failure value = %d, want %d", got, maxInt32)
	}
	if loaded := gauge.Load(); loaded != maxInt32 {
		t.Fatalf("Int32Gauge.Load() after failed TrySub overflow = %d, want %d", loaded, maxInt32)
	}
}

// TestInt32GaugeTrySubUnderflowLeavesStateUnchanged verifies subtracting
// positive deltas checks the lower bound and leaves the current state intact.
func TestInt32GaugeTrySubUnderflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(minInt32)

	got, ok := gauge.TrySub(1)
	if ok {
		t.Fatal("Int32Gauge.TrySub(1) at min ok = true, want false")
	}
	if got != minInt32 {
		t.Fatalf("Int32Gauge.TrySub(1) failure value = %d, want %d", got, minInt32)
	}
	if loaded := gauge.Load(); loaded != minInt32 {
		t.Fatalf("Int32Gauge.Load() after failed TrySub underflow = %d, want %d", loaded, minInt32)
	}
}

// TestInt32GaugeTrySubHandlesMinInt32Delta verifies TrySub never computes
// -minInt32. That negation is not representable and would turn an invariant
// check into undefined-looking signed wrap behavior.
func TestInt32GaugeTrySubHandlesMinInt32Delta(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge

	got, ok := gauge.TrySub(minInt32)
	if ok {
		t.Fatal("Int32Gauge.TrySub(minInt32) from 0 ok = true, want false")
	}
	if got != 0 {
		t.Fatalf("Int32Gauge.TrySub(minInt32) failure value = %d, want 0", got)
	}
	if loaded := gauge.Load(); loaded != 0 {
		t.Fatalf("Int32Gauge.Load() after failed TrySub(minInt32) = %d, want 0", loaded)
	}
}

// TestInt32GaugeIncAndDec verifies single-unit signed gauge arithmetic.
func TestInt32GaugeIncAndDec(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge

	if got := gauge.Inc(); got != 1 {
		t.Fatalf("Int32Gauge.Inc() = %d, want 1", got)
	}
	if got := gauge.Dec(); got != 0 {
		t.Fatalf("Int32Gauge.Dec() = %d, want 0", got)
	}
	if got := gauge.Dec(); got != -1 {
		t.Fatalf("second Int32Gauge.Dec() = %d, want -1", got)
	}
}

// TestInt32GaugeSwap verifies explicit owner-controlled replacement semantics.
func TestInt32GaugeSwap(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(-10)

	if old := gauge.Swap(25); old != -10 {
		t.Fatalf("Int32Gauge.Swap(25) old value = %d, want -10", old)
	}
	if got := gauge.Load(); got != 25 {
		t.Fatalf("Int32Gauge.Load() after Swap(25) = %d, want 25", got)
	}
}

// TestInt32GaugeCompareAndSwap verifies conditional owner-controlled transitions.
func TestInt32GaugeCompareAndSwap(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(25)

	if swapped := gauge.CompareAndSwap(-10, 40); swapped {
		t.Fatal("Int32Gauge.CompareAndSwap(-10, 40) = true, want false")
	}
	if got := gauge.Load(); got != 25 {
		t.Fatalf("Int32Gauge.Load() after failed CAS = %d, want 25", got)
	}
	if swapped := gauge.CompareAndSwap(25, 40); !swapped {
		t.Fatal("Int32Gauge.CompareAndSwap(25, 40) = false, want true")
	}
	if got := gauge.Load(); got != 40 {
		t.Fatalf("Int32Gauge.Load() after successful CAS = %d, want 40", got)
	}
}

// TestInt32GaugeExactBoundaryOperations verifies legal transitions at both
// signed limits. Reaching min or max is valid; crossing either limit is the
// invariant violation tested by panic and Try* failure cases.
func TestInt32GaugeExactBoundaryOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		start int32
		op    func(*Int32Gauge) int32
		want  int32
	}{
		{name: "AddToMax", start: maxInt32 - 1, op: func(g *Int32Gauge) int32 { return g.Add(1) }, want: maxInt32},
		{name: "AddToMin", start: minInt32 + 1, op: func(g *Int32Gauge) int32 { return g.Add(-1) }, want: minInt32},
		{name: "SubToMax", start: maxInt32 - 1, op: func(g *Int32Gauge) int32 { return g.Sub(-1) }, want: maxInt32},
		{name: "SubToMin", start: minInt32 + 1, op: func(g *Int32Gauge) int32 { return g.Sub(1) }, want: minInt32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var gauge Int32Gauge
			gauge.Store(tt.start)

			if got := tt.op(&gauge); got != tt.want {
				t.Fatalf("%s result = %d, want %d", tt.name, got, tt.want)
			}
			if got := gauge.Load(); got != tt.want {
				t.Fatalf("%s Load() = %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}

// TestInt32GaugePanicsOnAddOverflow verifies Add rejects positive overflow as
// an invariant violation instead of silently wrapping signed current state.
func TestInt32GaugePanicsOnAddOverflow(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(maxInt32)

	mustPanicWithValue(t, errInt32GaugeOverflow, func() {
		_ = gauge.Add(1)
	})
}

// TestInt32GaugePanicsOnAddUnderflow verifies Add rejects negative underflow as
// an invariant violation instead of silently wrapping signed current state.
func TestInt32GaugePanicsOnAddUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(minInt32)

	mustPanicWithValue(t, errInt32GaugeUnderflow, func() {
		_ = gauge.Add(-1)
	})
}

// TestInt32GaugePanicsOnSubOverflow verifies Sub rejects upward overflow when a
// negative delta would move the gauge beyond maxInt32.
func TestInt32GaugePanicsOnSubOverflow(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(maxInt32)

	mustPanicWithValue(t, errInt32GaugeOverflow, func() {
		_ = gauge.Sub(-1)
	})
}

// TestInt32GaugePanicsOnSubUnderflow verifies Sub rejects downward underflow
// when a positive delta would move the gauge below minInt32.
func TestInt32GaugePanicsOnSubUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(minInt32)

	mustPanicWithValue(t, errInt32GaugeUnderflow, func() {
		_ = gauge.Sub(1)
	})
}

// TestInt32GaugeIncPanicsOnOverflow verifies Inc preserves Add overflow checks.
func TestInt32GaugeIncPanicsOnOverflow(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(maxInt32)

	mustPanicWithValue(t, errInt32GaugeOverflow, func() {
		_ = gauge.Inc()
	})
}

// TestInt32GaugeDecPanicsOnUnderflow verifies Dec preserves Sub underflow checks.
func TestInt32GaugeDecPanicsOnUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Int32Gauge
	gauge.Store(minInt32)

	mustPanicWithValue(t, errInt32GaugeUnderflow, func() {
		_ = gauge.Dec()
	})
}

// TestInt32GaugeConcurrentSignedUpdates verifies deterministic signed accounting under contention.
func TestInt32GaugeConcurrentSignedUpdates(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const iterations = 10_000

	var gauge Int32Gauge
	runConcurrent(t, goroutines, func() {
		for range iterations {
			gauge.Add(2)
			gauge.Sub(1)
		}
	})

	want := int32(goroutines * iterations)
	if got := gauge.Load(); got != want {
		t.Fatalf("Int32Gauge.Load() after concurrent signed updates = %d, want %d", got, want)
	}
}
