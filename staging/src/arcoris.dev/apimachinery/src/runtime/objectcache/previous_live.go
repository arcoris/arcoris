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

// PreviousLive returns the newest retained live version observed before before.
//
// Revisions are not dense, so PreviousLive scans retained object versions
// instead of deriving a predecessor revision arithmetically.
func (c *Cache) PreviousLive(key objectstore.Key, before objectstore.Revision) (GetResult, error) {
	if c == nil {
		return GetResult{}, ErrInvalidCache
	}
	if !objectstore.KeyMatchesListRequest(key, c.collection) {
		return GetResult{}, ErrOutsideCollection
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	state, err := c.prepareHistoricalReadLocked(key, before)
	if err != nil {
		return GetResult{}, err
	}
	if state.record == nil {
		return GetResult{}, ErrHistoryUnavailable
	}

	var found objectVersion
	ok := false
	state.record.newestToOldest(func(version objectVersion) bool {
		if version.Live && version.Revision.Before(before) {
			found = version
			ok = true
			return false
		}
		return true
	})
	if !ok {
		return GetResult{}, ErrHistoryUnavailable
	}

	return liveResultAt(key, found.Revision, found.State), nil
}
