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

	"arcoris.dev/component-base/pkg/clock"
	"arcoris.dev/component-base/pkg/delay"
)

func TestRunSucceedsOnFirstAttempt(t *testing.T) {
	config := configOf()
	calls := 0

	value, err := run(context.Background(), func(context.Context) (int, error) {
		calls++
		return 42, nil
	}, config)

	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if value != 42 {
		t.Fatalf("run value = %d, want 42", value)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
}

func TestRunDefaultDoesNotRetry(t *testing.T) {
	config := configOf()
	errBoom := errors.New("boom")
	calls := 0

	value, err := run(context.Background(), func(context.Context) (int, error) {
		calls++
		return 99, errBoom
	}, config)

	if !errors.Is(err, errBoom) {
		t.Fatalf("run error = %v, want %v", err, errBoom)
	}
	if value != 0 {
		t.Fatalf("run value = %d, want zero value", value)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
}

func TestRunRetriesRetryableErrors(t *testing.T) {
	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithDelaySchedule(delay.Immediate()),
	)

	errBoom := errors.New("boom")
	calls := 0

	value, err := run(context.Background(), func(context.Context) (int, error) {
		calls++
		if calls < 3 {
			return 0, errBoom
		}
		return 7, nil
	}, config)

	if err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if value != 7 {
		t.Fatalf("run value = %d, want 7", value)
	}
	if calls != 3 {
		t.Fatalf("operation calls = %d, want 3", calls)
	}
}

func TestRunStopsAtMaxAttempts(t *testing.T) {
	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Immediate()),
	)

	errBoom := errors.New("boom")
	calls := 0

	_, err := run(context.Background(), func(context.Context) (int, error) {
		calls++
		return 0, errBoom
	}, config)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("run error does not match ErrExhausted: %v", err)
	}
	if !errors.Is(err, errBoom) {
		t.Fatalf("run error does not preserve last operation error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("operation calls = %d, want 2", calls)
	}

	outcome, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome returned ok=false")
	}
	if outcome.Reason != StopReasonMaxAttempts {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonMaxAttempts)
	}
	if outcome.Attempts != 2 {
		t.Fatalf("Outcome.Attempts = %d, want 2", outcome.Attempts)
	}
}

func TestRunReturnsOriginalNonRetryableError(t *testing.T) {
	config := configOf(
		WithClassifier(NeverRetry()),
		WithMaxAttempts(5),
		WithDelaySchedule(delay.Immediate()),
	)

	errBoom := errors.New("boom")
	calls := 0

	_, err := run(context.Background(), func(context.Context) (int, error) {
		calls++
		return 0, errBoom
	}, config)

	if !errors.Is(err, errBoom) {
		t.Fatalf("run error = %v, want %v", err, errBoom)
	}
	if errors.Is(err, ErrExhausted) {
		t.Fatalf("non-retryable error matched ErrExhausted")
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
}

func TestRunStopsWhenDelaySequenceIsExhausted(t *testing.T) {
	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(10),
		WithDelaySchedule(delay.Limit(delay.Immediate(), 1)),
	)

	errBoom := errors.New("boom")
	calls := 0

	_, err := run(context.Background(), func(context.Context) (int, error) {
		calls++
		return 0, errBoom
	}, config)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("run error does not match ErrExhausted: %v", err)
	}
	if calls != 2 {
		t.Fatalf("operation calls = %d, want 2", calls)
	}

	outcome, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome returned ok=false")
	}
	if outcome.Reason != StopReasonDelayExhausted {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonDelayExhausted)
	}
}

func TestRunStopsAtMaxElapsed(t *testing.T) {
	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(10),
		WithDelaySchedule(delay.Fixed(time.Hour)),
		WithMaxElapsed(time.Millisecond),
	)

	errBoom := errors.New("boom")
	calls := 0

	_, err := run(context.Background(), func(context.Context) (int, error) {
		calls++
		return 0, errBoom
	}, config)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("run error does not match ErrExhausted: %v", err)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}

	outcome, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome returned ok=false")
	}
	if outcome.Reason != StopReasonMaxElapsed {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonMaxElapsed)
	}
}

