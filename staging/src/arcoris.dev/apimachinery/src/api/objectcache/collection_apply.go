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

// validateApply checks cache-local preconditions without mutating collection
// state. The caller must already have validated the objectstore.Change shape and
// stale revision ordering.
func (col collection) validateApply(change objectstore.Change) error {
	switch change.Kind {
	case objectstore.ChangeCreated:
		return col.validateCreate(change)
	case objectstore.ChangeUpdated:
		return col.validateUpdate(change)
	case objectstore.ChangeDeleted:
		return col.validateDelete(change)
	default:
		return invalidChangeError(nil)
	}
}

// applyValidated dispatches a change that already passed validateApply.
func (col *collection) applyValidated(change objectstore.Change) {
	switch change.Kind {
	case objectstore.ChangeCreated:
		col.applyCreateValidated(change)
	case objectstore.ChangeUpdated:
		col.applyUpdateValidated(change)
	case objectstore.ChangeDeleted:
		col.applyDeleteValidated(change)
	}
}

func (col collection) validateCreate(change objectstore.Change) error {
	if _, exists := col.items[change.Key]; exists {
		return invalidChangeStateError("create key already exists", change.Key)
	}

	return nil
}

func (col collection) validateUpdate(change objectstore.Change) error {
	current, exists := col.items[change.Key]
	if !exists {
		return invalidChangeStateError("update key is missing", change.Key)
	}
	if current.State.Revision != change.Before.Revision {
		return invalidChangeStateError("update before revision does not match cache", change.Key)
	}

	return nil
}

func (col collection) validateDelete(change objectstore.Change) error {
	current, exists := col.items[change.Key]
	if !exists {
		return invalidChangeStateError("delete key is missing", change.Key)
	}
	if current.State.Revision != change.Before.Revision {
		return invalidChangeStateError("delete before revision does not match cache", change.Key)
	}

	return nil
}

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

// invalidChangeStateError reports a structurally valid change that does not
// match the cache's current materialized state.
func invalidChangeStateError(message string, key objectstore.Key) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidChange, message, key.String())
}
