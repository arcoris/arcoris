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

package valuemerge

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/listmapkey"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// listMapEntry stores one concrete ListMap item after selector extraction.
type listMapEntry struct {
	selector  fieldpath.Selector
	indexPath fieldpath.Path
	item      value.Value
}

// listMapIndex stores ListMap items by selector while preserving physical order.
type listMapIndex struct {
	order  []string
	lookup map[string]listMapEntry
}

// listMapIndex extracts stable selector identity for all present list items.
func (m *merger) listMapIndex(
	path fieldpath.Path,
	list operand,
	element types.Descriptor,
	keys []types.FieldName,
) (listMapIndex, error) {
	if len(keys) == 0 {
		return listMapIndex{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"list map has no key fields",
		)
	}
	if list.Absent() || list.Value().IsNull() {
		return emptyListMapIndex(), nil
	}

	items := listItems(list)
	index := listMapIndex{
		order:  make([]string, 0, len(items)),
		lookup: make(map[string]listMapEntry, len(items)),
	}

	for i, item := range items {
		entry, err := m.listMapEntryAt(path, i, item, element, keys)
		if err != nil {
			return listMapIndex{}, err
		}
		if err := index.insert(path, entry); err != nil {
			return listMapIndex{}, err
		}
	}

	return index, nil
}

// emptyListMapIndex returns an initialized empty index.
func emptyListMapIndex() listMapIndex {
	return listMapIndex{lookup: map[string]listMapEntry{}}
}

// insert adds entry or reports duplicate selector identity.
func (i *listMapIndex) insert(path fieldpath.Path, entry listMapEntry) error {
	key := entry.selector.String()
	if previous, exists := i.lookup[key]; exists {
		return duplicateListMapEntryError(
			path.Select(entry.selector),
			previous.indexPath,
			entry.indexPath,
		)
	}

	i.order = append(i.order, key)
	i.lookup[key] = entry

	return nil
}

// listMapEntryAt extracts one item selector using the shared key helper.
func (m *merger) listMapEntryAt(
	path fieldpath.Path,
	index int,
	item value.Value,
	element types.Descriptor,
	keys []types.FieldName,
) (listMapEntry, error) {
	indexPath := path.Index(index)
	selector, err := listmapkey.ExtractSelector(
		indexPath,
		item,
		element,
		keys,
		listmapkey.Options{Resolver: m.resolver, MaxDepth: m.maxDepth},
	)
	if err != nil {
		return listMapEntry{}, mergeListMapKeyError(indexPath, err)
	}

	return listMapEntry{
		selector:  selector,
		indexPath: indexPath,
		item:      item,
	}, nil
}
