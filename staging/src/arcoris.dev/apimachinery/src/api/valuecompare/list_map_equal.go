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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// equalListMap compares associative-list items by selector identity.
//
// Physical order is ignored, matching compareListMap. Iteration uses sorted
// selector strings so traversal errors stay deterministic.
func (c *comparer) equalListMap(
	path fieldpath.Path,
	oldList value.ListView,
	newList value.ListView,
	element types.Type,
	keys []types.FieldName,
	depth int,
) (bool, error) {
	oldEntries, err := c.listMapEntries(path, oldList, element, keys)
	if err != nil {
		return false, err
	}
	newEntries, err := c.listMapEntries(path, newList, element, keys)
	if err != nil {
		return false, err
	}
	if len(oldEntries) != len(newEntries) {
		return false, nil
	}

	for _, key := range sortedMapKeys(oldEntries) {
		oldEntry := oldEntries[key]
		newEntry, found := newEntries[key]
		if !found {
			return false, nil
		}

		equal, err := c.equalValue(path.Select(oldEntry.selector), oldEntry.item, newEntry.item, element, depth+1)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}
