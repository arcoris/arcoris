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
	"testing"

	"arcoris.dev/health"
)

// EvaluatorFunc adapts a function to health.Evaluator.
type EvaluatorFunc = SourceFunc

// RegistryEvaluator evaluates checkers from a health.Registry.
//
// RegistryEvaluator is a deterministic test fixture. It does not own timeout,
// parallelism, panic normalization, or execution-policy behavior.
type RegistryEvaluator struct {
	registry *health.Registry
}

var _ health.Evaluator = (*RegistryEvaluator)(nil)

// NewEvaluator returns a deterministic registry-backed test evaluator.
func NewEvaluator(t testing.TB, r *health.Registry) *RegistryEvaluator {
	t.Helper()

	if r == nil {
		t.Fatal("healthtest.NewEvaluator() received nil registry")
	}

	return &RegistryEvaluator{registry: r}
}

// Evaluate runs registered checks for target in registry order.
func (e *RegistryEvaluator) Evaluate(ctx context.Context, target health.Target) (health.Report, error) {
	if e == nil || e.registry == nil {
		return UnknownReport(target), nil
	}

	checks := e.registry.Checks(target)
	results := make([]health.Result, 0, len(checks))
	for _, checker := range checks {
		results = append(results, checker.Check(ctx).Normalize(checker.Name(), ObservedTime))
	}

	return reportFromResults(target, results...), nil
}

// NewEvaluatorForTarget returns an evaluator with checks under target.
//
// Checks are registered in the supplied order, preserving the order that later
// reports expose to adapter tests.
func NewEvaluatorForTarget(
	t testing.TB,
	target health.Target,
	checks ...health.Checker,
) health.Evaluator {
	t.Helper()

	return NewEvaluator(t, NewRegistry(t, ForTarget(target, checks...)))
}

// NewEvaluatorWithResults returns an evaluator backed by one report.
//
// The report status is aggregated from results while preserving result order.
// Use NewEvaluatorForTarget when a test needs custom checker behavior.
func NewEvaluatorWithResults(
	t testing.TB,
	target health.Target,
	results ...health.Result,
) health.Evaluator {
	t.Helper()

	return NewEvaluatorWithReport(reportFromResults(target, results...))
}

// NewEvaluatorWithReport returns an evaluator that always returns report.
func NewEvaluatorWithReport(report health.Report) health.Evaluator {
	return NewStaticSource(report)
}

// NewEvaluatorWithReports returns an evaluator backed by target-specific reports.
func NewEvaluatorWithReports(reports ...health.Report) health.Evaluator {
	byTarget := make(map[health.Target]health.Report, len(reports))
	for _, report := range reports {
		byTarget[report.Target] = report
	}

	return NewTargetSource(byTarget)
}

// NewErrorEvaluator returns an evaluator that always fails with err.
func NewErrorEvaluator(err error) health.Evaluator {
	return NewErrorSource(err)
}

// NewSequenceEvaluator returns an evaluator backed by a report sequence.
func NewSequenceEvaluator(target health.Target, reports ...health.Report) health.Evaluator {
	return NewSequenceSource(target, reports...)
}

// NewDefaultTargetsEvaluator returns an evaluator for startup, live, and ready.
//
// The helper mirrors the conventional ARCORIS health target surface used
// by transport adapter default installation tests.
func NewDefaultTargetsEvaluator(t testing.TB) health.Evaluator {
	t.Helper()

	return NewEvaluatorWithReports(
		Report(health.TargetStartup, health.StatusHealthy, HealthyResult("startup")),
		Report(health.TargetLive, health.StatusHealthy, HealthyResult("live")),
		Report(health.TargetReady, health.StatusHealthy, HealthyResult("ready")),
	)
}

func reportFromResults(target health.Target, results ...health.Result) health.Report {
	return Report(target, aggregateStatus(results), results...)
}

func aggregateStatus(results []health.Result) health.Status {
	if len(results) == 0 {
		return health.StatusUnknown
	}

	status := health.StatusHealthy
	for _, result := range results {
		if result.Status.MoreSevereThan(status) {
			status = result.Status
		}
	}

	return status
}
