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

	"arcoris.dev/apimachinery/api/objectstore"
)

// shard protects one portion of the key-to-slot index.
//
// The mutex protects only slots map structure. It is not held while a caller's
// object state is loaded or transitioned. A plain Mutex is intentional: the
// critical sections are tiny, and mixed get/create benchmarks favored it over
// an RWMutex for this map-structure-only lock.
type shard struct {
	// mu protects structural access to slots.
	mu sync.Mutex

	// slots maps validated keys to per-object publication slots.
	slots map[objectstore.Key]*slot
}

// init prepares an empty shard.
func (s *shard) init() {
	s.slots = make(map[objectstore.Key]*slot)
}

// get returns the existing slot for key.
func (s *shard) get(key objectstore.Key) *slot {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.slots[key]
}

// getOrCreate returns a stable per-object slot for key.
func (s *shard) getOrCreate(key objectstore.Key) *slot {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing := s.slots[key]; existing != nil {
		return existing
	}

	created := new(slot)
	s.slots[key] = created

	return created
}

// shardFor selects the deterministic shard for key.
func (s *Store) shardFor(key objectstore.Key) *shard {
	return &s.shards[hashKey(key)&s.mask]
}
