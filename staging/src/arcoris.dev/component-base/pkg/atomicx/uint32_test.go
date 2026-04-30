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

// TestPaddedUint32ZeroValue verifies the raw primitive starts as zero.
func TestPaddedUint32ZeroValue(t *testing.T) {
	t.Parallel()

	var value PaddedUint32

	if got := value.Load(); got != 0 {
		t.Fatalf("zero-value PaddedUint32.Load() = %d, want 0", got)
	}
}

// TestPaddedUint32StoreAndLoad verifies raw owner-controlled publication.
func TestPaddedUint32StoreAndLoad(t *testing.T) {
	t.Parallel()

	var value PaddedUint32

	value.Store(42)

	if got := value.Load(); got != 42 {
		t.Fatalf("PaddedUint32.Load() after Store(42) = %d, want 42", got)
	}
}

// TestPaddedUint32Add verifies raw unsigned addition without gauge invariants.
func TestPaddedUint32Add(t *testing.T) {
	t.Parallel()

	var value PaddedUint32

	if got := value.Add(10); got != 10 {
		t.Fatalf("PaddedUint32.Add(10) = %d, want 10", got)
	}

	if got := value.Add(5); got != 15 {
		t.Fatalf("PaddedUint32.Add(5) = %d, want 15", got)
	}

	if got := value.Load(); got != 15 {
		t.Fatalf("PaddedUint32.Load() = %d, want 15", got)
	}
}

// TestPaddedUint32AddZero verifies zero deltas return the current observed value.
func TestPaddedUint32AddZero(t *testing.T) {
	t.Parallel()

	var value PaddedUint32

	value.Store(17)

	if got := value.Add(0); got != 17 {
		t.Fatalf("PaddedUint32.Add(0) = %d, want 17", got)
	}

	if got := value.Load(); got != 17 {
		t.Fatalf("PaddedUint32.Load() after Add(0) = %d, want 17", got)
	}
}

// TestPaddedUint32Inc verifies Inc is raw Add(1) convenience behavior.
func TestPaddedUint32Inc(t *testing.T) {
	t.Parallel()

	var value PaddedUint32

	if got := value.Inc(); got != 1 {
		t.Fatalf("first PaddedUint32.Inc() = %d, want 1", got)
	}

	if got := value.Inc(); got != 2 {
		t.Fatalf("second PaddedUint32.Inc() = %d, want 2", got)
	}

	if got := value.Load(); got != 2 {
		t.Fatalf("PaddedUint32.Load() = %d, want 2", got)
	}
}

// TestPaddedUint32Swap verifies unconditional raw replacement semantics.
func TestPaddedUint32Swap(t *testing.T) {
	t.Parallel()

	var value PaddedUint32

	value.Store(10)

	if old := value.Swap(25); old != 10 {
		t.Fatalf("PaddedUint32.Swap(25) old value = %d, want 10", old)
	}

	if got := value.Load(); got != 25 {
		t.Fatalf("PaddedUint32.Load() after Swap(25) = %d, want 25", got)
	}
}

// TestPaddedUint32CompareAndSwap verifies expected-value raw transitions.
func TestPaddedUint32CompareAndSwap(t *testing.T) {
	t.Parallel()

	var value PaddedUint32

	value.Store(10)

	if swapped := value.CompareAndSwap(9, 20); swapped {
		t.Fatalf("PaddedUint32.CompareAndSwap(9, 20) = true, want false")
	}

	if got := value.Load(); got != 10 {
		t.Fatalf("PaddedUint32.Load() after failed CAS = %d, want 10", got)
	}

	if swapped := value.CompareAndSwap(10, 20); !swapped {
		t.Fatalf("PaddedUint32.CompareAndSwap(10, 20) = false, want true")
	}

	if got := value.Load(); got != 20 {
		t.Fatalf("PaddedUint32.Load() after successful CAS = %d, want 20", got)
	}
}

// TestPaddedUint32RawAddWraps verifies raw primitives allow unsigned wrap.
func TestPaddedUint32RawAddWraps(t *testing.T) {
	t.Parallel()

	var value PaddedUint32

	value.Store(^uint32(0))

	if got := value.Add(1); got != 0 {
		t.Fatalf("PaddedUint32.Add(1) from max uint32 = %d, want 0", got)
	}

	if got := value.Load(); got != 0 {
		t.Fatalf("PaddedUint32.Load() after wrap = %d, want 0", got)
	}
}

// TestPaddedUint32ConcurrentAdd verifies deterministic atomic increments under contention.
func TestPaddedUint32ConcurrentAdd(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const incrementsPerGoroutine = 10_000
	const want uint32 = goroutines * incrementsPerGoroutine

	var value PaddedUint32
	runConcurrent(t, goroutines, func() {
		for range incrementsPerGoroutine {
			value.Inc()
		}
	})

	if got := value.Load(); got != want {
		t.Fatalf("PaddedUint32.Load() after concurrent increments = %d, want %d", got, want)
	}
}
