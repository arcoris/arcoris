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

	"arcoris.dev/component-base/pkg/backoff"
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
		WithBackoff(backoff.Immediate()),
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
		WithBackoff(backoff.Immediate()),
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
		WithBackoff(backoff.Immediate()),
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

func TestRunStopsWhenBackoffSequenceIsExhausted(t *testing.T) {
	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(10),
		WithBackoff(backoff.Limit(backoff.Immediate(), 1)),
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
	if outcome.Reason != StopReasonBackoffExhausted {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonBackoffExhausted)
	}
}

func TestRunStopsAtMaxElapsed(t *testing.T) {
	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(10),
		WithBackoff(backoff.Fixed(time.Hour)),
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

	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithBackoff(backoff.Immediate()),
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
}

func TestRunEmitsObserverEvents(t *testing.T) {
	var events []Event

	config := configOf(
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithBackoff(backoff.Immediate()),
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
