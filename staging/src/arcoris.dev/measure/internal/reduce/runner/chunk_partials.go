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

	"arcoris.dev/measure/internal/reduce/core"
	"arcoris.dev/measure/internal/reduce/merge"
)

// reduceFixedChunkWorkerPartials maps fixed-size chunks to complete chunk
// partials, folds those chunk partials into worker-local partials, and merges
// active workers in deterministic worker order.
//
// This is the Reduce-family fixed-chunk path: mappers may assign each chunk
// partial because the runner owns per-chunk folding with mergeFn. Chunk
// ownership is balanced and contiguous; no atomic cursor is used.
func reduceFixedChunkWorkerPartials[T any](
	n int,
	opts core.Options,
	scratch *core.Scratch[T],
	mapChunk core.IndexedIntoMapper[T],
	mergeFn core.Merger[T],
) (T, bool) {
	var zero T
	if n <= 0 {
		return zero, false
	}

	opts = core.NormalizeOptions(opts)
	if shouldReduceSequentially(n, opts) {
		return reduceSequentiallyIndexed(n, mapChunk)
	}

	work := chunkWorkers(n, opts)
	if work.workers <= 1 {
		return reduceSequentiallyIndexed(n, mapChunk)
	}

	scratch = ensureScratch(scratch)
	partials := scratch.EnsurePartialsDirty(work.workers)
	used := scratch.EnsureUsed(work.workers)
	fillFixedChunkWorkerPartials(n, work.size, work.count, partials, used, mapChunk, mergeFn)
	active := compactUsedPartials(partials, used)
	return merge.Merge(active, opts.MergeMode, mergeFn)
}

// fillFixedChunkWorkerPartials maps statically assigned fixed chunks into
// worker-local partial slots.
//
// The partial and used slices are indexed by worker slot. A worker writes its
// partial only after it has processed at least one chunk, so idle workers leave
// dirty partial slots untouched and are removed by compactUsedPartials before
// merging. Each chunk maps into a fresh temporary partial before mergeFn folds
// it into the worker-local accumulator, so mappers may assign chunk results
// rather than manually accumulating across chunks. That extra mergeFn call per
// chunk is the main reason buffer-backed algorithms should prefer the
// Accumulate family when they can update worker-local state directly.
func fillFixedChunkWorkerPartials[T any](
	n int,
	chunk int,
	chunks int,
	partials []T,
	used []bool,
	mapChunk core.IndexedIntoMapper[T],
	mergeFn core.Merger[T],
) {
	var wg sync.WaitGroup
	wg.Add(len(partials))
	for worker := range partials {
		go func() {
			var local T
			var chunkPartial T
			var zero T
			active := false
			startChunk, endChunk := fixedChunkBlock(chunks, len(partials), worker)
			for chunkIndex := startChunk; chunkIndex < endChunk; chunkIndex++ {
				chunkPartial = zero
				mapChunk(worker, chunkRange(n, chunk, chunkIndex), &chunkPartial)
				if active {
					mergeFn(&local, chunkPartial)
				} else {
					local = chunkPartial
					active = true
				}
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
