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

func TestReduceIntoEmptyReturnsFalse(t *testing.T) {
	got, ok := ReduceInto[int](
		0,
		core.Options{},
		nil,
		func(core.Range, *int) { t.Fatal("mapper must not run for empty input") },
		func(*int, int) { t.Fatal("merge must not run for empty input") },
	)
	if ok {
		t.Fatal("ReduceInto returned ok for empty input")
	}
	if got != 0 {
		t.Fatalf("ReduceInto() = %d, want zero value", got)
	}
}

func TestReduceIntoSequentialMapsExactlyOnce(t *testing.T) {
	var calls atomic.Int64
	got, ok := ReduceInto[int](
		17,
		core.Options{Workers: 8, Strategy: core.StrategySequential},
		nil,
		func(r core.Range, dst *int) {
			calls.Add(1)
			if r.Start != 0 || r.End != 17 {
				t.Errorf("sequential range = [%d,%d), want [0,17)", r.Start, r.End)
			}
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("ReduceInto returned false for non-empty input")
	}
	if got != 17 {
		t.Fatalf("ReduceInto() = %d, want 17", got)
	}
	if calls.Load() != 1 {
		t.Fatalf("mapper calls = %d, want 1", calls.Load())
	}
}

func TestReduceIntoSmallInputFallsBackToSequential(t *testing.T) {
	var calls atomic.Int64
	got, ok := ReduceInto[int](
		10,
		core.Options{Workers: 8, MinItemsPerWorker: 100, Strategy: core.StrategyBalanced},
		nil,
		func(r core.Range, dst *int) {
			calls.Add(1)
			if r.Start != 0 || r.End != 10 {
				t.Errorf("fallback range = [%d,%d), want [0,10)", r.Start, r.End)
			}
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("ReduceInto returned false for non-empty input")
	}
	if got != 10 {
		t.Fatalf("ReduceInto() = %d, want 10", got)
	}
	if calls.Load() != 1 {
		t.Fatalf("mapper calls = %d, want 1", calls.Load())
	}
}

func TestReduceIntoBalancedProcessesEveryIndexOnce(t *testing.T) {
	const n = 1000
	seen := make([]atomic.Int64, n)
	got, ok := ReduceInto[int](
		n,
		core.Options{Workers: 4, MinItemsPerWorker: 100, Strategy: core.StrategyBalanced},
		nil,
		func(r core.Range, dst *int) {
			for i := r.Start; i < r.End; i++ {
				seen[i].Add(1)
				*dst += i
			}
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("ReduceInto returned false for balanced input")
	}
	if want := (n - 1) * n / 2; got != want {
		t.Fatalf("ReduceInto() = %d, want %d", got, want)
	}
	assertEveryIndexOnce(t, seen)
}

func TestReduceIntoFixedChunksUsesWorkerLocalPartials(t *testing.T) {
	const n = 1000
	var scratch core.Scratch[int]
	got, ok := ReduceInto[int](
		n,
		core.Options{Workers: 4, MinItemsPerWorker: 1, ChunkSize: 10, Strategy: core.StrategyFixedChunks},
		&scratch,
		func(r core.Range, dst *int) {
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("ReduceInto returned false for fixed input")
	}
	if got != n {
		t.Fatalf("ReduceInto() = %d, want %d", got, n)
	}
	if len(scratch.Partials) != 4 {
		t.Fatalf("partials = %d, want one slot per worker", len(scratch.Partials))
	}
	if chunks := chunkCount(n, 10); len(scratch.Partials) >= chunks {
		t.Fatalf("partials = %d, chunks = %d; fixed execution should not allocate per chunk", len(scratch.Partials), chunks)
	}
}

func TestReduceIntoDynamicChunksSkipsIdlePartials(t *testing.T) {
	var zeroMerged atomic.Int64
	got, ok := ReduceInto[nonNeutralPartial](
		10,
		core.Options{Workers: 8, MinItemsPerWorker: 1, ChunkSize: 100, Strategy: core.StrategyDynamicChunks},
		nil,
		func(r core.Range, dst *nonNeutralPartial) {
			dst.Value = r.Len()
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
		t.Fatal("ReduceInto returned false for dynamic input")
	}
	if got.Value != 10 || !got.Active {
		t.Fatalf("ReduceInto() = %#v, want active value 10", got)
	}
	if zeroMerged.Load() != 0 {
		t.Fatalf("merged %d inactive zero-value partials", zeroMerged.Load())
	}
}

// assertEveryIndexOnce verifies the core coverage invariant shared by runner
// strategy tests.
func assertEveryIndexOnce(t *testing.T, seen []atomic.Int64) {
	t.Helper()
	for i := range seen {
		if got := seen[i].Load(); got != 1 {
			t.Fatalf("index %d processed %d times, want 1", i, got)
		}
	}
}

// nonNeutralPartial makes zero-value partials observable in merge tests.
type nonNeutralPartial struct {
	Value  int
	Active bool
}
