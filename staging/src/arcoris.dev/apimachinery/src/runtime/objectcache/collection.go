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

import "arcoris.dev/apimachinery/api/objectstore"

// collection is the cache-owned materialized state.
//
// order defines public output order. items owns the current live item for every
// ordered key. The collection has no query indexes in this first runtime cache
// core; objectquery support is intentionally deferred.
//
// collection is not internally synchronized. Cache owns synchronization and
// swaps or mutates collection values only while holding Cache.mu.
type collection struct {
	// revision is the collection boundary installed by Replace or advanced by
	// ApplyChange. Revision zero is valid for an empty ready collection.
	revision objectstore.Revision
	// order is the deterministic output order for List. The cache preserves
	// replacement order, appends creates, preserves update slots, and removes
	// deletes without reordering survivors.
	order []objectstore.Key
	// items owns detached live states by key. Every key in order must exist in
	// items, and no item outside order is part of the public collection. Empty
	// collections may keep items nil to avoid needless maps.
	items map[objectstore.Key]objectstore.ListItem
}

// len returns the materialized live item count.
func (col collection) len() int {
	return len(col.order)
}

// item returns one detached item without exposing cache-owned state.
//
// The bool is only a latest-state presence signal. Historical reads use
// objectRecord tombstones to distinguish retained absence from unknown history.
func (col collection) item(key objectstore.Key) (objectstore.ListItem, bool) {
	item, ok := col.items[key]
	if !ok {
		return objectstore.ListItem{}, false
	}

	return item.Clone(), true
}

// listResult returns a detached ListResult at the collection boundary revision.
func (col collection) listResult() objectstore.ListResult {
	return objectstore.ListResult{Items: col.listItems(), Revision: col.revision}
}

// listItems returns detached items in deterministic cache order.
//
// The returned slice is safe for callers to mutate. It preserves replacement
// order plus committed create/update/delete order semantics.
func (col collection) listItems() []objectstore.ListItem {
	if len(col.order) == 0 {
		return nil
	}

	out := make([]objectstore.ListItem, 0, len(col.order))
	for _, key := range col.order {
		out = append(out, col.items[key].Clone())
	}

	return out
}
