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

func TestFillFixedChunkWorkerPartialsProcessesEveryIndexOnce(t *testing.T) {
	const n = 257
	partials := make([]int, 4)
	used := make([]bool, 4)
	seen := make([]atomic.Int64, n)
	fillFixedChunkWorkerPartials(
		n,
		13,
		chunkCount(n, 13),
		partials,
		used,
		func(_ int, r core.Range, dst *int) {
			count := 0
			for i := r.Start; i < r.End; i++ {
				seen[i].Add(1)
				count++
			}
			*dst = count
		},
		func(dst *int, src int) { *dst += src },
	)
	assertEveryIndexOnce(t, seen)
	if got := sumInts(compactUsedPartials(partials, used)); got != n {
		t.Fatalf("active partial sum = %d, want %d", got, n)
	}
}

func TestReduceFixedChunkWorkerPartialsDoesNotMergeInactiveWorkers(t *testing.T) {
	var zeroMerged atomic.Int64
	got, ok := reduceFixedChunkWorkerPartials[nonNeutralPartial](
		10,
		core.Options{
			Workers:           8,
			MinItemsPerWorker: 1,
			ChunkSize:         100,
			Strategy:          core.StrategyFixedChunks,
		},
		nil,
		func(_ int, r core.Range, dst *nonNeutralPartial) {
			dst.Value += r.Len()
			dst.Active = true
		},
		func(dst *nonNeutralPartial, src nonNeutralPartial) {
			if !src.Active {
				zeroMerged.Add(1)
			}
			dst.Value += src.Value
			dst.Active = dst.Active || src.Active
		},
	)
	if !ok {
		t.Fatal("reduceFixedChunkWorkerPartials returned false for non-empty input")
	}
	if got.Value != 10 || !got.Active {
		t.Fatalf("reduceFixedChunkWorkerPartials() = %#v, want active value 10", got)
	}
	if zeroMerged.Load() != 0 {
		t.Fatalf("merged %d inactive zero-value partials", zeroMerged.Load())
	}
}

func TestReduceFixedChunkWorkerPartialsUsesFewerPartialsThanChunks(t *testing.T) {
	var scratch core.Scratch[int]
	got, ok := reduceFixedChunkWorkerPartials[int](
		1000,
		core.Options{
			Workers:           4,
			MinItemsPerWorker: 1,
			ChunkSize:         10,
			Strategy:          core.StrategyFixedChunks,
		},
		&scratch,
		func(_ int, r core.Range, dst *int) { *dst += r.Len() },
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("reduceFixedChunkWorkerPartials returned false for non-empty input")
	}
	if got != 1000 {
		t.Fatalf("reduceFixedChunkWorkerPartials() = %d, want 1000", got)
	}
	if len(scratch.Partials) != 4 {
		t.Fatalf("partials = %d, want one slot per worker", len(scratch.Partials))
	}
	if chunks := chunkCount(1000, 10); len(scratch.Partials) >= chunks {
		t.Fatalf(
			"partials = %d, chunks = %d; want fewer partials than chunks",
			len(scratch.Partials),
			chunks,
		)
	}
}

func TestFillFixedChunkWorkerPartialsAccumulatesMultipleChunksPerWorker(t *testing.T) {
	partials := make([]int, 2)
	used := make([]bool, 2)
	fillFixedChunkWorkerPartials(
		100,
		1,
		chunkCount(100, 1),
		partials,
		used,
		func(_ int, r core.Range, dst *int) {
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	active := compactUsedPartials(partials, used)
	if len(active) == 0 {
		t.Fatal("expected at least one active worker")
	}
	for _, partial := range active {
		if partial > 1 {
			return
		}
	}
	t.Fatalf("active partials = %#v, want at least one worker to accumulate multiple chunks", active)
}

func TestReduceFixedChunkWorkerPartialsSequentialFallbackCanExceedChunkSize(t *testing.T) {
	var maxLen atomic.Int64
	got, ok := reduceFixedChunkWorkerPartials[int](
		10,
		core.Options{
			Workers:           8,
			MinItemsPerWorker: 100,
			ChunkSize:         3,
			Strategy:          core.StrategyFixedChunks,
		},
		nil,
		func(_ int, r core.Range, dst *int) {
			maxLen.Store(int64(r.Len()))
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok || got != 10 {
		t.Fatalf("reduceFixedChunkWorkerPartials() = %d ok=%v, want 10 true", got, ok)
	}
	if maxLen.Load() != 10 {
		t.Fatalf("max range len = %d, want full sequential range 10", maxLen.Load())
	}
}

func TestReduceFixedChunkWorkerPartialsParallelRespectsChunkSize(t *testing.T) {
	var maxLen atomic.Int64
	got, ok := reduceFixedChunkWorkerPartials[int](
		10,
		core.Options{
			Workers:           4,
			MinItemsPerWorker: 1,
			ChunkSize:         3,
			Strategy:          core.StrategyFixedChunks,
		},
		nil,
		func(_ int, r core.Range, dst *int) {
			for {
				old := maxLen.Load()
				if int64(r.Len()) <= old || maxLen.CompareAndSwap(old, int64(r.Len())) {
					break
				}
			}
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok || got != 10 {
		t.Fatalf("reduceFixedChunkWorkerPartials() = %d ok=%v, want 10 true", got, ok)
	}
	if maxLen.Load() > 3 {
		t.Fatalf("max range len = %d, want <= 3", maxLen.Load())
	}
}

func sumInts(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}
