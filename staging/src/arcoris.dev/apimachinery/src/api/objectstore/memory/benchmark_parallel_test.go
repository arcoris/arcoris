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

package memory

import (
	"context"
	"sync/atomic"
	"testing"
)

func BenchmarkMemoryParallelMixedDistinctObjects(b *testing.B) {
	store := benchmarkStore(b)
	var next atomic.Uint64
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := int(next.Add(1))
			key := testKey(index)
			created, err := store.Create(context.Background(), key, testState("created"))
			if err != nil {
				continue
			}
			updated, err := store.Update(context.Background(), key, created.Revision, testState("updated"))
			if err != nil {
				continue
			}
			_, _ = store.Delete(context.Background(), key, updated.Revision)
		}
	})
}
