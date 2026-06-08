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

	"arcoris.dev/chrono/delay"
	panicassert "arcoris.dev/testutil/panic"
)

func TestDoValueObservedReturnsValueAndOutcome(t *testing.T) {
	t.Parallel()

	value, outcome, err := DoValueObserved(context.Background(), func(context.Context) (int, error) {
		return 42, nil
	})

	if err != nil {
		t.Fatalf("DoValueObserved error = %v, want nil", err)
	}
	if value != 42 {
		t.Fatalf("DoValueObserved value = %d, want 42", value)
	}
	if !outcome.IsValid() {
		t.Fatalf("Outcome is invalid: %+v", outcome)
	}
	if outcome.Reason != StopReasonSucceeded {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonSucceeded)
	}
	if outcome.LastErr != nil {
		t.Fatalf("Outcome.LastErr = %v, want nil", outcome.LastErr)
	}
}

func TestDoValueObservedReturnsSuccessfulRetryValue(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("boom")
	calls := 0

	value, outcome, err := DoValueObserved(
		context.Background(),
		func(context.Context) (int, error) {
			calls++
			if calls == 1 {
				return 100, errBoom
			}
			return 42, nil
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Immediate()),
	)

	if err != nil {
		t.Fatalf("DoValueObserved error = %v, want nil", err)
	}
	if value != 42 {
		t.Fatalf("DoValueObserved value = %d, want final successful value", value)
	}
	if calls != 2 {
		t.Fatalf("operation calls = %d, want 2", calls)
	}
	if !outcome.IsValid() {
		t.Fatalf("Outcome is invalid: %+v", outcome)
	}
	if outcome.Attempts != 2 {
		t.Fatalf("Outcome.Attempts = %d, want 2", outcome.Attempts)
	}
	if outcome.Reason != StopReasonSucceeded {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonSucceeded)
	}
}

func TestDoValueObservedReturnsZeroValueOnTerminalFailure(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("boom")

	value, outcome, err := DoValueObserved(
		context.Background(),
		func(context.Context) (int, error) {
			return 99, errBoom
		},
		WithClassifier(RetryAll()),
		WithMaxAttempts(1),
	)

	if !errors.Is(err, ErrExhausted) {
		t.Fatalf("DoValueObserved error = %v, want ErrExhausted", err)
	}
	if value != 0 {
		t.Fatalf("DoValueObserved value = %d, want zero value", value)
	}
	if !outcome.IsValid() {
		t.Fatalf("Outcome is invalid: %+v", outcome)
	}
	if outcome.Reason != StopReasonMaxAttempts {
		t.Fatalf("Outcome.Reason = %s, want %s", outcome.Reason, StopReasonMaxAttempts)
	}

	exhausted, ok := ExhaustedOutcome(err)
	if !ok {
		t.Fatalf("ExhaustedOutcome ok=false")
	}
	if exhausted != outcome {
		t.Fatalf("ExhaustedOutcome = %+v, want %+v", exhausted, outcome)
	}
}

func TestDoValueObservedStopEventMatchesReturnedOutcome(t *testing.T) {
	t.Parallel()

	var stop Outcome
	stopEvents := 0

	value, outcome, err := DoValueObserved(
		context.Background(),
		func(context.Context) (string, error) {
			return "ok", nil
		},
		WithObserverFunc(func(_ context.Context, event Event) {
			if event.Kind == EventRetryStop {
				stop = event.Outcome
				stopEvents++
			}
		}),
	)

	if err != nil {
		t.Fatalf("DoValueObserved error = %v, want nil", err)
	}
	if value != "ok" {
		t.Fatalf("DoValueObserved value = %q, want ok", value)
	}
	if stopEvents != 1 {
		t.Fatalf("stop events = %d, want 1", stopEvents)
	}
	if stop != outcome {
		t.Fatalf("stop outcome = %+v, want returned %+v", stop, outcome)
	}
}

func TestDoValueReturnsSameValueAndErrorAsObservedConvenience(t *testing.T) {
	t.Parallel()

	value, err := DoValue(context.Background(), func(context.Context) (int, error) {
		return 42, nil
	})

	if err != nil {
		t.Fatalf("DoValue error = %v, want nil", err)
	}
	if value != 42 {
		t.Fatalf("DoValue value = %d, want 42", value)
	}
}

func TestDoValueObservedPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	panicassert.RequireErrorIs(t, ErrNilContext, func() {
		_, _, _ = DoValueObserved(nil, func(context.Context) (int, error) {
			return 0, nil
		})
	})
	panicassert.RequireErrorIs(t, ErrNilValueOperation, func() {
		_, _, _ = DoValueObserved[int](context.Background(), nil)
	})
	panicassert.RequireErrorIs(t, ErrNilOption, func() {
		_, _, _ = DoValueObserved(
			context.Background(),
			func(context.Context) (int, error) { return 0, nil },
			nil,
		)
	})
}
