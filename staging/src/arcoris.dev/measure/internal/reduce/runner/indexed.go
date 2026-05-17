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
	"arcoris.dev/measure/internal/reduce/core"
	"arcoris.dev/measure/internal/reduce/merge"
	"arcoris.dev/measure/internal/reduce/planner"
)

// DoIndexedInto executes a planned reduction while exposing the worker slot to
// the mapper.
//
// The worker index is an execution slot, not necessarily the range index for
// fixed plans with more chunks than workers. Merge order remains range order for
// planned strategies and worker order for dynamic strategy.
func DoIndexedInto[T any](n int, opts core.Options, scratch *core.Scratch[T], mapRange core.IndexedIntoMapper[T], mergeFn core.Merger[T]) (T, bool) {
	var zero T
	if n <= 0 {
		return zero, false
	}
	opts = core.NormalizeOptions(opts)
	if opts.Strategy == core.StrategyDynamic {
		return DoDynamicInto(n, opts, scratch, mapRange, mergeFn)
	}
	if shouldRunSequential(n, opts) {
		var partial T
		mapRange(0, core.Range{Start: 0, End: n}, &partial)
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
		mapRange(0, ranges[0], &partial)
		return partial, true
	}
	partials := scratch.EnsurePartialsDirty(len(ranges))
	runRangesInto(ranges, partials, opts.Workers, mapRange)
	return merge.Merge(partials, opts.MergeMode, mergeFn)
}
