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
	"sync"
	"sync/atomic"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func BenchmarkShardLockMutexParallelGetDistinctKeys(b *testing.B) {
	keys := benchmarkShardKeys(1024)
	var shard shard
	shard.init()
	for _, key := range keys {
		shard.getOrCreate(key)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var n atomic.Uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[int(n.Add(1))&(len(keys)-1)]
			if shard.get(key) == nil {
				b.Fatal("slot missing")
			}
		}
	})
}

func BenchmarkShardLockRWMutexParallelGetDistinctKeys(b *testing.B) {
	keys := benchmarkShardKeys(1024)
	shard := newBenchmarkRWShard(keys)

	b.ReportAllocs()
	b.ResetTimer()
	var n atomic.Uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[int(n.Add(1))&(len(keys)-1)]
			if shard.get(key) == nil {
				b.Fatal("slot missing")
			}
		}
	})
}

func BenchmarkShardLockMutexParallelMixedGetCreate(b *testing.B) {
	keys := benchmarkShardKeys(1024)
	var shard shard
	shard.init()
	for _, key := range keys[:len(keys)/2] {
		shard.getOrCreate(key)
	}

	b.ReportAllocs()
	b.ResetTimer()
	var n atomic.Uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := int(n.Add(1)) & (len(keys) - 1)
			if i&7 == 0 {
				shard.getOrCreate(keys[i])
				continue
			}
			_ = shard.get(keys[i])
		}
	})
}

func BenchmarkShardLockRWMutexParallelMixedGetCreate(b *testing.B) {
	keys := benchmarkShardKeys(1024)
	shard := newBenchmarkRWShard(keys[:len(keys)/2])

	b.ReportAllocs()
	b.ResetTimer()
	var n atomic.Uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := int(n.Add(1)) & (len(keys) - 1)
			if i&7 == 0 {
				shard.getOrCreate(keys[i])
				continue
			}
			_ = shard.get(keys[i])
		}
	})
}

type benchmarkRWShard struct {
	mu    sync.RWMutex
	slots map[objectstore.Key]*slot
}

func newBenchmarkRWShard(keys []objectstore.Key) *benchmarkRWShard {
	shard := &benchmarkRWShard{slots: make(map[objectstore.Key]*slot, len(keys))}
	for _, key := range keys {
		shard.slots[key] = new(slot)
	}

	return shard
}

func (s *benchmarkRWShard) get(key objectstore.Key) *slot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.slots[key]
}

func (s *benchmarkRWShard) getOrCreate(key objectstore.Key) *slot {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing := s.slots[key]; existing != nil {
		return existing
	}

	created := new(slot)
	s.slots[key] = created

	return created
}

func benchmarkShardKeys(n int) []objectstore.Key {
	keys := make([]objectstore.Key, n)
	for i := range keys {
		keys[i] = testKey(i + 1)
	}

	return keys
}
