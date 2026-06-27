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

// ensureItems initializes the item map for create-after-empty transitions.
//
// Empty collections keep nil maps until a create arrives; this avoids
// allocation for ready-empty caches while preserving simple mutation code.
func (col *collection) ensureItems() {
	if col.items == nil {
		col.items = map[objectstore.Key]objectstore.ListItem{}
	}
}

// removeOrderKey removes a known key while preserving survivor order.
//
// The removed slot is zeroed before reslicing so an obsolete key graph is not
// retained by the backing array.
func (col *collection) removeOrderKey(key objectstore.Key) {
	for i, existing := range col.order {
		if existing.Equal(key) {
			copy(col.order[i:], col.order[i+1:])
			last := len(col.order) - 1
			col.order[last] = objectstore.Key{}
			col.order = col.order[:last]
			return
		}
	}
}
