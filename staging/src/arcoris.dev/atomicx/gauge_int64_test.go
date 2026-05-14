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

// TestInt64GaugeZeroValueIsUsable verifies signed gauges are ready for use without initialization.
func TestInt64GaugeZeroValueIsUsable(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge

	if got := gauge.Load(); got != 0 {
		t.Fatalf("zero-value Int64Gauge.Load() = %d, want 0", got)
	}
}

// TestInt64GaugeStoreAndLoad verifies owner-controlled signed state publication.
func TestInt64GaugeStoreAndLoad(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(-42)

	if got := gauge.Load(); got != -42 {
		t.Fatalf("Int64Gauge.Load() after Store(-42) = %d, want -42", got)
	}
}

// TestInt64GaugeAddSubAndZeroDeltas verifies signed arithmetic in both directions.
func TestInt64GaugeAddSubAndZeroDeltas(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge

	if got := gauge.Add(10); got != 10 {
		t.Fatalf("Int64Gauge.Add(10) = %d, want 10", got)
	}
	if got := gauge.Add(0); got != 10 {
		t.Fatalf("Int64Gauge.Add(0) = %d, want 10", got)
	}
	if got := gauge.Add(-4); got != 6 {
		t.Fatalf("Int64Gauge.Add(-4) = %d, want 6", got)
	}
	if got := gauge.Sub(10); got != -4 {
		t.Fatalf("Int64Gauge.Sub(10) = %d, want -4", got)
	}
	if got := gauge.Sub(0); got != -4 {
		t.Fatalf("Int64Gauge.Sub(0) = %d, want -4", got)
	}
	if got := gauge.Sub(-6); got != 2 {
		t.Fatalf("Int64Gauge.Sub(-6) = %d, want 2", got)
	}
}

// TestInt64GaugeTryAddSuccess verifies checked signed addition updates state.
func TestInt64GaugeTryAddSuccess(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(-10)

	got, ok := gauge.TryAdd(15)
	if !ok {
		t.Fatal("Int64Gauge.TryAdd(15) ok = false, want true")
	}
	if got != 5 {
		t.Fatalf("Int64Gauge.TryAdd(15) value = %d, want 5", got)
	}
	if loaded := gauge.Load(); loaded != 5 {
		t.Fatalf("Int64Gauge.Load() after TryAdd(15) = %d, want 5", loaded)
	}
}

// TestInt64GaugeTryAddOverflowLeavesStateUnchanged verifies positive overflow
// is reported without mutating state. Failed checked arithmetic must be safe for
// callers that treat refusal as normal control flow.
func TestInt64GaugeTryAddOverflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(maxInt64)

	got, ok := gauge.TryAdd(1)
	if ok {
		t.Fatal("Int64Gauge.TryAdd(1) at max ok = true, want false")
	}
	if got != maxInt64 {
		t.Fatalf("Int64Gauge.TryAdd(1) failure value = %d, want %d", got, maxInt64)
	}
	if loaded := gauge.Load(); loaded != maxInt64 {
		t.Fatalf("Int64Gauge.Load() after failed TryAdd overflow = %d, want %d", loaded, maxInt64)
	}
}

// TestInt64GaugeTryAddUnderflowLeavesStateUnchanged verifies negative underflow
// is reported without mutating state. Signed gauges must not wrap past the lower
// boundary.
func TestInt64GaugeTryAddUnderflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(minInt64)

	got, ok := gauge.TryAdd(-1)
	if ok {
		t.Fatal("Int64Gauge.TryAdd(-1) at min ok = true, want false")
	}
	if got != minInt64 {
		t.Fatalf("Int64Gauge.TryAdd(-1) failure value = %d, want %d", got, minInt64)
	}
	if loaded := gauge.Load(); loaded != minInt64 {
		t.Fatalf("Int64Gauge.Load() after failed TryAdd underflow = %d, want %d", loaded, minInt64)
	}
}

