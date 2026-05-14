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

	"arcoris.dev/chrono/delay"
)

func TestDoSucceeds(t *testing.T) {
	calls := 0

	err := Do(context.Background(), func(context.Context) error {
		calls++
		return nil
	})

	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
}

func TestDoReturnsOriginalNonRetryableError(t *testing.T) {
	errBoom := errors.New("boom")
	calls := 0

	err := Do(context.Background(), func(context.Context) error {
		calls++
		return errBoom
	})

	if !errors.Is(err, errBoom) {
		t.Fatalf("Do error = %v, want %v", err, errBoom)
	}
	if errors.Is(err, ErrExhausted) {
		t.Fatalf("non-retryable error matched ErrExhausted")
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
}

func TestDoRetriesRetryableErrors(t *testing.T) {
	errBoom := errors.New("boom")
	calls := 0

	err := Do(
		context.Background(),
		func(context.Context) error {
			calls++
			if calls < 3 {
				return errBoom
			}
			return nil
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithDelaySchedule(delay.Immediate()),
	)

	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("operation calls = %d, want 3", calls)
	}
}

func TestDoReturnsExhaustedAtMaxAttempts(t *testing.T) {
	errBoom := errors.New("boom")
	calls := 0

	err := Do(
		context.Background(),
		func(context.Context) error {
			calls++
			return errBoom
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Immediate()),
	)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("Do error does not match ErrExhausted: %v", err)
	}
	if !errors.Is(err, errBoom) {
		t.Fatalf("Do error does not preserve last operation error: %v", err)
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

func TestDoReturnsInterruptedWhenContextAlreadyStopped(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0

	err := Do(
		ctx,
		func(context.Context) error {
			calls++
			return nil
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithDelaySchedule(delay.Immediate()),
	)

	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("Do error does not match ErrInterrupted: %v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Do error does not preserve context.Canceled: %v", err)
	}
	if calls != 0 {
		t.Fatalf("operation calls = %d, want 0", calls)
	}
}

func TestDoPanicsOnNilContext(t *testing.T) {
	expectPanic(t, panicNilContext, func() {
		_ = Do(nil, func(context.Context) error {
			return nil
		})
	})
}

func TestDoPanicsOnNilOperation(t *testing.T) {
	expectPanic(t, panicNilOperation, func() {
		_ = Do(context.Background(), nil)
	})
}

func TestDoPanicsOnNilOption(t *testing.T) {
	expectPanic(t, panicNilOption, func() {
		_ = Do(
			context.Background(),
			func(context.Context) error {
				return nil
			},
			nil,
		)
	})
}
