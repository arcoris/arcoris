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

// applyCreateValidated adds a new live item and appends its key to stable
// output order.
func (col *collection) applyCreateValidated(change objectstore.Change) {
	item := objectstore.ListItem{Key: change.Key, State: change.After.Clone()}
	col.ensureMaps()
	col.order = append(col.order, change.Key)
	col.items[change.Key] = item
	col.indexes.add(item)
	col.revision = change.Revision
}

// applyUpdateValidated replaces one live item while preserving its existing
// order slot.
func (col *collection) applyUpdateValidated(change objectstore.Change) {
	current := col.items[change.Key]
	next := objectstore.ListItem{Key: change.Key, State: change.After.Clone()}
	col.indexes.remove(current)
	col.items[change.Key] = next
	col.indexes.add(next)
	col.revision = change.Revision
}

// applyDeleteValidated removes a live item from items, order, and every index
// bucket.
func (col *collection) applyDeleteValidated(change objectstore.Change) {
	current := col.items[change.Key]
	col.indexes.remove(current)
	delete(col.items, change.Key)
	col.removeOrderKey(change.Key)
	col.revision = change.Revision
	if len(col.items) == 0 {
		col.items = nil
		col.indexes = indexes{}
	}
}

// ensureMaps initializes lazy maps for create-after-empty transitions.
func (col *collection) ensureMaps() {
	if col.items == nil {
		col.items = map[objectstore.Key]objectstore.ListItem{}
	}
	if col.indexes.byNamespace == nil {
		col.indexes = newIndexes()
	}
}

// removeOrderKey removes a known key from the deterministic order slice.
func (col *collection) removeOrderKey(key objectstore.Key) {
	for i, existing := range col.order {
		if existing.Equal(key) {
			col.order = append(col.order[:i], col.order[i+1:]...)
			return
		}
	}
}
