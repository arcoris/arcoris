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

import (
	"sync"
	"testing"
)

const (
	// testMaxUint64 is the test-owned uint64 upper bound used outside gauge
	// tests. Keeping it here avoids coupling raw primitive, counter, snapshot,
	// and delta tests to gauge-local production constants.
	testMaxUint64 = ^uint64(0)

	// testMaxUint32 is the test-owned uint32 upper bound used outside gauge
	// tests. It keeps 32-bit counter and raw primitive boundary checks explicit.
	testMaxUint32 = ^uint32(0)

	// testMaxInt64 and testMinInt64 are test-owned signed int64 boundaries for
	// raw primitive tests. Raw padded primitives wrap; gauge constants belong to
	// checked gauge implementations.
	testMaxInt64 = int64(1<<63 - 1)
	testMinInt64 = -testMaxInt64 - 1

	// testMaxInt32 and testMinInt32 are test-owned signed int32 boundaries for
	// raw primitive tests. They keep raw atomic semantics separate from gauge
	// invariant checks.
	testMaxInt32 = int32(1<<31 - 1)
	testMinInt32 = -testMaxInt32 - 1
)

// mustPanicWithValue verifies that fn panics with exactly the expected value.
//
// Gauge tests use exact panic values because those values are part of the
// package's debugging contract for invariant violations.
func mustPanicWithValue(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("function did not panic, want panic value %#v", want)
		}
		if got != want {
			t.Fatalf("panic value = %#v, want %#v", got, want)
		}
	}()

	fn()
}

// runConcurrent runs the same deterministic test body in several goroutines.
//
// It centralizes WaitGroup setup while keeping the actual accounting operations
// visible in the calling test.
func runConcurrent(t *testing.T, goroutines int, fn func()) {
	t.Helper()

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()
}

// runConcurrentIndexed runs one deterministic test body per goroutine index.
//
// Tests use it when each goroutine needs a stable role, such as balancing
// positive and negative signed primitive updates.
func runConcurrentIndexed(t *testing.T, goroutines int, fn func(index int)) {
	t.Helper()

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		go func() {
			defer wg.Done()
			fn(i)
		}()
	}
	wg.Wait()
}
