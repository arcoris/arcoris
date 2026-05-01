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

package wait

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestNewTimerReturnsUsableTimer verifies that NewTimer initializes the runtime
// timer wrapper and exposes a stable delivery channel.
func TestNewTimerReturnsUsableTimer(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)
	defer timer.StopAndDrain()

	first := timer.C()
	second := timer.C()

	if first == nil {
		t.Fatal("Timer.C() = nil, want non-nil channel")
	}
	if first != second {
		t.Fatal("Timer.C() returned different channels, want stable channel")
	}
}

// TestTimerWaitReturnsNilWhenTimerFires verifies successful timer completion.
func TestTimerWaitReturnsNilWhenTimerFires(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Nanosecond)

	if err := timer.Wait(context.Background()); err != nil {
		t.Fatalf("Timer.Wait(...) = %v, want nil", err)
	}
}

// TestTimerWaitReturnsNilForImmediateTimer verifies that non-positive timer
// durations are immediately ready according to NewTimer semantics.
func TestTimerWaitReturnsNilForImmediateTimer(t *testing.T) {
	t.Parallel()

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

			timer := NewTimer(tt.duration)

			if err := timer.Wait(context.Background()); err != nil {
				t.Fatalf("Timer.Wait(...) = %v, want nil", err)
			}
		})
	}
}

// TestTimerWaitReturnsInterruptedWhenContextCancelledBeforeWait verifies that an
// already-cancelled context is classified as a wait-owned interruption.
func TestTimerWaitReturnsInterruptedWhenContextCancelledBeforeWait(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	timer := NewTimer(time.Hour)

	err := timer.Wait(ctx)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
}

// TestTimerWaitPreservesCancellationCause verifies that cancellation causes are
// visible through the wait-owned interruption wrapper returned by Wait.
func TestTimerWaitPreservesCancellationCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("timer owner stopped")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(cause)
	timer := NewTimer(time.Hour)

	err := timer.Wait(ctx)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
	mustMatch(t, err, cause)
}

// TestTimerWaitReturnsTimeoutWhenContextDeadlineExceededBeforeWait verifies
// wait-owned timeout classification for an expired context deadline.
func TestTimerWaitReturnsTimeoutWhenContextDeadlineExceededBeforeWait(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	timer := NewTimer(time.Hour)

	err := timer.Wait(ctx)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, context.DeadlineExceeded)
}

// TestTimerWaitPreservesDeadlineCause verifies that timeout causes remain
// visible through the wait-owned timeout wrapper returned by Wait.
func TestTimerWaitPreservesDeadlineCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("timer wait budget exhausted")
	ctx, cancel := context.WithTimeoutCause(context.Background(), time.Nanosecond, cause)
	defer cancel()
	<-ctx.Done()
	timer := NewTimer(time.Hour)

	err := timer.Wait(ctx)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, context.DeadlineExceeded)
	mustMatch(t, err, cause)
}

// TestTimerWaitReturnsInterruptedWhenContextCancelledDuringWait verifies that
// context cancellation stops an in-flight timer wait.
func TestTimerWaitReturnsInterruptedWhenContextCancelledDuringWait(t *testing.T) {
	t.Parallel()

	cause := errors.New("shutdown")
	ctx, cancel := context.WithCancelCause(context.Background())
	timer := NewTimer(time.Hour)
	errCh := make(chan error, 1)
	started := make(chan struct{})

	go func() {
		close(started)
		errCh <- timer.Wait(ctx)
	}()

	<-started
	cancel(cause)

	err := mustReceiveError(t, errCh)
	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
	mustMatch(t, err, cause)
}

// TestTimerStopPreventsDelivery verifies that Stop prevents an active long timer
// from delivering a value.
func TestTimerStopPreventsDelivery(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)

	if !timer.Stop() {
		t.Fatal("Timer.Stop() = false, want true for active timer")
	}
	mustNotReceiveTimerValue(t, timer.C())
}

