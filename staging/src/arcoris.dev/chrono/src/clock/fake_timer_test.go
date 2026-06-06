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

var _ Timer = (*fakeTimer)(nil)

// TestFakeTimerChannelIsStable verifies that C exposes the same delivery channel
// throughout the timer lifecycle.
func TestFakeTimerChannelIsStable(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	timer := clk.NewTimer(time.Hour)

	first := timer.C()
	second := timer.C()

	if first == nil {
		t.Fatal("fakeTimer.C() returned nil channel")
	}

	if first != second {
		t.Fatal("fakeTimer.C() returned different channels across calls")
	}
}

// TestFakeTimerChannelIsStableAcrossLifecycle verifies that Stop and Reset do
// not replace the timer channel.
func TestFakeTimerChannelIsStableAcrossLifecycle(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	timer := clk.NewTimer(time.Hour)
	ch := timer.C()

	_ = timer.Stop()
	_ = timer.Reset(time.Hour)
	clk.Step(time.Hour)
	_ = channelassert.RequireReceive(t, ch, clockTestTimeout)

	if timer.C() != ch {
		t.Fatal("fakeTimer.C() changed after Stop/Reset/fire lifecycle")
	}
}

// TestFakeTimerDoesNotFireBeforeDeadline verifies one-shot timer deadline
// behavior before fake time reaches the timer deadline.
func TestFakeTimerDoesNotFireBeforeDeadline(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	timer := clk.NewTimer(10 * time.Second)

	clk.Step(9 * time.Second)

	channelassert.RequireNoReceive(t, timer.C())
}

// TestFakeTimerFiresWhenDeadlineIsReached verifies that a fake timer delivers
// once when fake time reaches its deadline.
func TestFakeTimerFiresWhenDeadlineIsReached(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(10 * time.Second)

	clk.Step(10 * time.Second)

	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "timer delivery", got, start.Add(10*time.Second))

	channelassert.RequireNoReceive(t, timer.C())
}

// TestFakeTimerFiresWhenDeadlineIsPassed verifies that fake timers fire when
// fake time advances beyond their deadline.
func TestFakeTimerFiresWhenDeadlineIsPassed(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(10 * time.Second)

	clk.Step(30 * time.Second)

	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "timer delivery", got, start.Add(30*time.Second))

	channelassert.RequireNoReceive(t, timer.C())
}

// TestFakeTimerNonPositiveDurationIsImmediatelyReady verifies that fake timers
// preserve immediate one-shot semantics for zero and negative durations.
func TestFakeTimerNonPositiveDurationIsImmediatelyReady(t *testing.T) {
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
			timer := clk.NewTimer(tc.d)

			got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
			mustEqualTime(t, "timer immediate delivery", got, start)

			channelassert.RequireNoReceive(t, timer.C())
		})
	}
}

// TestFakeTimerStopPreventsDelivery verifies that Stop removes an active timer
// from the fake clock registry.
func TestFakeTimerStopPreventsDelivery(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	timer := clk.NewTimer(10 * time.Second)

	if stopped := timer.Stop(); !stopped {
		t.Fatal("fakeTimer.Stop() = false for active timer, want true")
	}

	clk.Step(10 * time.Second)

	channelassert.RequireNoReceive(t, timer.C())

	if stopped := timer.Stop(); stopped {
		t.Fatal("second fakeTimer.Stop() = true, want false")
	}
}

// TestFakeTimerStopDoesNotCloseChannel verifies that Stop is not a channel
// lifecycle signal.
func TestFakeTimerStopDoesNotCloseChannel(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	timer := clk.NewTimer(time.Hour)

	_ = timer.Stop()

	select {
	case _, ok := <-timer.C():
		if !ok {
			t.Fatal("fakeTimer.Stop closed the timer channel")
		}
	default:
	}
}

// TestFakeTimerStopAfterFireReportsInactive verifies that Stop reports false
// after the timer has already fired.
func TestFakeTimerStopAfterFireReportsInactive(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	timer := clk.NewTimer(5 * time.Second)

	clk.Step(5 * time.Second)
	_ = channelassert.RequireReceive(t, timer.C(), clockTestTimeout)

	if stopped := timer.Stop(); stopped {
		t.Fatal("fakeTimer.Stop() after fired timer = true, want false")
	}
}

// TestFakeTimerResetActiveTimerMovesDeadline verifies that Reset on an active
// timer returns true and schedules the timer relative to current fake time.
func TestFakeTimerResetActiveTimerMovesDeadline(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(10 * time.Second)

	if wasActive := timer.Reset(20 * time.Second); !wasActive {
		t.Fatal("fakeTimer.Reset() for active timer = false, want true")
	}

	clk.Step(19 * time.Second)
	channelassert.RequireNoReceive(t, timer.C())

	clk.Step(time.Second)

	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "timer delivery after Reset", got, start.Add(20*time.Second))
}

// TestFakeTimerResetActiveTimerLaterAfterPartialAdvance verifies that Reset
// reschedules an active timer relative to the current fake time, not relative to
// the timer's original creation time.
func TestFakeTimerResetActiveTimerLaterAfterPartialAdvance(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(10 * time.Second)

	clk.Step(4 * time.Second)

	if wasActive := timer.Reset(20 * time.Second); !wasActive {
		t.Fatal("fakeTimer.Reset() for active timer = false, want true")
	}

	clk.Step(19 * time.Second)
	channelassert.RequireNoReceive(t, timer.C())

	clk.Step(time.Second)

	want := start.Add(24 * time.Second)
	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "timer delivery after later Reset", got, want)
}

