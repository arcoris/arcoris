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

func BenchmarkBalanced_1K(b *testing.B) {
	benchmarkBalanced(b, 1_024)
}

func BenchmarkBalanced_64K(b *testing.B) {
	benchmarkBalanced(b, 64*1024)
}

func BenchmarkBalanced_1M(b *testing.B) {
	benchmarkBalanced(b, 1_000_000)
}

func benchmarkBalanced(b *testing.B, n int) {
	opts := core.Options{Workers: 8, MinItemsPerWorker: 1024, Strategy: core.StrategyBalanced}
	var ranges []core.Range
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ranges = Balanced(n, opts, ranges)
	}
	_ = ranges
}
