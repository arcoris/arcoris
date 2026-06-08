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

package bulkhead

import "testing"

func BenchmarkSetLimitSame(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(8)
	for b.Loop() {
		benchmarkSnapshot = bulkhead.SetLimit(8)
	}
}

func BenchmarkSetLimitChanged(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(8)
	limit := Amount(8)
	for b.Loop() {
		if limit == 8 {
			limit = 9
		} else {
			limit = 8
		}
		benchmarkSnapshot = bulkhead.SetLimit(limit)
	}
}

func BenchmarkSetLimitDebt(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(8)
	lease, _, ok := bulkhead.TryAcquireAmount(8)
	if !ok {
		b.Fatal("initial TryAcquireAmount failed")
	}
	defer lease.Release()

	limit := Amount(4)
	for b.Loop() {
		if limit == 4 {
			limit = 3
		} else {
			limit = 4
		}
		benchmarkSnapshot = bulkhead.SetLimit(limit)
	}
}
