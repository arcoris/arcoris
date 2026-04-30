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

// TestPaddedInt64ZeroValue verifies the raw signed primitive starts as zero.
func TestPaddedInt64ZeroValue(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	if got := value.Load(); got != 0 {
		t.Fatalf("zero-value PaddedInt64.Load() = %d, want 0", got)
	}
}

// TestPaddedInt64StoreAndLoad verifies raw owner-controlled signed publication.
func TestPaddedInt64StoreAndLoad(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	value.Store(-42)

	if got := value.Load(); got != -42 {
		t.Fatalf("PaddedInt64.Load() after Store(-42) = %d, want -42", got)
	}
}

// TestPaddedInt64AddPositiveDelta verifies raw positive signed addition.
func TestPaddedInt64AddPositiveDelta(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	if got := value.Add(10); got != 10 {
		t.Fatalf("PaddedInt64.Add(10) = %d, want 10", got)
	}

	if got := value.Add(5); got != 15 {
		t.Fatalf("PaddedInt64.Add(5) = %d, want 15", got)
	}

	if got := value.Load(); got != 15 {
		t.Fatalf("PaddedInt64.Load() = %d, want 15", got)
	}
}

// TestPaddedInt64AddNegativeDelta verifies raw negative signed addition.
func TestPaddedInt64AddNegativeDelta(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	value.Store(10)

	if got := value.Add(-3); got != 7 {
		t.Fatalf("PaddedInt64.Add(-3) = %d, want 7", got)
	}

	if got := value.Add(-10); got != -3 {
		t.Fatalf("PaddedInt64.Add(-10) = %d, want -3", got)
	}

	if got := value.Load(); got != -3 {
		t.Fatalf("PaddedInt64.Load() = %d, want -3", got)
	}
}

// TestPaddedInt64AddZero verifies zero deltas return the current observed value.
func TestPaddedInt64AddZero(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	value.Store(-17)

	if got := value.Add(0); got != -17 {
		t.Fatalf("PaddedInt64.Add(0) = %d, want -17", got)
	}

	if got := value.Load(); got != -17 {
		t.Fatalf("PaddedInt64.Load() after Add(0) = %d, want -17", got)
	}
}

// TestPaddedInt64IncAndDec verifies raw single-unit signed updates.
func TestPaddedInt64IncAndDec(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	if got := value.Inc(); got != 1 {
		t.Fatalf("PaddedInt64.Inc() = %d, want 1", got)
	}

	if got := value.Dec(); got != 0 {
		t.Fatalf("PaddedInt64.Dec() = %d, want 0", got)
	}

	if got := value.Dec(); got != -1 {
		t.Fatalf("second PaddedInt64.Dec() = %d, want -1", got)
	}

	if got := value.Load(); got != -1 {
		t.Fatalf("PaddedInt64.Load() = %d, want -1", got)
	}
}

// TestPaddedInt64Swap verifies unconditional raw replacement semantics.
func TestPaddedInt64Swap(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	value.Store(-10)

	if old := value.Swap(25); old != -10 {
		t.Fatalf("PaddedInt64.Swap(25) old value = %d, want -10", old)
	}

	if got := value.Load(); got != 25 {
		t.Fatalf("PaddedInt64.Load() after Swap(25) = %d, want 25", got)
	}
}

// TestPaddedInt64CompareAndSwap verifies expected-value raw transitions.
func TestPaddedInt64CompareAndSwap(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	value.Store(-10)

	if swapped := value.CompareAndSwap(9, 20); swapped {
		t.Fatalf("PaddedInt64.CompareAndSwap(9, 20) = true, want false")
	}

	if got := value.Load(); got != -10 {
		t.Fatalf("PaddedInt64.Load() after failed CAS = %d, want -10", got)
	}

	if swapped := value.CompareAndSwap(-10, 20); !swapped {
		t.Fatalf("PaddedInt64.CompareAndSwap(-10, 20) = false, want true")
	}

	if got := value.Load(); got != 20 {
		t.Fatalf("PaddedInt64.Load() after successful CAS = %d, want 20", got)
	}
}

// TestPaddedInt64RawAddWraps verifies raw signed primitives follow atomic wrap semantics.
func TestPaddedInt64RawAddWraps(t *testing.T) {
	t.Parallel()

	var value PaddedInt64

	value.Store(maxInt64)

	if got := value.Add(1); got != minInt64 {
		t.Fatalf("PaddedInt64.Add(1) from max int64 = %d, want %d", got, minInt64)
	}

	if got := value.Load(); got != minInt64 {
		t.Fatalf("PaddedInt64.Load() after wrap = %d, want %d", got, minInt64)
	}
}

// TestPaddedInt64ConcurrentBalancedAdd verifies deterministic signed updates under contention.
func TestPaddedInt64ConcurrentBalancedAdd(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const updatesPerGoroutine = 10_000

	var value PaddedInt64

	runConcurrentIndexed(t, goroutines, func(i int) {
		delta := int64(1)
		if i%2 != 0 {
			delta = -1
		}

		for range updatesPerGoroutine {
			value.Add(delta)
		}
	})

	if got := value.Load(); got != 0 {
		t.Fatalf("PaddedInt64.Load() after balanced concurrent updates = %d, want 0", got)
	}
}