// TestFakeTimerResetActiveTimerEarlierAfterPartialAdvance covers the opposite
// reschedule direction: an active timer can be moved closer to the current fake
// time after fake time has already advanced.
func TestFakeTimerResetActiveTimerEarlierAfterPartialAdvance(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(30 * time.Second)

	clk.Step(10 * time.Second)

	if wasActive := timer.Reset(5 * time.Second); !wasActive {
		t.Fatal("fakeTimer.Reset() for active timer = false, want true")
	}

	clk.Step(4 * time.Second)
	channelassert.RequireNoReceive(t, timer.C())

	clk.Step(time.Second)

	want := start.Add(15 * time.Second)
	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "timer delivery after earlier Reset", got, want)
}

// TestFakeTimerResetStoppedTimerReactivates verifies that Reset can reuse a
// stopped timer and that the return value reports the previous inactive state.
func TestFakeTimerResetStoppedTimerReactivates(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(10 * time.Second)

	if stopped := timer.Stop(); !stopped {
		t.Fatal("fakeTimer.Stop() = false for active timer, want true")
	}

	if wasActive := timer.Reset(5 * time.Second); wasActive {
		t.Fatal("fakeTimer.Reset() after Stop = true, want false")
	}

	clk.Step(5 * time.Second)

	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "timer delivery after Reset from stopped state", got, start.Add(5*time.Second))
}

// TestFakeTimerResetFiredTimerReactivates verifies timer reuse after an earlier
// delivery has already been consumed.
func TestFakeTimerResetFiredTimerReactivates(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(5 * time.Second)

	clk.Step(5 * time.Second)
	_ = channelassert.RequireReceive(t, timer.C(), clockTestTimeout)

	if wasActive := timer.Reset(3 * time.Second); wasActive {
		t.Fatal("fakeTimer.Reset() after fire = true, want false")
	}

	clk.Step(3 * time.Second)

	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "timer delivery after Reset from fired state", got, start.Add(8*time.Second))
}

// TestFakeTimerResetNonPositiveDurationDeliversImmediately verifies immediate
// Reset behavior for zero and negative durations.
func TestFakeTimerResetNonPositiveDurationDeliversImmediately(t *testing.T) {
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
			timer := clk.NewTimer(time.Hour)

			if wasActive := timer.Reset(tc.d); !wasActive {
				t.Fatal("fakeTimer.Reset(non-positive) for active timer = false, want true")
			}

			got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
			mustEqualTime(t, "timer immediate delivery after Reset", got, start)

			channelassert.RequireNoReceive(t, timer.C())
		})
	}
}

// TestFakeTimerResetZeroDeliversImmediatelyWithoutStep documents the direct
// immediate-delivery path used by Reset(0). It does not require a following
// Step(0) to make the timer ready.
func TestFakeTimerResetZeroDeliversImmediatelyWithoutStep(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(time.Hour)

	if wasActive := timer.Reset(0); !wasActive {
		t.Fatal("fakeTimer.Reset(0) for active timer = false, want true")
	}

	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "timer delivery after Reset(0)", got, start)
}

// TestFakeTimerResetDropsImmediateDeliveryWhenChannelIsFull documents the
// non-blocking fake timer delivery policy. Reset(0) must not block when the
// timer channel still contains a stale unread value; the new immediate delivery
// is dropped.
func TestFakeTimerResetDropsImmediateDeliveryWhenChannelIsFull(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(5 * time.Second)

	clk.Step(5 * time.Second)

	done := make(chan struct{})
	go func() {
		timer.Reset(0)
		close(done)
	}()
	channelassert.RequireSignal(t, done, clockTestTimeout)

	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "stale timer delivery", got, start.Add(5*time.Second))

	channelassert.RequireNoReceive(t, timer.C())
}

// TestFakeTimerStepDropsDeliveryWhenChannelIsFull verifies that fake-time
// advancement never blocks on an unread timer value.
func TestFakeTimerStepDropsDeliveryWhenChannelIsFull(t *testing.T) {
	t.Parallel()

	start := fakeClockTestTime()
	clk := NewFakeClock(start)
	timer := clk.NewTimer(5 * time.Second)

	clk.Step(5 * time.Second)

	if wasActive := timer.Reset(5 * time.Second); wasActive {
		t.Fatal("fakeTimer.Reset() after fire = true, want false")
	}

	done := make(chan struct{})
	go func() {
		clk.Step(5 * time.Second)
		close(done)
	}()
	channelassert.RequireSignal(t, done, clockTestTimeout)

	got := channelassert.RequireReceive(t, timer.C(), clockTestTimeout)
	mustEqualTime(t, "stale timer delivery", got, start.Add(5*time.Second))

	channelassert.RequireNoReceive(t, timer.C())
}

// TestFakeTimerConcurrentStopResetAndStepIsRaceSafe is a lifecycle race smoke
// test. It intentionally does not assert an ordering between concurrent owner
// operations; the race detector is the useful oracle here.
func TestFakeTimerConcurrentStopResetAndStepIsRaceSafe(t *testing.T) {
	t.Parallel()

	clk := NewFakeClock(fakeClockTestTime())
	timer := clk.NewTimer(time.Hour)

	var wg sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()

			for i := 0; i < 50; i++ {
				d := time.Duration(worker+i+1) * time.Nanosecond

				_ = timer.Stop()
				_ = timer.Reset(d)
				clk.Step(time.Nanosecond)

				select {
				case <-timer.C():
				default:
				}
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
