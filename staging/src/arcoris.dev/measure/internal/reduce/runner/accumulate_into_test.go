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

func TestAccumulateIntoSequentialMapsOnce(t *testing.T) {
	var calls atomic.Int64
	got, ok := AccumulateInto[int](
		25,
		core.Options{Workers: 8, Strategy: core.StrategySequential},
		nil,
		func(r core.Range, dst *int) {
			calls.Add(1)
			*dst += r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("AccumulateInto returned false for non-empty input")
	}
	if got != 25 {
		t.Fatalf("AccumulateInto() = %d, want 25", got)
	}
	if calls.Load() != 1 {
		t.Fatalf("accumulator calls = %d, want 1", calls.Load())
	}
}

func TestAccumulateIntoFixedChunksPreservesAccumulatorState(t *testing.T) {
	got, ok := AccumulateInto[int](
		100,
		core.Options{Workers: 2, MinItemsPerWorker: 1, ChunkSize: 1, Strategy: core.StrategyFixedChunks},
		nil,
		func(_ core.Range, dst *int) {
			// Intentionally ignores range length: this verifies that repeated
			// calls receive the same worker-local dst and preserve prior state.
			*dst += 1
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("AccumulateInto returned false for non-empty input")
	}
	if got != 100 {
		t.Fatalf("AccumulateInto() = %d, want one increment per chunk", got)
	}
}
