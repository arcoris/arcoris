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

			if err := Delay(context.Background(), tc.d); err != nil {
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

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
}

// TestDelayReturnsInterruptedWhenContextCancelledBeforeDelay verifies
// wait-owned cancellation classification before a positive delay starts.
func TestDelayReturnsInterruptedWhenContextCancelledBeforeDelay(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Delay(ctx, time.Hour)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
}

// TestDelayPreservesCancellationCause verifies that context cancellation causes
// remain visible through the wait-owned interruption wrapper.
func TestDelayPreservesCancellationCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("shutdown requested")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(cause)

	err := Delay(ctx, time.Hour)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
	mustMatch(t, err, cause)
}

// TestDelayReturnsTimeoutWhenContextDeadlineExceededBeforeDelay verifies
// wait-owned timeout classification for an already-expired context deadline.
func TestDelayReturnsTimeoutWhenContextDeadlineExceededBeforeDelay(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	err := Delay(ctx, time.Hour)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, context.DeadlineExceeded)
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

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, context.DeadlineExceeded)
	mustMatch(t, err, cause)
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

	err := mustReceiveError(t, errCh)
	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
	mustMatch(t, err, cause)
}

// TestDelayPanicsOnNilContext verifies invalid context validation at the public
// delay boundary.
func TestDelayPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilContext, func() {
		_ = Delay(nil, time.Second)
	})
}
