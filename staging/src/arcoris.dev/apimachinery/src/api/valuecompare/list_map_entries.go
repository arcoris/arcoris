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

// listMapEntries indexes list items by canonical selector string.
func (c *comparer) listMapEntries(
	path fieldpath.Path,
	list value.ListView,
	element types.Type,
	keys []types.FieldName,
) (map[string]listMapEntry, error) {
	if len(keys) == 0 {
		return nil, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list map has no key fields")
	}

	n := list.Len()
	if n == 0 {
		return nil, nil
	}

	entries := make(map[string]listMapEntry, n)
	for i := 0; i < n; i++ {
		entry, err := c.listMapEntryAt(path, list, i, element, keys)
		if err != nil {
			return nil, err
		}

		key := entry.selector.String()
		if previous, exists := entries[key]; exists {
			return nil, duplicateListMapEntryError(path.Select(entry.selector), previous.indexPath, entry.indexPath)
		}

		entries[key] = entry
	}

	return entries, nil
}

// listMapEntryAt extracts selector identity for one physical list item.
func (c *comparer) listMapEntryAt(
	path fieldpath.Path,
	list value.ListView,
	index int,
	element types.Type,
	keys []types.FieldName,
) (listMapEntry, error) {
	item, _ := list.At(index)
	indexPath := path.Index(index)
	selector, err := c.extractListMapSelector(indexPath, item, element, keys)
	if err != nil {
		return listMapEntry{}, err
	}

	return listMapEntry{
		selector:  selector,
		indexPath: indexPath,
		item:      item,
	}, nil
}
