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

	"arcoris.dev/health"
)

// NewEvaluator returns an evaluator with test-safe timeout defaults.
//
// A zero default timeout avoids timing flakes in adapter tests. Caller-supplied
// options are applied afterward so tests can still opt into explicit timeout or
// execution-policy behavior.
//
// This helper is for successful fixture construction. Tests that need to assert
// NewEvaluator error behavior should call health.NewEvaluator directly.
func NewEvaluator(t testing.TB, registry *health.Registry, opts ...health.EvaluatorOption) *health.Evaluator {
	t.Helper()

	allOptions := []health.EvaluatorOption{health.WithDefaultTimeout(0)}
	allOptions = append(allOptions, opts...)

	evaluator, err := health.NewEvaluator(registry, allOptions...)
	if err != nil {
		t.Fatalf("health.NewEvaluator() = %v, want nil", err)
	}

	return evaluator
}

// NewEvaluatorForTarget returns an evaluator with checks under target.
//
// Checks are registered in the supplied order, preserving the order that later
// reports expose to adapter tests.
func NewEvaluatorForTarget(
	t testing.TB,
	target health.Target,
	checks ...health.Checker,
) *health.Evaluator {
	t.Helper()

	return NewEvaluator(t, NewRegistry(t, ForTarget(target, checks...)))
}

// NewEvaluatorWithResults returns an evaluator whose checkers return results.
//
// Each result becomes one checker named after result.Name. This keeps adapter
// fixtures compact while still exercising the real Registry and Evaluator path.
// Use NewEvaluatorForTarget when a test needs custom checker behavior.
func NewEvaluatorWithResults(
	t testing.TB,
	target health.Target,
	results ...health.Result,
) *health.Evaluator {
	t.Helper()

	checks := make([]health.Checker, 0, len(results))
	for _, result := range results {
		result := result
		checks = append(checks, FuncChecker(result.Name, func(context.Context) health.Result {
			return result
		}))
	}

	return NewEvaluatorForTarget(t, target, checks...)
}

// NewDefaultTargetsEvaluator returns an evaluator for startup, live, and ready.
//
// The helper mirrors the conventional ARCORIS health target surface used
// by transport adapter default installation tests.
func NewDefaultTargetsEvaluator(t testing.TB) *health.Evaluator {
	t.Helper()

	return NewEvaluator(t, NewRegistry(
		t,
		ForTarget(health.TargetStartup, HealthyChecker("startup")),
		ForTarget(health.TargetLive, HealthyChecker("live")),
		ForTarget(health.TargetReady, HealthyChecker("ready")),
	))
}
