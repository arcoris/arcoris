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

// applyValidated dispatches a change that has already passed validateApply.
//
// Callers must hold Cache.mu through the collection that owns col. The method
// intentionally performs no validation; keeping validation in validateApply
// makes failed changes leave collection state untouched.
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

// applyCreateValidated adds a live item and appends the key to output order.
//
// Creates are the only mutation that grows order. Updates keep the existing
// slot and deletes remove a slot without reordering survivors.
func (col *collection) applyCreateValidated(change objectstore.Change) {
	item := objectstore.ListItem{Key: change.Key, State: change.After.Clone()}
	col.ensureItems()
	col.order = append(col.order, change.Key)
	col.items[change.Key] = item
	col.revision = change.Revision
}

// applyUpdateValidated replaces an existing item without moving its order slot.
//
// The replacement state is cloned because objectstore.Change belongs to the
// caller; collection owns everything it retains.
func (col *collection) applyUpdateValidated(change objectstore.Change) {
	col.items[change.Key] = objectstore.ListItem{Key: change.Key, State: change.After.Clone()}
	col.revision = change.Revision
}

// applyDeleteValidated removes one item and keeps remaining order stable.
//
// Delete removes only latest live state. Optional tombstone retention is handled
// by Cache.appendHistoryLocked after latest mutation succeeds.
func (col *collection) applyDeleteValidated(change objectstore.Change) {
	delete(col.items, change.Key)
	col.removeOrderKey(change.Key)
	if len(col.items) == 0 {
		col.items = nil
	}
	col.revision = change.Revision
}
