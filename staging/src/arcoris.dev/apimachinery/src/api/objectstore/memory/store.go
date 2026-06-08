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
	"sync/atomic"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Store is a concurrent in-memory objectstore.Store implementation.
//
// Store owns a fixed sharded key index and a store-local revision counter.
// Shards protect map structure only. Per-object slots own the committed record
// pointer and are updated with atomic compare-and-swap.
type Store struct {
	// revision is the global store-local monotonic commit counter.
	revision atomic.Uint64

	// shards is fixed at construction and never mutated after publication.
	shards []shard

	// mask is shardCount-1 and is valid because shard count is a power of two.
	mask uint64
}

var _ objectstore.Store = (*Store)(nil)
