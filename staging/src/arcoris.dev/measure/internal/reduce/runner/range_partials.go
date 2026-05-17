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

package runner

import (
	"sync"
	"sync/atomic"

	"arcoris.dev/measure/internal/reduce/core"
	"arcoris.dev/measure/internal/reduce/merge"
	"arcoris.dev/measure/internal/reduce/planner"
)

// reduceBalancedRangePartials plans balanced contiguous ranges, fills one partial
// per planned range, and merges those partials in range order.
//
// Range-local partial execution is deterministic and matches StrategyBalanced's
// balanced planning model. It is intentionally not the chunk-strategy default:
// fine-grained chunk plans can create far more ranges than workers, which would
// inflate scratch storage and merge work.
func reduceBalancedRangePartials[T any](n int, opts core.Options, scratch *core.Scratch[T], mapRange core.IndexedIntoMapper[T], mergeFn core.Merger[T]) (T, bool) {
	var zero T
	if n <= 0 {
		return zero, false
	}
	opts = core.NormalizeOptions(opts)
	if scratch == nil {
		scratch = new(core.Scratch[T])
	}
	scratch.Ranges = planner.Balanced(n, opts, scratch.Ranges)
	ranges := scratch.Ranges
	if len(ranges) == 0 {
		return zero, false
	}
	if len(ranges) == 1 {
		var partial T
		mapRange(0, ranges[0], &partial)
		return partial, true
	}
	partials := scratch.EnsurePartialsDirty(len(ranges))
	fillRangePartialsQueued(ranges, partials, opts.Workers, mapRange)
	return merge.Merge(partials, opts.MergeMode, mergeFn)
}

// fillRangePartialsQueued maps ranges into range-indexed partial slots while
// capping goroutine count.
//
// Range-local partial execution is deterministic because partials stay indexed
// by range order, but it can allocate and merge too many partials for
// fine-grained chunk plans. The queued path is useful when a deterministic
// balanced plan has more ranges than available workers.
func fillRangePartialsQueued[T any](ranges []core.Range, partials []T, workers int, mapRange core.IndexedIntoMapper[T]) {
	workers = activeWorkers(workers, len(ranges))
	if workers == len(ranges) {
		fillRangePartialsOneToOne(ranges, partials, mapRange)
		return
	}

	var next atomic.Int64
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		go func() {
			for {
				i := int(next.Add(1) - 1)
				if i >= len(ranges) {
					break
				}
				var local T
				mapRange(worker, ranges[i], &local)
				partials[i] = local
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// fillRangePartialsOneToOne maps each range in its own goroutine.
//
// Each goroutine computes into a stack-local partial before publishing to the
// range-indexed partial slot. That keeps writes out of shared storage while the
// mapper's hot loop is running and preserves deterministic range-order merge
// input.
func fillRangePartialsOneToOne[T any](ranges []core.Range, partials []T, mapRange core.IndexedIntoMapper[T]) {
	var wg sync.WaitGroup
	wg.Add(len(ranges))
	for worker, r := range ranges {
		go func() {
			var local T
			mapRange(worker, r, &local)
			partials[worker] = local
			wg.Done()
		}()
	}
	wg.Wait()
}
