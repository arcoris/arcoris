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
	"fmt"

	"arcoris.dev/apimachinery/api/objectstore"
)

// apply dispatches a validated, newer objectstore.Change into collection-local
// mutation logic. The caller is responsible for atomic cache publication.
func (col *collection) apply(change objectstore.Change) error {
	switch change.Kind {
	case objectstore.ChangeCreated:
		return col.applyCreate(change)
	case objectstore.ChangeUpdated:
		return col.applyUpdate(change)
	case objectstore.ChangeDeleted:
		return col.applyDelete(change)
	default:
		return invalidChangeError(nil)
	}
}

// applyCreate adds a new live item and appends its key to stable output order.
func (col *collection) applyCreate(change objectstore.Change) error {
	if _, exists := col.items[change.Key]; exists {
		return invalidChangeStateError("create key already exists", change.Key)
	}

	item := objectstore.ListItem{Key: change.Key, State: change.After.Clone()}
	col.ensureMaps()
	col.order = append(col.order, change.Key)
	col.items[change.Key] = item
	col.indexes.add(item)
	col.revision = change.Revision

	return nil
}

// applyUpdate replaces one live item while preserving its existing order slot.
func (col *collection) applyUpdate(change objectstore.Change) error {
	current, exists := col.items[change.Key]
	if !exists {
		return invalidChangeStateError("update key is missing", change.Key)
	}
	if current.State.Revision != change.Before.Revision {
		return invalidChangeStateError("update before revision does not match cache", change.Key)
	}

	next := objectstore.ListItem{Key: change.Key, State: change.After.Clone()}
	col.indexes.remove(current)
	col.items[change.Key] = next
	col.indexes.add(next)
	col.revision = change.Revision

	return nil
}

// applyDelete removes a live item from items, order, and every index bucket.
func (col *collection) applyDelete(change objectstore.Change) error {
	current, exists := col.items[change.Key]
	if !exists {
		return invalidChangeStateError("delete key is missing", change.Key)
	}
	if current.State.Revision != change.Before.Revision {
		return invalidChangeStateError("delete before revision does not match cache", change.Key)
	}

	col.indexes.remove(current)
	delete(col.items, change.Key)
	col.removeOrderKey(change.Key)
	col.revision = change.Revision
	if len(col.items) == 0 {
		col.items = nil
		col.indexes = indexes{}
	}

	return nil
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

// invalidChangeStateError reports a structurally valid change that does not
// match the cache's current materialized state.
func invalidChangeStateError(message string, key objectstore.Key) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidChange, message, key.String())
}
