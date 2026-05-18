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

package merge

import "testing"

type smallStruct struct {
	A, B int64
}

type largeStruct struct {
	V [16]int64
}

func BenchmarkMergeLinear_SmallStruct(b *testing.B) {
	partials := make([]smallStruct, 64)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Linear(partials, func(dst *smallStruct, src smallStruct) {
			dst.A += src.A
			dst.B += src.B
		})
	}
}

func BenchmarkMergePairwise_SmallStruct_IncludingCopy(b *testing.B) {
	partials := make([]smallStruct, 64)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		work := append([]smallStruct(nil), partials...)
		_, _ = PairwiseInPlace(work, func(dst *smallStruct, src smallStruct) {
			dst.A += src.A
			dst.B += src.B
		})
	}
}

func BenchmarkMergeLinear_LargeStruct(b *testing.B) {
	partials := make([]largeStruct, 64)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Linear(partials, func(dst *largeStruct, src largeStruct) {
			for j := range dst.V {
				dst.V[j] += src.V[j]
			}
		})
	}
}

func BenchmarkMergePairwise_LargeStruct_IncludingCopy(b *testing.B) {
	partials := make([]largeStruct, 64)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		work := append([]largeStruct(nil), partials...)
		_, _ = PairwiseInPlace(work, func(dst *largeStruct, src largeStruct) {
			for j := range dst.V {
				dst.V[j] += src.V[j]
			}
		})
	}
}
