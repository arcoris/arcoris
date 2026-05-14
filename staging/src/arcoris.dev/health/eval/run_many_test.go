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

func TestEvaluateChecksSequential(t *testing.T) {
	t.Parallel()

	evaluator := mustExecutionEvaluator(t, health.NewRegistry(), WithDefaultTimeout(0))
	order := make([]string, 0, 2)
	checks := []health.Checker{
		health.MustCheck("first", func(context.Context) health.Result {
			order = append(order, "first")
			return health.Healthy("first")
		}),
		health.MustCheck("second", func(context.Context) health.Result {
			order = append(order, "second")
			return health.Healthy("second")
		}),
	}

	results := evaluator.evaluateChecksSequential(context.Background(), checks, 0)

	if got, want := executionResultNames(results), []string{"first", "second"}; !sameStrings(got, want) {
		t.Fatalf("result names = %v, want %v", got, want)
	}
	if !sameStrings(order, []string{"first", "second"}) {
		t.Fatalf("execution order = %v, want [first second]", order)
	}
}

func TestEvaluateChecksFallsBackToSequentialForParallelLimitOne(t *testing.T) {
	t.Parallel()

	evaluator := mustExecutionEvaluator(t, health.NewRegistry(), WithDefaultTimeout(0))
	blocked := false
	checks := []health.Checker{
		health.MustCheck("first", func(context.Context) health.Result {
			blocked = true
			return health.Healthy("first")
		}),
		health.MustCheck("second", func(context.Context) health.Result {
			if !blocked {
				return health.Unhealthy("second", health.ReasonMisconfigured, "parallel execution was observed")
			}

			return health.Healthy("second")
		}),
	}

	results := evaluator.evaluateChecks(
		context.Background(),
		checks,
		0,
		ExecutionPolicy{Mode: ExecutionParallel, MaxConcurrency: 1},
	)

	if got := aggregateStatus(results); got != health.StatusHealthy {
		t.Fatalf("aggregateStatus() = %s, want healthy", got)
	}
}

func TestAggregateStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		results []health.Result
		want    health.Status
	}{
		{name: "empty", want: health.StatusUnknown},
		{name: "healthy", results: []health.Result{health.Healthy("healthy")}, want: health.StatusHealthy},
		{name: "degraded", results: []health.Result{health.Healthy("healthy"), health.Degraded("degraded", health.ReasonOverloaded, "degraded")}, want: health.StatusDegraded},
		{name: "unknown", results: []health.Result{health.Degraded("degraded", health.ReasonOverloaded, "degraded"), health.Unknown("unknown", health.ReasonNotObserved, "unknown")}, want: health.StatusUnknown},
		{name: "unhealthy", results: []health.Result{health.Unknown("unknown", health.ReasonNotObserved, "unknown"), health.Unhealthy("unhealthy", health.ReasonFatal, "unhealthy")}, want: health.StatusUnhealthy},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := aggregateStatus(tc.results); got != tc.want {
				t.Fatalf("aggregateStatus() = %s, want %s", got, tc.want)
			}
		})
	}
}
