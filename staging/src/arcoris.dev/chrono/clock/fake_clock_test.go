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

package clock

import "testing"

// TestNewFakeClockInitializesCurrentTime verifies that the constructor stores
// the exact initial fake time chosen by the test.
func TestNewFakeClockInitializesCurrentTime(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clock := NewFakeClock(start)

	mustEqualTime(t, "FakeClock.Now()", clock.Now(), start)
}

// TestNewFakeClockInitializesRegistries verifies that fake waiters, timers, and
// tickers can be registered without relying on lazy nil-map behavior.
//
// Lazy initialization still exists defensively, but the constructor should
// produce a fully ready fake clock.
func TestNewFakeClockInitializesRegistries(t *testing.T) {
	t.Parallel()

	clock := NewFakeClock(fakeClockTestTime())

	clock.mu.Lock()
	defer clock.mu.Unlock()

	if clock.waiters == nil {
		t.Fatal("FakeClock.waiters is nil")
	}

	if clock.timers == nil {
		t.Fatal("FakeClock.timers is nil")
	}

	if clock.tickers == nil {
		t.Fatal("FakeClock.tickers is nil")
	}

	if len(clock.waiters) != 0 {
		t.Fatalf("len(FakeClock.waiters) = %d, want 0", len(clock.waiters))
	}

	if len(clock.timers) != 0 {
		t.Fatalf("len(FakeClock.timers) = %d, want 0", len(clock.timers))
	}

	if len(clock.tickers) != 0 {
		t.Fatalf("len(FakeClock.tickers) = %d, want 0", len(clock.tickers))
	}
}

// TestFakeClockDoesNotAdvanceWithRealTime verifies the core fake-clock
// invariant: real runtime time must not change fake time.
func TestFakeClockDoesNotAdvanceWithRealTime(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clock := NewFakeClock(start)

	first := clock.Now()
	second := clock.Now()

	mustEqualTime(t, "first FakeClock.Now()", first, start)
	mustEqualTime(t, "second FakeClock.Now()", second, start)
}

// TestFakeClockRegistriesBecomeEmptyAfterDueDelivery verifies that due waiters,
// timers, and stopped tickers do not remain registered forever.
func TestFakeClockRegistriesBecomeEmptyAfterDueDelivery(t *testing.T) {
	t.Parallel()

	clock := NewFakeClock(fakeClockTestTime())

	waiter := clock.After(5)
	timer := clock.NewTimer(5)
	ticker := clock.NewTicker(5)

	clock.Step(5)

	_ = mustReceiveTime(t, waiter)
	_ = mustReceiveTime(t, timer.C())
	_ = mustReceiveTime(t, ticker.C())

	ticker.Stop()

	clock.mu.Lock()
	defer clock.mu.Unlock()

	if len(clock.waiters) != 0 {
		t.Fatalf("len(FakeClock.waiters) = %d, want 0", len(clock.waiters))
	}

	if len(clock.timers) != 0 {
		t.Fatalf("len(FakeClock.timers) = %d, want 0", len(clock.timers))
	}

	if len(clock.tickers) != 0 {
		t.Fatalf("len(FakeClock.tickers) = %d, want 0", len(clock.tickers))
	}
}
