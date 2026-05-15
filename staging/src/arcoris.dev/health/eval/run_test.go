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

package eval

import (
	"context"
	"errors"
	"testing"
	"time"

	"arcoris.dev/health"
)

func TestEvaluateCheckHandlesNilChecker(t *testing.T) {
	t.Parallel()

	evaluator := mustEvaluator(t, health.NewRegistry(), WithClock(newStepClock(testObserved, testObserved)))

	res := evaluator.evaluateCheck(context.Background(), nil, 0)
	if res.Status != health.StatusUnknown || res.Reason != health.ReasonNotObserved {
		t.Fatalf("nil checker result = %+v", res)
	}
	if !errors.Is(res.Cause, health.ErrNilChecker) {
		t.Fatalf("nil checker cause = %v, want health.ErrNilChecker", res.Cause)
	}
}

func TestRunCheckTimeout(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	defer close(release)

	checker := checkerFunc{
		name: "blocked",
		fn: func(context.Context) health.Result {
			<-release
			return health.Healthy("blocked")
		},
	}
	evaluator := mustEvaluator(t, health.NewRegistry())

	res := evaluator.runCheck(context.Background(), checker, time.Nanosecond)
	if res.Status != health.StatusUnknown || res.Reason != health.ReasonTimeout {
		t.Fatalf("timeout result = %+v, want unknown timeout", res)
	}
	if !errors.Is(res.Cause, context.DeadlineExceeded) {
		t.Fatalf("timeout cause = %v, want context deadline", res.Cause)
	}
}

func TestRunCheckReturnsResultBeforeTimeout(t *testing.T) {
	t.Parallel()

	checker := checkerFunc{
		name: "fast",
		fn: func(context.Context) health.Result {
			return health.Healthy("fast")
		},
	}
	evaluator := mustEvaluator(t, health.NewRegistry())

	res := evaluator.runCheck(context.Background(), checker, time.Second)
	if res.Status != health.StatusHealthy {
		t.Fatalf("status = %s, want healthy", res.Status)
	}
}

func TestRunCheckParentCancellation(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	defer close(release)

	ctx, cancel := context.WithCancelCause(context.Background())
	cause := errors.New("owner canceled")
	cancel(cause)

	checker := checkerFunc{
		name: "blocked",
		fn: func(context.Context) health.Result {
			<-release
			return health.Healthy("blocked")
		},
	}
	evaluator := mustEvaluator(t, health.NewRegistry())

	res := evaluator.runCheck(ctx, checker, time.Second)
	if res.Status != health.StatusUnknown || res.Reason != health.ReasonCanceled {
		t.Fatalf("cancel result = %+v, want unknown canceled", res)
	}
	if !errors.Is(res.Cause, cause) {
		t.Fatalf("cancel cause = %v, want custom cause", res.Cause)
	}
}

func TestRunCheckWithZeroTimeoutRunsInline(t *testing.T) {
	t.Parallel()

	called := false
	checker := checkerFunc{
		name: "inline",
		fn: func(context.Context) health.Result {
			called = true
			return health.Healthy("inline")
		},
	}
	evaluator := mustEvaluator(t, health.NewRegistry())

	res := evaluator.runCheck(context.Background(), checker, 0)
	if !called {
		t.Fatal("checker was not called")
	}
	if res.Status != health.StatusHealthy {
		t.Fatalf("status = %s, want healthy", res.Status)
	}
}

func TestCallCheckRecoversPanic(t *testing.T) {
	t.Parallel()

	checker := checkerFunc{
		name: "panic_check",
		fn: func(context.Context) health.Result {
			panic("boom")
		},
	}

	res := callCheck(context.Background(), checker)
	if res.Status != health.StatusUnhealthy || res.Reason != health.ReasonPanic {
		t.Fatalf("panic result = %+v, want unhealthy panic", res)
	}

	var panicErr PanicError
	if !errors.As(res.Cause, &panicErr) {
		t.Fatalf("panic cause = %T, want PanicError", res.Cause)
	}
	if panicErr.Value != "boom" || len(panicErr.Stack) == 0 {
		t.Fatalf("panic details = %+v", panicErr)
	}
}

func TestNormalizeEvaluatedResult(t *testing.T) {
	t.Parallel()

	res := normalizeEvaluatedResult(
		health.Result{Status: health.StatusHealthy, Duration: -time.Second},
		"storage",
		testObserved,
		time.Second,
	)
	if res.Name != "storage" {
		t.Fatalf("name = %q, want storage", res.Name)
	}
	if res.Observed != testObserved {
		t.Fatalf("observed = %v, want %v", res.Observed, testObserved)
	}
	if res.Duration != time.Second {
		t.Fatalf("duration = %s, want 1s", res.Duration)
	}

	res = normalizeEvaluatedResult(
		health.Result{Status: health.StatusHealthy, Duration: -time.Second},
		"storage",
		testObserved,
		-time.Second,
	)
	if res.Duration != 0 {
		t.Fatalf("negative fallback duration = %s, want 0", res.Duration)
	}
}

func TestNormalizeEvaluatedResultRejectsMismatchedResultName(t *testing.T) {
	t.Parallel()

	res := normalizeEvaluatedResult(
		health.Healthy("database"),
		"storage",
		testObserved,
		time.Second,
	)

	if res.Name != "storage" || res.Status != health.StatusUnknown || res.Reason != health.ReasonMisconfigured {
		t.Fatalf("mismatched result normalization = %+v, want unknown misconfigured storage", res)
	}
	if !errors.Is(res.Cause, ErrMismatchedCheckResult) {
		t.Fatalf("cause = %v, want ErrMismatchedCheckResult", res.Cause)
	}
}

func TestNormalizeEvaluatedResultRejectsInvalidReason(t *testing.T) {
	t.Parallel()

	res := normalizeEvaluatedResult(
		health.Result{Status: health.StatusHealthy, Reason: health.Reason("bad-reason")},
		"storage",
		testObserved,
		time.Second,
	)

	if res.Name != "storage" || res.Status != health.StatusUnknown || res.Reason != health.ReasonMisconfigured {
		t.Fatalf("invalid reason normalization = %+v, want unknown misconfigured storage", res)
	}
	if !errors.Is(res.Cause, ErrInvalidCheckResult) {
		t.Fatalf("cause = %v, want ErrInvalidCheckResult", res.Cause)
	}
}

func TestInterruptedResult(t *testing.T) {
	t.Parallel()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 0)
	defer timeoutCancel()

	timeout := interruptedResult("storage", timeoutCtx)
	if timeout.Reason != health.ReasonTimeout || !errors.Is(timeout.Cause, context.DeadlineExceeded) {
		t.Fatalf("timeout result = %+v", timeout)
	}

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	canceled := interruptedResult("storage", canceledCtx)
	if canceled.Reason != health.ReasonCanceled || !errors.Is(canceled.Cause, context.Canceled) {
		t.Fatalf("canceled result = %+v", canceled)
	}

	active := interruptedResult("storage", context.Background())
	if active.Reason != health.ReasonCanceled || active.Cause != nil {
		t.Fatalf("active context result = %+v, want canceled reason with nil cause", active)
	}
}
