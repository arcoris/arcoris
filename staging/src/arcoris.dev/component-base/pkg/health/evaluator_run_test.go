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

package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestEvaluateCheckHandlesNilChecker(t *testing.T) {
	t.Parallel()

	evaluator := mustEvaluator(t, NewRegistry(), WithClock(newStepClock(testObserved, testObserved)))

	result := evaluator.evaluateCheck(context.Background(), nil, 0)
	if result.Status != StatusUnknown || result.Reason != ReasonNotObserved {
		t.Fatalf("nil checker result = %+v", result)
	}
	if !errors.Is(result.Cause, ErrNilChecker) {
		t.Fatalf("nil checker cause = %v, want ErrNilChecker", result.Cause)
	}
}

func TestRunCheckTimeout(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	defer close(release)

	checker := checkerFunc{
		name: "blocked",
		fn: func(context.Context) Result {
			<-release
			return Healthy("blocked")
		},
	}
	evaluator := mustEvaluator(t, NewRegistry())

	result := evaluator.runCheck(context.Background(), checker, time.Nanosecond)
	if result.Status != StatusUnknown || result.Reason != ReasonTimeout {
		t.Fatalf("timeout result = %+v, want unknown timeout", result)
	}
	if !errors.Is(result.Cause, context.DeadlineExceeded) {
		t.Fatalf("timeout cause = %v, want context deadline", result.Cause)
	}
}

func TestRunCheckReturnsResultBeforeTimeout(t *testing.T) {
	t.Parallel()

	checker := checkerFunc{
		name: "fast",
		fn: func(context.Context) Result {
			return Healthy("fast")
		},
	}
	evaluator := mustEvaluator(t, NewRegistry())

	result := evaluator.runCheck(context.Background(), checker, time.Second)
	if result.Status != StatusHealthy {
		t.Fatalf("status = %s, want healthy", result.Status)
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
		fn: func(context.Context) Result {
			<-release
			return Healthy("blocked")
		},
	}
	evaluator := mustEvaluator(t, NewRegistry())

	result := evaluator.runCheck(ctx, checker, time.Second)
	if result.Status != StatusUnknown || result.Reason != ReasonCanceled {
		t.Fatalf("cancel result = %+v, want unknown canceled", result)
	}
	if !errors.Is(result.Cause, cause) {
		t.Fatalf("cancel cause = %v, want custom cause", result.Cause)
	}
}

func TestRunCheckWithZeroTimeoutRunsInline(t *testing.T) {
	t.Parallel()

	called := false
	checker := checkerFunc{
		name: "inline",
		fn: func(context.Context) Result {
			called = true
			return Healthy("inline")
		},
	}
	evaluator := mustEvaluator(t, NewRegistry())

	result := evaluator.runCheck(context.Background(), checker, 0)
	if !called {
		t.Fatal("checker was not called")
	}
	if result.Status != StatusHealthy {
		t.Fatalf("status = %s, want healthy", result.Status)
	}
}

func TestCallCheckRecoversPanic(t *testing.T) {
	t.Parallel()

	checker := checkerFunc{
		name: "panic_check",
		fn: func(context.Context) Result {
			panic("boom")
		},
	}

	result := callCheck(context.Background(), checker)
	if result.Status != StatusUnhealthy || result.Reason != ReasonPanic {
		t.Fatalf("panic result = %+v, want unhealthy panic", result)
	}

	var panicErr PanicError
	if !errors.As(result.Cause, &panicErr) {
		t.Fatalf("panic cause = %T, want PanicError", result.Cause)
	}
	if panicErr.Value != "boom" || len(panicErr.Stack) == 0 {
		t.Fatalf("panic details = %+v", panicErr)
	}
}

func TestNormalizeEvaluatedResult(t *testing.T) {
	t.Parallel()

	result := normalizeEvaluatedResult(
		Result{Status: StatusHealthy, Duration: -time.Second},
		"storage",
		testObserved,
		time.Second,
	)
	if result.Name != "storage" {
		t.Fatalf("name = %q, want storage", result.Name)
	}
	if result.Observed != testObserved {
		t.Fatalf("observed = %v, want %v", result.Observed, testObserved)
	}
	if result.Duration != time.Second {
		t.Fatalf("duration = %s, want 1s", result.Duration)
	}

	result = normalizeEvaluatedResult(
		Result{Status: StatusHealthy, Duration: -time.Second},
		"storage",
		testObserved,
		-time.Second,
	)
	if result.Duration != 0 {
		t.Fatalf("negative fallback duration = %s, want 0", result.Duration)
	}
}

func TestInterruptedResult(t *testing.T) {
	t.Parallel()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 0)
	defer timeoutCancel()

	timeout := interruptedResult("storage", timeoutCtx)
	if timeout.Reason != ReasonTimeout || !errors.Is(timeout.Cause, context.DeadlineExceeded) {
		t.Fatalf("timeout result = %+v", timeout)
	}

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	canceled := interruptedResult("storage", canceledCtx)
	if canceled.Reason != ReasonCanceled || !errors.Is(canceled.Cause, context.Canceled) {
		t.Fatalf("canceled result = %+v", canceled)
	}

	active := interruptedResult("storage", context.Background())
	if active.Reason != ReasonCanceled || active.Cause != nil {
		t.Fatalf("active context result = %+v, want canceled reason with nil cause", active)
	}
}
