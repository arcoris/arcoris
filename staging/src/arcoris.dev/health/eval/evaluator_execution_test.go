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
	"testing"

	"arcoris.dev/health"
)

func TestEvaluatorDefaultExecutionIsSequential(t *testing.T) {
	t.Parallel()

	registry := health.NewRegistry()
	order := make(chan string, 2)

	mustRegisterExecutionCheck(t, registry, health.TargetReady, "first", func(context.Context) health.Result {
		order <- "first"
		return health.Healthy("first")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "second", func(context.Context) health.Result {
		order <- "second"
		return health.Healthy("second")
	})

	evaluator := mustExecutionEvaluator(t, registry, WithDefaultTimeout(0))

	report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if report.Status != health.StatusHealthy {
		t.Fatalf("health.Status = %s, want healthy", report.Status)
	}
	if got := <-order; got != "first" {
		t.Fatalf("first observed check = %q, want first", got)
	}
	if got := <-order; got != "second" {
		t.Fatalf("second observed check = %q, want second", got)
	}
	if got, want := evaluator.executionPolicyFor(health.TargetReady), DefaultExecutionPolicy(); got != want {
		t.Fatalf("executionPolicyFor(ready) = %+v, want %+v", got, want)
	}
}

func TestEvaluatorTargetExecutionPolicyOverrideWins(t *testing.T) {
	t.Parallel()

	registry := health.NewRegistry()
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "ready", func(context.Context) health.Result {
		return health.Healthy("ready")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetLive, "live", func(context.Context) health.Result {
		return health.Healthy("live")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithParallelChecks(4),
		WithTargetSequentialChecks(health.TargetLive),
	)

	if got, want := evaluator.executionPolicyFor(health.TargetReady), ParallelExecutionPolicy(4); got != want {
		t.Fatalf("executionPolicyFor(ready) = %+v, want %+v", got, want)
	}
	if got, want := evaluator.executionPolicyFor(health.TargetLive), DefaultExecutionPolicy(); got != want {
		t.Fatalf("executionPolicyFor(live) = %+v, want %+v", got, want)
	}
}

func TestEvaluatorTargetExecutionPolicyAppliesOnlyToConfiguredTarget(t *testing.T) {
	t.Parallel()

	registry := health.NewRegistry()
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "ready", func(context.Context) health.Result {
		return health.Healthy("ready")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetLive, "live", func(context.Context) health.Result {
		return health.Healthy("live")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(health.TargetReady, 4),
	)

	if got, want := evaluator.executionPolicyFor(health.TargetReady), ParallelExecutionPolicy(4); got != want {
		t.Fatalf("executionPolicyFor(ready) = %+v, want %+v", got, want)
	}
	if got, want := evaluator.executionPolicyFor(health.TargetLive), DefaultExecutionPolicy(); got != want {
		t.Fatalf("executionPolicyFor(live) = %+v, want %+v", got, want)
	}
}