func TestRunReturnsInterruptedWhenContextAlreadyStopped(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var stop Event
	stopEvents := 0
	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithDelaySchedule(delay.Immediate()),
		WithObserverFunc(func(_ context.Context, event Event) {
			if event.Kind != EventRetryStop {
				return
			}
			stop = event
			stopEvents++
		}),
	)

	calls := 0

	_, err := run(ctx, func(context.Context) (int, error) {
		calls++
		return 0, errors.New("boom")
	}, config)

	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("run error does not match ErrInterrupted: %v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error does not preserve context.Canceled: %v", err)
	}
	if calls != 0 {
		t.Fatalf("operation calls = %d, want 0", calls)
	}
	if stopEvents != 1 {
		t.Fatalf("stop events = %d, want 1", stopEvents)
	}
	if !stop.IsValid() {
		t.Fatalf("stop event is invalid: %+v", stop)
	}
	if stop.Outcome.Reason != StopReasonInterrupted {
		t.Fatalf("stop reason = %s, want %s", stop.Outcome.Reason, StopReasonInterrupted)
	}
	if stop.Outcome.Attempts != 0 {
		t.Fatalf("stop attempts = %d, want 0", stop.Outcome.Attempts)
	}
}

func TestRunEmitsObserverEvents(t *testing.T) {
	var events []Event

	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Immediate()),
		WithObserverFunc(func(_ context.Context, event Event) {
			if !event.IsValid() {
				t.Fatalf("observer received invalid event: %+v", event)
			}
			events = append(events, event)
		}),
	)

	errBoom := errors.New("boom")
	calls := 0

	_, err := run(context.Background(), func(context.Context) (int, error) {
		calls++
		return 0, errBoom
	}, config)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("run error does not match ErrExhausted: %v", err)
	}

	wantKinds := []EventKind{
		EventAttemptStart,
		EventAttemptFailure,
		EventRetryDelay,
		EventAttemptStart,
		EventAttemptFailure,
		EventRetryStop,
	}

	if len(events) != len(wantKinds) {
		t.Fatalf("events len = %d, want %d: %+v", len(events), len(wantKinds), events)
	}

	for i, want := range wantKinds {
		if events[i].Kind != want {
			t.Fatalf("event[%d].Kind = %s, want %s", i, events[i].Kind, want)
		}
	}
}

func TestRunOperationOwnedContextErrorIsNotInterrupted(t *testing.T) {
	config := configOf()

	var stop Event
	observedStop := false
	config.observers = append(config.observers, ObserverFunc(func(_ context.Context, event Event) {
		if event.Kind == EventRetryStop {
			stop = event
			observedStop = true
		}
	}))

	_, err := run(context.Background(), func(context.Context) (int, error) {
		return 0, context.Canceled
	}, config)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error = %v, want context.Canceled", err)
	}
	if errors.Is(err, ErrInterrupted) {
		t.Fatalf("operation-owned context error matched ErrInterrupted")
	}
	if !observedStop {
		t.Fatalf("retry stop event was not observed")
	}
	if stop.Outcome.Reason != StopReasonNonRetryable {
		t.Fatalf("stop reason = %s, want %s", stop.Outcome.Reason, StopReasonNonRetryable)
	}
	if !stop.IsValid() {
		t.Fatalf("stop event is invalid: %+v", stop)
	}
}

func TestRunReturnsInterruptedWhenContextStopsDuringDelay(t *testing.T) {
	errBoom := errors.New("boom")
	fake := clock.NewFakeClock(time.Unix(100, 0))
	signalingClock := newRetryTimerSignalClock(fake)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopEvents := make(chan Event, 1)
	config := configOf(
		WithClock(signalingClock),
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Fixed(time.Hour)),
		WithObserverFunc(func(_ context.Context, event Event) {
			if event.Kind == EventRetryStop {
				stopEvents <- event
			}
		}),
	)

	errs := make(chan error, 1)
	go func() {
		_, err := run(ctx, func(context.Context) (int, error) {
			return 0, errBoom
		}, config)
		errs <- err
	}()

	// The timer-created signal proves retry has crossed into its own delay wait.
	// Cancelling after this point distinguishes retry-owned interruption from an
	// operation returning a raw context error.
	<-signalingClock.timerCreated
	cancel()

	err := <-errs
	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("run error does not match ErrInterrupted: %v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("run error does not preserve context.Canceled: %v", err)
	}

	stop := <-stopEvents
	if !stop.IsValid() {
		t.Fatalf("stop event is invalid: %+v", stop)
	}
	if stop.Outcome.Reason != StopReasonInterrupted {
		t.Fatalf("stop reason = %s, want %s", stop.Outcome.Reason, StopReasonInterrupted)
	}
	if stop.Outcome.Attempts != 1 {
		t.Fatalf("stop attempts = %d, want 1", stop.Outcome.Attempts)
	}
	if !errors.Is(stop.Outcome.LastErr, errBoom) {
		t.Fatalf("stop last error = %v, want %v", stop.Outcome.LastErr, errBoom)
	}
}

