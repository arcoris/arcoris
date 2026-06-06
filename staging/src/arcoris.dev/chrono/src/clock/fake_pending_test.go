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
	"testing"
	"time"
)

func TestFakeClockPendingReportsZeroForNewAndZeroValueClocks(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	requirePending(t, clk.Pending(), Pending{})

	var zero FakeClock
	requirePending(t, zero.Pending(), Pending{})
}

func TestFakeClockPendingTracksAfterAndSleepWaiters(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	// After registers synchronously. Sleep registers from its goroutine, so the
	// test waits for Pending to observe it before advancing fake time.
	after := clk.After(time.Hour)
	requirePending(t, clk.Pending(), Pending{Waiters: 1})

	done := make(chan struct{})
	go func() {
		clk.Sleep(time.Hour)
		close(done)
	}()
	waitUntil(t, "Sleep waiter is registered", func() bool {
		return clk.Pending().Waiters == 2
	})

	clk.Step(time.Hour)

	_ = channelassert.RequireReceive(t, after, clockTestTimeout)
	channelassert.RequireSignal(t, done, clockTestTimeout)
	requirePending(t, clk.Pending(), Pending{})
}

func TestFakeClockPendingTracksTimerLifecycle(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	timer := clk.NewTimer(time.Hour)
	requirePending(t, clk.Pending(), Pending{Timers: 1})

	if !timer.Stop() {
		t.Fatal("timer.Stop() = false, want true")
	}
	requirePending(t, clk.Pending(), Pending{})

	timer.Reset(time.Hour)
	requirePending(t, clk.Pending(), Pending{Timers: 1})

	clk.Step(time.Hour)
	_ = channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	requirePending(t, clk.Pending(), Pending{})
}

func TestFakeClockPendingTracksTickerLifecycle(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	ticker := clk.NewTicker(time.Hour)
	requirePending(t, clk.Pending(), Pending{Tickers: 1})

	ticker.Stop()
	requirePending(t, clk.Pending(), Pending{})

	ticker.Reset(time.Hour)
	defer ticker.Stop()
	requirePending(t, clk.Pending(), Pending{Tickers: 1})
}

func TestFakeClockPendingDoesNotMutateOrDeliver(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)

	// Pending must be a read-only inspection. The due waiter should stay queued
	// until Set(current) explicitly runs the delivery path.
	waiter := registerDueWaiter(t, clk)
	requirePending(t, clk.Pending(), Pending{Waiters: 1})
	channelassert.RequireNoReceive(t, waiter)

	clk.Set(start)

	got := channelassert.RequireReceive(t, waiter, clockTestTimeout)
	mustEqualTime(t, "waiter delivery", got, start)

	requirePending(t, clk.Pending(), Pending{})
}

func TestFakeClockPendingDoesNotCountUnreadDeliveredValues(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())

	waiter := clk.After(time.Second)
	timer := clk.NewTimer(time.Second)
	ticker := clk.NewTicker(time.Second)
	defer ticker.Stop()

	requirePending(t, clk.Pending(), Pending{
		Waiters: 1,
		Timers:  1,
		Tickers: 1,
	})

	clk.Step(time.Second)

	requirePending(t, clk.Pending(), Pending{Tickers: 1})

	_ = channelassert.RequireReceive(t, waiter, clockTestTimeout)
	_ = channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	_ = channelassert.RequireReceive(t, ticker.C(), clockTestTimeout)

	requirePending(t, clk.Pending(), Pending{Tickers: 1})

	ticker.Stop()
	requirePending(t, clk.Pending(), Pending{})
}

func requirePending(t *testing.T, got Pending, want Pending) {
	t.Helper()

	if got != want {
		t.Fatalf("Pending() = %#v, want %#v", got, want)
	}
}
