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

	"arcoris.dev/measure/internal/reduce"
	"arcoris.dev/measure/internal/reduce/merge"
)

// DoDynamicInto executes a dynamic chunk-claiming reduction.
//
// Dynamic execution is intended for variable-cost chunks. Each worker keeps one
// local partial result and repeatedly maps chunks claimed from an atomic cursor.
// Merge order is worker-slot order, not input-range order, so floating-point
// callers should expect different grouping from static execution.
func DoDynamicInto[T any](n int, opts reduce.Options, scratch *reduce.Scratch[T], mapChunk reduce.IndexedIntoMapper[T], mergeFn reduce.Merger[T]) (T, bool) {
	var zero T
	if n <= 0 {
		return zero, false
	}
	opts = reduce.NormalizeOptions(opts)
	if shouldRunSequential(n, opts) {
		var partial T
		mapChunk(0, reduce.Range{Start: 0, End: n}, &partial)
		return partial, true
	}
	chunk := opts.ChunkSize
	if chunk <= 0 {
		chunk = reduce.DefaultChunkSize
	}
	workers := activeWorkers(opts.Workers, chunkCount(n, chunk))
	if workers <= 1 {
		var partial T
		mapChunk(0, reduce.Range{Start: 0, End: n}, &partial)
		return partial, true
	}
	if scratch == nil {
		scratch = new(reduce.Scratch[T])
	}
	partials := scratch.EnsurePartialsDirty(workers)
	var next atomic.Int64
	var wg sync.WaitGroup
	wg.Add(workers)
	for worker := 0; worker < workers; worker++ {
		worker := worker
		go func() {
			var local T
			for {
				start := int(next.Add(int64(chunk)) - int64(chunk))
				if start >= n {
					break
				}
				end := start + chunk
				if end > n {
					end = n
				}
				mapChunk(worker, reduce.Range{Start: start, End: end}, &local)
			}
			partials[worker] = local
			wg.Done()
		}()
	}
	wg.Wait()
	return merge.Merge(partials, opts.MergeMode, mergeFn)
}

// chunkCount returns the number of chunks needed to cover n items with chunk
// size chunk. The formula avoids n+chunk overflow for large inputs.
func chunkCount(n, chunk int) int {
	if n <= 0 {
		return 0
	}
	return 1 + (n-1)/chunk
}
