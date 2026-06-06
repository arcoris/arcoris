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

package clock

import (
	channelassert "arcoris.dev/testutil/channel"

	"sync"
	"testing"
	"time"
)

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

// TestFakeClockZeroValueReportsZeroTime verifies the documented zero-value
// contract for tests that intentionally use a zero-time fake clock.
func TestFakeClockZeroValueReportsZeroTime(t *testing.T) {
	t.Parallel()

	var clk FakeClock

	mustEqualTime(t, "zero FakeClock.Now()", clk.Now(), time.Time{})
	if got := clk.Since(time.Time{}); got != 0 {
		t.Fatalf("zero FakeClock.Since(zero) = %s, want 0", got)
	}
	if got := clk.Until(time.Time{}); got != 0 {
		t.Fatalf("zero FakeClock.Until(zero) = %s, want 0", got)
	}
}

// TestFakeClockZeroValueCanRegisterWaitersTimersAndTickers verifies that lazy
// registry initialization makes the zero value usable for all fake primitives.
func TestFakeClockZeroValueCanRegisterWaitersTimersAndTickers(t *testing.T) {
	t.Parallel()

	var clk FakeClock

	waiter := clk.After(time.Second)
	timer := clk.NewTimer(time.Second)
	ticker := clk.NewTicker(time.Second)
	defer ticker.Stop()

	clk.Step(time.Second)

	want := time.Time{}.Add(time.Second)

	gotWaiter := channelassert.RequireReceive(t, waiter, clockTestTimeout)
	gotTimer := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	gotTicker := channelassert.RequireReceive(t, ticker.C(), clockTestTimeout)

	mustEqualTime(t, "zero FakeClock After delivery", gotWaiter, want)
	mustEqualTime(t, "zero FakeClock timer delivery", gotTimer, want)
	mustEqualTime(t, "zero FakeClock ticker delivery", gotTicker, want)
}

// TestFakeClockConcurrentOperationsAreRaceSafe exercises the fake clock's core
// operations together. Its main value is under go test -race.
func TestFakeClockConcurrentOperationsAreRaceSafe(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()

			for i := 0; i < 25; i++ {
				d := time.Duration(worker+i+1) * time.Nanosecond

				_ = clk.Now()
				_ = clk.Since(start)
				_ = clk.After(d)

				timer := clk.NewTimer(d)
				_ = timer.C()
				_ = timer.Reset(d + time.Nanosecond)
				_ = timer.Stop()

				ticker := clk.NewTicker(d + time.Nanosecond)
				_ = ticker.C()
				ticker.Reset(d + 2*time.Nanosecond)
				ticker.Stop()

				clk.Step(time.Nanosecond)
			}
		}(worker)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	channelassert.RequireSignal(t, done, clockTestTimeout)
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

	_ = channelassert.RequireReceive(t, waiter, clockTestTimeout)
	_ = channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	_ = channelassert.RequireReceive(t, ticker.C(), clockTestTimeout)

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
