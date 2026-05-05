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
	"testing"
)

func TestEvaluatorDefaultExecutionIsSequential(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	order := make(chan string, 2)

	mustRegisterExecutionCheck(t, registry, TargetReady, "first", func(context.Context) Result {
		order <- "first"
		return Healthy("first")
	})
	mustRegisterExecutionCheck(t, registry, TargetReady, "second", func(context.Context) Result {
		order <- "second"
		return Healthy("second")
	})

	evaluator := mustExecutionEvaluator(t, registry, WithDefaultTimeout(0))

	report, err := evaluator.Evaluate(context.Background(), TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if report.Status != StatusHealthy {
		t.Fatalf("Status = %s, want healthy", report.Status)
	}
	if got := <-order; got != "first" {
		t.Fatalf("first observed check = %q, want first", got)
	}
	if got := <-order; got != "second" {
		t.Fatalf("second observed check = %q, want second", got)
	}
	if got, want := evaluator.executionPolicyFor(TargetReady), DefaultExecutionPolicy(); got != want {
		t.Fatalf("executionPolicyFor(ready) = %+v, want %+v", got, want)
	}
}

func TestEvaluatorTargetExecutionPolicyOverrideWins(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	mustRegisterExecutionCheck(t, registry, TargetReady, "ready", func(context.Context) Result {
		return Healthy("ready")
	})
	mustRegisterExecutionCheck(t, registry, TargetLive, "live", func(context.Context) Result {
		return Healthy("live")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithParallelChecks(4),
		WithTargetSequentialChecks(TargetLive),
	)

	if got, want := evaluator.executionPolicyFor(TargetReady), ParallelExecutionPolicy(4); got != want {
		t.Fatalf("executionPolicyFor(ready) = %+v, want %+v", got, want)
	}
	if got, want := evaluator.executionPolicyFor(TargetLive), DefaultExecutionPolicy(); got != want {
		t.Fatalf("executionPolicyFor(live) = %+v, want %+v", got, want)
	}
}

func TestEvaluatorTargetExecutionPolicyAppliesOnlyToConfiguredTarget(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	mustRegisterExecutionCheck(t, registry, TargetReady, "ready", func(context.Context) Result {
		return Healthy("ready")
	})
	mustRegisterExecutionCheck(t, registry, TargetLive, "live", func(context.Context) Result {
		return Healthy("live")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(TargetReady, 4),
	)

	if got, want := evaluator.executionPolicyFor(TargetReady), ParallelExecutionPolicy(4); got != want {
		t.Fatalf("executionPolicyFor(ready) = %+v, want %+v", got, want)
	}
	if got, want := evaluator.executionPolicyFor(TargetLive), DefaultExecutionPolicy(); got != want {
		t.Fatalf("executionPolicyFor(live) = %+v, want %+v", got, want)
	}
}
