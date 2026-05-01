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

var _ Timer = (*fakeTimer)(nil)

// TestFakeTimerChannelIsStable verifies that C exposes the same delivery channel
// throughout the timer lifecycle.
func TestFakeTimerChannelIsStable(t *testing.T) {
	t.Parallel()

	clock := NewFakeClock(fakeClockTestTime())
	timer := clock.NewTimer(time.Hour)

	first := timer.C()
	second := timer.C()

	if first == nil {
		t.Fatal("fakeTimer.C() returned nil channel")
	}

	if first != second {
		t.Fatal("fakeTimer.C() returned different channels across calls")
	}
}

// TestFakeTimerDoesNotFireBeforeDeadline verifies one-shot timer deadline
// behavior before fake time reaches the timer deadline.
func TestFakeTimerDoesNotFireBeforeDeadline(t *testing.T) {
	t.Parallel()

	clock := NewFakeClock(fakeClockTestTime())
	timer := clock.NewTimer(10 * time.Second)

	clock.Step(9 * time.Second)

	mustNotReceiveTime(t, timer.C())
}

// TestFakeTimerFiresWhenDeadlineIsReached verifies that a fake timer delivers
// once when fake time reaches its deadline.
func TestFakeTimerFiresWhenDeadlineIsReached(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clock := NewFakeClock(start)
	timer := clock.NewTimer(10 * time.Second)

	clock.Step(10 * time.Second)

	mustEqualTime(t, "timer delivery", mustReceiveTime(t, timer.C()), start.Add(10*time.Second))
	mustNotReceiveTime(t, timer.C())
}

// TestFakeTimerFiresWhenDeadlineIsPassed verifies that fake timers fire when
// fake time advances beyond their deadline.
func TestFakeTimerFiresWhenDeadlineIsPassed(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clock := NewFakeClock(start)
	timer := clock.NewTimer(10 * time.Second)

	clock.Step(30 * time.Second)

	mustEqualTime(t, "timer delivery", mustReceiveTime(t, timer.C()), start.Add(30*time.Second))
	mustNotReceiveTime(t, timer.C())
}

// TestFakeTimerNonPositiveDurationIsImmediatelyReady verifies that fake timers
// preserve immediate one-shot semantics for zero and negative durations.
func TestFakeTimerNonPositiveDurationIsImmediatelyReady(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "zero",
			duration: 0,
		},
		{
			name:     "negative",
			duration: -time.Second,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clock := NewFakeClock(start)
			timer := clock.NewTimer(tt.duration)

			mustEqualTime(t, "timer immediate delivery", mustReceiveTime(t, timer.C()), start)
			mustNotReceiveTime(t, timer.C())
		})
	}
}

// TestFakeTimerStopPreventsDelivery verifies that Stop removes an active timer
// from the fake clock registry.
func TestFakeTimerStopPreventsDelivery(t *testing.T) {
	t.Parallel()

	clock := NewFakeClock(fakeClockTestTime())
	timer := clock.NewTimer(10 * time.Second)

	if stopped := timer.Stop(); !stopped {
		t.Fatal("fakeTimer.Stop() = false for active timer, want true")
	}

	clock.Step(10 * time.Second)

	mustNotReceiveTime(t, timer.C())

	if stopped := timer.Stop(); stopped {
		t.Fatal("second fakeTimer.Stop() = true, want false")
	}
}

// TestFakeTimerStopAfterFireReportsInactive verifies that Stop reports false
// after the timer has already fired.
func TestFakeTimerStopAfterFireReportsInactive(t *testing.T) {
	t.Parallel()

	clock := NewFakeClock(fakeClockTestTime())
	timer := clock.NewTimer(5 * time.Second)

	clock.Step(5 * time.Second)
	_ = mustReceiveTime(t, timer.C())

	if stopped := timer.Stop(); stopped {
		t.Fatal("fakeTimer.Stop() after fired timer = true, want false")
	}
}

// TestFakeTimerResetActiveTimerMovesDeadline verifies that Reset on an active
// timer returns true and schedules the timer relative to current fake time.
func TestFakeTimerResetActiveTimerMovesDeadline(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clock := NewFakeClock(start)
	timer := clock.NewTimer(10 * time.Second)

	if wasActive := timer.Reset(20 * time.Second); !wasActive {
		t.Fatal("fakeTimer.Reset() for active timer = false, want true")
	}

	clock.Step(19 * time.Second)
	mustNotReceiveTime(t, timer.C())

	clock.Step(time.Second)
	mustEqualTime(t, "timer delivery after Reset", mustReceiveTime(t, timer.C()), start.Add(20*time.Second))
}

// TestFakeTimerResetStoppedTimerReactivates verifies that Reset can reuse a
// stopped timer and that the return value reports the previous inactive state.
func TestFakeTimerResetStoppedTimerReactivates(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clock := NewFakeClock(start)
	timer := clock.NewTimer(10 * time.Second)

	if stopped := timer.Stop(); !stopped {
		t.Fatal("fakeTimer.Stop() = false for active timer, want true")
	}

	if wasActive := timer.Reset(5 * time.Second); wasActive {
		t.Fatal("fakeTimer.Reset() after Stop = true, want false")
	}

	clock.Step(5 * time.Second)

	mustEqualTime(t, "timer delivery after Reset from stopped state", mustReceiveTime(t, timer.C()), start.Add(5*time.Second))
}

// TestFakeTimerResetFiredTimerReactivates verifies timer reuse after an earlier
// delivery has already been consumed.
func TestFakeTimerResetFiredTimerReactivates(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clock := NewFakeClock(start)
	timer := clock.NewTimer(5 * time.Second)

	clock.Step(5 * time.Second)
	_ = mustReceiveTime(t, timer.C())

	if wasActive := timer.Reset(3 * time.Second); wasActive {
		t.Fatal("fakeTimer.Reset() after fire = true, want false")
	}

	clock.Step(3 * time.Second)

	mustEqualTime(t, "timer delivery after Reset from fired state", mustReceiveTime(t, timer.C()), start.Add(8*time.Second))
}

// TestFakeTimerResetNonPositiveDurationDeliversImmediately verifies immediate
// Reset behavior for zero and negative durations.
func TestFakeTimerResetNonPositiveDurationDeliversImmediately(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "zero",
			duration: 0,
		},
		{
			name:     "negative",
			duration: -time.Second,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clock := NewFakeClock(start)
			timer := clock.NewTimer(time.Hour)

			if wasActive := timer.Reset(tt.duration); !wasActive {
				t.Fatal("fakeTimer.Reset(non-positive) for active timer = false, want true")
			}

			mustEqualTime(t, "timer immediate delivery after Reset", mustReceiveTime(t, timer.C()), start)
			mustNotReceiveTime(t, timer.C())
		})
	}
}

// TestFakeTimerResetDropsImmediateDeliveryWhenChannelIsFull documents the
// non-blocking fake timer delivery policy. Reset(0) must not block when the
// timer channel still contains a stale unread value; the new immediate delivery
// is dropped.
func TestFakeTimerResetDropsImmediateDeliveryWhenChannelIsFull(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clock := NewFakeClock(start)
	timer := clock.NewTimer(5 * time.Second)

	clock.Step(5 * time.Second)

	done := make(chan struct{})
	go func() {
		timer.Reset(0)
		close(done)
	}()
	mustReceiveSignal(t, done)

	mustEqualTime(t, "stale timer delivery", mustReceiveTime(t, timer.C()), start.Add(5*time.Second))
	mustNotReceiveTime(t, timer.C())
}
