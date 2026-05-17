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

	"arcoris.dev/measure/internal/reduce"
)

func BenchmarkDoIntoStaticSum(b *testing.B) {
	values := make([]int, 1_000_000)
	for i := range values {
		values[i] = i
	}
	opts := reduce.Options{Workers: 8, MinItemsPerWorker: 1024, Strategy: reduce.StrategyStatic}
	var scratch reduce.Scratch[int]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := DoInto[int](len(values), opts, &scratch, func(r reduce.Range, dst *int) {
			chunk := values[r.Start:r.End]
			for j := 0; j < len(chunk); j++ {
				*dst += chunk[j]
			}
		}, func(dst *int, src int) { *dst += src })
		if !ok || got == 0 {
			b.Fatalf("unexpected result: %d ok=%v", got, ok)
		}
	}
}

func BenchmarkDoIntoFixedSmallChunks(b *testing.B) {
	values := make([]int, 65_536)
	for i := range values {
		values[i] = i
	}
	opts := reduce.Options{Workers: 8, MinItemsPerWorker: 1, ChunkSize: 32, Strategy: reduce.StrategyFixed}
	var scratch reduce.Scratch[int]
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		got, ok := DoInto[int](len(values), opts, &scratch, func(r reduce.Range, dst *int) {
			for _, value := range values[r.Start:r.End] {
				*dst += value
			}
		}, func(dst *int, src int) { *dst += src })
		if !ok || got == 0 {
			b.Fatalf("unexpected result: %d ok=%v", got, ok)
		}
	}
}