// TestTimerStopAndDrainPreventsDelivery verifies the combined stop/drain helper
// for an active timer.
func TestTimerStopAndDrainPreventsDelivery(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)

	if !timer.StopAndDrain() {
		t.Fatal("Timer.StopAndDrain() = false, want true for active timer")
	}
	mustNotReceiveTimerValue(t, timer.C())
}

// TestTimerResetReschedulesStoppedTimer verifies that Reset can reuse a stopped
// timer and make it fire again.
func TestTimerResetReschedulesStoppedTimer(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)
	if !timer.StopAndDrain() {
		t.Fatal("Timer.StopAndDrain() = false, want true for active timer")
	}

	if active := timer.Reset(time.Nanosecond); active {
		t.Fatal("Timer.Reset(...) = true after stopped timer, want false")
	}
	if err := timer.Wait(context.Background()); err != nil {
		t.Fatalf("Timer.Wait(...) after Reset = %v, want nil", err)
	}
}

// TestTimerResetReschedulesActiveTimer verifies that Reset reports an active
// timer and replaces its deadline.
func TestTimerResetReschedulesActiveTimer(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)

	if active := timer.Reset(time.Nanosecond); !active {
		t.Fatal("Timer.Reset(...) = false for active timer, want true")
	}
	if err := timer.Wait(context.Background()); err != nil {
		t.Fatalf("Timer.Wait(...) after active Reset = %v, want nil", err)
	}
}

// TestTimerResetAllowsImmediateDuration verifies that Reset preserves the
// low-level timer rule that non-positive durations are immediately ready.
func TestTimerResetAllowsImmediateDuration(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)
	defer timer.StopAndDrain()

	_ = timer.Reset(0)
	if err := timer.Wait(context.Background()); err != nil {
		t.Fatalf("Timer.Wait(...) after immediate Reset = %v, want nil", err)
	}
}

// TestTimerWaitPanicsOnNilContext verifies context validation at the Timer.Wait
// boundary.
func TestTimerWaitPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	timer := NewTimer(time.Hour)
	defer timer.StopAndDrain()

	mustPanicWith(t, errNilContext, func() {
		_ = timer.Wait(nil)
	})
}

// TestTimerMethodsPanicOnNilReceiver verifies that Timer methods reject nil
// receiver usage consistently.
func TestTimerMethodsPanicOnNilReceiver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(timer *Timer)
	}{
		{
			name: "C",
			call: func(timer *Timer) { _ = timer.C() },
		},
		{
			name: "Wait",
			call: func(timer *Timer) { _ = timer.Wait(context.Background()) },
		},
		{
			name: "Stop",
			call: func(timer *Timer) { _ = timer.Stop() },
		},
		{
			name: "StopAndDrain",
			call: func(timer *Timer) { _ = timer.StopAndDrain() },
		},
		{
			name: "Reset",
			call: func(timer *Timer) { _ = timer.Reset(time.Second) },
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilTimer, func() {
				tt.call(nil)
			})
		})
	}
}

// TestTimerMethodsPanicOnZeroValue verifies that Timer methods reject zero-value
// usage consistently.
func TestTimerMethodsPanicOnZeroValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(timer *Timer)
	}{
		{
			name: "C",
			call: func(timer *Timer) { _ = timer.C() },
		},
		{
			name: "Wait",
			call: func(timer *Timer) { _ = timer.Wait(context.Background()) },
		},
		{
			name: "Stop",
			call: func(timer *Timer) { _ = timer.Stop() },
		},
		{
			name: "StopAndDrain",
			call: func(timer *Timer) { _ = timer.StopAndDrain() },
		},
		{
			name: "Reset",
			call: func(timer *Timer) { _ = timer.Reset(time.Second) },
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var timer Timer
			mustPanicWith(t, errNilTimer, func() {
				tt.call(&timer)
			})
		})
	}
}
