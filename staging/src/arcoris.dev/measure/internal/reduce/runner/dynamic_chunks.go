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
)

// reduceDynamicChunkWorkerPartials executes dynamic chunk-claiming reduction with
// one partial per active worker.
//
// Dynamic execution is intended for variable-cost chunks. Each worker keeps one
// local partial result and repeatedly maps chunks claimed from an atomic cursor.
// Only workers that claimed at least one chunk are merged, because generic
// partial values do not have a guaranteed zero-value identity. Merge order is
// active worker-slot order, not input-range order, so floating-point callers
// should expect different grouping from balanced execution.
func reduceDynamicChunkWorkerPartials[T any](n int, opts core.Options, scratch *core.Scratch[T], mapChunk core.IndexedIntoMapper[T], mergeFn core.Merger[T]) (T, bool) {
	var zero T
	if n <= 0 {
		return zero, false
	}
	opts = core.NormalizeOptions(opts)
	if shouldReduceSequentially(n, opts) {
		return reduceSequentiallyIndexed(n, mapChunk)
	}
	chunk := opts.ChunkSize
	if chunk <= 0 {
		chunk = core.DefaultChunkSize
	}
	chunks := chunkCount(n, chunk)
	workers := activeWorkers(opts.Workers, chunks)
	if workers <= 1 {
		return reduceSequentiallyIndexed(n, mapChunk)
	}
	if scratch == nil {
		scratch = new(core.Scratch[T])
	}
	partials := scratch.EnsurePartialsDirty(workers)
	used := scratch.EnsureUsed(workers)
	fillDynamicChunkWorkerPartials(n, chunk, chunks, partials, used, mapChunk, mergeFn)
	active := compactUsedPartials(partials, used)
	return merge.Merge(active, opts.MergeMode, mergeFn)
}

// fillDynamicChunkWorkerPartials maps dynamically claimed chunks into worker-local
// partial slots.
//
// The cursor advances by chunk index, so each atomic operation assigns one
// whole chunk. That keeps scheduling cost out of the per-element loop and lets
// idle worker slots remain inactive instead of contributing dirty or zero-value
// partials to the final merge. Each chunk maps into a temporary partial and is
// folded into the worker-local accumulator with mergeFn.
func fillDynamicChunkWorkerPartials[T any](n int, chunk int, chunks int, partials []T, used []bool, mapChunk core.IndexedIntoMapper[T], mergeFn core.Merger[T]) {
	var next atomic.Int64
	var wg sync.WaitGroup
	wg.Add(len(partials))
	for worker := range partials {
		go func() {
			var local T
			var chunkPartial T
			var zero T
			active := false
			for {
				chunkIndex := int(next.Add(1) - 1)
				if chunkIndex >= chunks {
					break
				}
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

// chunkCount returns the number of chunks needed to cover n items with chunk
// size chunk. The formula avoids n+chunk overflow for large inputs.
func chunkCount(n, chunk int) int {
	if n <= 0 || chunk <= 0 {
		return 0
	}
	return 1 + (n-1)/chunk
}
