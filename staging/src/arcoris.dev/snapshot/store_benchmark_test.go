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

package snapshot

import "testing"

func BenchmarkStoreSnapshotSmallValue(b *testing.B) {
	store := NewStore(42, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = store.Snapshot()
	}
}

func BenchmarkStoreReplaceSmallValue(b *testing.B) {
	store := NewStore(0, Identity[int])

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = store.Replace(i)
	}
}

func BenchmarkStoreSnapshotSlice100(b *testing.B) {
	val := make([]string, 100)
	store := NewStore(val, cloneStrings)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = store.Snapshot()
	}
}

func BenchmarkStoreUpdateSlice100(b *testing.B) {
	val := make([]string, 100)
	store := NewStore(val, cloneStrings)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = store.Update(func(v []string) []string {
			v[0] = "updated"
			return v
		})
	}
}
