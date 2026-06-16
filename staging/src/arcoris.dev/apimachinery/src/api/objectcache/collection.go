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

// collection is the cache-owned materialized state shared by Snapshot and Cache.
//
// order defines public output order. items owns the current live item for every
// ordered key. indexes are acceleration only and must never be used as a
// semantic substitute for objectquery.Predicate.Match.
type collection struct {
	revision objectstore.Revision
	order    []objectstore.Key
	items    map[objectstore.Key]objectstore.ListItem
	indexes  indexes
}

func (col collection) isZero() bool {
	return len(col.order) == 0 && col.revision.IsZero()
}

func (col collection) len() int {
	return len(col.order)
}

func (col collection) item(key objectstore.Key) (objectstore.ListItem, bool) {
	item, ok := col.items[key]
	if !ok {
		return objectstore.ListItem{}, false
	}

	return item.Clone(), true
}

func (col collection) listItems() []objectstore.ListItem {
	if len(col.order) == 0 {
		return nil
	}

	out := make([]objectstore.ListItem, 0, len(col.order))
	for _, key := range col.order {
		out = append(out, col.items[key].Clone())
	}

	return out
}
