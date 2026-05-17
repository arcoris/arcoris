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

// Do maps each planned range to a complete partial result and merges partials.
//
// Do is the simplest entry point when the mapper naturally returns a complete
// partial value. It allocates temporary scratch as needed; repeated callers
// should prefer DoInto with caller-owned Scratch or Runner.
func Do[T any](n int, opts core.Options, mapRange core.Mapper[T], mergeFn core.Merger[T]) (T, bool) {
	return DoInto(n, opts, nil, func(r core.Range, dst *T) { *dst = mapRange(r) }, mergeFn)
}

// DoInto executes a planned reduction using caller-owned scratch.
//
// DoInto handles sequential, static, fixed, and dynamic strategies. For n <= 0
// it returns the zero value and false. For non-empty input it returns the merged
// partial and true. The scratch argument may be nil, but passing one lets callers
// reuse range and partial buffers across calls.
func DoInto[T any](n int, opts core.Options, scratch *core.Scratch[T], mapRange core.IntoMapper[T], mergeFn core.Merger[T]) (T, bool) {
	var zero T
	if n <= 0 {
		return zero, false
	}
	opts = core.NormalizeOptions(opts)
	if opts.Strategy == core.StrategyDynamic {
		return DoDynamicInto(n, opts, scratch, func(_ int, r core.Range, dst *T) { mapRange(r, dst) }, mergeFn)
	}
	if shouldRunSequential(n, opts) {
		var partial T
		mapRange(core.Range{Start: 0, End: n}, &partial)
		return partial, true
	}
	if scratch == nil {
		scratch = new(core.Scratch[T])
	}
	scratch.Ranges = planner.Plan(n, opts, scratch.Ranges)
	ranges := scratch.Ranges
	if len(ranges) == 0 {
		return zero, false
	}
	if len(ranges) == 1 {
		var partial T
		mapRange(ranges[0], &partial)
		return partial, true
	}
	partials := scratch.EnsurePartialsDirty(len(ranges))
	runRangesInto(ranges, partials, opts.Workers, func(_ int, r core.Range, dst *T) { mapRange(r, dst) })
	return merge.Merge(partials, opts.MergeMode, mergeFn)
}

// runRangesInto stores one complete partial per range while capping goroutine
// count at workers.
//
// The partial slice remains indexed by range, so merge order stays deterministic
// even when a small worker set consumes many fixed chunks.
func runRangesInto[T any](ranges []core.Range, partials []T, workers int, mapRange core.IndexedIntoMapper[T]) {
	workers = activeWorkers(workers, len(ranges))
	if workers == len(ranges) {
		runRangesOneToOne(ranges, partials, mapRange)
		return
	}

	var next atomic.Int64
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		worker := worker
		go func() {
			for {
				i := int(next.Add(1) - 1)
				if i >= len(ranges) {
					break
				}
				clearPartial(partials, i)
				mapRange(worker, ranges[i], &partials[i])
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// runRangesOneToOne is the low-overhead path for static plans where each range
// has its own worker goroutine. Workers compute into stack-local partials to
// avoid false sharing while the hot loop runs.
func runRangesOneToOne[T any](ranges []core.Range, partials []T, mapRange core.IndexedIntoMapper[T]) {
	var wg sync.WaitGroup
	wg.Add(len(ranges))
	for worker, r := range ranges {
		worker, r := worker, r
		go func() {
			var local T
			mapRange(worker, r, &local)
			partials[worker] = local
			wg.Done()
		}()
	}
	wg.Wait()
}

// clearPartial removes stale scratch state before a mapper writes directly into
// a reused partial slot.
func clearPartial[T any](partials []T, i int) {
	var zero T
	partials[i] = zero
}

// activeWorkers caps worker count to available jobs while preserving a valid
// positive worker count for non-empty job sets.
func activeWorkers(workers, jobs int) int {
	if jobs <= 0 {
		return 0
	}
	if workers <= 0 || workers > jobs {
		return jobs
	}
	return workers
}

// shouldRunSequential reports whether n and opts should bypass goroutine
// startup and execute as one full-input range.
func shouldRunSequential(n int, opts core.Options) bool {
	return opts.Strategy == core.StrategySequential || opts.Workers <= 1 || n/2 < opts.MinItemsPerWorker
}
