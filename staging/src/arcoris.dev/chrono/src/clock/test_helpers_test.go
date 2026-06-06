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
	"runtime"
	"testing"
	"time"
)

const (
	// clockTestTimeout is a safety guard for tests that involve goroutines or
	// real runtime timer channels.
	//
	// Clock tests must not depend on real-time sleeping for correctness. This
	// timeout is used only to prevent a broken test or implementation from
	// hanging the test process indefinitely.
	clockTestTimeout = 500 * time.Millisecond
)

func fakeClockTestTime() time.Time {
	return time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
}

func waitUntil(t *testing.T, description string, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(clockTestTimeout)

	for time.Now().Before(deadline) {
		if condition() {
			return
		}

		runtime.Gosched()
	}

	t.Fatalf("condition did not become true before safety timeout %s: %s", clockTestTimeout, description)
}

func mustEqualTime(t *testing.T, name string, got, want time.Time) {
	t.Helper()

	if !got.Equal(want) {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}

// registerDueWaiter installs a waiter that is already due but not yet delivered.
//
// Production callers create waiters through After or Sleep. These tests need a
// pending due entry so Set(current), Step(0), and Pending can prove how the
// registry behaves before delivery runs.
func registerDueWaiter(t *testing.T, clk *FakeClock) <-chan time.Time {
	t.Helper()

	ch := make(chan time.Time, 1)

	clk.mu.Lock()
	defer clk.mu.Unlock()

	clk.ensureWaiterStoreLocked()

	waiter := &fakeWaiter{
		deadline: clk.now,
		sequence: clk.nextSequenceLocked(),
		ch:       ch,
		active:   true,
	}
	clk.waiters[waiter] = struct{}{}

	return ch
}

// forceTimerDue marks a fake timer as due without delivering its channel value.
//
// The helper keeps already-due delivery tests focused on Set(current) and
// Step(0), rather than on NewTimer or Reset immediate-delivery paths.
func forceTimerDue(t *testing.T, clk *FakeClock, timer Timer) {
	t.Helper()

	fake, ok := timer.(*fakeTimer)
	if !ok {
		t.Fatalf("timer has type %T, want *fakeTimer", timer)
	}
	if fake.clock != clk {
		t.Fatal("timer belongs to a different fake clock")
	}

	clk.mu.Lock()
	defer clk.mu.Unlock()

	clk.ensureTimerStoreLocked()

	fake.deadline = clk.now
	fake.sequence = clk.nextSequenceLocked()
	fake.active = true
	clk.timers[fake] = struct{}{}
}

// forceTickerDue marks a fake ticker as due without sending a tick.
//
// Tests use this to exercise the advancement delivery path with a ticker whose
// next tick is exactly the current fake time.
func forceTickerDue(t *testing.T, clk *FakeClock, ticker Ticker) {
	t.Helper()

	fake, ok := ticker.(*fakeTicker)
	if !ok {
		t.Fatalf("ticker has type %T, want *fakeTicker", ticker)
	}
	if fake.clock != clk {
		t.Fatal("ticker belongs to a different fake clock")
	}

	clk.mu.Lock()
	defer clk.mu.Unlock()

	clk.ensureTickerStoreLocked()

	fake.next = clk.now
	fake.sequence = clk.nextSequenceLocked()
	fake.active = true
	clk.tickers[fake] = struct{}{}
}
