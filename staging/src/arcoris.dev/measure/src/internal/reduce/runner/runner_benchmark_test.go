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

import "testing"

import "arcoris.dev/measure/internal/reduce/core"

func BenchmarkReduceInto_Sequential(b *testing.B) {
	values := benchmarkValues(64 * 1024)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1, Strategy: core.StrategySequential}
	benchmarkReduceIntoSum(b, values, opts)
}

func BenchmarkReduceInto_Balanced(b *testing.B) {
	values := benchmarkValues(1_000_000)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1024, Strategy: core.StrategyBalanced}
	benchmarkReduceIntoSum(b, values, opts)
}

func BenchmarkReduceInto_FixedChunks(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         32,
		Strategy:          core.StrategyFixedChunks,
	}
	benchmarkReduceIntoSum(b, values, opts)
}

func BenchmarkReduceIndexedInto_Balanced(b *testing.B) {
	values := benchmarkValues(1_000_000)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1024, Strategy: core.StrategyBalanced}
	benchmarkReduceIndexedIntoSum(b, values, opts)
}

func BenchmarkReduceIndexedInto_FixedChunks(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         32,
		Strategy:          core.StrategyFixedChunks,
	}
	benchmarkReduceIndexedIntoSum(b, values, opts)
}

func BenchmarkReduceInto_DynamicChunks(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         256,
		Strategy:          core.StrategyDynamicChunks,
	}
	benchmarkReduceIntoSum(b, values, opts)
}

func BenchmarkAccumulateInto_Sequential(b *testing.B) {
	values := benchmarkValues(64 * 1024)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1, Strategy: core.StrategySequential}
	benchmarkAccumulateIntoSum(b, values, opts)
}

func BenchmarkAccumulateInto_Balanced(b *testing.B) {
	values := benchmarkValues(1_000_000)
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1024, Strategy: core.StrategyBalanced}
	benchmarkAccumulateIntoSum(b, values, opts)
}

func BenchmarkAccumulateInto_FixedChunks(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         32,
		Strategy:          core.StrategyFixedChunks,
	}
	benchmarkAccumulateIntoSum(b, values, opts)
}

func BenchmarkAccumulateInto_DynamicChunks(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         256,
		Strategy:          core.StrategyDynamicChunks,
	}
	benchmarkAccumulateIntoSum(b, values, opts)
}

func BenchmarkReduceVsAccumulate_FixedChunks(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         32,
		Strategy:          core.StrategyFixedChunks,
	}
	b.Run("reduce", func(b *testing.B) { benchmarkReduceIntoSum(b, values, opts) })
	b.Run("accumulate", func(b *testing.B) { benchmarkAccumulateIntoSum(b, values, opts) })
}

func BenchmarkReduceVsAccumulate_DynamicChunks(b *testing.B) {
	values := benchmarkValues(65_536)
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         256,
		Strategy:          core.StrategyDynamicChunks,
	}
	b.Run("reduce", func(b *testing.B) { benchmarkReduceIntoSum(b, values, opts) })
	b.Run("accumulate", func(b *testing.B) { benchmarkAccumulateIntoSum(b, values, opts) })
}

func BenchmarkReduceInto_BucketPartial_FixedChunks(b *testing.B) {
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         32,
		Strategy:          core.StrategyFixedChunks,
	}
	benchmarkReduceIntoBucketPartial(b, 65_536, opts)
}

func BenchmarkAccumulateInto_BucketPartial_FixedChunks(b *testing.B) {
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         32,
		Strategy:          core.StrategyFixedChunks,
	}
	benchmarkAccumulateIntoBucketPartial(b, 65_536, opts)
}

func BenchmarkReduceVsAccumulate_BufferBackedPartial_FixedChunks(b *testing.B) {
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         32,
		Strategy:          core.StrategyFixedChunks,
	}
	b.Run("reduce", func(b *testing.B) {
		benchmarkReduceIntoBucketPartial(b, 65_536, opts)
	})
	b.Run("accumulate", func(b *testing.B) {
		benchmarkAccumulateIntoBucketPartial(b, 65_536, opts)
	})
}

func BenchmarkReduceVsAccumulate_BufferBackedPartial_DynamicChunks(b *testing.B) {
	opts := core.Options{
		Workers:           8,
		MinItemsPerWorker: 1,
		ChunkSize:         256,
		Strategy:          core.StrategyDynamicChunks,
	}
	b.Run("reduce", func(b *testing.B) {
		benchmarkReduceIntoBucketPartial(b, 65_536, opts)
	})
	b.Run("accumulate", func(b *testing.B) {
		benchmarkAccumulateIntoBucketPartial(b, 65_536, opts)
	})
}

