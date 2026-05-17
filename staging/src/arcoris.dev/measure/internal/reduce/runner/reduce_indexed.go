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

import "arcoris.dev/measure/internal/reduce/core"

// ReduceIndexedInto executes a reduction while exposing the worker slot to the
// mapper.
//
// The worker index is an execution slot, not necessarily the range index for
// fixed plans with more chunks than workers. Merge order remains range order
// for balanced range-local execution and active worker order for fixed or dynamic
// worker-local execution.
func ReduceIndexedInto[T any](
	n int,
	opts core.Options,
	scratch *core.Scratch[T],
	mapRange core.IndexedIntoMapper[T],
	mergeFn core.Merger[T],
) (T, bool) {
	var zero T
	if n <= 0 {
		return zero, false
	}

	opts = core.NormalizeOptions(opts)

	switch {
	case opts.Strategy == core.StrategySequential:
		return reduceSequentiallyIndexed(n, mapRange)
	case shouldReduceSequentially(n, opts):
		return reduceSequentiallyIndexed(n, mapRange)
	}

	switch opts.Strategy {
	case core.StrategyDynamicChunks:
		return reduceDynamicChunkWorkerPartials(
			n,
			opts,
			scratch,
			mapRange,
			mergeFn,
		)
	case core.StrategyFixedChunks:
		return reduceFixedChunkWorkerPartials(
			n,
			opts,
			scratch,
			mapRange,
			mergeFn,
		)
	default:
		return reduceBalancedRangePartials(
			n,
			opts,
			scratch,
			mapRange,
			mergeFn,
		)
	}
}
