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

// TestUint64CounterZeroValue verifies lifetime counters start at zero.
func TestUint64CounterZeroValue(t *testing.T) {
	t.Parallel()

	var counter Uint64Counter

	if got := counter.Load(); got != 0 {
		t.Fatalf("zero-value Uint64Counter.Load() = %d, want 0", got)
	}
}

// TestUint64CounterAddAndInc verifies mutable counter increments only move forward.
func TestUint64CounterAddAndInc(t *testing.T) {
	t.Parallel()

	var counter Uint64Counter

	if got := counter.Add(0); got != 0 {
		t.Fatalf("Uint64Counter.Add(0) = %d, want 0", got)
	}

	if got := counter.Inc(); got != 1 {
		t.Fatalf("Uint64Counter.Inc() = %d, want 1", got)
	}

	if got := counter.Add(9); got != 10 {
		t.Fatalf("Uint64Counter.Add(9) = %d, want 10", got)
	}

	if got := counter.Load(); got != 10 {
		t.Fatalf("Uint64Counter.Load() = %d, want 10", got)
	}
}

// TestUint64CounterAddZeroDoesNotChangeValue verifies zero batches are no-ops.
func TestUint64CounterAddZeroDoesNotChangeValue(t *testing.T) {
	t.Parallel()

	var counter Uint64Counter

	counter.Add(42)

	if got := counter.Add(0); got != 42 {
		t.Fatalf("Uint64Counter.Add(0) = %d, want 42", got)
	}

	if got := counter.Load(); got != 42 {
		t.Fatalf("Uint64Counter.Load() after Add(0) = %d, want 42", got)
	}
}

// TestUint64CounterWrapsLikeUint64 verifies lifetime counters allow raw unsigned wrap.
func TestUint64CounterWrapsLikeUint64(t *testing.T) {
	t.Parallel()

	var counter Uint64Counter

	counter.Add(testMaxUint64)

	if got := counter.Load(); got != testMaxUint64 {
		t.Fatalf("Uint64Counter.Load() after Add(max) = %d, want %d", got, testMaxUint64)
	}

	if got := counter.Inc(); got != 0 {
		t.Fatalf("Uint64Counter.Inc() after max = %d, want 0", got)
	}

	if got := counter.Load(); got != 0 {
		t.Fatalf("Uint64Counter.Load() after wrap = %d, want 0", got)
	}
}

// TestUint64CounterConcurrentInc verifies deterministic single-event accounting under contention.
func TestUint64CounterConcurrentInc(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const incrementsPerGoroutine = 10_000

	var counter Uint64Counter
	runConcurrent(t, goroutines, func() {
		for range incrementsPerGoroutine {
			counter.Inc()
		}
	})

	want := uint64(goroutines * incrementsPerGoroutine)
	if got := counter.Load(); got != want {
		t.Fatalf("Uint64Counter.Load() after concurrent increments = %d, want %d", got, want)
	}
}

// TestUint64CounterConcurrentAdd verifies deterministic batched accounting under contention.
func TestUint64CounterConcurrentAdd(t *testing.T) {
	t.Parallel()

	const goroutines = 16
	const additionsPerGoroutine = 5_000
	const delta = 3

	var counter Uint64Counter
	runConcurrent(t, goroutines, func() {
		for range additionsPerGoroutine {
			counter.Add(delta)
		}
	})

	want := uint64(goroutines * additionsPerGoroutine * delta)
	if got := counter.Load(); got != want {
		t.Fatalf("Uint64Counter.Load() after concurrent Add = %d, want %d", got, want)
	}
}
