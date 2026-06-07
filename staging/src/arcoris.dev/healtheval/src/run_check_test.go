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

package eval

import (
	"context"
	"errors"
	"testing"
	"time"

	"arcoris.dev/health"
)

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
	evaluator := mustEvaluator(t, emptyRegistry(t))

	res := evaluator.runCheck(context.Background(), checker, "blocked", time.Nanosecond)
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
	evaluator := mustEvaluator(t, emptyRegistry(t))

	res := evaluator.runCheck(context.Background(), checker, "fast", time.Second)
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
	evaluator := mustEvaluator(t, emptyRegistry(t))

	res := evaluator.runCheck(ctx, checker, "blocked", time.Second)
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
	evaluator := mustEvaluator(t, emptyRegistry(t))

	res := evaluator.runCheck(context.Background(), checker, "inline", 0)
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

	res := callCheck(context.Background(), checker, "panic_check")
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
