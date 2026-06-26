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

// buildCollection converts a validated ListResult into cache-owned state.
//
// It still checks duplicate keys because objectstore.ListResult validation
// verifies item shape and collection membership, but duplicate materialized
// keys are a cache read-model invariant.
func buildCollection(result objectstore.ListResult) (collection, error) {
	col := collection{revision: result.Revision}
	if len(result.Items) == 0 {
		return col, nil
	}

	col.order = make([]objectstore.Key, 0, len(result.Items))
	col.items = make(map[objectstore.Key]objectstore.ListItem, len(result.Items))

	for _, item := range result.Items {
		if _, exists := col.items[item.Key]; exists {
			return collection{}, errorWith(ErrInvalidRead, ErrDuplicateKey)
		}
		item = item.Clone()
		col.order = append(col.order, item.Key)
		col.items[item.Key] = item
	}

	return col, nil
}
