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

package valuefieldset

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractList interprets value.KindList through a list descriptor.
func (e *extractor) extractList(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Type,
	depth int,
) (fieldpath.Set, error) {
	if err := requireKind(path, val, value.KindList, descriptor.Code()); err != nil {
		return fieldpath.Set{}, err
	}

	listView, ok := descriptor.List()
	if !ok {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not a list",
		)
	}

	element := listView.Element()
	if !element.IsValid() {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"list element descriptor is invalid",
		)
	}

	valueView, _ := val.List()
	if valueView.IsEmpty() {
		return setAt(path)
	}

	switch listView.Semantics() {
	case types.ListAtomic,
		types.ListSet:
		return setAt(path)
	case types.ListOrdered:
		return e.extractIndexedList(path, valueView, element, depth)
	case types.ListMap:
		return e.extractListMap(path, valueView, element, listView.MapKeys(), depth)
	default:
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"list semantics are invalid",
		)
	}
}

// extractIndexedList extracts ordered-list item paths by physical index.
//
// Ordered list semantics explicitly make item position part of semantic
// addressing. Atomic and set-like lists intentionally do not use this helper
// because their future ownership/apply behavior treats the complete list as one
// field until a stable non-index identity model exists.
func (e *extractor) extractIndexedList(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Type,
	depth int,
) (fieldpath.Set, error) {
	out := fieldpath.EmptySet()
	for index := 0; index < valueView.Len(); index++ {
		item, _ := valueView.At(index)

		itemSet, err := e.extract(path.Index(index), item, element, depth+1)
		if err != nil {
			return fieldpath.Set{}, err
		}

		out = out.Union(itemSet)
	}

	return out, nil
}

// extractListMap extracts ListMap items through stable selector paths.
func (e *extractor) extractListMap(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Type,
	keys []types.FieldName,
	depth int,
) (fieldpath.Set, error) {
	if len(keys) == 0 {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"list map has no key fields",
		)
	}

	out := fieldpath.EmptySet()
	seen := make(map[string]fieldpath.Path, valueView.Len())

	for index := 0; index < valueView.Len(); index++ {
		item, _ := valueView.At(index)
		indexPath := path.Index(index)

		selector, err := e.listMapSelector(indexPath, item, element, keys)
		if err != nil {
			return fieldpath.Set{}, err
		}

		selectorPath := path.Select(selector)
		selectorKey := selector.String()
		if previous, exists := seen[selectorKey]; exists {
			return fieldpath.Set{}, errorfAt(
				selectorPath,
				ErrDuplicateListKey,
				ErrorReasonDuplicateListKey,
				"duplicate list map key; first occurrence at %s, duplicate at %s",
				previous,
				indexPath,
			)
		}
		seen[selectorKey] = indexPath

		itemSet, err := e.extract(selectorPath, item, element, depth+1)
		if err != nil {
			return fieldpath.Set{}, err
		}

		out = out.Union(itemSet)
	}

	return out, nil
}
