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

// compareListMap compares items by selector identity, ignoring physical order.
func (c *comparer) compareListMap(
	path fieldpath.Path,
	oldList value.ListView,
	newList value.ListView,
	element types.Type,
	keys []types.FieldName,
	depth int,
) (Result, error) {
	oldEntries, err := c.listMapEntries(path, oldList, element, keys)
	if err != nil {
		return Result{}, err
	}
	newEntries, err := c.listMapEntries(path, newList, element, keys)
	if err != nil {
		return Result{}, err
	}

	result := EmptyResult()
	for _, key := range unionSortedListMapKeys(oldEntries, newEntries) {
		oldEntry, oldFound := oldEntries[key]
		newEntry, newFound := newEntries[key]
		selector := oldEntry.selector
		if !oldFound {
			selector = newEntry.selector
		}

		child, err := c.compare(
			path.Select(selector),
			listMapOperand(oldEntry, oldFound),
			listMapOperand(newEntry, newFound),
			element,
			depth+1,
		)
		if err != nil {
			return Result{}, err
		}

		result = result.merge(child)
	}

	return result, nil
}
