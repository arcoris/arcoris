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

import "arcoris.dev/measure/internal/reduce/core"

// AccumulateIndexedInto executes a reduction with direct worker-local
// accumulation.
//
// This is the performance-oriented API family. The accumulator may be called
// repeatedly with the same worker-local dst, so it must preserve and extend
// existing partial state. Strategy-specific paths remain private so all dispatch
// policy stays auditable here.
func AccumulateIndexedInto[T any](
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

	switch {
	case opts.Strategy == core.StrategySequential:
		return accumulateSequentiallyIndexed(n, accumulate)
	case shouldReduceSequentially(n, opts):
		return accumulateSequentiallyIndexed(n, accumulate)
	}

	switch opts.Strategy {
	case core.StrategyFixedChunks:
		return accumulateFixedChunkWorkerPartials(
			n,
			opts,
			scratch,
			accumulate,
			mergeFn,
		)
	case core.StrategyDynamicChunks:
		return accumulateDynamicChunkWorkerPartials(
			n,
			opts,
			scratch,
			accumulate,
			mergeFn,
		)
	default:
		return accumulateBalancedWorkerPartials(
			n,
			opts,
			scratch,
			accumulate,
			mergeFn,
		)
	}
}
