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

// list evaluates predicate over collection order.
//
// Indexes only reduce the candidate key set. The predicate remains the final
// semantic authority and every returned item is cloned.
func (col collection) list(predicate objectquery.Predicate) []objectstore.ListItem {
	if len(col.order) == 0 {
		return nil
	}

	plan := col.indexes.plan(predicate)
	out := make([]objectstore.ListItem, 0, len(col.order))
	for _, key := range col.order {
		if !plan.includes(key) {
			continue
		}
		item := col.items[key]
		if predicate.Match(item) {
			out = append(out, item.Clone())
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
