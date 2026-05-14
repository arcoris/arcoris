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

var _ Timer = (*realTimer)(nil)

// TestRealTimerChannelIsStable verifies that the adapter exposes the underlying
// timer channel consistently and does not allocate or replace the channel on
// repeated C calls.
func TestRealTimerChannelIsStable(t *testing.T) {
	t.Parallel()

	timer := &realTimer{
		timer: time.NewTimer(time.Hour),
	}
	defer timer.Stop()

	first := timer.C()
	second := timer.C()

	if first == nil {
		t.Fatal("realTimer.C() returned nil channel")
	}

	if first != second {
		t.Fatal("realTimer.C() returned different channels across calls")
	}
}

// TestRealTimerDelivers verifies that realTimer forwards delivery from the
// underlying standard library timer.
func TestRealTimerDelivers(t *testing.T) {
	t.Parallel()

	timer := &realTimer{
		timer: time.NewTimer(realClockTestDelay),
	}
	defer timer.Stop()

	if got := mustReceiveTime(t, timer.C()); got.IsZero() {
		t.Fatal("realTimer delivered zero time")
	}
}

// TestRealTimerStopPreventsLongTimerDelivery verifies the basic active-timer
// Stop contract without depending on precise runtime timing.
func TestRealTimerStopPreventsLongTimerDelivery(t *testing.T) {
	t.Parallel()

	timer := &realTimer{
		timer: time.NewTimer(time.Hour),
	}

	if stopped := timer.Stop(); !stopped {
		t.Fatal("realTimer.Stop() = false for active long-running timer, want true")
	}

	select {
	case got := <-timer.C():
		t.Fatalf("realTimer delivered after Stop: %v", got)
	default:
	}
}

// TestRealTimerStopAfterDeliveryReportsInactive verifies the adapter preserves
// the standard library inactive-timer result after a timer has fired and its
// value has been received.
func TestRealTimerStopAfterDeliveryReportsInactive(t *testing.T) {
	t.Parallel()

	timer := &realTimer{
		timer: time.NewTimer(realClockTestDelay),
	}

	_ = mustReceiveTime(t, timer.C())

	if stopped := timer.Stop(); stopped {
		t.Fatal("realTimer.Stop() after delivered timer = true, want false")
	}
}

// TestRealTimerResetActiveTimer verifies that Reset delegates the active-timer
// lifecycle to the underlying standard library timer and schedules a new
// delivery.
func TestRealTimerResetActiveTimer(t *testing.T) {
	t.Parallel()

	timer := &realTimer{
		timer: time.NewTimer(time.Hour),
	}
	defer timer.Stop()

	if wasActive := timer.Reset(realClockTestDelay); !wasActive {
		t.Fatal("realTimer.Reset() for active timer = false, want true")
	}

	if got := mustReceiveTime(t, timer.C()); got.IsZero() {
		t.Fatal("realTimer delivered zero time after Reset")
	}
}

// TestRealTimerResetStoppedTimerReactivates verifies that Reset can schedule a
// stopped timer again. The returned active-state flag must reflect the previous
// stopped state.
func TestRealTimerResetStoppedTimerReactivates(t *testing.T) {
	t.Parallel()

	timer := &realTimer{
		timer: time.NewTimer(time.Hour),
	}

	if stopped := timer.Stop(); !stopped {
		t.Fatal("realTimer.Stop() = false for active long-running timer, want true")
	}

	if wasActive := timer.Reset(realClockTestDelay); wasActive {
		t.Fatal("realTimer.Reset() after Stop = true, want false")
	}
	defer timer.Stop()

	if got := mustReceiveTime(t, timer.C()); got.IsZero() {
		t.Fatal("realTimer delivered zero time after Reset from stopped state")
	}
}

// TestRealTimerResetNonPositiveDurationDelivers verifies that the adapter
// preserves standard library behavior for immediate timers.
func TestRealTimerResetNonPositiveDurationDelivers(t *testing.T) {
	t.Parallel()

	timer := &realTimer{
		timer: time.NewTimer(time.Hour),
	}
	defer timer.Stop()

	_ = timer.Reset(0)

	if got := mustReceiveTime(t, timer.C()); got.IsZero() {
		t.Fatal("realTimer delivered zero time after Reset(0)")
	}
}
