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

// resetHistoryLocked clears all retained history after a replacement boundary.
//
// Cache.mu must be held. Replace may follow continuity loss, so old history
// cannot prove versions across the new boundary and is intentionally discarded.
func (c *Cache) resetHistoryLocked(col collection) {
	if !c.historyEnabled {
		c.records = nil
		return
	}

	c.records = make(map[objectstore.Key]*objectRecord, len(col.order))
	for _, key := range col.order {
		item := col.items[key]
		record := newObjectRecord(key, c.retainedVersionsPerObject)
		record.append(liveVersion(col.revision, item.State))
		c.records[key] = record
	}
}

// appendHistoryLocked records the cache-observed version produced by change.
//
// Cache.mu must be held. The caller has already validated and applied the
// change to latest; this helper only mirrors the successful mutation into the
// optional per-object history record.
func (c *Cache) appendHistoryLocked(change objectstore.Change) {
	if !c.historyEnabled {
		return
	}

	record := c.records[change.Key]
	if record == nil {
		record = newObjectRecord(change.Key, c.retainedVersionsPerObject)
		c.records[change.Key] = record
	}

	switch change.Kind {
	case objectstore.ChangeCreated, objectstore.ChangeUpdated:
		record.append(liveVersion(change.Revision, change.After))
	case objectstore.ChangeDeleted:
		record.append(tombstoneVersion(change.Revision))
	}
}
