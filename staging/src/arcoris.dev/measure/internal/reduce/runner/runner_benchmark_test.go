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
	"arcoris.dev/measure/internal/reduce/planner"
)

func BenchmarkReduceInto_Sequential(b *testing.B) {
	values := benchmarkValues(64 * 1024)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1, Strategy: core.StrategySequential}
	benchmarkReduceIntoSum(b, values, opts)
}

func BenchmarkReduceInto_StaticRangePartials(b *testing.B) {
	values := benchmarkValues(1_000_000)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1024, Strategy: core.StrategyStatic}
	benchmarkReduceIntoSum(b, values, opts)
}

func BenchmarkReduceInto_FixedWorkerPartials(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1, ChunkSize: 32, Strategy: core.StrategyFixed}
	benchmarkReduceIntoSum(b, values, opts)
}

func BenchmarkReduceIndexedInto_StaticRangePartials(b *testing.B) {
	values := benchmarkValues(1_000_000)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1024, Strategy: core.StrategyStatic}
	benchmarkReduceIndexedIntoSum(b, values, opts)
}

func BenchmarkReduceIndexedInto_FixedWorkerPartials(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1, ChunkSize: 32, Strategy: core.StrategyFixed}
	benchmarkReduceIndexedIntoSum(b, values, opts)
}

func BenchmarkDynamicWorkerPartials(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1, ChunkSize: 256, Strategy: core.StrategyDynamic}
	var scratch core.Scratch[int]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := reduceDynamicWorkerPartials[int](len(values), opts, &scratch, func(_ int, r core.Range, dst *int) {
			sumRange(values, r, dst)
		}, func(dst *int, src int) { *dst += src })
		if !ok || got == 0 {
			b.Fatalf("unexpected result: %d ok=%v", got, ok)
		}
	}
}

func BenchmarkFillRangePartialsQueued(b *testing.B) {
	ranges := planner.Static(1_000_000, core.Options{Workers: 64, MinItemsPerWorker: 1, Strategy: core.StrategyStatic}, nil)
	partials := make([]int, len(ranges))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fillRangePartialsQueued(ranges, partials, 8, func(_ int, r core.Range, dst *int) {
			*dst = r.Len()
		})
	}
}

func BenchmarkFillWorkerPartialsFixed(b *testing.B) {
	partials := make([]int, 8)
	used := make([]bool, 8)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := range used {
			used[j] = false
		}
		fillWorkerPartialsFixed(1_000_000, 1024, partials, used, func(_ int, r core.Range, dst *int) {
			*dst = r.Len()
		}, func(dst *int, src int) { *dst += src })
	}
}

func BenchmarkCompactUsedPartials(b *testing.B) {
	partials := []int{1, 2, 3, 4, 5, 6, 7, 8}
	used := []bool{true, false, true, true, false, true, false, true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = compactUsedPartials(partials, used)
	}
}

func benchmarkReduceIntoSum(b *testing.B, values []int, opts core.Options) {
	var scratch core.Scratch[int]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := ReduceInto[int](len(values), opts, &scratch, func(r core.Range, dst *int) {
			sumRange(values, r, dst)
		}, func(dst *int, src int) { *dst += src })
		if !ok || got == 0 {
			b.Fatalf("unexpected result: %d ok=%v", got, ok)
		}
	}
}

func benchmarkReduceIndexedIntoSum(b *testing.B, values []int, opts core.Options) {
	var scratch core.Scratch[int]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := ReduceIndexedInto[int](len(values), opts, &scratch, func(_ int, r core.Range, dst *int) {
			sumRange(values, r, dst)
		}, func(dst *int, src int) { *dst += src })
		if !ok || got == 0 {
			b.Fatalf("unexpected result: %d ok=%v", got, ok)
		}
	}
}

func benchmarkValues(n int) []int {
	values := make([]int, n)
	for i := range values {
		values[i] = i + 1
	}
	return values
}

func sumRange(values []int, r core.Range, dst *int) {
	for _, value := range values[r.Start:r.End] {
		*dst += value
	}
}
