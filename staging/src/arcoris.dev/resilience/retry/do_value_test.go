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

func TestDoValueSucceeds(t *testing.T) {
	calls := 0

	value, err := DoValue(context.Background(), func(context.Context) (int, error) {
		calls++
		return 42, nil
	})

	if err != nil {
		t.Fatalf("DoValue returned error: %v", err)
	}
	if value != 42 {
		t.Fatalf("DoValue value = %d, want 42", value)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
}

func TestDoValueReturnsOriginalNonRetryableErrorAndZeroValue(t *testing.T) {
	errBoom := errors.New("boom")
	calls := 0

	value, err := DoValue(context.Background(), func(context.Context) (int, error) {
		calls++
		return 99, errBoom
	})

	if !errors.Is(err, errBoom) {
		t.Fatalf("DoValue error = %v, want %v", err, errBoom)
	}
	if errors.Is(err, ErrExhausted) {
		t.Fatalf("non-retryable error matched ErrExhausted")
	}
	if value != 0 {
		t.Fatalf("DoValue value = %d, want zero value", value)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
}

func TestDoValueRetriesAndReturnsSuccessfulValue(t *testing.T) {
	errBoom := errors.New("boom")
	calls := 0

	value, err := DoValue(
		context.Background(),
		func(context.Context) (int, error) {
			calls++
			if calls < 3 {
				return calls, errBoom
			}
			return 42, nil
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithDelaySchedule(delay.Immediate()),
	)

	if err != nil {
		t.Fatalf("DoValue returned error: %v", err)
	}
	if value != 42 {
		t.Fatalf("DoValue value = %d, want 42", value)
	}
	if calls != 3 {
		t.Fatalf("operation calls = %d, want 3", calls)
	}
}

func TestDoValueReturnsExhaustedAndZeroValue(t *testing.T) {
	errBoom := errors.New("boom")
	calls := 0

	value, err := DoValue(
		context.Background(),
		func(context.Context) (int, error) {
			calls++
			return 99, errBoom
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Immediate()),
	)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("DoValue error does not match ErrExhausted: %v", err)
	}
	if !errors.Is(err, errBoom) {
		t.Fatalf("DoValue error does not preserve last operation error: %v", err)
	}
	if value != 0 {
		t.Fatalf("DoValue value = %d, want zero value", value)
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

func TestDoValueReturnsInterruptedAndZeroValueWhenContextAlreadyStopped(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0

	value, err := DoValue(
		ctx,
		func(context.Context) (int, error) {
			calls++
			return 42, nil
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(3),
		WithDelaySchedule(delay.Immediate()),
	)

	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("DoValue error does not match ErrInterrupted: %v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("DoValue error does not preserve context.Canceled: %v", err)
	}
	if value != 0 {
		t.Fatalf("DoValue value = %d, want zero value", value)
	}
	if calls != 0 {
		t.Fatalf("operation calls = %d, want 0", calls)
	}
}

func TestDoValuePanicsOnNilContext(t *testing.T) {
	expectPanic(t, panicNilContext, func() {
		_, _ = DoValue(nil, func(context.Context) (int, error) {
			return 0, nil
		})
	})
}

func TestDoValuePanicsOnNilOperation(t *testing.T) {
	expectPanic(t, panicNilValueOperation, func() {
		_, _ = DoValue[int](context.Background(), nil)
	})
}

func TestDoValuePanicsOnNilOption(t *testing.T) {
	expectPanic(t, panicNilOption, func() {
		_, _ = DoValue(
			context.Background(),
			func(context.Context) (int, error) {
				return 0, nil
			},
			nil,
		)
	})
}
