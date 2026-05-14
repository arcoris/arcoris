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
	"sync"
	"time"
)

// evaluateChecks evaluates checks with the supplied execution policy.
//
// The helper owns multi-check scheduling only. It deliberately delegates every
// individual check to evaluateCheck so timeout handling, panic recovery,
// cancellation normalization, mismatched-name normalization, invalid-reason
// normalization, observed timestamps, and durations remain centralized in
// evaluator_run.go.
func (e *Evaluator) evaluateChecks(
	ctx context.Context,
	checks []Checker,
	timeout time.Duration,
	executionPolicy ExecutionPolicy,
) []Result {
	executionPolicy = executionPolicy.Normalize()

	switch executionPolicy.Mode {
	case ExecutionParallel:
		return e.evaluateChecksParallel(ctx, checks, timeout, executionPolicy.MaxConcurrency)
	default:
		return e.evaluateChecksSequential(ctx, checks, timeout)
	}
}

// evaluateChecksSequential evaluates checks one by one in registry order.
func (e *Evaluator) evaluateChecksSequential(
	ctx context.Context,
	checks []Checker,
	timeout time.Duration,
) []Result {
	results := make([]Result, 0, len(checks))
	for _, chk := range checks {
		results = append(results, e.evaluateCheck(ctx, chk, timeout))
	}

	return results
}

// evaluateChecksParallel evaluates checks with bounded concurrency while
// preserving registry order in the returned results.
//
// The implementation preallocates the result slice and assigns exactly one index
// from exactly one goroutine. It never appends concurrently. This preserves
// deterministic Report.Checks order even when checks complete out of order.
//
// Parallel execution is intentionally not fail-fast. Every registered check is
// still evaluated through evaluateCheck so timeout, panic, cancellation,
// mismatched-name, invalid-reason, timestamp, and duration normalization remains
// centralized and consistent with sequential execution.
func (e *Evaluator) evaluateChecksParallel(
	ctx context.Context,
	checks []Checker,
	timeout time.Duration,
	maxConcurrency int,
) []Result {
	if len(checks) == 0 {
		return nil
	}
	if maxConcurrency <= 1 || len(checks) == 1 {
		return e.evaluateChecksSequential(ctx, checks, timeout)
	}
	if maxConcurrency > len(checks) {
		maxConcurrency = len(checks)
	}

	results := make([]Result, len(checks))
	sem := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	for i, chk := range checks {
		// Acquire before spawning the goroutine so the evaluator bounds both
		// active checks and spawned goroutines. Go 1.22+ gives range variables
		// per-iteration scope, so defensive loop-variable shadowing is not needed.
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			results[i] = e.evaluateCheck(ctx, chk, timeout)
		}()
	}

	wg.Wait()
	return results
}

// aggregateStatus returns the most severe status in results.
//
// Empty result slices are treated as StatusUnknown. Evaluate handles the normal
// no-check case before calling this helper, but the defensive fallback keeps the
// aggregation boundary conservative for tests and future internal callers.
func aggregateStatus(results []Result) Status {
	if len(results) == 0 {
		return StatusUnknown
	}

	status := StatusHealthy
	for _, res := range results {
		if res.Status.MoreSevereThan(status) {
			status = res.Status
		}
	}

	return status
}
