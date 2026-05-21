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


package runner

import (
	"sync"

	"arcoris.dev/measure/internal/reduce/core"
	"arcoris.dev/measure/internal/reduce/merge"
	"arcoris.dev/measure/internal/reduce/planner"
)

// accumulateBalancedWorkerPartials assigns each balanced range to one
// worker-local partial and merges those active partials in range order.
//
// Balanced accumulation still starts each worker-local partial from the zero
// value of T. Accumulators that need maps, slices, or other internal buffers
// must lazily initialize them before their first update.
func accumulateBalancedWorkerPartials[T any](
	n int,
	opts core.Options,
	scratch *core.Scratch[T],
	accumulate core.IndexedAccumulator[T],
	mergeFn core.Merger[T],
) (T, bool) {
	var zero T
	if n <= 0 {
		return zero, false
	}

	opts = core.NormalizeOptions(opts)
	scratch = ensureScratch(scratch)
	scratch.Ranges = planner.Balanced(n, opts, scratch.Ranges)
	ranges := scratch.Ranges

	if len(ranges) == 0 {
		return zero, false
	}
	if len(ranges) == 1 {
		return accumulateSequentiallyIndexed(n, accumulate)
	}

	partials := scratch.EnsurePartialsDirty(len(ranges))
	used := scratch.EnsureUsed(len(ranges))
	fillBalancedWorkerAccumulators(ranges, partials, used, accumulate)
	active := compactUsedPartials(partials, used)
	return merge.Merge(active, opts.MergeMode, mergeFn)
}

// fillBalancedWorkerAccumulators assigns one balanced range to each worker
// slot.
//
// The current balanced planner never returns more ranges than worker slots, so
// this path stays one-to-one and deterministic.
func fillBalancedWorkerAccumulators[T any](
	ranges []core.Range,
	partials []T,
	used []bool,
	accumulate core.IndexedAccumulator[T],
) {
	var wg sync.WaitGroup
	wg.Add(len(ranges))
	for worker, r := range ranges {
		go func() {
			var local T
			accumulate(worker, r, &local)
			partials[worker] = local
			used[worker] = true
			wg.Done()
		}()
	}
	wg.Wait()
}