// TestInt64GaugeTrySubSuccess verifies checked signed subtraction updates state.
func TestInt64GaugeTrySubSuccess(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(10)

	got, ok := gauge.TrySub(15)
	if !ok {
		t.Fatal("Int64Gauge.TrySub(15) ok = false, want true")
	}
	if got != -5 {
		t.Fatalf("Int64Gauge.TrySub(15) value = %d, want -5", got)
	}
	if loaded := gauge.Load(); loaded != -5 {
		t.Fatalf("Int64Gauge.Load() after TrySub(15) = %d, want -5", loaded)
	}
}

// TestInt64GaugeTrySubOverflowLeavesStateUnchanged verifies subtracting negative
// deltas checks the upper bound without computing through a wrapped value.
func TestInt64GaugeTrySubOverflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(maxInt64)

	got, ok := gauge.TrySub(-1)
	if ok {
		t.Fatal("Int64Gauge.TrySub(-1) at max ok = true, want false")
	}
	if got != maxInt64 {
		t.Fatalf("Int64Gauge.TrySub(-1) failure value = %d, want %d", got, maxInt64)
	}
	if loaded := gauge.Load(); loaded != maxInt64 {
		t.Fatalf("Int64Gauge.Load() after failed TrySub overflow = %d, want %d", loaded, maxInt64)
	}
}

// TestInt64GaugeTrySubUnderflowLeavesStateUnchanged verifies subtracting
// positive deltas checks the lower bound and leaves the current state intact.
func TestInt64GaugeTrySubUnderflowLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(minInt64)

	got, ok := gauge.TrySub(1)
	if ok {
		t.Fatal("Int64Gauge.TrySub(1) at min ok = true, want false")
	}
	if got != minInt64 {
		t.Fatalf("Int64Gauge.TrySub(1) failure value = %d, want %d", got, minInt64)
	}
	if loaded := gauge.Load(); loaded != minInt64 {
		t.Fatalf("Int64Gauge.Load() after failed TrySub underflow = %d, want %d", loaded, minInt64)
	}
}

// TestInt64GaugeTrySubHandlesMinInt64Delta verifies TrySub never computes
// -minInt64. That negation is not representable and would turn an invariant
// check into undefined-looking signed wrap behavior.
func TestInt64GaugeTrySubHandlesMinInt64Delta(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge

	got, ok := gauge.TrySub(minInt64)
	if ok {
		t.Fatal("Int64Gauge.TrySub(minInt64) from 0 ok = true, want false")
	}
	if got != 0 {
		t.Fatalf("Int64Gauge.TrySub(minInt64) failure value = %d, want 0", got)
	}
	if loaded := gauge.Load(); loaded != 0 {
		t.Fatalf("Int64Gauge.Load() after failed TrySub(minInt64) = %d, want 0", loaded)
	}
}

// TestInt64GaugeIncAndDec verifies single-unit signed gauge arithmetic.
func TestInt64GaugeIncAndDec(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge

	if got := gauge.Inc(); got != 1 {
		t.Fatalf("Int64Gauge.Inc() = %d, want 1", got)
	}
	if got := gauge.Dec(); got != 0 {
		t.Fatalf("Int64Gauge.Dec() = %d, want 0", got)
	}
	if got := gauge.Dec(); got != -1 {
		t.Fatalf("second Int64Gauge.Dec() = %d, want -1", got)
	}
}

// TestInt64GaugeSwap verifies explicit owner-controlled replacement semantics.
func TestInt64GaugeSwap(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(-10)

	if old := gauge.Swap(25); old != -10 {
		t.Fatalf("Int64Gauge.Swap(25) old value = %d, want -10", old)
	}
	if got := gauge.Load(); got != 25 {
		t.Fatalf("Int64Gauge.Load() after Swap(25) = %d, want 25", got)
	}
}

// TestInt64GaugeCompareAndSwap verifies conditional owner-controlled transitions.
func TestInt64GaugeCompareAndSwap(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(25)

	if swapped := gauge.CompareAndSwap(-10, 40); swapped {
		t.Fatal("Int64Gauge.CompareAndSwap(-10, 40) = true, want false")
	}
	if got := gauge.Load(); got != 25 {
		t.Fatalf("Int64Gauge.Load() after failed CAS = %d, want 25", got)
	}
	if swapped := gauge.CompareAndSwap(25, 40); !swapped {
		t.Fatal("Int64Gauge.CompareAndSwap(25, 40) = false, want true")
	}
	if got := gauge.Load(); got != 40 {
		t.Fatalf("Int64Gauge.Load() after successful CAS = %d, want 40", got)
	}
}