func TestRunStopEventsAreValidForTerminalReasons(t *testing.T) {
	errBoom := errors.New("boom")

	tests := []struct {
		name    string
		ctx     func() context.Context
		config  config
		op      ValueOperation[int]
		wantErr error
		reason  StopReason
	}{
		{
			name:   "succeeded",
			ctx:    context.Background,
			config: configOf(),
			op: func(context.Context) (int, error) {
				return 1, nil
			},
			reason: StopReasonSucceeded,
		},
		{
			name:   "non retryable",
			ctx:    context.Background,
			config: configOf(),
			op: func(context.Context) (int, error) {
				return 0, errBoom
			},
			wantErr: errBoom,
			reason:  StopReasonNonRetryable,
		},
		{
			name: "max attempts",
			ctx:  context.Background,
			config: configOf(
				WithClassifier(RetryAll()),
				WithMaxAttempts(1),
			),
			op: func(context.Context) (int, error) {
				return 0, errBoom
			},
			wantErr: ErrExhausted,
			reason:  StopReasonMaxAttempts,
		},
		{
			name: "max elapsed",
			ctx:  context.Background,
			config: configOf(
				WithClassifier(RetryAll()),
				WithMaxAttempts(2),
				WithDelaySchedule(delay.Fixed(time.Second)),
				WithMaxElapsed(time.Nanosecond),
			),
			op: func(context.Context) (int, error) {
				return 0, errBoom
			},
			wantErr: ErrExhausted,
			reason:  StopReasonMaxElapsed,
		},
		{
			name: "delay exhausted",
			ctx:  context.Background,
			config: configOf(
				WithClassifier(RetryAll()),
				WithMaxAttempts(2),
			),
			op: func(context.Context) (int, error) {
				return 0, errBoom
			},
			wantErr: ErrExhausted,
			reason:  StopReasonDelayExhausted,
		},
		{
			name: "interrupted before attempt",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			config:  configOf(),
			op:      func(context.Context) (int, error) { return 0, nil },
			wantErr: ErrInterrupted,
			reason:  StopReasonInterrupted,
		},
		{
			name: "interrupted after failed attempt",
			ctx:  context.Background,
			config: configOf(
				WithClassifier(RetryAll()),
				WithMaxAttempts(2),
			),
			op: func(context.Context) (int, error) {
				return 0, errBoom
			},
			wantErr: ErrInterrupted,
			reason:  StopReasonInterrupted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopEvents := 0
			var stop Event
			config := tt.config
			ctx := tt.ctx()
			if tt.reason == StopReasonDelayExhausted {
				config.delay = retryTestSchedule{sequence: &retryTestSequence{}}
			}
			if tt.name == "interrupted after failed attempt" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(context.Background())
				// Cancelling from the failure observer stops retry at its next
				// context boundary, after LastErr has been set but before delay
				// selection can happen.
				config.observers = append(config.observers, ObserverFunc(func(_ context.Context, event Event) {
					if event.Kind == EventAttemptFailure {
						cancel()
					}
				}))
			}
			config.observers = append(config.observers, ObserverFunc(func(_ context.Context, event Event) {
				if event.Kind != EventRetryStop {
					return
				}
				stopEvents++
				stop = event
			}))

			_, err := run(ctx, tt.op, config)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("run error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("run error = %v, want %v", err, tt.wantErr)
			}

			if stopEvents != 1 {
				t.Fatalf("stop events = %d, want 1", stopEvents)
			}
			if !stop.IsValid() {
				t.Fatalf("stop event is invalid: %+v", stop)
			}
			if stop.Outcome.Reason != tt.reason {
				t.Fatalf("stop reason = %s, want %s", stop.Outcome.Reason, tt.reason)
			}
		})
	}
}
