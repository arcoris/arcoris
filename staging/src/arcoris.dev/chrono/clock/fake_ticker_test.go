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

var _ Ticker = (*fakeTicker)(nil)

// TestFakeTickerChannelIsStable verifies that C exposes the same delivery
// channel throughout the ticker lifecycle.
func TestFakeTickerChannelIsStable(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	ticker := clk.NewTicker(time.Hour)
	defer ticker.Stop()

	first := ticker.C()
	second := ticker.C()

	if first == nil {
		t.Fatal("fakeTicker.C() returned nil channel")
	}

	if first != second {
		t.Fatal("fakeTicker.C() returned different channels across calls")
	}
}

// TestFakeTickerPanicsForNonPositiveInterval verifies that FakeClock.NewTicker
// matches the standard library ticker contract for invalid intervals.
func TestFakeTickerPanicsForNonPositiveInterval(t *testing.T) {
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

			mustPanicWithValue(t, errFakeTickerNonPositiveInterval, func() {
				_ = clk.NewTicker(tc.d)
			})
		})
	}
}

// TestFakeTickerDoesNotTickBeforeInterval verifies that tickers are not delivered
// before fake time reaches the first interval.
func TestFakeTickerDoesNotTickBeforeInterval(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	ticker := clk.NewTicker(10 * time.Second)
	defer ticker.Stop()

	clk.Step(9 * time.Second)

	mustNotReceiveTime(t, ticker.C())
}

// TestFakeTickerTicksWhenIntervalIsReached verifies the first periodic tick.
func TestFakeTickerTicksWhenIntervalIsReached(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	ticker := clk.NewTicker(10 * time.Second)
	defer ticker.Stop()

	clk.Step(10 * time.Second)

	mustEqualTime(t, "ticker delivery", mustReceiveTime(t, ticker.C()), start.Add(10*time.Second))
	mustNotReceiveTime(t, ticker.C())
}

// TestFakeTickerTicksRepeatedlyAcrossAdvanceCalls verifies repeated periodic
// behavior when fake time advances in separate operations.
func TestFakeTickerTicksRepeatedlyAcrossAdvanceCalls(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	ticker := clk.NewTicker(5 * time.Second)
	defer ticker.Stop()

	clk.Step(5 * time.Second)
	mustEqualTime(t, "first ticker delivery", mustReceiveTime(t, ticker.C()), start.Add(5*time.Second))

	clk.Step(5 * time.Second)
	mustEqualTime(t, "second ticker delivery", mustReceiveTime(t, ticker.C()), start.Add(10*time.Second))

	clk.Step(5 * time.Second)
	mustEqualTime(t, "third ticker delivery", mustReceiveTime(t, ticker.C()), start.Add(15*time.Second))
}

// TestFakeTickerDeliversAtMostOneTickPerAdvance verifies the documented fake
// ticker policy. Large fake-time jumps must not produce an unbounded burst of
// missed ticks.
func TestFakeTickerDeliversAtMostOneTickPerAdvance(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	ticker := clk.NewTicker(5 * time.Second)
	defer ticker.Stop()

	clk.Step(30 * time.Second)

	mustEqualTime(t, "ticker delivery after large Step", mustReceiveTime(t, ticker.C()), start.Add(30*time.Second))
	mustNotReceiveTime(t, ticker.C())

	clk.Step(5 * time.Second)
	mustEqualTime(t, "ticker delivery after next Step", mustReceiveTime(t, ticker.C()), start.Add(35*time.Second))
}

// TestFakeTickerStopPreventsFutureTicks verifies that Stop removes the ticker
// from the fake clock registry.
func TestFakeTickerStopPreventsFutureTicks(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	ticker := clk.NewTicker(5 * time.Second)

	ticker.Stop()

	clk.Step(5 * time.Second)
	mustNotReceiveTime(t, ticker.C())

	ticker.Stop()
}

// TestFakeTickerResetChangesInterval verifies that Reset schedules the next tick
// relative to the current fake time.
func TestFakeTickerResetChangesInterval(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	ticker := clk.NewTicker(10 * time.Second)
	defer ticker.Stop()

	clk.Step(4 * time.Second)

	ticker.Reset(3 * time.Second)

	clk.Step(2 * time.Second)
	mustNotReceiveTime(t, ticker.C())

	clk.Step(time.Second)
	mustEqualTime(t, "ticker delivery after Reset", mustReceiveTime(t, ticker.C()), start.Add(7*time.Second))
}

// TestFakeTickerResetStoppedTickerReactivates verifies that Reset can reactivate
// a stopped fake ticker.
func TestFakeTickerResetStoppedTickerReactivates(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	ticker := clk.NewTicker(10 * time.Second)

	ticker.Stop()
	ticker.Reset(3 * time.Second)
	defer ticker.Stop()

	clk.Step(3 * time.Second)

	mustEqualTime(t, "ticker delivery after Reset from stopped state", mustReceiveTime(t, ticker.C()), start.Add(3*time.Second))
}

// TestFakeTickerResetDropsDeliveryWhenChannelIsFull documents that Reset does
// not drain stale ticks. A later tick is dropped when the ticker channel is
// already full.
func TestFakeTickerResetDropsDeliveryWhenChannelIsFull(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	ticker := clk.NewTicker(5 * time.Second)
	defer ticker.Stop()

	clk.Step(5 * time.Second)

	ticker.Reset(3 * time.Second)
	clk.Step(3 * time.Second)

	mustEqualTime(t, "stale ticker delivery", mustReceiveTime(t, ticker.C()), start.Add(5*time.Second))
	mustNotReceiveTime(t, ticker.C())
}

// TestFakeTickerResetPanicsForNonPositiveInterval verifies Reset uses the same
// interval validation rule as NewTicker.
func TestFakeTickerResetPanicsForNonPositiveInterval(t *testing.T) {
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
			ticker := clk.NewTicker(time.Hour)
			defer ticker.Stop()

			mustPanicWithValue(t, errFakeTickerNonPositiveInterval, func() {
				ticker.Reset(tc.d)
			})
		})
	}
}

// TestFakeTickerDropsDeliveryWhenChannelIsFull verifies non-blocking delivery.
// A full ticker channel must not block fake-clock advancement.
func TestFakeTickerDropsDeliveryWhenChannelIsFull(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	ticker := clk.NewTicker(5 * time.Second)
	defer ticker.Stop()

	clk.Step(5 * time.Second)
	mustEqualTime(t, "first ticker delivery", mustReceiveTime(t, ticker.C()), start.Add(5*time.Second))

	clk.Step(5 * time.Second)

	done := make(chan struct{})
	go func() {
		clk.Step(5 * time.Second)
		close(done)
	}()
	mustReceiveSignal(t, done)

	mustEqualTime(t, "second ticker delivery", mustReceiveTime(t, ticker.C()), start.Add(10*time.Second))
	mustNotReceiveTime(t, ticker.C())
}
