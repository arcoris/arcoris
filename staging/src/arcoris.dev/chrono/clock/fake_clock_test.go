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
	clk := NewFakeClock(start)

	mustEqualTime(t, "FakeClock.Now()", clk.Now(), start)
}

// TestNewFakeClockInitializesRegistries verifies that fake waiters, timers, and
// tickers can be registered without relying on lazy nil-map behavior.
//
// Lazy initialization still exists defensively, but the constructor should
// produce a fully ready fake clock.
func TestNewFakeClockInitializesRegistries(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	clk.mu.Lock()
	defer clk.mu.Unlock()

	if clk.waiters == nil {
		t.Fatal("FakeClock.waiters is nil")
	}

	if clk.timers == nil {
		t.Fatal("FakeClock.timers is nil")
	}

	if clk.tickers == nil {
		t.Fatal("FakeClock.tickers is nil")
	}

	if len(clk.waiters) != 0 {
		t.Fatalf("len(FakeClock.waiters) = %d, want 0", len(clk.waiters))
	}

	if len(clk.timers) != 0 {
		t.Fatalf("len(FakeClock.timers) = %d, want 0", len(clk.timers))
	}

	if len(clk.tickers) != 0 {
		t.Fatalf("len(FakeClock.tickers) = %d, want 0", len(clk.tickers))
	}
}

// TestFakeClockDoesNotAdvanceWithRealTime verifies the core fake-clock
// invariant: real runtime time must not change fake time.
func TestFakeClockDoesNotAdvanceWithRealTime(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	first := clk.Now()
	second := clk.Now()

	mustEqualTime(t, "first FakeClock.Now()", first, start)
	mustEqualTime(t, "second FakeClock.Now()", second, start)
}

// TestFakeClockRegistriesBecomeEmptyAfterDueDelivery verifies that due waiters,
// timers, and stopped tickers do not remain registered forever.
func TestFakeClockRegistriesBecomeEmptyAfterDueDelivery(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	waiter := clk.After(5)
	timer := clk.NewTimer(5)
	ticker := clk.NewTicker(5)

	clk.Step(5)

	_ = mustReceiveTime(t, waiter)
	_ = mustReceiveTime(t, timer.C())
	_ = mustReceiveTime(t, ticker.C())

	ticker.Stop()

	clk.mu.Lock()
	defer clk.mu.Unlock()

	if len(clk.waiters) != 0 {
		t.Fatalf("len(FakeClock.waiters) = %d, want 0", len(clk.waiters))
	}

	if len(clk.timers) != 0 {
		t.Fatalf("len(FakeClock.timers) = %d, want 0", len(clk.timers))
	}

	if len(clk.tickers) != 0 {
		t.Fatalf("len(FakeClock.tickers) = %d, want 0", len(clk.tickers))
	}
}
