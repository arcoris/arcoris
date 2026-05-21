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


package planner

import (
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func BenchmarkFixedChunks_1M_Chunk1K(b *testing.B) {
	benchmarkFixedChunks(b, 1_000_000, 1_024)
}

func BenchmarkFixedChunks_1M_Chunk64K(b *testing.B) {
	benchmarkFixedChunks(b, 1_000_000, 64*1024)
}

func benchmarkFixedChunks(
	b *testing.B,
	n int,
	chunk int,
) {
	opts := core.Options{ChunkSize: chunk, Strategy: core.StrategyFixedChunks}
	var ranges []core.Range
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ranges = FixedChunks(n, opts, ranges)
	}
	_ = ranges
}
