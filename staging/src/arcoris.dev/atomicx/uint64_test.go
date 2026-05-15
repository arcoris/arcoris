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

// TestPaddedUint64ZeroValue verifies the raw primitive starts as zero.
func TestPaddedUint64ZeroValue(t *testing.T) {
	t.Parallel()

	var val PaddedUint64

	if got := val.Load(); got != 0 {
		t.Fatalf("zero-value PaddedUint64.Load() = %d, want 0", got)
	}
}

// TestPaddedUint64StoreAndLoad verifies raw owner-controlled publication.
func TestPaddedUint64StoreAndLoad(t *testing.T) {
	t.Parallel()

	var val PaddedUint64

	val.Store(42)

	if got := val.Load(); got != 42 {
		t.Fatalf("PaddedUint64.Load() after Store(42) = %d, want 42", got)
	}
}

// TestPaddedUint64Add verifies raw unsigned addition without gauge invariants.
func TestPaddedUint64Add(t *testing.T) {
	t.Parallel()

	var val PaddedUint64

	if got := val.Add(10); got != 10 {
		t.Fatalf("PaddedUint64.Add(10) = %d, want 10", got)
	}

	if got := val.Add(5); got != 15 {
		t.Fatalf("PaddedUint64.Add(5) = %d, want 15", got)
	}

	if got := val.Load(); got != 15 {
		t.Fatalf("PaddedUint64.Load() = %d, want 15", got)
	}
}

// TestPaddedUint64AddZero verifies zero deltas return the current observed value.
func TestPaddedUint64AddZero(t *testing.T) {
	t.Parallel()

	var val PaddedUint64

	val.Store(17)

	if got := val.Add(0); got != 17 {
		t.Fatalf("PaddedUint64.Add(0) = %d, want 17", got)
	}

	if got := val.Load(); got != 17 {
		t.Fatalf("PaddedUint64.Load() after Add(0) = %d, want 17", got)
	}
}

// TestPaddedUint64Inc verifies Inc is raw Add(1) convenience behavior.
func TestPaddedUint64Inc(t *testing.T) {
	t.Parallel()

	var val PaddedUint64

	if got := val.Inc(); got != 1 {
		t.Fatalf("first PaddedUint64.Inc() = %d, want 1", got)
	}

	if got := val.Inc(); got != 2 {
		t.Fatalf("second PaddedUint64.Inc() = %d, want 2", got)
	}

	if got := val.Load(); got != 2 {
		t.Fatalf("PaddedUint64.Load() = %d, want 2", got)
	}
}

// TestPaddedUint64Swap verifies unconditional raw replacement semantics.
func TestPaddedUint64Swap(t *testing.T) {
	t.Parallel()

	var val PaddedUint64

	val.Store(10)

	if old := val.Swap(25); old != 10 {
		t.Fatalf("PaddedUint64.Swap(25) old value = %d, want 10", old)
	}

	if got := val.Load(); got != 25 {
		t.Fatalf("PaddedUint64.Load() after Swap(25) = %d, want 25", got)
	}
}

// TestPaddedUint64CompareAndSwap verifies expected-value raw transitions.
func TestPaddedUint64CompareAndSwap(t *testing.T) {
	t.Parallel()

	var val PaddedUint64

	val.Store(10)

	if swapped := val.CompareAndSwap(9, 20); swapped {
		t.Fatalf("PaddedUint64.CompareAndSwap(9, 20) = true, want false")
	}

	if got := val.Load(); got != 10 {
		t.Fatalf("PaddedUint64.Load() after failed CAS = %d, want 10", got)
	}

	if swapped := val.CompareAndSwap(10, 20); !swapped {
		t.Fatalf("PaddedUint64.CompareAndSwap(10, 20) = false, want true")
	}

	if got := val.Load(); got != 20 {
		t.Fatalf("PaddedUint64.Load() after successful CAS = %d, want 20", got)
	}
}

// TestPaddedUint64RawAddWraps verifies raw primitives allow unsigned wrap.
func TestPaddedUint64RawAddWraps(t *testing.T) {
	t.Parallel()

	var val PaddedUint64

	val.Store(testMaxUint64)

	if got := val.Add(1); got != 0 {
		t.Fatalf("PaddedUint64.Add(1) from max uint64 = %d, want 0", got)
	}

	if got := val.Load(); got != 0 {
		t.Fatalf("PaddedUint64.Load() after wrap = %d, want 0", got)
	}
}

// TestPaddedUint64ConcurrentAdd verifies deterministic atomic increments under contention.
func TestPaddedUint64ConcurrentAdd(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const incrementsPerGoroutine = 10_000
	const want = goroutines * incrementsPerGoroutine

	var val PaddedUint64
	runConcurrent(t, goroutines, func() {
		for range incrementsPerGoroutine {
			val.Inc()
		}
	})

	if got := val.Load(); got != want {
		t.Fatalf("PaddedUint64.Load() after concurrent increments = %d, want %d", got, want)
	}
}
