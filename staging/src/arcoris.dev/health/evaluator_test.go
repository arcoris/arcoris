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

func TestNewEvaluatorRejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	if _, err := NewEvaluator(nil); !errors.Is(err, ErrNilRegistry) {
		t.Fatalf("NewEvaluator(nil) = %v, want ErrNilRegistry", err)
	}
	if _, err := NewEvaluator(NewRegistry(), nil); !errors.Is(err, ErrNilEvaluatorOption) {
		t.Fatalf("NewEvaluator(nil option) = %v, want ErrNilEvaluatorOption", err)
	}
	if _, err := NewEvaluator(NewRegistry(), WithDefaultTimeout(-time.Second)); !errors.Is(err, ErrInvalidTimeout) {
		t.Fatalf("NewEvaluator(invalid option) = %v, want ErrInvalidTimeout", err)
	}
}

func TestEvaluatorTimeoutForUsesTargetOverride(t *testing.T) {
	t.Parallel()

	evaluator := mustEvaluator(
		t,
		NewRegistry(),
		WithDefaultTimeout(time.Second),
		WithTargetTimeout(TargetReady, 2*time.Second),
	)

	if got := evaluator.timeoutFor(TargetReady); got != 2*time.Second {
		t.Fatalf("ready timeout = %s, want 2s", got)
	}
	if got := evaluator.timeoutFor(TargetLive); got != time.Second {
		t.Fatalf("live timeout = %s, want 1s", got)
	}
}

func TestEvaluateRejectsInvalidTarget(t *testing.T) {
	t.Parallel()

	evaluator := mustEvaluator(t, NewRegistry(), WithClock(newStepClock(testObserved)))

	report, err := evaluator.Evaluate(context.Background(), TargetUnknown)
	if !errors.Is(err, ErrInvalidTarget) {
		t.Fatalf("Evaluate(invalid target) = %v, want ErrInvalidTarget", err)
	}
	if report.Target != TargetUnknown || report.Status != StatusUnknown || report.Observed != testObserved {
		t.Fatalf("invalid target report = %+v", report)
	}
}

func TestEvaluateReturnsUnknownForTargetWithoutChecks(t *testing.T) {
	t.Parallel()

	evaluator := mustEvaluator(t, NewRegistry(), WithClock(newStepClock(testObserved)))

	report, err := evaluator.Evaluate(nil, TargetReady)
	if err != nil {
		t.Fatalf("Evaluate(empty target) = %v, want nil", err)
	}
	if report.Target != TargetReady || report.Status != StatusUnknown || len(report.Checks) != 0 {
		t.Fatalf("empty target report = %+v", report)
	}
}

func TestEvaluateRunsChecksInRegistryOrderAndAggregatesSeverity(t *testing.T) {
	t.Parallel()

	registry := mustRegistry(
		t,
		TargetReady,
		mustCheck(t, "first", Healthy("")),
		mustCheck(t, "second", Degraded("", ReasonOverloaded, "overloaded")),
	)
	start := testObserved
	evaluator := mustEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithClock(newStepClock(
			start,
			start,
			start.Add(10*time.Millisecond),
			start.Add(10*time.Millisecond),
			start.Add(30*time.Millisecond),
			start.Add(30*time.Millisecond),
		)),
	)

	report, err := evaluator.Evaluate(context.Background(), TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if report.Status != StatusDegraded {
		t.Fatalf("status = %s, want degraded", report.Status)
	}
	if report.Duration != 30*time.Millisecond {
		t.Fatalf("duration = %s, want 30ms", report.Duration)
	}
	if len(report.Checks) != 2 || report.Checks[0].Name != "first" || report.Checks[1].Name != "second" {
		t.Fatalf("checks order = %+v", report.Checks)
	}
	if report.Checks[0].Duration != 10*time.Millisecond || report.Checks[1].Duration != 20*time.Millisecond {
		t.Fatalf("check durations = %s, %s; want 10ms, 20ms", report.Checks[0].Duration, report.Checks[1].Duration)
	}
}

func TestEvaluateClampsNegativeReportDuration(t *testing.T) {
	t.Parallel()

	registry := mustRegistry(t, TargetReady, mustCheck(t, "storage", Healthy("storage")))
	evaluator := mustEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithClock(newStepClock(
			testObserved,
			testObserved,
			testObserved.Add(time.Millisecond),
			testObserved.Add(-time.Second),
		)),
	)

	report, err := evaluator.Evaluate(context.Background(), TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if report.Duration != 0 {
		t.Fatalf("report duration = %s, want 0", report.Duration)
	}
}

func TestEvaluateDoesNotRetainReport(t *testing.T) {
	t.Parallel()

	registry := mustRegistry(t, TargetReady, mustCheck(t, "storage", Healthy("storage")))
	evaluator := mustEvaluator(t, registry, WithDefaultTimeout(0))

	report, err := evaluator.Evaluate(context.Background(), TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	report.Checks[0] = Unhealthy("storage", ReasonFatal, "mutated")

	again, err := evaluator.Evaluate(context.Background(), TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() second = %v, want nil", err)
	}
	if again.Checks[0].Status != StatusHealthy {
		t.Fatalf("retained mutated report status = %s, want healthy", again.Checks[0].Status)
	}
}
