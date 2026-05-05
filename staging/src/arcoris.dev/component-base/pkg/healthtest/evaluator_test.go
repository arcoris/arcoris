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

package healthtest

import (
	"context"
	"testing"

	"arcoris.dev/component-base/pkg/health"
)

func TestEvaluatorHelpers(t *testing.T) {
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
