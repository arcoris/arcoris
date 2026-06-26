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

// Get returns a detached state for key when the cache is ready and contains it.
func (c *Cache) Get(key objectstore.Key) (objectstore.State, bool) {
	if c == nil || !objectstore.KeyMatchesListRequest(key, c.collection) {
		return objectstore.State{}, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.ready {
		return objectstore.State{}, false
	}
	item, ok := c.col.item(key)
	if !ok {
		return objectstore.State{}, false
	}

	return item.State, true
}
