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

// accumulateBalancedWorkerPartials assigns each balanced range to one
// worker-local partial and merges those active partials in range/worker order.
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

// accumulateFixedChunkWorkerPartials accumulates statically assigned chunks
// directly into worker-local partials.
func accumulateFixedChunkWorkerPartials[T any](
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
	if shouldReduceSequentially(n, opts) {
		return accumulateSequentiallyIndexed(n, accumulate)
	}

	work := chunkWorkers(n, opts)
	if work.workers <= 1 {
		return accumulateSequentiallyIndexed(n, accumulate)
	}

	scratch = ensureScratch(scratch)
	partials := scratch.EnsurePartialsDirty(work.workers)
	used := scratch.EnsureUsed(work.workers)
	fillFixedChunkWorkerAccumulators(n, work.size, work.count, partials, used, accumulate)
	active := compactUsedPartials(partials, used)
	return merge.Merge(active, opts.MergeMode, mergeFn)
}

// accumulateDynamicChunkWorkerPartials accumulates dynamically claimed chunks
// directly into worker-local partials.
func accumulateDynamicChunkWorkerPartials[T any](
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
	if shouldReduceSequentially(n, opts) {
		return accumulateSequentiallyIndexed(n, accumulate)
	}

	work := chunkWorkers(n, opts)
	if work.workers <= 1 {
		return accumulateSequentiallyIndexed(n, accumulate)
	}

	scratch = ensureScratch(scratch)
	partials := scratch.EnsurePartialsDirty(work.workers)
	used := scratch.EnsureUsed(work.workers)
	fillDynamicChunkWorkerAccumulators(n, work.size, work.count, partials, used, accumulate)
	active := compactUsedPartials(partials, used)
	return merge.Merge(active, opts.MergeMode, mergeFn)
}

// fillBalancedWorkerAccumulators assigns one balanced range to each worker slot.
//
// Balanced accumulation preserves the same range order as the planned ranges,
// but it calls the Accumulate-family callback directly into a worker-local
// partial instead of asking for a complete temporary partial first.
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

// fillFixedChunkWorkerAccumulators calls accumulate directly on each
// worker-local partial for a deterministic contiguous block of chunks.
func fillFixedChunkWorkerAccumulators[T any](
	n int,
	chunk int,
	chunks int,
	partials []T,
	used []bool,
	accumulate core.IndexedAccumulator[T],
) {
	var wg sync.WaitGroup
	wg.Add(len(partials))
	for worker := range partials {
		go func() {
			var local T
			active := false
			startChunk, endChunk := fixedChunkBlock(chunks, len(partials), worker)
			for chunkIndex := startChunk; chunkIndex < endChunk; chunkIndex++ {
				accumulate(worker, chunkRange(n, chunk, chunkIndex), &local)
				active = true
			}
			if active {
				partials[worker] = local
				used[worker] = true
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// fillDynamicChunkWorkerAccumulators uses one atomic operation per chunk and
// accumulates claimed chunks directly into worker-local partials.
func fillDynamicChunkWorkerAccumulators[T any](
	n int,
	chunk int,
	chunks int,
	partials []T,
	used []bool,
	accumulate core.IndexedAccumulator[T],
) {
	var next atomic.Int64
	var wg sync.WaitGroup
	wg.Add(len(partials))
	for worker := range partials {
		go func() {
			var local T
			active := false
			for {
				chunkIndex := int(next.Add(1) - 1)
				if chunkIndex >= chunks {
					break
				}
				accumulate(worker, chunkRange(n, chunk, chunkIndex), &local)
				active = true
			}
			if active {
				partials[worker] = local
				used[worker] = true
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
