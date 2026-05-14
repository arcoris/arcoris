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

var _ Ticker = (*realTicker)(nil)

// TestRealTickerChannelIsStable verifies that the adapter exposes the underlying
// ticker channel consistently and does not allocate or replace the channel on
// repeated C calls.
func TestRealTickerChannelIsStable(t *testing.T) {
	t.Parallel()

	ticker := &realTicker{
		ticker: time.NewTicker(time.Hour),
	}
	defer ticker.Stop()

	first := ticker.C()
	second := ticker.C()

	if first == nil {
		t.Fatal("realTicker.C() returned nil channel")
	}

	if first != second {
		t.Fatal("realTicker.C() returned different channels across calls")
	}
}

// TestRealTickerDelivers verifies that realTicker forwards ticks from the
// underlying standard library ticker.
func TestRealTickerDelivers(t *testing.T) {
	t.Parallel()

	ticker := &realTicker{
		ticker: time.NewTicker(realClockTestDelay),
	}
	defer ticker.Stop()

	if got := mustReceiveTime(t, ticker.C()); got.IsZero() {
		t.Fatal("realTicker delivered zero time")
	}
}

// TestRealTickerStopPreventsLongTickerDelivery verifies that Stop turns off a
// ticker that has not yet reached its first tick.
func TestRealTickerStopPreventsLongTickerDelivery(t *testing.T) {
	t.Parallel()

	ticker := &realTicker{
		ticker: time.NewTicker(time.Hour),
	}

	ticker.Stop()

	select {
	case got := <-ticker.C():
		t.Fatalf("realTicker delivered after Stop: %v", got)
	default:
	}
}

// TestRealTickerStopIsIdempotent verifies that repeated Stop calls are safe.
// This mirrors the lifecycle expectation of the standard library ticker.
func TestRealTickerStopIsIdempotent(t *testing.T) {
	t.Parallel()

	ticker := &realTicker{
		ticker: time.NewTicker(time.Hour),
	}

	ticker.Stop()
	ticker.Stop()
}

// TestRealTickerResetChangesPeriod verifies that Reset delegates to the
// underlying standard library ticker and can move a long ticker to a short
// period.
func TestRealTickerResetChangesPeriod(t *testing.T) {
	t.Parallel()

	ticker := &realTicker{
		ticker: time.NewTicker(time.Hour),
	}
	defer ticker.Stop()

	ticker.Reset(realClockTestDelay)

	if got := mustReceiveTime(t, ticker.C()); got.IsZero() {
		t.Fatal("realTicker delivered zero time after Reset")
	}
}

// TestRealTickerResetPanicsForNonPositiveDuration verifies that the adapter
// intentionally preserves time.Ticker.Reset behavior instead of translating the
// panic.
func TestRealTickerResetPanicsForNonPositiveDuration(t *testing.T) {
	t.Parallel()

	ticker := &realTicker{
		ticker: time.NewTicker(time.Hour),
	}
	defer ticker.Stop()

	defer func() {
		if got := recover(); got == nil {
			t.Fatal("realTicker.Reset(0) did not panic")
		}
	}()

	ticker.Reset(0)
}
