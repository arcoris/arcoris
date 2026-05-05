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
	"sync/atomic"
	"testing"
	"time"
)

func TestEvaluatorParallelPreservesRegistryOrderWhenChecksFinishOutOfOrder(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	releaseFirst := make(chan struct{})
	firstStarted := make(chan struct{})
	secondDone := make(chan struct{})
	thirdDone := make(chan struct{})

	mustRegisterExecutionCheck(t, registry, TargetReady, "first", func(context.Context) Result {
		close(firstStarted)
		<-releaseFirst
		return Healthy("first")
	})
	mustRegisterExecutionCheck(t, registry, TargetReady, "second", func(context.Context) Result {
		close(secondDone)
		return Healthy("second")
	})
	mustRegisterExecutionCheck(t, registry, TargetReady, "third", func(context.Context) Result {
		close(thirdDone)
		return Healthy("third")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(TargetReady, 3),
	)

	done := make(chan Report, 1)
	go func() {
		report, err := evaluator.Evaluate(context.Background(), TargetReady)
		if err != nil {
			t.Errorf("Evaluate() = %v, want nil", err)
		}
		done <- report
	}()

	<-firstStarted
	<-secondDone
	<-thirdDone
	close(releaseFirst)

	var report Report
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

	registry := NewRegistry()
	release := make(chan struct{})
	started := make(chan struct{}, checkCount)

	var active atomic.Int64
	var maxSeen atomic.Int64

	for i := 0; i < checkCount; i++ {
		name := executionCheckName(i)
		mustRegisterExecutionCheck(t, registry, TargetReady, name, func(context.Context) Result {
			current := active.Add(1)
			updateMaxInt64(&maxSeen, current)
			started <- struct{}{}

			<-release

			active.Add(-1)
			return Healthy(name)
		})
	}

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(TargetReady, limit),
	)

	done := make(chan Report, 1)
	go func() {
		report, err := evaluator.Evaluate(context.Background(), TargetReady)
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

	var report Report
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

	registry := NewRegistry()
	mustRegisterExecutionCheck(t, registry, TargetReady, "healthy", func(context.Context) Result {
		return Healthy("healthy")
	})
	mustRegisterExecutionCheck(t, registry, TargetReady, "degraded", func(context.Context) Result {
		return Degraded("degraded", ReasonOverloaded, "degraded")
	})
	mustRegisterExecutionCheck(t, registry, TargetReady, "unhealthy", func(context.Context) Result {
		return Unhealthy("unhealthy", ReasonFatal, "unhealthy")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(TargetReady, 3),
	)

	report, err := evaluator.Evaluate(context.Background(), TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if report.Status != StatusUnhealthy {
		t.Fatalf("Status = %s, want unhealthy", report.Status)
	}
}

func TestEvaluatorParallelPreservesNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		checkName  string
		fn         CheckFunc
		wantStatus Status
		wantReason Reason
	}{
		{
			name:       "panic",
			checkName:  "panic_check",
			fn:         func(context.Context) Result { panic("boom") },
			wantStatus: StatusUnhealthy,
			wantReason: ReasonPanic,
		},
		{
			name:       "invalid reason",
			checkName:  "invalid_reason",
			fn:         func(context.Context) Result { return Unknown("invalid_reason", Reason("bad-reason"), "bad") },
			wantStatus: StatusUnknown,
			wantReason: ReasonMisconfigured,
		},
		{
			name:       "mismatched name",
			checkName:  "mismatched_name",
			fn:         func(context.Context) Result { return Healthy("other_name") },
			wantStatus: StatusUnknown,
			wantReason: ReasonMisconfigured,
		},
		{
			name:      "cause preserved internally",
			checkName: "cause_check",
			fn: func(context.Context) Result {
				return Unhealthy("cause_check", ReasonFatal, "failed").WithCause(errors.New("private"))
			},
			wantStatus: StatusUnhealthy,
			wantReason: ReasonFatal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			registry := NewRegistry()
			mustRegisterExecutionCheck(t, registry, TargetReady, tc.checkName, tc.fn)

			evaluator := mustExecutionEvaluator(
				t,
				registry,
				WithDefaultTimeout(0),
				WithTargetParallelChecks(TargetReady, 2),
			)

			report, err := evaluator.Evaluate(context.Background(), TargetReady)
			if err != nil {
				t.Fatalf("Evaluate() = %v, want nil", err)
			}
			if len(report.Checks) != 1 {
				t.Fatalf("checks = %d, want 1", len(report.Checks))
			}

			result := report.Checks[0]
			if result.Status != tc.wantStatus {
				t.Fatalf("Status = %s, want %s", result.Status, tc.wantStatus)
			}
			if result.Reason != tc.wantReason {
				t.Fatalf("Reason = %s, want %s", result.Reason, tc.wantReason)
			}
			if result.Observed.IsZero() {
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
		wantReason Reason
	}{
		{
			name:       "timeout",
			ctx:        context.Background(),
			timeout:    time.Nanosecond,
			wantReason: ReasonTimeout,
		},
		{
			name: "canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			timeout:    time.Second,
			wantReason: ReasonCanceled,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			release := make(chan struct{})
			defer close(release)

			registry := NewRegistry()
			mustRegisterExecutionCheck(t, registry, TargetReady, "blocking_one", blockingAfterContextDone(release))
			mustRegisterExecutionCheck(t, registry, TargetReady, "blocking_two", blockingAfterContextDone(release))

			evaluator := mustExecutionEvaluator(
				t,
				registry,
				WithTargetTimeout(TargetReady, tc.timeout),
				WithTargetParallelChecks(TargetReady, 2),
			)

			done := make(chan Report, 1)
			go func() {
				report, err := evaluator.Evaluate(tc.ctx, TargetReady)
				if err != nil {
					t.Errorf("Evaluate() = %v, want nil", err)
				}
				done <- report
			}()

			var report Report
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
					t.Fatalf("Reason for %s = %s, want %s", res.Name, res.Reason, tc.wantReason)
				}
			}
		})
	}
}
