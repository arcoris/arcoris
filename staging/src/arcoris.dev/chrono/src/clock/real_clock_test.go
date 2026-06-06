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

const (
	// realClockTestDelay is intentionally small because these tests verify
	// contract-level behavior, not timing precision.
	realClockTestDelay = time.Millisecond
)

// TestRealClockZeroValueIsUsable verifies that RealClock can be used directly
// without construction, configuration, or injected state.
func TestRealClockZeroValueIsUsable(t *testing.T) {
	t.Parallel()

	var clk RealClock

	if got := clk.Now(); got.IsZero() {
		t.Fatal("RealClock.Now() returned zero time")
	}
}

// TestRealClockSinceReturnsNonNegativeElapsedDuration verifies the basic elapsed
// time contract without asserting precise wall-clock timing.
func TestRealClockSinceReturnsNonNegativeElapsedDuration(t *testing.T) {
	t.Parallel()

	var clk RealClock

	start := clk.Now()
	elapsed := clk.Since(start)

	if elapsed < 0 {
		t.Fatalf("RealClock.Since(start) = %s, want non-negative duration", elapsed)
	}
}

// TestRealClockUntilReturnsDeadlineDuration verifies the basic deadline
// calculation contract without asserting precise wall-clock timing.
func TestRealClockUntilReturnsDeadlineDuration(t *testing.T) {
	t.Parallel()

	var clk RealClock

	deadline := clk.Now().Add(time.Hour)
	remaining := clk.Until(deadline)

	if remaining <= 0 {
		t.Fatalf("RealClock.Until(deadline) = %s, want positive duration", remaining)
	}
	if remaining > time.Hour {
		t.Fatalf("RealClock.Until(deadline) = %s, want at most 1h", remaining)
	}
}

// TestRealClockAfterDelivers verifies that RealClock.After is wired to a real
// runtime wait source. The timeout is only a test guard, not part of the clock
// contract.
func TestRealClockAfterDelivers(t *testing.T) {
	t.Parallel()

	var clk RealClock

	if got := channelassert.RequireReceive(t, clk.After(realClockTestDelay), clockTestTimeout); got.IsZero() {
		t.Fatal("RealClock.After delivered zero time")
	}
}

// TestRealClockSleepAcceptsNonPositiveDurations verifies that RealClock preserves
// the standard library sleep behavior for non-positive durations.
func TestRealClockSleepAcceptsNonPositiveDurations(t *testing.T) {
	t.Parallel()

	var clk RealClock

	done := make(chan struct{})
	go func() {
		clk.Sleep(0)
		clk.Sleep(-time.Second)
		close(done)
	}()

	channelassert.RequireSignal(t, done, clockTestTimeout)
}

// TestRealClockNewTimerReturnsUsableTimer verifies that RealClock creates a
// Timer contract backed by a real runtime timer.
func TestRealClockNewTimerReturnsUsableTimer(t *testing.T) {
	t.Parallel()

	var clk RealClock

	timer := clk.NewTimer(realClockTestDelay)
	if timer == nil {
		t.Fatal("RealClock.NewTimer returned nil")
	}
	defer timer.Stop()

	if timer.C() == nil {
		t.Fatal("RealClock.NewTimer returned timer with nil channel")
	}

	if got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout); got.IsZero() {
		t.Fatal("RealClock.NewTimer timer delivered zero time")
	}
}

// TestRealClockNewTickerReturnsUsableTicker verifies that RealClock creates a
// Ticker contract backed by a real runtime ticker.
func TestRealClockNewTickerReturnsUsableTicker(t *testing.T) {
	t.Parallel()

	var clk RealClock

	ticker := clk.NewTicker(realClockTestDelay)
	if ticker == nil {
		t.Fatal("RealClock.NewTicker returned nil")
	}
	defer ticker.Stop()

	if ticker.C() == nil {
		t.Fatal("RealClock.NewTicker returned ticker with nil channel")
	}

	if got := channelassert.RequireReceive(t, ticker.C(), clockTestTimeout); got.IsZero() {
		t.Fatal("RealClock.NewTicker ticker delivered zero time")
	}
}
