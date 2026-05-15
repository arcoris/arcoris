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
	"sync/atomic"
	"testing"
	"time"

	"arcoris.dev/health"
)

func TestEvaluatorParallelPreservesRegistryOrderWhenChecksFinishOutOfOrder(t *testing.T) {
	t.Parallel()

	registry := health.NewRegistry()
	releaseFirst := make(chan struct{})
	firstStarted := make(chan struct{})
	secondDone := make(chan struct{})
	thirdDone := make(chan struct{})

	mustRegisterExecutionCheck(t, registry, health.TargetReady, "first", func(context.Context) health.Result {
		close(firstStarted)
		<-releaseFirst
		return health.Healthy("first")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "second", func(context.Context) health.Result {
		close(secondDone)
		return health.Healthy("second")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "third", func(context.Context) health.Result {
		close(thirdDone)
		return health.Healthy("third")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(health.TargetReady, 3),
	)

	done := make(chan health.Report, 1)
	go func() {
		report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
		if err != nil {
			t.Errorf("Evaluate() = %v, want nil", err)
		}
		done <- report
	}()

	<-firstStarted
	<-secondDone
	<-thirdDone
	close(releaseFirst)

	var report health.Report
	select {
	case report = <-done:
	case <-time.After(executionTestTimeout):
		t.Fatal("parallel evaluation did not finish")
	}

	names := executionResultNames(report.Checks)
	want := []string{"first", "second", "third"}

	if !sameStrings(names, want) {
		t.Fatalf("result names = %v, want %v", names, want)
	}
}

func TestEvaluatorParallelRespectsMaxConcurrency(t *testing.T) {
	t.Parallel()

	const checkCount = 8
	const limit = 3

	registry := health.NewRegistry()
	release := make(chan struct{})
	started := make(chan struct{}, checkCount)

	var active atomic.Int64
	var maxSeen atomic.Int64

	for i := 0; i < checkCount; i++ {
		name := executionCheckName(i)
		mustRegisterExecutionCheck(t, registry, health.TargetReady, name, func(context.Context) health.Result {
			cur := active.Add(1)
			updateMaxInt64(&maxSeen, cur)
			started <- struct{}{}

			<-release

			active.Add(-1)
			return health.Healthy(name)
		})
	}

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(health.TargetReady, limit),
	)

	done := make(chan health.Report, 1)
	go func() {
		report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
		if err != nil {
			t.Errorf("Evaluate() = %v, want nil", err)
		}
		done <- report
	}()

	for i := 0; i < limit; i++ {
		<-started
	}

	if got := maxSeen.Load(); got != limit {
		t.Fatalf("max concurrency = %d, want exactly %d", got, limit)
	}

	close(release)

	var report health.Report
	select {
	case report = <-done:
	case <-time.After(executionTestTimeout):
		t.Fatal("parallel evaluation did not finish")
	}

	if got := maxSeen.Load(); got > limit {
		t.Fatalf("max concurrency after completion = %d, want <= %d", got, limit)
	}
	if len(report.Checks) != checkCount {
		t.Fatalf("checks = %d, want %d", len(report.Checks), checkCount)
	}
}

func TestEvaluatorParallelAggregatesMostSevereStatus(t *testing.T) {
	t.Parallel()

	registry := health.NewRegistry()
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "healthy", func(context.Context) health.Result {
		return health.Healthy("healthy")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "degraded", func(context.Context) health.Result {
		return health.Degraded("degraded", health.ReasonOverloaded, "degraded")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "unhealthy", func(context.Context) health.Result {
		return health.Unhealthy("unhealthy", health.ReasonFatal, "unhealthy")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(health.TargetReady, 3),
	)

	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if report.Status != health.StatusUnhealthy {
		t.Fatalf("health.Status = %s, want unhealthy", report.Status)
	}
}

func TestEvaluatorParallelPreservesNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		checkName  string
		fn         health.CheckFunc
		wantStatus health.Status
		wantReason health.Reason
	}{
		{
			name:       "panic",
			checkName:  "panic_check",
			fn:         func(context.Context) health.Result { panic("boom") },
			wantStatus: health.StatusUnhealthy,
			wantReason: health.ReasonPanic,
		},
		{
			name:      "invalid reason",
			checkName: "invalid_reason",
			fn: func(context.Context) health.Result {
				return health.Unknown("invalid_reason", health.Reason("bad-reason"), "bad")
			},
			wantStatus: health.StatusUnknown,
			wantReason: health.ReasonMisconfigured,
		},
		{
			name:       "mismatched name",
			checkName:  "mismatched_name",
			fn:         func(context.Context) health.Result { return health.Healthy("other_name") },
			wantStatus: health.StatusUnknown,
			wantReason: health.ReasonMisconfigured,
		},
		{
			name:      "cause preserved internally",
			checkName: "cause_check",
			fn: func(context.Context) health.Result {
				return health.Unhealthy("cause_check", health.ReasonFatal, "failed").WithCause(errors.New("private"))
			},
			wantStatus: health.StatusUnhealthy,
			wantReason: health.ReasonFatal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			registry := health.NewRegistry()
			mustRegisterExecutionCheck(t, registry, health.TargetReady, tc.checkName, tc.fn)

			evaluator := mustExecutionEvaluator(
				t,
				registry,
				WithDefaultTimeout(0),
				WithTargetParallelChecks(health.TargetReady, 2),
			)

			report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
			if err != nil {
				t.Fatalf("Evaluate() = %v, want nil", err)
			}
			if len(report.Checks) != 1 {
				t.Fatalf("checks = %d, want 1", len(report.Checks))
			}

			res := report.Checks[0]
			if res.Status != tc.wantStatus {
				t.Fatalf("health.Status = %s, want %s", res.Status, tc.wantStatus)
			}
			if res.Reason != tc.wantReason {
				t.Fatalf("health.Reason = %s, want %s", res.Reason, tc.wantReason)
			}
			if res.Observed.IsZero() {
				t.Fatal("Observed is zero")
			}
		})
	}
}

func TestEvaluatorParallelTimeoutAndCancel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		ctx        context.Context
		timeout    time.Duration
		wantReason health.Reason
	}{
		{
			name:       "timeout",
			ctx:        context.Background(),
			timeout:    time.Nanosecond,
			wantReason: health.ReasonTimeout,
		},
		{
			name: "canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			timeout:    time.Second,
			wantReason: health.ReasonCanceled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			release := make(chan struct{})
			defer close(release)

			registry := health.NewRegistry()
			mustRegisterExecutionCheck(t, registry, health.TargetReady, "blocking_one", blockingAfterContextDone(release))
			mustRegisterExecutionCheck(t, registry, health.TargetReady, "blocking_two", blockingAfterContextDone(release))

			evaluator := mustExecutionEvaluator(
				t,
				registry,
				WithTargetTimeout(health.TargetReady, tc.timeout),
				WithTargetParallelChecks(health.TargetReady, 2),
			)

			done := make(chan health.Report, 1)
			go func() {
				report, err := evaluator.Evaluate(tc.ctx, health.TargetReady)
				if err != nil {
					t.Errorf("Evaluate() = %v, want nil", err)
				}
				done <- report
			}()

			var report health.Report
			select {
			case report = <-done:
			case <-time.After(executionTestTimeout):
				t.Fatal("parallel evaluation did not finish")
			}

			if len(report.Checks) != 2 {
				t.Fatalf("checks = %d, want 2", len(report.Checks))
			}
			for _, res := range report.Checks {
				if res.Reason != tc.wantReason {
					t.Fatalf("health.Reason for %s = %s, want %s", res.Name, res.Reason, tc.wantReason)
				}
			}
		})
	}
}