// TestInt64GaugeExactBoundaryOperations verifies legal transitions at both
// signed limits. Reaching min or max is valid; crossing either limit is the
// invariant violation tested by panic and Try* failure cases.
func TestInt64GaugeExactBoundaryOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		start int64
		op    func(*Int64Gauge) int64
		want  int64
	}{
		{name: "AddToMax", start: maxInt64 - 1, op: func(g *Int64Gauge) int64 { return g.Add(1) }, want: maxInt64},
		{name: "AddToMin", start: minInt64 + 1, op: func(g *Int64Gauge) int64 { return g.Add(-1) }, want: minInt64},
		{name: "SubToMax", start: maxInt64 - 1, op: func(g *Int64Gauge) int64 { return g.Sub(-1) }, want: maxInt64},
		{name: "SubToMin", start: minInt64 + 1, op: func(g *Int64Gauge) int64 { return g.Sub(1) }, want: minInt64},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gauge Int64Gauge
			gauge.Store(tc.start)

			if got := tc.op(&gauge); got != tc.want {
				t.Fatalf("%s result = %d, want %d", tc.name, got, tc.want)
			}
			if got := gauge.Load(); got != tc.want {
				t.Fatalf("%s Load() = %d, want %d", tc.name, got, tc.want)
			}
		})
	}
}

// TestInt64GaugePanicsOnAddOverflow verifies Add rejects positive overflow as
// an invariant violation instead of silently wrapping signed current state.
func TestInt64GaugePanicsOnAddOverflow(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(maxInt64)

	mustPanicWithValue(t, errInt64GaugeOverflow, func() {
		_ = gauge.Add(1)
	})
}

// TestInt64GaugePanicsOnAddUnderflow verifies Add rejects negative underflow as
// an invariant violation instead of silently wrapping signed current state.
func TestInt64GaugePanicsOnAddUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(minInt64)

	mustPanicWithValue(t, errInt64GaugeUnderflow, func() {
		_ = gauge.Add(-1)
	})
}

// TestInt64GaugePanicsOnSubOverflow verifies Sub rejects upward overflow when a
// negative delta would move the gauge beyond maxInt64.
func TestInt64GaugePanicsOnSubOverflow(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(maxInt64)

	mustPanicWithValue(t, errInt64GaugeOverflow, func() {
		_ = gauge.Sub(-1)
	})
}

// TestInt64GaugePanicsOnSubUnderflow verifies Sub rejects downward underflow
// when a positive delta would move the gauge below minInt64.
func TestInt64GaugePanicsOnSubUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(minInt64)

	mustPanicWithValue(t, errInt64GaugeUnderflow, func() {
		_ = gauge.Sub(1)
	})
}

// TestInt64GaugeIncPanicsOnOverflow verifies Inc preserves Add overflow checks.
func TestInt64GaugeIncPanicsOnOverflow(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(maxInt64)

	mustPanicWithValue(t, errInt64GaugeOverflow, func() {
		_ = gauge.Inc()
	})
}

// TestInt64GaugeDecPanicsOnUnderflow verifies Dec preserves Sub underflow checks.
func TestInt64GaugeDecPanicsOnUnderflow(t *testing.T) {
	t.Parallel()

	var gauge Int64Gauge
	gauge.Store(minInt64)

	mustPanicWithValue(t, errInt64GaugeUnderflow, func() {
		_ = gauge.Dec()
	})
}

// TestInt64GaugeConcurrentSignedUpdates verifies deterministic signed accounting under contention.
func TestInt64GaugeConcurrentSignedUpdates(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const iterations = 10_000

	var gauge Int64Gauge
	runConcurrent(t, goroutines, func() {
		for range iterations {
			gauge.Add(2)
			gauge.Sub(1)
		}
	})

	want := int64(goroutines * iterations)
	if got := gauge.Load(); got != want {
		t.Fatalf("Int64Gauge.Load() after concurrent signed updates = %d, want %d", got, want)
	}
}
