// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package channelassert

import (
	"testing"
	"time"
)

// RequireReceive waits for a value until timeout and returns it.
func RequireReceive[T any](t testing.TB, ch <-chan T, timeout time.Duration) T {
	t.Helper()

	requirePositiveTimeout(t, timeout)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case got, ok := <-ch:
		if !ok {
			t.Fatal("channel closed before value was received")
		}
		return got
	case <-timer.C:
		t.Fatalf("did not receive value before safety timeout %s", timeout)
		var zero T
		return zero
	}
}

// RequireNoReceive fails if ch can be received from immediately.
func RequireNoReceive[T any](t testing.TB, ch <-chan T) {
	t.Helper()

	select {
	case _, ok := <-ch:
		if !ok {
			t.Fatal("channel closed, want open and empty")
		}
		t.Fatal("received value, want no immediate receive")
	default:
	}
}

// RequireClosed waits until ch is closed before timeout.
func RequireClosed[T any](t testing.TB, ch <-chan T, timeout time.Duration) {
	t.Helper()

	requirePositiveTimeout(t, timeout)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("received value before channel close")
		}
	case <-timer.C:
		t.Fatalf("channel did not close before safety timeout %s", timeout)
	}
}

// RequireSignal waits for a struct{} signal until timeout.
func RequireSignal(t testing.TB, ch <-chan struct{}, timeout time.Duration) {
	t.Helper()

	requirePositiveTimeout(t, timeout)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-ch:
	case <-timer.C:
		t.Fatalf("did not receive signal before safety timeout %s", timeout)
	}
}

// RequireNoSignal fails if ch can be received from immediately.
func RequireNoSignal(t testing.TB, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
		t.Fatal("received signal, want no immediate signal")
	default:
	}
}

func requirePositiveTimeout(t testing.TB, timeout time.Duration) {
	t.Helper()

	if timeout <= 0 {
		t.Fatalf("timeout = %s, want positive", timeout)
	}
}
