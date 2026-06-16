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

// NewSnapshot builds an immutable snapshot from a store list result.
//
// The result is cloned before being retained. Duplicate storage keys are
// rejected because a snapshot represents one materialized live collection, not
// an arbitrary multiset of list items.
func NewSnapshot(result objectstore.ListResult) (Snapshot, error) {
	col, err := buildCollection(result, ErrInvalidSnapshot)
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{col: col}, nil
}

func cloneListItems(items []objectstore.ListItem) []objectstore.ListItem {
	if len(items) == 0 {
		return nil
	}

	out := make([]objectstore.ListItem, len(items))
	for i, item := range items {
		out[i] = item.Clone()
	}

	return out
}
