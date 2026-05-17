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

func TestAccumulateIndexedIntoFixedChunksDoesNotMergePerChunk(t *testing.T) {
	var mergeCalls atomic.Int64
	got, ok := AccumulateIndexedInto[int](
		100,
		core.Options{Workers: 2, MinItemsPerWorker: 1, ChunkSize: 1, Strategy: core.StrategyFixedChunks},
		nil,
		func(_ int, _ core.Range, dst *int) {
			*dst += 1
		},
		func(dst *int, src int) {
			mergeCalls.Add(1)
			*dst += src
		},
	)
	if !ok {
		t.Fatal("AccumulateIndexedInto returned false for non-empty input")
	}
	if got != 100 {
		t.Fatalf("AccumulateIndexedInto() = %d, want 100", got)
	}
	if mergeCalls.Load() != 1 {
		t.Fatalf("merge calls = %d, want one final worker merge", mergeCalls.Load())
	}
}

func TestAccumulateIndexedIntoDynamicChunksAccumulatesWorkerLocalState(t *testing.T) {
	var calls [2]atomic.Int64
	got, ok := AccumulateIndexedInto[int](
		100,
		core.Options{Workers: 2, MinItemsPerWorker: 1, ChunkSize: 1, Strategy: core.StrategyDynamicChunks},
		nil,
		func(worker int, _ core.Range, dst *int) {
			calls[worker].Add(1)
			*dst += 1
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("AccumulateIndexedInto returned false for non-empty input")
	}
	if got != 100 {
		t.Fatalf("AccumulateIndexedInto() = %d, want 100", got)
	}
	if calls[0].Load() <= 1 && calls[1].Load() <= 1 {
		t.Fatalf("worker calls = [%d %d], want repeated calls into at least one worker accumulator", calls[0].Load(), calls[1].Load())
	}
}

func TestAccumulateIndexedIntoDoesNotMergeInactiveWorkers(t *testing.T) {
	var inactiveMerged atomic.Int64
	got, ok := AccumulateIndexedInto[nonNeutralPartial](
		10,
		core.Options{Workers: 8, MinItemsPerWorker: 1, ChunkSize: 100, Strategy: core.StrategyDynamicChunks},
		nil,
		func(_ int, r core.Range, dst *nonNeutralPartial) {
			dst.Value += r.Len()
			dst.Active = true
		},
		func(dst *nonNeutralPartial, src nonNeutralPartial) {
			if !src.Active {
				inactiveMerged.Add(1)
			}
			dst.Value += src.Value
			dst.Active = dst.Active || src.Active
		},
	)
	if !ok {
		t.Fatal("AccumulateIndexedInto returned false for non-empty input")
	}
	if got.Value != 10 || !got.Active {
		t.Fatalf("AccumulateIndexedInto() = %#v, want active value 10", got)
	}
	if inactiveMerged.Load() != 0 {
		t.Fatalf("merged %d inactive worker partials", inactiveMerged.Load())
	}
}
