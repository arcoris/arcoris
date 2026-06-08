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
	"sync"
	"sync/atomic"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func BenchmarkMemoryUpdateSingleObject(b *testing.B) {
	store := benchmarkStore(b)
	key := testKey(1)
	current := benchmarkCreatedState(b, store, key)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		next, err := store.Update(context.Background(), key, current.Revision, testState("updated"))
		if err != nil {
			b.Fatalf("Update() error: %v", err)
		}
		current = next
	}
}

func BenchmarkMemoryUpdateDistinctObjects(b *testing.B) {
	store := benchmarkStore(b)
	const objects = 1024
	keys := make([]objectstore.Key, objects)
	states := make([]objectstore.State, objects)
	for i := 0; i < objects; i++ {
		keys[i] = testKey(i)
		states[i] = benchmarkCreatedState(b, store, keys[i])
	}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		index := i % objects
		next, err := store.Update(context.Background(), keys[index], states[index].Revision, testState("updated"))
		if err != nil {
			b.Fatalf("Update() error: %v", err)
		}
		states[index] = next
	}
}

func BenchmarkMemoryParallelUpdateSameObject(b *testing.B) {
	store := benchmarkStore(b)
	key := testKey(1)
	current := benchmarkCreatedState(b, store, key)
	var mu sync.Mutex
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			expected := current.Revision
			mu.Unlock()

			next, err := store.Update(context.Background(), key, expected, testState("updated"))
			if err == nil {
				mu.Lock()
				current = next
				mu.Unlock()
			}
		}
	})
}

func BenchmarkMemoryParallelUpdateDistinctObjects(b *testing.B) {
	store := benchmarkStore(b)
	const objects = 1024
	keys := make([]objectstore.Key, objects)
	states := make([]objectstore.State, objects)
	var locks [objects]sync.Mutex
	for i := 0; i < objects; i++ {
		keys[i] = testKey(i)
		states[i] = benchmarkCreatedState(b, store, keys[i])
	}

	var nextIndex atomic.Uint64
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			index := int(nextIndex.Add(1) % objects)
			locks[index].Lock()
			expected := states[index].Revision
			updated, err := store.Update(context.Background(), keys[index], expected, testState("updated"))
			if err == nil {
				states[index] = updated
			}
			locks[index].Unlock()
		}
	})
}
