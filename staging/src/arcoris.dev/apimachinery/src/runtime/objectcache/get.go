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

// Get returns the latest known live state for key.
//
// A missing key is a successful known absence answer when the cache is ready
// and the key belongs to the cache collection. Get never consults historical
// records; it reads only the latest live materialized collection. The served
// revision is the current collection boundary whether the key is live or absent.
func (c *Cache) Get(key objectstore.Key) (GetResult, error) {
	if c == nil {
		return GetResult{}, ErrInvalidCache
	}
	if !objectstore.KeyMatchesListRequest(key, c.collection) {
		return GetResult{}, ErrOutsideCollection
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.ready {
		return GetResult{}, ErrNotReady
	}
	result := GetResult{Key: key, Revision: c.latest.revision}
	item, ok := c.latest.item(key)
	if !ok {
		return result, nil
	}

	result.State = item.State
	result.Found = true

	return result, nil
}
