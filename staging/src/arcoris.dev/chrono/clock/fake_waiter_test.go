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

// TestFakeClockAfterDoesNotDeliverBeforeDeadline verifies one-shot waiter
// deadline behavior.
func TestFakeClockAfterDoesNotDeliverBeforeDeadline(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	ch := clk.After(10 * time.Second)

	clk.Step(9 * time.Second)
	mustNotReceiveTime(t, ch)
}

// TestFakeClockAfterDeliversWhenDeadlineIsReached verifies that After releases
// exactly when fake time reaches its deadline.
func TestFakeClockAfterDeliversWhenDeadlineIsReached(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	ch := clk.After(10 * time.Second)

	clk.Step(10 * time.Second)

	mustEqualTime(t, "After delivery", mustReceiveTime(t, ch), start.Add(10*time.Second))
	mustNotReceiveTime(t, ch)
}

// TestFakeClockAfterDeliversWhenDeadlineIsPassed verifies that waiters are
// released when fake time advances beyond their deadline.
func TestFakeClockAfterDeliversWhenDeadlineIsPassed(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	ch := clk.After(10 * time.Second)

	clk.Step(30 * time.Second)

	mustEqualTime(t, "After delivery", mustReceiveTime(t, ch), start.Add(30*time.Second))
	mustNotReceiveTime(t, ch)
}

// TestFakeClockAfterNonPositiveDurationIsImmediatelyReady verifies immediate
// waiter semantics for zero and negative durations.
func TestFakeClockAfterNonPositiveDurationIsImmediatelyReady(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()

	tests := []struct {
		name string
		d    time.Duration
	}{
		{
			name: "zero",
			d:    0,
		},
		{
			name: "negative",
			d:    -time.Second,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clk := NewFakeClock(start)

			ch := clk.After(tc.d)

			mustEqualTime(t, "After immediate delivery", mustReceiveTime(t, ch), start)
			mustNotReceiveTime(t, ch)
		})
	}
}

// TestFakeClockSleepBlocksUntilStep verifies that Sleep is controlled by fake
// time and does not observe real runtime time.
func TestFakeClockSleepBlocksUntilStep(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	done := make(chan struct{})

	go func() {
		clk.Sleep(10 * time.Second)
		close(done)
	}()

	waitUntil(t, "Sleep waiter is registered", clk.HasWaiters)

	mustNotReceiveSignal(t, done)

	clk.Step(9 * time.Second)
	mustNotReceiveSignal(t, done)

	clk.Step(time.Second)
	mustReceiveSignal(t, done)
}

// TestFakeClockSleepNonPositiveDurationReturns verifies immediate Sleep
// behavior for zero and negative durations.
func TestFakeClockSleepNonPositiveDurationReturns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    time.Duration
	}{
		{
			name: "zero",
			d:    0,
		},
		{
			name: "negative",
			d:    -time.Second,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clk := NewFakeClock(fakeClockTestTime())
			done := make(chan struct{})

			go func() {
				clk.Sleep(tc.d)
				close(done)
			}()

			mustReceiveSignal(t, done)
		})
	}
}

// TestFakeClockHasWaitersReportsPendingAfterAndSleep verifies that HasWaiters is
// a coordination helper for pending fake-time waiters only.
func TestFakeClockHasWaitersReportsPendingAfterAndSleep(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	if clk.HasWaiters() {
		t.Fatal("FakeClock.HasWaiters() = true before any waiters, want false")
	}

	ch := clk.After(5 * time.Second)

	if !clk.HasWaiters() {
		t.Fatal("FakeClock.HasWaiters() = false after After registration, want true")
	}

	clk.Step(5 * time.Second)
	_ = mustReceiveTime(t, ch)

	if clk.HasWaiters() {
		t.Fatal("FakeClock.HasWaiters() = true after waiter delivery, want false")
	}
}

// TestFakeClockHasWaitersDoesNotReportTimersOrTickers verifies that timer and
// ticker registries are separate from waiter tracking.
func TestFakeClockHasWaitersDoesNotReportTimersOrTickers(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	timer := clk.NewTimer(time.Hour)
	ticker := clk.NewTicker(time.Hour)
	defer timer.Stop()
	defer ticker.Stop()

	if clk.HasWaiters() {
		t.Fatal("FakeClock.HasWaiters() = true with only timer/ticker registered, want false")
	}
}
