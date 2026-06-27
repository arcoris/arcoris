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
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

// Query evaluates predicate over the latest live collection.
//
// Query accepts an already-compiled objectquery.Predicate so callers retain
// responsibility for query construction, selectable-field registration, and
// validation through objectquery.Compile. It scans latest live items in cache
// order, uses Predicate.Match as the semantic source of truth, and returns
// detached matching items at the current cache revision. Historical records,
// tombstones, and future index planning hints are not consulted.
func (c *Cache) Query(predicate objectquery.Predicate) (objectstore.ListResult, error) {
	if c == nil {
		return objectstore.ListResult{}, ErrInvalidCache
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.ready {
		return objectstore.ListResult{}, ErrNotReady
	}

	var items []objectstore.ListItem
	for _, key := range c.latest.order {
		item := c.latest.items[key]
		if predicate.Match(item) {
			items = append(items, item.Clone())
		}
	}

	return objectstore.ListResult{Items: items, Revision: c.latest.revision}, nil
}
