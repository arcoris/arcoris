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

	"arcoris.dev/measure/internal/reduce"
)

func TestDoIndexedIntoPassesWorkerIndexes(t *testing.T) {
	var scratch reduce.Scratch[int]
	_, ok := DoIndexedInto[int](
		1000,
		reduce.Options{Workers: 4, MinItemsPerWorker: 100, Strategy: reduce.StrategyStatic},
		&scratch,
		func(worker int, r reduce.Range, dst *int) {
			*dst = worker + r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("expected ok")
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

func TestDoIndexedIntoBoundsFixedWorkerSlots(t *testing.T) {
	var maxWorker atomic.Int64
	_, ok := DoIndexedInto[int](
		100,
		reduce.Options{Workers: 2, MinItemsPerWorker: 1, ChunkSize: 10, Strategy: reduce.StrategyFixed},
		nil,
		func(worker int, r reduce.Range, dst *int) {
			for {
				old := maxWorker.Load()
				if int64(worker) <= old || maxWorker.CompareAndSwap(old, int64(worker)) {
					break
				}
			}
			*dst = r.Len()
		},
		func(dst *int, src int) { *dst += src },
	)
	if !ok {
		t.Fatal("expected ok")
	}
	if got := maxWorker.Load(); got >= 2 {
		t.Fatalf("max worker = %d, want < 2", got)
	}
}
