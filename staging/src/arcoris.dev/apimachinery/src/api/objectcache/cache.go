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

package objectcache

import (
	"sync"

	"arcoris.dev/apimachinery/api/objectstore"
)

// Cache is a mutable, concurrency-safe materialized object collection.
//
// Cache owns detached state and private indexes. Public reads return clones.
// Replace rebuilds the collection atomically. Apply consumes validated
// objectstore.Change values and advances the cache revision when the change is
// newer than the current watermark.
type Cache struct {
	mu  sync.RWMutex
	col collection
}

// IsZero reports whether c has no cached items and no observed revision.
func (c *Cache) IsZero() bool {
	if c == nil {
		return true
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.col.isZero()
}

// Len returns the number of currently cached live items.
func (c *Cache) Len() int {
	if c == nil {
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.col.len()
}

// Revision returns the current cache revision watermark.
func (c *Cache) Revision() objectstore.Revision {
	if c == nil {
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.col.revision
}
