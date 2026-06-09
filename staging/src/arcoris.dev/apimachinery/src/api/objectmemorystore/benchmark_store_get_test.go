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

package objectmemorystore

import (
	"context"
	"sync/atomic"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func BenchmarkMemoryGetSingleObject(b *testing.B) {
	store := benchmarkStore(b)
	key := testKey(1)
	benchmarkCreatedState(b, store, key)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, _ = store.Get(context.Background(), key)
	}
}

func BenchmarkMemoryParallelGetSameObject(b *testing.B) {
	store := benchmarkStore(b)
	key := testKey(1)
	benchmarkCreatedState(b, store, key)
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = store.Get(context.Background(), key)
		}
	})
}

func BenchmarkMemoryParallelGetDistinctObjects(b *testing.B) {
	store := benchmarkStore(b)
	const objects = 1024
	keys := make([]objectstore.Key, objects)
	for i := 0; i < objects; i++ {
		key := testKey(i)
		benchmarkCreatedState(b, store, key)
		keys[i] = key
	}

	var next atomic.Uint64
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := int(next.Add(1) % objects)
			_, _, _ = store.Get(context.Background(), keys[index])
		}
	})
}
