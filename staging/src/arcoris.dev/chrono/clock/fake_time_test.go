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

import (
	"testing"
	"time"
)

// TestFakeClockSinceUsesFakeTime verifies that elapsed duration is computed from
// fake time, not from the real process clock.
func TestFakeClockSinceUsesFakeTime(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	clk.Step(10 * time.Second)

	if got := clk.Since(start); got != 10*time.Second {
		t.Fatalf("FakeClock.Since(start) = %s, want %s", got, 10*time.Second)
	}
}

// TestFakeClockSetMovesToFutureTime verifies exact Set behavior for monotonic
// fake time movement.
func TestFakeClockSetMovesToFutureTime(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	next := start.Add(30 * time.Second)
	clk := NewFakeClock(start)

	clk.Set(next)

	mustEqualTime(t, "FakeClock.Now() after Set", clk.Now(), next)
}

// TestFakeClockSetToCurrentTimeIsAllowed verifies that Set may be used to
// process already-due work without changing fake time.
func TestFakeClockSetToCurrentTimeIsAllowed(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	clk.Set(start)

	mustEqualTime(t, "FakeClock.Now() after Set(current)", clk.Now(), start)
}

// TestFakeClockSetBackwardsPanics verifies the monotonic fake-time invariant.
//
// Backwards movement would make already-fired waiters, timers, and ticker
// deadlines ambiguous, so FakeClock rejects it.
func TestFakeClockSetBackwardsPanics(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	mustPanicWithValue(t, errFakeClockBackwardSet, func() {
		clk.Set(start.Add(-time.Nanosecond))
	})

	mustEqualTime(t, "FakeClock.Now() after failed backward Set", clk.Now(), start)
}

// TestFakeClockStepMovesByDuration verifies that Step advances fake time by an
// exact duration.
func TestFakeClockStepMovesByDuration(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	clk.Step(15 * time.Second)

	mustEqualTime(t, "FakeClock.Now() after Step", clk.Now(), start.Add(15*time.Second))
}

// TestFakeClockStepZeroIsAllowed verifies that Step(0) keeps time stable while
// still exercising the due-delivery path.
func TestFakeClockStepZeroIsAllowed(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	clk.Step(0)

	mustEqualTime(t, "FakeClock.Now() after Step(0)", clk.Now(), start)
}

// TestFakeClockStepNegativePanics verifies that fake time cannot move backwards
// through duration-based advancement.
func TestFakeClockStepNegativePanics(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	mustPanicWithValue(t, errFakeClockNegativeStep, func() {
		clk.Step(-time.Nanosecond)
	})

	mustEqualTime(t, "FakeClock.Now() after failed negative Step", clk.Now(), start)
}

// TestFakeClockPrivateStepNegativePanics verifies the monotonic fake-time
// invariant at the private mutation layer used by package-internal tests.
func TestFakeClockPrivateStepNegativePanics(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	mustPanicWithValue(t, errFakeClockNegativeStep, func() {
		clk.step(-time.Nanosecond)
	})

	mustEqualTime(t, "FakeClock.Now() after failed private negative step", clk.Now(), start)
}

// TestFakeClockStepDeliversDueWaitersTimersAndTickers verifies that time
// advancement coordinates all fake delivery sources through one path.
func TestFakeClockStepDeliversDueWaitersTimersAndTickers(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	waiter := clk.After(10 * time.Second)
	timer := clk.NewTimer(10 * time.Second)
	ticker := clk.NewTicker(10 * time.Second)
	defer ticker.Stop()

	clk.Step(9 * time.Second)

	mustNotReceiveTime(t, waiter)
	mustNotReceiveTime(t, timer.C())
	mustNotReceiveTime(t, ticker.C())

	clk.Step(time.Second)

	want := start.Add(10 * time.Second)

	mustEqualTime(t, "waiter delivery", mustReceiveTime(t, waiter), want)
	mustEqualTime(t, "timer delivery", mustReceiveTime(t, timer.C()), want)
	mustEqualTime(t, "ticker delivery", mustReceiveTime(t, ticker.C()), want)
}

// TestFakeClockCollectsDueDeliveriesInStableOrder verifies the documented fake
// delivery order without relying on receiver scheduling.
//
// This test intentionally calls the package-private step method so it can
// inspect the delivery plan before channel delivery and receiver scheduling.
func TestFakeClockCollectsDueDeliveriesInStableOrder(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	lateWaiter := clk.After(20 * time.Second)
	earlyWaiter := clk.After(10 * time.Second)
	firstSameDeadlineWaiter := clk.After(15 * time.Second)
	secondSameDeadlineWaiter := clk.After(15 * time.Second)

	lateTimer := clk.NewTimer(20 * time.Second)
	earlyTimer := clk.NewTimer(10 * time.Second)
	firstSameDeadlineTimer := clk.NewTimer(15 * time.Second)
	secondSameDeadlineTimer := clk.NewTimer(15 * time.Second)

	lateTicker := clk.NewTicker(20 * time.Second)
	defer lateTicker.Stop()
	earlyTicker := clk.NewTicker(10 * time.Second)
	defer earlyTicker.Stop()
	firstSameDeadlineTicker := clk.NewTicker(15 * time.Second)
	defer firstSameDeadlineTicker.Stop()
	secondSameDeadlineTicker := clk.NewTicker(15 * time.Second)
	defer secondSameDeadlineTicker.Stop()

	deliveries := clk.step(20 * time.Second)

	if len(deliveries.waiters) != 4 {
		t.Fatalf("len(waiter deliveries) = %d, want 4", len(deliveries.waiters))
	}
	if len(deliveries.timers) != 4 {
		t.Fatalf("len(timer deliveries) = %d, want 4", len(deliveries.timers))
	}
	if len(deliveries.tickers) != 4 {
		t.Fatalf("len(ticker deliveries) = %d, want 4", len(deliveries.tickers))
	}

	wantWaiters := []<-chan time.Time{
		earlyWaiter,
		firstSameDeadlineWaiter,
		secondSameDeadlineWaiter,
		lateWaiter,
	}
	for i, want := range wantWaiters {
		if (<-chan time.Time)(deliveries.waiters[i].ch) != want {
			t.Fatalf("waiter delivery %d uses channel %p, want %p", i, deliveries.waiters[i].ch, want)
		}
	}

	wantTimers := []<-chan time.Time{
		earlyTimer.C(),
		firstSameDeadlineTimer.C(),
		secondSameDeadlineTimer.C(),
		lateTimer.C(),
	}
	for i, want := range wantTimers {
		if (<-chan time.Time)(deliveries.timers[i].ch) != want {
			t.Fatalf("timer delivery %d uses channel %p, want %p", i, deliveries.timers[i].ch, want)
		}
	}

	wantTickers := []<-chan time.Time{
		earlyTicker.C(),
		firstSameDeadlineTicker.C(),
		secondSameDeadlineTicker.C(),
		lateTicker.C(),
	}
	for i, want := range wantTickers {
		if (<-chan time.Time)(deliveries.tickers[i].ch) != want {
			t.Fatalf("ticker delivery %d uses channel %p, want %p", i, deliveries.tickers[i].ch, want)
		}
	}

	deliveries.deliver()
}
