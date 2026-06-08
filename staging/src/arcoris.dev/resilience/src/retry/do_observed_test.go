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

package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
	panicassert "arcoris.dev/testutil/panic"
)

func TestDoObservedReturnsTerminalOutcome(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("boom")

	tests := []struct {
		name    string
		ctx     func(*testing.T) context.Context
		op      func(*testing.T) Operation
		opts    func(*testing.T) []Option
		wantErr error
		reason  StopReason
		calls   int
	}{
		{
			name: "success",
			ctx:  func(*testing.T) context.Context { return context.Background() },
			op: func(*testing.T) Operation {
				return func(context.Context) error { return nil }
			},
			opts:   func(*testing.T) []Option { return nil },
			reason: StopReasonSucceeded,
			calls:  1,
		},
		{
			name: "non retryable",
			ctx:  func(*testing.T) context.Context { return context.Background() },
			op: func(*testing.T) Operation {
				return func(context.Context) error { return errBoom }
			},
			opts:    func(*testing.T) []Option { return nil },
			wantErr: errBoom,
			reason:  StopReasonNonRetryable,
			calls:   1,
		},
		{
			name: "max attempts",
			ctx:  func(*testing.T) context.Context { return context.Background() },
			op: func(*testing.T) Operation {
				return func(context.Context) error { return errBoom }
			},
			opts: func(*testing.T) []Option {
				return []Option{
					WithClassifier(RetryAll()),
					WithMaxAttempts(2),
					WithDelaySchedule(delay.Immediate()),
				}
			},
			wantErr: ErrExhausted,
			reason:  StopReasonMaxAttempts,
			calls:   2,
		},
		{
			name: "max elapsed",
			ctx:  func(*testing.T) context.Context { return context.Background() },
			op: func(*testing.T) Operation {
				return func(context.Context) error { return errBoom }
			},
			opts: func(*testing.T) []Option {
				return []Option{
					WithClassifier(RetryAll()),
					WithMaxAttempts(3),
					WithDelaySchedule(delay.Fixed(time.Second)),
					WithMaxElapsed(time.Nanosecond),
				}
			},
			wantErr: ErrExhausted,
			reason:  StopReasonMaxElapsed,
			calls:   1,
		},
		{
			name: "deadline",
			ctx: func(t *testing.T) context.Context {
				now := retryFutureNow()
				ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Millisecond))
				t.Cleanup(cancel)
				return ctx
			},
			op: func(*testing.T) Operation {
				return func(context.Context) error { return errBoom }
			},
			opts: func(*testing.T) []Option {
				now := retryFutureNow()
				return []Option{
					WithClock(clock.NewFakeClock(now)),
					WithClassifier(RetryAll()),
					WithMaxAttempts(3),
					WithDelaySchedule(delay.Fixed(time.Second)),
				}
			},
			wantErr: ErrExhausted,
			reason:  StopReasonDeadline,
			calls:   1,
		},
		{
			name: "delay exhausted",
			ctx:  func(*testing.T) context.Context { return context.Background() },
			op: func(*testing.T) Operation {
				return func(context.Context) error { return errBoom }
			},
			opts: func(*testing.T) []Option {
				return []Option{
					WithClassifier(RetryAll()),
					WithMaxAttempts(3),
					WithDelaySchedule(retryTestSchedule{sequence: &retryTestSequence{}}),
				}
			},
			wantErr: ErrExhausted,
			reason:  StopReasonDelayExhausted,
			calls:   1,
		},
		{
			name: "context canceled before attempt",
			ctx: func(*testing.T) context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			op: func(*testing.T) Operation {
				return func(context.Context) error { return nil }
			},
			opts:    func(*testing.T) []Option { return nil },
			wantErr: ErrInterrupted,
			reason:  StopReasonInterrupted,
			calls:   0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			calls := 0
			op := tt.op(t)
			wrapped := func(ctx context.Context) error {
				calls++
				return op(ctx)
			}

			var stop Outcome
			stopEvents := 0
			opts := append([]Option(nil), tt.opts(t)...)
			opts = append(opts, WithObserverFunc(func(_ context.Context, event Event) {
				if event.Kind == EventRetryStop {
					stop = event.Outcome
					stopEvents++
				}
			}))

			outcome, err := DoObserved(tt.ctx(t), wrapped, opts...)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("DoObserved error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DoObserved error = %v, want %v", err, tt.wantErr)
			}
			if !outcome.IsValid() {
				t.Fatalf("Outcome is invalid: %+v", outcome)
			}
			if outcome.Reason != tt.reason {
				t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, tt.reason)
			}
			if calls != tt.calls {
				t.Fatalf("operation calls = %d, want %d", calls, tt.calls)
			}
			if stopEvents != 1 {
				t.Fatalf("stop events = %d, want 1", stopEvents)
			}
			if stop != outcome {
				t.Fatalf("stop outcome = %+v, want returned %+v", stop, outcome)
			}
			if errors.Is(err, ErrExhausted) {
				exhausted, ok := ExhaustedOutcome(err)
				if !ok {
					t.Fatalf("ExhaustedOutcome ok=false")
				}
				if exhausted != outcome {
					t.Fatalf("ExhaustedOutcome = %+v, want %+v", exhausted, outcome)
				}
			}
		})
	}
}

