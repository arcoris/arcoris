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

// TestUntilReturnsNilWhenConditionIsImmediatelySatisfied verifies immediate
// success semantics.
func TestUntilReturnsNilWhenConditionIsImmediatelySatisfied(t *testing.T) {
	t.Parallel()

	calls := 0
	condition := func(context.Context) (bool, error) {
		calls++
		return true, nil
	}

	if err := Until(context.Background(), time.Hour, condition); err != nil {
		t.Fatalf("Until(...) = %v, want nil", err)
	}
	if calls != 1 {
		t.Fatalf("condition calls = %d, want 1", calls)
	}
}

// TestUntilRepeatsUntilConditionIsSatisfied verifies fixed-interval repeated
// evaluation after unsatisfied condition results.
func TestUntilRepeatsUntilConditionIsSatisfied(t *testing.T) {
	t.Parallel()

	calls := 0
	condition := func(context.Context) (bool, error) {
		calls++
		return calls == 3, nil
	}

	if err := Until(context.Background(), time.Millisecond, condition); err != nil {
		t.Fatalf("Until(...) = %v, want nil", err)
	}
	if calls != 3 {
		t.Fatalf("condition calls = %d, want 3", calls)
	}
}

// TestUntilReturnsConditionErrorUnchanged verifies that condition-owned errors
// are terminal and are not wrapped by the wait loop.
func TestUntilReturnsConditionErrorUnchanged(t *testing.T) {
	t.Parallel()

	conditionErr := errors.New("condition failed")
	condition := func(context.Context) (bool, error) {
		return false, conditionErr
	}

	err := Until(context.Background(), time.Hour, condition)

	if err != conditionErr {
		t.Fatalf("Until(...) = %v, want exact condition error %v", err, conditionErr)
	}
	mustNotBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
}

// TestUntilReturnsRawConditionContextErrorUnchanged verifies that raw context
// errors returned by a condition remain condition-owned errors.
func TestUntilReturnsRawConditionContextErrorUnchanged(t *testing.T) {
	t.Parallel()

	condition := func(context.Context) (bool, error) {
		return false, context.Canceled
	}

	err := Until(context.Background(), time.Hour, condition)

	if err != context.Canceled {
		t.Fatalf("Until(...) = %v, want raw context.Canceled", err)
	}
	mustNotBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
}

// TestUntilDoesNotEvaluateWhenContextIsCancelledBeforeStart verifies that a
// stopped context wins before the first condition evaluation.
func TestUntilDoesNotEvaluateWhenContextIsCancelledBeforeStart(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	condition := func(context.Context) (bool, error) {
		calls++
		return true, nil
	}

	err := Until(ctx, time.Hour, condition)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
	if calls != 0 {
		t.Fatalf("condition calls = %d, want 0", calls)
	}
}

// TestUntilReturnsTimeoutWhenContextDeadlineExceededBeforeStart verifies that an
// already-expired deadline is classified as a wait-owned timeout.
func TestUntilReturnsTimeoutWhenContextDeadlineExceededBeforeStart(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	calls := 0
	condition := func(context.Context) (bool, error) {
		calls++
		return true, nil
	}

	err := Until(ctx, time.Hour, condition)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, context.DeadlineExceeded)
	if calls != 0 {
		t.Fatalf("condition calls = %d, want 0", calls)
	}
}

// TestUntilReturnsInterruptedWhenContextIsCancelledDuringCondition verifies that
// cancellation after an unsatisfied evaluation stops the wait without another
// sleep.
func TestUntilReturnsInterruptedWhenContextIsCancelledDuringCondition(t *testing.T) {
	t.Parallel()

	cause := errors.New("shutdown requested")
	ctx, cancel := context.WithCancelCause(context.Background())
	calls := 0
	condition := func(context.Context) (bool, error) {
		calls++
		cancel(cause)
		return false, nil
	}

	err := Until(ctx, time.Hour, condition)

	mustBeInterrupted(t, err)
	mustNotBeTimedOut(t, err)
	mustMatch(t, err, context.Canceled)
	mustMatch(t, err, cause)
	if calls != 1 {
		t.Fatalf("condition calls = %d, want 1", calls)
	}
}

// TestUntilReturnsNilWhenConditionSucceedsWhileCancellingContext verifies that a
// successful condition result wins for the evaluation that produced it.
func TestUntilReturnsNilWhenConditionSucceedsWhileCancellingContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	condition := func(context.Context) (bool, error) {
		cancel()
		return true, nil
	}

	if err := Until(ctx, time.Hour, condition); err != nil {
		t.Fatalf("Until(...) = %v, want nil", err)
	}
}

// TestUntilReturnsTimeoutWhenContextDeadlineExpiresDuringWait verifies context
// stop behavior while Until is sleeping between evaluations.
func TestUntilReturnsTimeoutWhenContextDeadlineExpiresDuringWait(t *testing.T) {
	t.Parallel()

	cause := errors.New("finite wait budget exhausted")
	ctx, cancel := context.WithTimeoutCause(context.Background(), time.Millisecond, cause)
	defer cancel()

	condition := func(context.Context) (bool, error) {
		return false, nil
	}

	err := Until(ctx, time.Hour, condition)

	mustBeTimedOut(t, err)
	mustBeInterrupted(t, err)
	mustMatch(t, err, context.DeadlineExceeded)
	mustMatch(t, err, cause)
}

// TestUntilPanicsOnNilContext verifies invalid context validation at the public
// wait-loop boundary.
func TestUntilPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilContext, func() {
		_ = Until(nil, time.Second, Satisfied)
	})
}

// TestUntilPanicsOnNonPositiveInterval verifies that fixed-cadence waits reject
// intervals that would produce busy-loop semantics.
func TestUntilPanicsOnNonPositiveInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		interval time.Duration
	}{
		{
			name:     "zero",
			interval: 0,
		},
		{
			name:     "negative",
			interval: -time.Second,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNonPositiveInterval, func() {
				_ = Until(context.Background(), tc.interval, Satisfied)
			})
		})
	}
}

// TestUntilPanicsOnNilCondition verifies that Until reuses the package-wide nil
// condition policy.
func TestUntilPanicsOnNilCondition(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, errNilCondition, func() {
		_ = Until(context.Background(), time.Second, nil)
	})
}

// TestUntilDoesNotRecoverConditionPanics verifies that panic recovery is not a
// responsibility of the low-level fixed-cadence wait loop.
func TestUntilDoesNotRecoverConditionPanics(t *testing.T) {
	t.Parallel()

	panicValue := "condition panic"
	condition := func(context.Context) (bool, error) {
		panic(panicValue)
	}

	mustPanicWith(t, panicValue, func() {
		_ = Until(context.Background(), time.Second, condition)
	})
}
