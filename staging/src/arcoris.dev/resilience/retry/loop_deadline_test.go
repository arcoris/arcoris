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

package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
)

func TestRetryExecutionContextDeadlineWouldBeExceeded(t *testing.T) {
	now := retryFutureNow()
	fake := clock.NewFakeClock(now)
	execution := &retryExecution{config: config{clock: fake}}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name string
		ctx  context.Context
		d    time.Duration
		want bool
	}{
		{
			name: "no deadline",
			ctx:  context.Background(),
			d:    time.Second,
			want: false,
		},
		{
			name: "already canceled",
			ctx:  canceled,
			d:    time.Second,
			want: false,
		},
		{
			name: "delay below remaining",
			ctx:  retryContextWithDeadline(t, now.Add(100*time.Millisecond)),
			d:    50 * time.Millisecond,
			want: false,
		},
		{
			name: "delay equals remaining",
			ctx:  retryContextWithDeadline(t, now.Add(100*time.Millisecond)),
			d:    100 * time.Millisecond,
			want: true,
		},
		{
			name: "delay above remaining",
			ctx:  retryContextWithDeadline(t, now.Add(100*time.Millisecond)),
			d:    150 * time.Millisecond,
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := execution.contextDeadlineWouldBeExceeded(tc.ctx, tc.d)
			if got != tc.want {
				t.Fatalf("contextDeadlineWouldBeExceeded(%v) = %v, want %v", tc.d, got, tc.want)
			}
		})
	}
}

func TestRunStopsWhenDelayExceedsContextDeadlineBudget(t *testing.T) {
	errBoom := errors.New("boom")
	now := retryFutureNow()
	fake := clock.NewFakeClock(now)
	ctx := retryContextWithDeadline(t, now.Add(50*time.Millisecond))
	recorder := &retryObserverRecorder{}
	cfg := configOf(
		WithClock(fake),
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Fixed(100*time.Millisecond)),
		WithObserver(recorder),
	)

	calls := 0
	_, err := run(ctx, func(context.Context) (int, error) {
		calls++
		return 0, errBoom
	}, cfg)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("run error = %v, want ErrExhausted", err)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}

	outcome, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome returned ok=false")
	}
	if outcome.Reason != StopReasonDeadline {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonDeadline)
	}
	if outcome.Attempts != 1 {
		t.Fatalf("Outcome.Attempts = %d, want 1", outcome.Attempts)
	}
	if !errors.Is(outcome.LastErr, errBoom) {
		t.Fatalf("Outcome.LastErr = %v, want %v", outcome.LastErr, errBoom)
	}
	if retryCountEvents(recorder.events, EventRetryDelay) != 0 {
		t.Fatalf("retry delay events = %d, want 0", retryCountEvents(recorder.events, EventRetryDelay))
	}

	stop := retryLastEvent(t, recorder.events)
	if stop.Kind != EventRetryStop {
		t.Fatalf("last event kind = %s, want %s", stop.Kind, EventRetryStop)
	}
	if stop.Outcome.Reason != StopReasonDeadline {
		t.Fatalf("stop reason = %s, want %s", stop.Outcome.Reason, StopReasonDeadline)
	}
}

func TestRunStopsWhenDelayEqualsContextDeadlineBudget(t *testing.T) {
	errBoom := errors.New("boom")
	now := retryFutureNow()
	fake := clock.NewFakeClock(now)
	ctx := retryContextWithDeadline(t, now.Add(50*time.Millisecond))
	recorder := &retryObserverRecorder{}
	cfg := configOf(
		WithClock(fake),
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Fixed(50*time.Millisecond)),
		WithObserver(recorder),
	)

	calls := 0
	_, err := run(ctx, func(context.Context) (int, error) {
		calls++
		return 0, errBoom
	}, cfg)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("run error = %v, want ErrExhausted", err)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}

	outcome, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome returned ok=false")
	}
	if outcome.Reason != StopReasonDeadline {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonDeadline)
	}
	if retryCountEvents(recorder.events, EventRetryDelay) != 0 {
		t.Fatalf("retry delay events = %d, want 0", retryCountEvents(recorder.events, EventRetryDelay))
	}
}

func TestRunAllowsDelayBelowContextDeadlineBudget(t *testing.T) {
	errBoom := errors.New("boom")
	now := retryFutureNow()
	fake := clock.NewFakeClock(now)
	signalingClock := newRetryTimerSignalClock(fake)
	ctx := retryContextWithDeadline(t, now.Add(100*time.Millisecond))
	recorder := &retryObserverRecorder{}
	cfg := configOf(
		WithClock(signalingClock),
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Fixed(50*time.Millisecond)),
		WithObserver(recorder),
	)

	calls := 0
	done := make(chan error, 1)
	go func() {
		_, err := run(ctx, func(context.Context) (int, error) {
			calls++
			if calls == 1 {
				return 0, errBoom
			}
			return 7, nil
		}, cfg)
		done <- err
	}()

	<-signalingClock.timerCreated
	fake.Step(50 * time.Millisecond)

	err := <-done
	if err != nil {
		t.Fatalf("run error = %v, want nil", err)
	}
	if calls != 2 {
		t.Fatalf("operation calls = %d, want 2", calls)
	}
	if retryCountEvents(recorder.events, EventRetryDelay) != 1 {
		t.Fatalf("retry delay events = %d, want 1", retryCountEvents(recorder.events, EventRetryDelay))
	}

	stop := retryLastEvent(t, recorder.events)
	if stop.Outcome.Reason != StopReasonSucceeded {
		t.Fatalf("stop reason = %s, want %s", stop.Outcome.Reason, StopReasonSucceeded)
	}
}

func TestRunWithoutContextDeadlineDoesNotRestrictDelay(t *testing.T) {
	errBoom := errors.New("boom")
	now := retryFutureNow()
	fake := clock.NewFakeClock(now)
	signalingClock := newRetryTimerSignalClock(fake)
	recorder := &retryObserverRecorder{}
	cfg := configOf(
		WithClock(signalingClock),
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Fixed(50*time.Millisecond)),
		WithObserver(recorder),
	)

	calls := 0
	done := make(chan error, 1)
	go func() {
		_, err := run(context.Background(), func(context.Context) (int, error) {
			calls++
			if calls == 1 {
				return 0, errBoom
			}
			return 1, nil
		}, cfg)
		done <- err
	}()

	<-signalingClock.timerCreated
	fake.Step(50 * time.Millisecond)

	err := <-done
	if err != nil {
		t.Fatalf("run error = %v, want nil", err)
	}
	if calls != 2 {
		t.Fatalf("operation calls = %d, want 2", calls)
	}
	if retryCountEvents(recorder.events, EventRetryDelay) != 1 {
		t.Fatalf("retry delay events = %d, want 1", retryCountEvents(recorder.events, EventRetryDelay))
	}
}

func retryContextWithDeadline(t *testing.T, dl time.Time) context.Context {
	t.Helper()

	ctx, cancel := context.WithDeadline(context.Background(), dl)
	t.Cleanup(cancel)
	return ctx
}

func retryCountEvents(events []Event, kind EventKind) int {
	count := 0
	for _, event := range events {
		if event.Kind == kind {
			count++
		}
	}

	return count
}

func retryLastEvent(t *testing.T, events []Event) Event {
	t.Helper()

	if len(events) == 0 {
		t.Fatal("events len = 0, want at least one event")
	}

	return events[len(events)-1]
}
