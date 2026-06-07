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

package healthtest

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestHealthtestEvaluatorFunc(t *testing.T) {
	t.Parallel()

	evaluator := EvaluatorFunc(func(_ context.Context, target health.Target) (health.Report, error) {
		return HealthyReport(target), nil
	})

	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, report, health.StatusHealthy)
}

func TestHealthtestStaticEvaluator(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluatorForTarget(
		t,
		health.TargetReady,
		HealthyChecker("storage"),
		DegradedChecker("queue", health.ReasonOverloaded),
	)

	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}

	AssertReportStatus(t, report, health.StatusDegraded)
	AssertCheckOrder(t, report, "storage", "queue")
}

func TestHealthtestRegistryEvaluatorPreservesResolverOrder(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluatorForTarget(
		t,
		health.TargetReady,
		HealthyChecker("first"),
		HealthyChecker("second"),
		HealthyChecker("third"),
	)
	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	AssertCheckOrder(t, report, "first", "second", "third")
}

func TestHealthtestRegistryEvaluatorDoesNotOwnEvalPolicy(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluatorForTarget(t, health.TargetReady, NewPanicChecker("panic_check", "boom"))
	defer func() {
		if recover() == nil {
			t.Fatal("RegistryEvaluator recovered panic, want panic to remain test-visible")
		}
	}()

	_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
}

func TestHealthtestNewEvaluatorWithResultsAggregatesStatus(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluatorWithResults(
		t,
		health.TargetReady,
		HealthyResult("storage"),
		DegradedResult("queue", health.ReasonOverloaded),
		UnhealthyResult("database", health.ReasonFatal),
	)
	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, report, health.StatusUnhealthy)
	AssertCheckOrder(t, report, "storage", "queue", "database")
}

func TestHealthtestEvaluatorWithReports(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluatorWithReports(
		HealthyReport(health.TargetReady),
		UnhealthyReport(health.TargetLive),
	)

	ready, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate(ready) = %v, want nil", err)
	}
	AssertReportStatus(t, ready, health.StatusHealthy)

	live, err := evaluator.Evaluate(context.Background(), health.TargetLive)
	if err != nil {
		t.Fatalf("Evaluate(live) = %v, want nil", err)
	}
	AssertReportStatus(t, live, health.StatusUnhealthy)
}

func TestNewEvaluatorWithResults(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluatorWithResults(
		t,
		health.TargetReady,
		HealthyResult("storage"),
		UnhealthyResult("database", health.ReasonFatal),
	)
	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}

	AssertReportStatus(t, report, health.StatusUnhealthy)
	AssertCheckOrder(t, report, "storage", "database")
}

func TestHealthtestErrorEvaluator(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("private evaluator failure")
	evaluator := NewErrorEvaluator(wantErr)

	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Evaluate() = %v, want %v", err, wantErr)
	}
	if report.Target != health.TargetUnknown ||
		report.Status != health.StatusUnknown ||
		report.IsObserved() ||
		len(report.Checks) != 0 {
		t.Fatalf("report = %#v, want zero", report)
	}
}

func TestHealthtestNewSequenceEvaluatorReturnsSequence(t *testing.T) {
	t.Parallel()

	evaluator := NewSequenceEvaluator(
		health.TargetReady,
		StartingReport(health.TargetReady),
		HealthyReport(health.TargetReady),
	)

	first, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("first Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, first, health.StatusStarting)

	second, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("second Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, second, health.StatusHealthy)

	third, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("third Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, third, health.StatusHealthy)
}

func TestNewDefaultTargetsEvaluator(t *testing.T) {
	t.Parallel()

	evaluator := NewDefaultTargetsEvaluator(t)
	for _, target := range []health.Target{
		health.TargetStartup,
		health.TargetLive,
		health.TargetReady,
	} {
		report, err := evaluator.Evaluate(context.Background(), target)
		if err != nil {
			t.Fatalf("Evaluate(%s) = %v, want nil", target, err)
		}
		AssertReportStatus(t, report, health.StatusHealthy)
	}
}

func TestHealthtestFixturesAreValid(t *testing.T) {
	t.Parallel()

	for _, report := range []health.Report{
		HealthyReport(health.TargetReady),
		StartingReport(health.TargetStartup),
		DegradedReport(health.TargetLive),
		UnhealthyReport(health.TargetReady),
		UnknownReport(health.TargetReady),
		MixedReport(health.TargetReady),
	} {
		if !report.IsValid() {
			t.Fatalf("fixture report invalid: %+v", report)
		}
	}
}
