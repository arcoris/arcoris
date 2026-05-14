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

package probe

import (
	"context"
	"time"

	"arcoris.dev/health"
)

// runCycle evaluates configured targets once in configured order.
//
// A cycle is synchronous and non-overlapping. Runner deliberately evaluates
// targets sequentially; check-level parallelism belongs to the evaluator source.
//
// The store update no longer receives an explicit timestamp. The per-target
// snapshot.Store records Updated at commit time using the same clock configured
// on Runner. This keeps write metadata ownership inside the storage primitive
// and keeps runCycle focused on evaluator orchestration.
func (r *Runner) runCycle(ctx context.Context) {
	for _, target := range r.targets {
		report, ok := r.evaluateTarget(ctx, target)
		if !ok {
			return
		}
		if ctx.Err() != nil {
			return
		}

		r.store.update(target, report)
	}
}

// evaluateTarget evaluates one target and reports whether the result should be
// stored.
//
// If the runner context stops before or during evaluation, the result is not
// stored. This prevents normal runner shutdown from overwriting the last useful
// cached snapshot with cancellation artifacts.
//
// Evaluator errors are converted into an unknown report because probe is a cache
// layer over health reports, not an error-reporting transport. Context
// cancellation remains special: normal shutdown should not publish a synthetic
// canceled observation.
func (r *Runner) evaluateTarget(ctx context.Context, target health.Target) (health.Report, bool) {
	if ctx.Err() != nil {
		return health.Report{}, false
	}

	report, err := r.evaluator.Evaluate(ctx, target)
	if ctx.Err() != nil {
		return health.Report{}, false
	}
	if err != nil {
		return unknownReport(target, r.clock.Now()), true
	}

	return report, true
}

// unknownReport returns a defensive unknown report for unexpected evaluator
// errors.
//
// Targets are validated during Runner construction, so evaluator errors should
// be rare. The runner stores an unknown report instead of stopping the loop or
// exposing a raw internal error through Snapshot.
//
// observed is the health report observation time, not the cache commit time. The
// cache commit time is supplied later by snapshot.Store and exposed as
// Snapshot.Updated.
func unknownReport(target health.Target, observed time.Time) health.Report {
	return health.Report{
		Target:   target,
		Status:   health.StatusUnknown,
		Observed: observed,
	}
}
