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

import (
	"sync/atomic"
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestReduceIndexedIntoPassesWorkerSlots(t *testing.T) {
	var scratch core.Scratch[int]
	_, ok := ReduceIndexedInto[int](
		1000,
		core.Options{Workers: 4, MinItemsPerWorker: 100, Strategy: core.StrategyBalanced},
		&scratch,
		func(worker int, r core.Range, dst *int) {
			*dst = worker + r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("ReduceIndexedInto returned false for non-empty input")
	}
	if len(scratch.Partials) != 4 {
		t.Fatalf("partials = %d, want 4", len(scratch.Partials))
	}
	for worker, partial := range scratch.Partials {
		if partial < worker {
			t.Fatalf("partial[%d] = %d, expected worker contribution", worker, partial)
		}
	}
}

func TestReduceIndexedIntoBalancedProcessesEveryIndexOnce(t *testing.T) {
	assertReduceIndexedCoversInput(t, core.Options{
		Workers:           4,
		MinItemsPerWorker: 10,
		Strategy:          core.StrategyBalanced,
	})
}

func TestReduceIndexedIntoFixedChunksProcessesEveryIndexOnce(t *testing.T) {
	assertReduceIndexedCoversInput(t, core.Options{
		Workers:           3,
		MinItemsPerWorker: 1,
		ChunkSize:         7,
		Strategy:          core.StrategyFixedChunks,
	})
}

func TestReduceIndexedIntoDynamicChunksProcessesEveryIndexOnce(t *testing.T) {
	assertReduceIndexedCoversInput(t, core.Options{
		Workers:           5,
		MinItemsPerWorker: 1,
		ChunkSize:         11,
		Strategy:          core.StrategyDynamicChunks,
	})
}

func TestReduceIndexedIntoBoundsWorkerSlotsForFixedChunks(t *testing.T) {
	var maxWorker atomic.Int64
	_, ok := ReduceIndexedInto[int](
		100,
		core.Options{
			Workers:           2,
			MinItemsPerWorker: 1,
			ChunkSize:         10,
			Strategy:          core.StrategyFixedChunks,
		},
		nil,
		func(worker int, r core.Range, dst *int) {
			for {
				old := maxWorker.Load()
				if int64(worker) <= old || maxWorker.CompareAndSwap(old, int64(worker)) {
					break
				}
			}
			*dst += r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("ReduceIndexedInto returned false for non-empty input")
	}
	if got := maxWorker.Load(); got >= 2 {
		t.Fatalf("max worker = %d, want < 2", got)
	}
}

func assertReduceIndexedCoversInput(t *testing.T, opts core.Options) {
	t.Helper()
	const n = 257
	seen := make([]atomic.Int64, n)
	got, ok := ReduceIndexedInto[int](
		n,
		opts,
		nil,
		func(_ int, r core.Range, dst *int) {
			for i := r.Start; i < r.End; i++ {
				seen[i].Add(1)
				*dst += 1
			}
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("ReduceIndexedInto returned false for non-empty input")
	}
	if got != n {
		t.Fatalf("ReduceIndexedInto() = %d, want %d", got, n)
	}
	assertEveryIndexOnce(t, seen)
}
