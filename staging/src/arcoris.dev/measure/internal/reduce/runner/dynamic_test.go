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
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestReduceDynamicWorkerPartialsProcessesEveryIndexOnce(t *testing.T) {
	const n = 1000
	seen := make([]atomic.Int64, n)
	got, ok := reduceDynamicWorkerPartials[int](
		n,
		core.Options{Workers: 4, MinItemsPerWorker: 1, ChunkSize: 17, Strategy: core.StrategyDynamic},
		nil,
		func(_ int, r core.Range, dst *int) {
			for i := r.Start; i < r.End; i++ {
				seen[i].Add(1)
				*dst += i
			}
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("reduceDynamicWorkerPartials returned false for non-empty input")
	}
	if want := (n - 1) * n / 2; got != want {
		t.Fatalf("reduceDynamicWorkerPartials() = %d, want %d", got, want)
	}
	assertEveryIndexOnce(t, seen)
}

func TestReduceDynamicWorkerPartialsDoesNotMergeInactiveWorkers(t *testing.T) {
	var zeroMerged atomic.Int64
	got, ok := reduceDynamicWorkerPartials[nonNeutralPartial](
		10,
		core.Options{Workers: 8, MinItemsPerWorker: 1, ChunkSize: 100, Strategy: core.StrategyDynamic},
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
		t.Fatal("reduceDynamicWorkerPartials returned false for non-empty input")
	}
	if got.Value != 10 || !got.Active {
		t.Fatalf("reduceDynamicWorkerPartials() = %#v, want active value 10", got)
	}
	if zeroMerged.Load() != 0 {
		t.Fatalf("merged %d inactive zero-value partials", zeroMerged.Load())
	}
}

func TestReduceDynamicWorkerPartialsUsesChunkSize(t *testing.T) {
	var mu sync.Mutex
	var lengths []int
	got, ok := reduceDynamicWorkerPartials[int](
		10,
		core.Options{Workers: 3, MinItemsPerWorker: 1, ChunkSize: 3, Strategy: core.StrategyDynamic},
		nil,
		func(_ int, r core.Range, dst *int) {
			mu.Lock()
			lengths = append(lengths, r.Len())
			mu.Unlock()
			*dst += r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("reduceDynamicWorkerPartials returned false for non-empty input")
	}
	if got != 10 {
		t.Fatalf("reduceDynamicWorkerPartials() = %d, want 10", got)
	}
	sort.Ints(lengths)
	want := []int{1, 3, 3, 3}
	if len(lengths) != len(want) {
		t.Fatalf("chunk lengths = %v, want %v", lengths, want)
	}
	for i := range want {
		if lengths[i] != want[i] {
			t.Fatalf("chunk lengths = %v, want %v", lengths, want)
		}
	}
}

func TestReduceDynamicWorkerPartialsSequentialFallbackMapsOnce(t *testing.T) {
	var calls atomic.Int64
	got, ok := reduceDynamicWorkerPartials[int](
		10,
		core.Options{Workers: 8, MinItemsPerWorker: 100, ChunkSize: 1, Strategy: core.StrategyDynamic},
		nil,
		func(_ int, r core.Range, dst *int) {
			calls.Add(1)
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("reduceDynamicWorkerPartials returned false for non-empty input")
	}
	if got != 10 {
		t.Fatalf("reduceDynamicWorkerPartials() = %d, want 10", got)
	}
	if calls.Load() != 1 {
		t.Fatalf("mapper calls = %d, want 1", calls.Load())
	}
}

func TestReduceDynamicWorkerPartialsEmptyReturnsFalse(t *testing.T) {
	_, ok := reduceDynamicWorkerPartials[int](
		0,
		core.Options{Workers: 4, Strategy: core.StrategyDynamic},
		nil,
		func(int, core.Range, *int) { t.Fatal("mapper must not run for empty input") },
		func(*int, int) { t.Fatal("merge must not run for empty input") },
	)
	if ok {
		t.Fatal("reduceDynamicWorkerPartials returned ok for empty input")
	}
}

func TestFillWorkerPartialsDynamicMarksOnlyActiveWorkers(t *testing.T) {
	partials := make([]int, 4)
	used := make([]bool, 4)
	fillWorkerPartialsDynamic(5, 100, partials, used, func(_ int, r core.Range, dst *int) {
		*dst = r.Len()
	}, func(dst *int, src int) { *dst += src })
	active := compactUsedPartials(partials, used)
	if len(active) != 1 {
		t.Fatalf("active partials = %d, want 1", len(active))
	}
	if active[0] != 5 {
		t.Fatalf("active[0] = %d, want 5", active[0])
	}
}
