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

package eval

import (
	"context"
	"time"

	"arcoris.dev/health"
)

// Evaluate runs all checks registered for target and returns an aggregated
// report.
//
// target MUST be concrete. Invalid or non-concrete targets return a
// health.StatusUnknown report and an error classified as health.ErrInvalidTarget.
//
// A nil ctx is treated as context.Background. This mirrors defensive boundaries
// in other adjacent ARCORIS packages and avoids panics in diagnostics and tests.
//
// If target has no registered checks, Evaluate returns a health.StatusUnknown
// report. Absence of checks is not treated as healthy because health requires an
// affirmative observation.
//
// Evaluate is not fail-fast after checks are resolved. It attempts to produce
// one result per resolved check. If ctx is already canceled, each check receives
// that canceled context and should return quickly; evaluator-owned normalization
// turns cooperative interruption into unknown canceled results.
//
// Evaluate is synchronous regardless of execution policy. Parallel execution
// affects only how checks are scheduled inside this call; the caller still
// receives one complete health.Report after all scheduled checks have produced
// caller-visible Results.
func (e *Evaluator) Evaluate(ctx context.Context, target health.Target) (health.Report, error) {
	started := e.clock.Now()

	if !target.IsConcrete() {
		return unknownReport(target, started), health.InvalidTargetError{Target: target}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	checks, err := e.resolveTargetChecks(target)
	if err != nil {
		return unknownReport(target, started), err
	}
	if len(checks) == 0 {
		return unknownReport(target, started), nil
	}

	results := e.evaluateTargetChecks(ctx, target, checks)
	finished := e.clock.Now()

	return health.Report{
		Target:   target,
		Status:   health.AggregateStatus(results),
		Observed: finished,
		Duration: nonNegativeDuration(e.clock.Since(started)),
		Checks:   results,
	}, nil
}

// resolveTargetChecks returns checks for target after resolver contract checks.
func (e *Evaluator) resolveTargetChecks(target health.Target) ([]health.Checker, error) {
	set, err := e.resolver.ResolveChecks(target)
	if err != nil {
		return nil, err
	}
	if set.Target() != target {
		return nil, MismatchedResolvedTargetError{
			Requested: target,
			Resolved:  set.Target(),
		}
	}

	return set.Checks(), nil
}

// evaluateTargetChecks evaluates checks with target-specific runtime settings.
func (e *Evaluator) evaluateTargetChecks(ctx context.Context, target health.Target, checks []health.Checker) []health.Result {
	timeout := e.timeoutFor(target)
	executionPolicy := e.executionPolicyFor(target)

	return e.evaluateChecks(ctx, checks, timeout, executionPolicy)
}

// unknownReport returns a conservative report for evaluation boundaries that
// could not produce affirmative check observations.
func unknownReport(target health.Target, observed time.Time) health.Report {
	return health.Report{
		Target:   target,
		Status:   health.StatusUnknown,
		Observed: observed,
	}
}