func TestDoObservedReturnsInterruptedAfterFailedAttempt(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("boom")
	ctx, cancel := context.WithCancel(context.Background())

	calls := 0
	outcome, err := DoObserved(
		ctx,
		func(context.Context) error {
			calls++
			return errBoom
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithDelaySchedule(delay.Immediate()),
		WithObserverFunc(func(_ context.Context, event Event) {
			if event.Kind == EventAttemptFailure {
				cancel()
			}
		}),
	)

	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("DoObserved error = %v, want ErrInterrupted", err)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
	if !outcome.IsValid() {
		t.Fatalf("Outcome is invalid: %+v", outcome)
	}
	if outcome.Reason != StopReasonInterrupted {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonInterrupted)
	}
	if outcome.Attempts != 1 {
		t.Fatalf("Outcome.Attempts = %d, want 1", outcome.Attempts)
	}
	if !errors.Is(outcome.LastErr, errBoom) {
		t.Fatalf("Outcome.LastErr = %v, want %v", outcome.LastErr, errBoom)
	}
}

func TestDoObservedKeepsOperationOwnedContextErrors(t *testing.T) {
	t.Parallel()

	outcome, err := DoObserved(context.Background(), func(context.Context) error {
		return context.DeadlineExceeded
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("DoObserved error = %v, want context deadline error", err)
	}
	if errors.Is(err, ErrInterrupted) {
		t.Fatalf("operation-owned context error matched ErrInterrupted")
	}
	if outcome.Reason != StopReasonNonRetryable {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonNonRetryable)
	}
}

func TestDoReturnsSameErrorAsDoObserved(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("boom")

	err := Do(
		context.Background(),
		func(context.Context) error { return errBoom },
		WithClassifier(RetryAll()),
		WithMaxAttempts(1),
	)
	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("Do error = %v, want ErrExhausted", err)
	}

	outcome, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("Do exhausted error did not preserve Outcome")
	}
	if outcome.Reason != StopReasonMaxAttempts {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonMaxAttempts)
	}
}

func TestDoObservedPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	panicassert.RequireErrorIs(t, ErrNilContext, func() {
		_, _ = DoObserved(nil, func(context.Context) error {
			return nil
		})
	})
	panicassert.RequireErrorIs(t, ErrNilOperation, func() {
		_, _ = DoObserved(context.Background(), nil)
	})
	panicassert.RequireErrorIs(t, ErrNilOption, func() {
		_, _ = DoObserved(
			context.Background(),
			func(context.Context) error { return nil },
			nil,
		)
	})
}
