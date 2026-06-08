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

package snapshot

import "testing"

func BenchmarkStoreSnapshotSmallValue(b *testing.B) {
	store := NewStore(42, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntSnapshotSink = store.Snapshot()
	}
}

func BenchmarkStoreStampedSmallValue(b *testing.B) {
	store := NewStore(42, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntStampedSink = store.Stamped()
	}
}

func BenchmarkStoreRevision(b *testing.B) {
	store := NewStore(42, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkRevisionSink = store.Revision()
	}
}

func BenchmarkStoreSnapshotParallel(b *testing.B) {
	store := NewStore(42, Identity[int])

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var snap Snapshot[int]
		for pb.Next() {
			snap = store.Snapshot()
		}
		benchmarkSinkMu.Lock()
		benchmarkIntSnapshotSink = snap
		benchmarkSinkMu.Unlock()
	})
}

func BenchmarkStoreStampedParallel(b *testing.B) {
	store := NewStore(42, Identity[int])

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var stamped Stamped[int]
		for pb.Next() {
			stamped = store.Stamped()
		}
		benchmarkSinkMu.Lock()
		benchmarkIntStampedSink = stamped
		benchmarkSinkMu.Unlock()
	})
}

func BenchmarkStoreReplaceSmallValue(b *testing.B) {
	store := NewStore(0, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntSnapshotSink = store.Replace(i)
	}
}

func BenchmarkStoreReplaceStampedSmallValue(b *testing.B) {
	store := NewStore(0, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntStampedSink = store.ReplaceStamped(i)
	}
}

func BenchmarkStoreReplaceParallel(b *testing.B) {
	store := NewStore(0, Identity[int])

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var snap Snapshot[int]
		for pb.Next() {
			snap = store.Replace(1)
		}
		benchmarkSinkMu.Lock()
		benchmarkIntSnapshotSink = snap
		benchmarkSinkMu.Unlock()
	})
}

func BenchmarkStoreUpdateSmallValue(b *testing.B) {
	store := NewStore(0, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntSnapshotSink = store.Update(func(v int) int {
			return v + 1
		})
	}
}

func BenchmarkStoreUpdateStampedSmallValue(b *testing.B) {
	store := NewStore(0, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkIntStampedSink = store.UpdateStamped(func(v int) int {
			return v + 1
		})
	}
}

func BenchmarkStoreSnapshotSlice100(b *testing.B) {
	val := make([]string, 100)
	store := NewStore(val, cloneStrings)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkSliceSnapshotSink = store.Snapshot()
	}
}

func BenchmarkStoreStampedSlice100(b *testing.B) {
	val := make([]string, 100)
	store := NewStore(val, cloneStrings)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkSliceStampedSink = store.Stamped()
	}
}

func BenchmarkStoreReplaceSlice100(b *testing.B) {
	val := make([]string, 100)
	store := NewStore(val, cloneStrings)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkSliceSnapshotSink = store.Replace(val)
	}
}

func BenchmarkStoreUpdateParallel(b *testing.B) {
	store := NewStore(0, Identity[int])

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var snap Snapshot[int]
		for pb.Next() {
			snap = store.Update(func(v int) int {
				return v + 1
			})
		}
		benchmarkSinkMu.Lock()
		benchmarkIntSnapshotSink = snap
		benchmarkSinkMu.Unlock()
	})
}

func BenchmarkStoreUpdateSlice100(b *testing.B) {
	val := make([]string, 100)
	store := NewStore(val, cloneStrings)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkSliceSnapshotSink = store.Update(func(v []string) []string {
			v[0] = "updated"
			return v
		})
	}
}
