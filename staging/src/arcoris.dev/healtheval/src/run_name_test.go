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
	"sync/atomic"
	"testing"
	"time"

	"arcoris.dev/health"
)

func TestEvaluateCheckReadsCheckerNameOnce(t *testing.T) {
	t.Parallel()

	checker := &countingNameChecker{
		name: "single_name",
		fn: func(context.Context) health.Result {
			return health.Healthy("")
		},
	}
	evaluator := mustEvaluator(t, emptyRegistry(t), WithDefaultTimeout(0))

	res := evaluator.evaluateCheck(context.Background(), checker, 0)

	if got := checker.nameCalls.Load(); got != 1 {
		t.Fatalf("Name() calls = %d, want 1", got)
	}
	if res.Name != "single_name" || res.Status != health.StatusHealthy {
		t.Fatalf("result = %+v, want healthy single_name", res)
	}
}

func TestEvaluateCheckUsesResolvedNameForPanicRecovery(t *testing.T) {
	t.Parallel()

	checker := &countingNameChecker{
		name:            "panic_check",
		panicAfterFirst: true,
		fn: func(context.Context) health.Result {
			panic("boom")
		},
	}
	evaluator := mustEvaluator(t, emptyRegistry(t), WithDefaultTimeout(0))

	res := evaluator.evaluateCheck(context.Background(), checker, 0)

	if got := checker.nameCalls.Load(); got != 1 {
		t.Fatalf("Name() calls = %d, want 1", got)
	}
	if res.Name != "panic_check" || res.Reason != health.ReasonPanic {
		t.Fatalf("panic result = %+v, want panic_check panic", res)
	}
}

func TestEvaluateCheckUsesResolvedNameForTimeout(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	defer close(release)

	checker := &countingNameChecker{
		name:            "timeout_check",
		panicAfterFirst: true,
		fn: func(context.Context) health.Result {
			<-release
			return health.Healthy("timeout_check")
		},
	}
	evaluator := mustEvaluator(t, emptyRegistry(t))

	res := evaluator.evaluateCheck(context.Background(), checker, time.Nanosecond)

	if got := checker.nameCalls.Load(); got != 1 {
		t.Fatalf("Name() calls = %d, want 1", got)
	}
	if res.Name != "timeout_check" || res.Reason != health.ReasonTimeout {
		t.Fatalf("timeout result = %+v, want timeout_check timeout", res)
	}
}

func TestEvaluateCheckInvalidCheckerNameStaysValid(t *testing.T) {
	t.Parallel()

	checker := checkerFunc{
		name: "bad-name",
		fn: func(context.Context) health.Result {
			t.Fatal("checker with invalid name should not execute")
			return health.Healthy("")
		},
	}
	evaluator := mustEvaluator(t, emptyRegistry(t), WithClock(newStepClock(testObserved, testObserved)))

	res := evaluator.evaluateCheck(context.Background(), checker, 0)

	if res.Name != "" ||
		res.Status != health.StatusUnknown ||
		res.Reason != health.ReasonMisconfigured ||
		!res.IsValid() {
		t.Fatalf("invalid checker result = %+v, want valid unknown misconfigured unnamed result", res)
	}
	if !errors.Is(res.Cause, health.ErrInvalidCheckName) {
		t.Fatalf("cause = %v, want health.ErrInvalidCheckName", res.Cause)
	}
}

type countingNameChecker struct {
	name            string
	panicAfterFirst bool
	nameCalls       atomic.Int64
	fn              func(context.Context) health.Result
}

func (checker *countingNameChecker) Name() string {
	if checker.nameCalls.Add(1) > 1 && checker.panicAfterFirst {
		panic("Name called more than once")
	}

	return checker.name
}

func (checker *countingNameChecker) Check(ctx context.Context) health.Result {
	return checker.fn(ctx)
}
