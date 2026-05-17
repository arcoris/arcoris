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

func TestFillFixedChunkWorkerAccumulatorsUsesChunkCountForWorkers(t *testing.T) {
	partials := make([]int, activeWorkers(8, chunkCount(10, 100)))
	used := make([]bool, len(partials))
	fillFixedChunkWorkerAccumulators(
		10,
		100,
		chunkCount(10, 100),
		partials,
		used,
		func(_ int, r core.Range, dst *int) {
			*dst += r.Len()
		},
	)
	active := compactUsedPartials(partials, used)
	if len(active) != 1 {
		t.Fatalf("active workers = %d, want 1", len(active))
	}
	if active[0] != 10 {
		t.Fatalf("active[0] = %d, want 10", active[0])
	}
}

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