func BenchmarkFillFixedChunkWorkerPartials(b *testing.B) {
	partials := make([]int, 8)
	used := make([]bool, 8)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := range used {
			used[j] = false
		}
		fillFixedChunkWorkerPartials(
			1_000_000,
			1024,
			chunkCount(1_000_000, 1024),
			partials,
			used,
			func(_ int, r core.Range, dst *int) {
				*dst = r.Len()
			},
			func(dst *int, src int) { *dst += src },
		)
	}
}

func BenchmarkCompactUsedPartials(b *testing.B) {
	basePartials := []int{1, 2, 3, 4, 5, 6, 7, 8}
	baseUsed := []bool{true, false, true, true, false, true, false, true}
	partials := append([]int(nil), basePartials...)
	used := append([]bool(nil), baseUsed...)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		copy(partials, basePartials)
		copy(used, baseUsed)
		_ = compactUsedPartials(partials, used)
	}
}

type bucketPartial struct {
	Buckets []uint64
}

const bucketCount = 64

func benchmarkReduceIntoSum(
	b *testing.B,
	values []int,
	opts core.Options,
) {
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

func benchmarkReduceIntoBucketPartial(
	b *testing.B,
	n int,
	opts core.Options,
) {
	var scratch core.Scratch[bucketPartial]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := ReduceInto[bucketPartial](
			n,
			opts,
			&scratch,
			func(r core.Range, dst *bucketPartial) {
				dst.Buckets = make([]uint64, bucketCount)
				fillBucketRange(r, dst.Buckets)
			},
			func(dst *bucketPartial, src bucketPartial) {
				if dst.Buckets == nil {
					dst.Buckets = make([]uint64, len(src.Buckets))
				}
				for i := range src.Buckets {
					dst.Buckets[i] += src.Buckets[i]
				}
			},
		)
		if !ok || len(got.Buckets) != bucketCount {
			b.Fatalf("unexpected result: len=%d ok=%v", len(got.Buckets), ok)
		}
	}
}

func benchmarkReduceIndexedIntoSum(
	b *testing.B,
	values []int,
	opts core.Options,
) {
	var scratch core.Scratch[int]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := ReduceIndexedInto[int](
			len(values),
			opts,
			&scratch,
			func(_ int, r core.Range, dst *int) {
				sumRange(values, r, dst)
			},
			func(dst *int, src int) { *dst += src },
		)
		if !ok || got == 0 {
			b.Fatalf("unexpected result: %d ok=%v", got, ok)
		}
	}
}

func benchmarkAccumulateIntoSum(
	b *testing.B,
	values []int,
	opts core.Options,
) {
	var scratch core.Scratch[int]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := AccumulateInto[int](len(values), opts, &scratch, func(r core.Range, dst *int) {
			sumRange(values, r, dst)
		}, func(dst *int, src int) { *dst += src })
		if !ok || got == 0 {
			b.Fatalf("unexpected result: %d ok=%v", got, ok)
		}
	}
}

func benchmarkAccumulateIntoBucketPartial(
	b *testing.B,
	n int,
	opts core.Options,
) {
	var scratch core.Scratch[bucketPartial]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := AccumulateInto[bucketPartial](
			n,
			opts,
			&scratch,
			func(r core.Range, dst *bucketPartial) {
				if dst.Buckets == nil {
					dst.Buckets = make([]uint64, bucketCount)
				}
				fillBucketRange(r, dst.Buckets)
			},
			func(dst *bucketPartial, src bucketPartial) {
				if dst.Buckets == nil {
					dst.Buckets = make([]uint64, len(src.Buckets))
				}
				for i := range src.Buckets {
					dst.Buckets[i] += src.Buckets[i]
				}
			},
		)
		if !ok || len(got.Buckets) != bucketCount {
			b.Fatalf("unexpected result: len=%d ok=%v", len(got.Buckets), ok)
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

func sumRange(
	values []int,
	r core.Range,
	dst *int,
) {
	for _, value := range values[r.Start:r.End] {
		*dst += value
	}
}

func fillBucketRange(r core.Range, buckets []uint64) {
	for i := r.Start; i < r.End; i++ {
		buckets[i%len(buckets)]++
	}
}
