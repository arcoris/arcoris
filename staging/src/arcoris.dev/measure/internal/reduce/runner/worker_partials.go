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

// reduceFixedWorkerPartials executes fixed-size chunks with one accumulator per
// active worker and then merges only those active worker-local partials.
//
// Fixed worker-local execution is the default fixed strategy because a tiny
// chunk size can produce many more chunks than workers. Keeping partials at
// worker cardinality bounds scratch use and merge cost without assuming that a
// zero-value partial is a merge identity.
func reduceFixedWorkerPartials[T any](n int, opts core.Options, scratch *core.Scratch[T], mapChunk core.IndexedIntoMapper[T], mergeFn core.Merger[T]) (T, bool) {
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
	workers := activeWorkers(opts.Workers, n)
	if workers <= 1 {
		return reduceSequentiallyIndexed(n, mapChunk)
	}
	if scratch == nil {
		scratch = new(core.Scratch[T])
	}
	partials := scratch.EnsurePartialsDirty(workers)
	used := scratch.EnsureUsed(workers)
	fillWorkerPartialsFixed(n, chunk, partials, used, mapChunk, mergeFn)
	active := compactUsedPartials(partials, used)
	return merge.Merge(active, opts.MergeMode, mergeFn)
}

// fillWorkerPartialsFixed maps fixed-size chunks into worker-local partial
// slots.
//
// The partial and used slices are indexed by worker slot. A worker writes its
// partial only after it has processed at least one chunk, so idle workers leave
// dirty partial slots untouched and are removed by compactUsedPartials before
// merging. Each chunk maps into a fresh temporary partial before mergeFn folds
// it into the worker-local accumulator, so mappers may assign chunk results
// rather than manually accumulating across chunks. The atomic cursor advances
// once per chunk, never per element.
func fillWorkerPartialsFixed[T any](n int, chunk int, partials []T, used []bool, mapChunk core.IndexedIntoMapper[T], mergeFn core.Merger[T]) {
	chunks := chunkCount(n, chunk)
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
				start := chunkIndex * chunk
				end := start + chunk
				if end > n {
					end = n
				}
				chunkPartial = zero
				mapChunk(worker, core.Range{Start: start, End: end}, &chunkPartial)
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
