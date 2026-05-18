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
	"sync/atomic"
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestFillDynamicChunkWorkerAccumulatorsProcessesEveryIndexOnce(t *testing.T) {
	const n = 257
	partials := make([]int, 4)
	used := make([]bool, 4)
	seen := make([]atomic.Int64, n)
	fillDynamicChunkWorkerAccumulators(
		n,
		13,
		chunkCount(n, 13),
		partials,
		used,
		func(_ int, r core.Range, dst *int) {
			for i := r.Start; i < r.End; i++ {
				seen[i].Add(1)
				*dst += 1
			}
		},
	)
	assertEveryIndexOnce(t, seen)
	if got := sumInts(compactUsedPartials(partials, used)); got != n {
		t.Fatalf("active partial sum = %d, want %d", got, n)
	}
}

func TestAccumulateDynamicChunkWorkerPartialsSequentialFallbackCanExceedChunkSize(t *testing.T) {
	var maxLen atomic.Int64
	got, ok := accumulateDynamicChunkWorkerPartials(
		10,
		core.Options{
			Workers:           8,
			MinItemsPerWorker: 100,
			ChunkSize:         3,
			Strategy:          core.StrategyDynamicChunks,
		},
		nil,
		func(_ int, r core.Range, dst *int) {
			maxLen.Store(int64(r.Len()))
			*dst += r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok || got != 10 {
		t.Fatalf("accumulateDynamicChunkWorkerPartials() = %d ok=%v, want 10 true", got, ok)
	}
	if maxLen.Load() != 10 {
		t.Fatalf("max range len = %d, want full sequential range 10", maxLen.Load())
	}
}

func TestAccumulateDynamicChunkWorkerPartialsParallelRespectsChunkSize(t *testing.T) {
	var maxLen atomic.Int64
	got, ok := accumulateDynamicChunkWorkerPartials(
		10,
		core.Options{
			Workers:           4,
			MinItemsPerWorker: 1,
			ChunkSize:         3,
			Strategy:          core.StrategyDynamicChunks,
		},
		nil,
		func(_ int, r core.Range, dst *int) {
			for {
				old := maxLen.Load()
				if int64(r.Len()) <= old || maxLen.CompareAndSwap(old, int64(r.Len())) {
					break
				}
			}
			*dst += r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok || got != 10 {
		t.Fatalf("accumulateDynamicChunkWorkerPartials() = %d ok=%v, want 10 true", got, ok)
	}
	if maxLen.Load() > 3 {
		t.Fatalf("max range len = %d, want <= 3", maxLen.Load())
	}
}
