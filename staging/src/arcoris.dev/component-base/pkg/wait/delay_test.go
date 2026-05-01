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
	"fmt"
	"testing"
	"time"
)

// TestDelayReturnsNilAfterPositiveDuration verifies successful one-shot delay
// completion for an active context.
func TestDelayReturnsNilAfterPositiveDuration(t *testing.T) {
	t.Parallel()

	if err := Delay(context.Background(), time.Nanosecond); err != nil {
		t.Fatalf("Delay(...) = %v, want nil", err)
	}
}

// TestDelayReturnsNilForNonPositiveDuration verifies that non-positive delays
// are immediate no-ops when the context is still active.
func TestDelayReturnsNilForNonPositiveDuration(t *testing.T) {
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

			if err := Delay(context.Background(), tt.duration); err != nil {
				t.Fatalf("Delay(...) = %v, want nil", err)
			}
		})
	}
}

// TestDelayChecksContextBeforeImmediateDuration verifies that an already-stopped
// context still wins before a non-positive duration is treated as complete.
func TestDelayChecksContextBeforeImmediateDuration(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Delay(ctx, 0)

	mustDelayBeInterrupted(t, err)
	mustDelayNotBeTimedOut(t, err)
	mustDelayMatch(t, err, context.Canceled)
}

// TestDelayReturnsInterruptedWhenContextCancelledBeforeDelay verifies
// wait-owned cancellation classification before a positive delay starts.
func TestDelayReturnsInterruptedWhenContextCancelledBeforeDelay(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Delay(ctx, time.Hour)

	mustDelayBeInterrupted(t, err)
	mustDelayNotBeTimedOut(t, err)
	mustDelayMatch(t, err, context.Canceled)
}

// TestDelayPreservesCancellationCause verifies that context cancellation causes
// remain visible through the wait-owned interruption wrapper.
func TestDelayPreservesCancellationCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("shutdown requested")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(cause)

	err := Delay(ctx, time.Hour)

	mustDelayBeInterrupted(t, err)
	mustDelayNotBeTimedOut(t, err)
	mustDelayMatch(t, err, cause)
}

// TestDelayReturnsTimeoutWhenContextDeadlineExceededBeforeDelay verifies
// wait-owned timeout classification for an already-expired context deadline.
func TestDelayReturnsTimeoutWhenContextDeadlineExceededBeforeDelay(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	err := Delay(ctx, time.Hour)

	mustDelayBeTimedOut(t, err)
	mustDelayBeInterrupted(t, err)
	mustDelayMatch(t, err, context.DeadlineExceeded)
}

// TestDelayPreservesDeadlineCause verifies that timeout causes remain visible
// through the wait-owned timeout wrapper.
func TestDelayPreservesDeadlineCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("wait budget exhausted")
	ctx, cancel := context.WithTimeoutCause(context.Background(), time.Nanosecond, cause)
	defer cancel()

	<-ctx.Done()
	err := Delay(ctx, time.Hour)

	mustDelayBeTimedOut(t, err)
	mustDelayBeInterrupted(t, err)
	mustDelayMatch(t, err, cause)
}

// TestDelayReturnsInterruptedWhenContextCancelledDuringDelay verifies that
// cancellation stops an in-flight delay without waiting for the timer duration.
func TestDelayReturnsInterruptedWhenContextCancelledDuringDelay(t *testing.T) {
	t.Parallel()

	cause := errors.New("stop delay")
	ctx, cancel := context.WithCancelCause(context.Background())
	errCh := make(chan error, 1)
	started := make(chan struct{})

	go func() {
		close(started)
		errCh <- Delay(ctx, time.Hour)
	}()

	<-started
	cancel(cause)

	err := mustReceiveDelayError(t, errCh)
	mustDelayBeInterrupted(t, err)
	mustDelayNotBeTimedOut(t, err)
	mustDelayMatch(t, err, cause)
}

// TestDelayPanicsOnNilContext verifies invalid context validation at the public
// delay boundary.
func TestDelayPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	mustDelayPanicWith(t, errNilContext, func() {
		_ = Delay(nil, time.Second)
	})
}

// mustDelayBeInterrupted fails the test unless err is a wait-owned
// interruption.
func mustDelayBeInterrupted(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("err is nil, want non-nil")
	}
	if !Interrupted(err) {
		t.Fatal("Interrupted(err) = false, want true")
	}
}

// mustDelayNotBeInterrupted fails the test if err is classified as a wait-owned
// interruption.
func mustDelayNotBeInterrupted(t *testing.T, err error) {
	t.Helper()

	if Interrupted(err) {
		t.Fatal("Interrupted(err) = true, want false")
	}
}

// mustDelayBeTimedOut fails the test unless err is a wait-owned timeout.
func mustDelayBeTimedOut(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("err is nil, want non-nil")
	}
	if !TimedOut(err) {
		t.Fatal("TimedOut(err) = false, want true")
	}
}

// mustDelayNotBeTimedOut fails the test if err is classified as a wait-owned
// timeout.
func mustDelayNotBeTimedOut(t *testing.T, err error) {
	t.Helper()

	if TimedOut(err) {
		t.Fatal("TimedOut(err) = true, want false")
	}
}

// mustDelayMatch fails the test unless err matches target through errors.Is.
func mustDelayMatch(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(err, %v) = false, want true", target)
	}
}

// mustReceiveDelayError waits for Delay running in another goroutine and fails
// the test if it does not return promptly.
func mustReceiveDelayError(t *testing.T, ch <-chan error) error {
	t.Helper()

	select {
	case err := <-ch:
		return err
	case <-time.After(time.Second):
		t.Fatal("Delay did not return after context cancellation")
		return nil
	}
}

// mustDelayPanicWith fails the test unless fn panics with want.
func mustDelayPanicWith(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		if got != want {
			t.Fatalf("panic = %s, want %s", fmt.Sprint(got), fmt.Sprint(want))
		}
	}()

	fn()
}
