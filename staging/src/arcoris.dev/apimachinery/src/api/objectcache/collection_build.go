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
	"fmt"

	"arcoris.dev/apimachinery/api/objectstore"
)

// buildCollection clones a store list result into cache-owned collection state.
// The invalid sentinel distinguishes snapshot construction errors from mutable
// cache construction or replacement errors.
func buildCollection(result objectstore.ListResult, invalid error) (collection, error) {
	cloned := result.Clone()
	items := make(map[objectstore.Key]objectstore.ListItem, len(cloned.Items))
	order := make([]objectstore.Key, 0, len(cloned.Items))
	idx := newIndexes()

	for _, item := range cloned.Items {
		if _, exists := items[item.Key]; exists {
			return collection{}, duplicateKeyError(invalid, item.Key)
		}
		items[item.Key] = item
		order = append(order, item.Key)
		idx.add(item)
	}

	if len(items) == 0 {
		items = nil
		idx = indexes{}
	}

	return collection{
		revision: cloned.Revision,
		order:    order,
		items:    items,
		indexes:  idx,
	}, nil
}

// duplicateKeyError preserves both the broad invalid-input sentinel and the
// duplicate-key sentinel for callers using errors.Is.
func duplicateKeyError(invalid error, key objectstore.Key) error {
	return fmt.Errorf(
		"%w: %w: %s",
		invalid,
		ErrDuplicateKey,
		key.String(),
	)
}
