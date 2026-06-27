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

// GetAt returns the retained object state as of a cache observation revision.
//
// Found=false is a proven absence answer. ErrHistoryUnavailable means the cache
// cannot prove the answer inside its per-object retention window or across the
// last replacement boundary. The newest retained version at or before revision
// defines the retained answer; retained tombstones prove absence.
func (c *Cache) GetAt(key objectstore.Key, revision objectstore.Revision) (GetResult, error) {
	if c == nil {
		return GetResult{}, ErrInvalidCache
	}
	if !objectstore.KeyMatchesListRequest(key, c.collection) {
		return GetResult{}, ErrOutsideCollection
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	state, err := c.prepareHistoricalReadLocked(key, revision)
	if err != nil {
		return GetResult{}, err
	}
	if state.record == nil {
		if revision == state.current {
			return knownAbsenceAt(key, revision), nil
		}
		return GetResult{}, ErrHistoryUnavailable
	}

	var found objectVersion
	ok := false
	state.record.newestToOldest(func(version objectVersion) bool {
		if !revision.Before(version.Revision) {
			found = version
			ok = true
			return false
		}
		return true
	})
	if !ok {
		return GetResult{}, ErrHistoryUnavailable
	}
	if !found.Live {
		return knownAbsenceAt(key, revision), nil
	}

	return liveResultAt(key, revision, found.State), nil
}
