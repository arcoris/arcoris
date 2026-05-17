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
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestRunnerReusesScratch(t *testing.T) {
	r := New[int](core.Options{Workers: 4, MinItemsPerWorker: 10, Strategy: core.StrategyBalanced})
	for i := 0; i < 3; i++ {
		got, ok := r.ReduceInto(100, func(rng core.Range, dst *int) {
			for x := rng.Start; x < rng.End; x++ {
				*dst += x
			}
		}, func(dst *int, src int) { *dst += src })
		if !ok {
			t.Fatal("expected ok")
		}
		if got != 4950 {
			t.Fatalf("got %d, want 4950", got)
		}
	}
}

func TestRunnerReduceIndexedIntoReusesScratch(t *testing.T) {
	r := New[int](core.Options{
		Workers:           3,
		MinItemsPerWorker: 1,
		ChunkSize:         7,
		Strategy:          core.StrategyFixedChunks,
	})
	got, ok := r.ReduceIndexedInto(42, func(_ int, rng core.Range, dst *int) {
		*dst += rng.Len()
	}, func(dst *int, src int) { *dst += src })
	if !ok {
		t.Fatal("expected ok")
	}
	if got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
}

func TestRunnerAccumulateIntoReusesScratch(t *testing.T) {
	r := New[int](core.Options{
		Workers:           3,
		MinItemsPerWorker: 1,
		ChunkSize:         7,
		Strategy:          core.StrategyFixedChunks,
	})
	got, ok := r.AccumulateInto(42, func(rng core.Range, dst *int) {
		*dst += rng.Len()
	}, func(dst *int, src int) { *dst += src })
	if !ok {
		t.Fatal("expected ok")
	}
	if got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
}

func TestRunnerClearAndReleaseMemoryManagement(t *testing.T) {
	r := New[*int](core.Options{
		Workers:           2,
		MinItemsPerWorker: 1,
		ChunkSize:         1,
		Strategy:          core.StrategyFixedChunks,
	})
	value := 1
	_, ok := r.ReduceInto(4, func(_ core.Range, dst **int) {
		*dst = &value
	}, func(dst **int, src *int) {
		*dst = src
	})
	if !ok {
		t.Fatal("expected ok")
	}
	if cap(r.scratch.Partials) == 0 {
		t.Fatal("expected retained partial storage")
	}

	r.Clear()
	retained := r.scratch.Partials[:cap(r.scratch.Partials)]
	for i, ptr := range retained {
		if ptr != nil {
			t.Fatalf("retained partial[%d] kept pointer after Clear", i)
		}
	}

	got, ok := r.ReduceInto(4, func(rng core.Range, dst **int) {
		*dst = &value
	}, func(dst **int, src *int) {
		*dst = src
	})
	if !ok || got != &value {
		t.Fatalf("ReduceInto after Clear = %v ok=%v, want value pointer true", got, ok)
	}

	r.Release()
	if r.scratch.Ranges != nil || r.scratch.Partials != nil || r.scratch.Used != nil {
		t.Fatal("Release did not drop scratch storage")
	}
	got, ok = r.ReduceInto(4, func(_ core.Range, dst **int) {
		*dst = &value
	}, func(dst **int, src *int) {
		*dst = src
	})
	if !ok || got != &value {
		t.Fatalf("ReduceInto after Release = %v ok=%v, want value pointer true", got, ok)
	}
}
