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

func TestEvaluateChecksSequential(t *testing.T) {
	t.Parallel()

	evaluator := mustExecutionEvaluator(t, NewRegistry(), WithDefaultTimeout(0))
	order := make([]string, 0, 2)
	checks := []Checker{
		MustCheck("first", func(context.Context) Result {
			order = append(order, "first")
			return Healthy("first")
		}),
		MustCheck("second", func(context.Context) Result {
			order = append(order, "second")
			return Healthy("second")
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

	evaluator := mustExecutionEvaluator(t, NewRegistry(), WithDefaultTimeout(0))
	blocked := false
	checks := []Checker{
		MustCheck("first", func(context.Context) Result {
			blocked = true
			return Healthy("first")
		}),
		MustCheck("second", func(context.Context) Result {
			if !blocked {
				return Unhealthy("second", ReasonMisconfigured, "parallel execution was observed")
			}

			return Healthy("second")
		}),
	}

	results := evaluator.evaluateChecks(
		context.Background(),
		checks,
		0,
		ExecutionPolicy{Mode: ExecutionParallel, MaxConcurrency: 1},
	)

	if got := aggregateStatus(results); got != StatusHealthy {
		t.Fatalf("aggregateStatus() = %s, want healthy", got)
	}
}

func TestAggregateStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		results []Result
		want    Status
	}{
		{name: "empty", want: StatusUnknown},
		{name: "healthy", results: []Result{Healthy("healthy")}, want: StatusHealthy},
		{name: "degraded", results: []Result{Healthy("healthy"), Degraded("degraded", ReasonOverloaded, "degraded")}, want: StatusDegraded},
		{name: "unknown", results: []Result{Degraded("degraded", ReasonOverloaded, "degraded"), Unknown("unknown", ReasonNotObserved, "unknown")}, want: StatusUnknown},
		{name: "unhealthy", results: []Result{Unknown("unknown", ReasonNotObserved, "unknown"), Unhealthy("unhealthy", ReasonFatal, "unhealthy")}, want: StatusUnhealthy},
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
