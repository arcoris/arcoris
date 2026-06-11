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
//
// Field-set extraction follows ownership/merge semantics, not diagnostic
// precision. Atomic and set-like lists produce the list path. Ordered lists
// produce index paths because index position is part of the API contract.
// ListMap values produce selector paths from declared key fields.
func (e *extractor) extractList(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) (fieldpath.Set, error) {
	if err := requireKind(path, val, value.KindList, descriptor.Code()); err != nil {
		return fieldpath.Set{}, err
	}

	listView, ok := descriptor.AsList()
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

	valueView, _ := val.AsList()
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
// addressing. Atomic lists are one field. Set-like lists also stay at the list
// path here: valuevalidation owns concrete scalar-set uniqueness, while field
// ownership and apply remain conservative until set item paths are designed.
func (e *extractor) extractIndexedList(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Descriptor,
	depth int,
) (fieldpath.Set, error) {
	var out setBuilder
	for i := 0; i < valueView.Len(); i++ {
		item, _ := valueView.At(i)

		itemSet, err := e.extract(path.Index(i), item, element, depth+1)
		if err != nil {
			return fieldpath.Set{}, err
		}

		out.AddSet(itemSet)
	}

	return out.Build(path)
}

// extractListMap extracts ListMap items through stable selector paths.
func (e *extractor) extractListMap(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Descriptor,
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

	var out setBuilder
	seen := make(map[string]fieldpath.Path, valueView.Len())

	for i := 0; i < valueView.Len(); i++ {
		item, _ := valueView.At(i)
		indexPath := path.Index(i)

		selector, err := e.listMapSelector(indexPath, item, element, keys)
		if err != nil {
			return fieldpath.Set{}, err
		}

		selectorPath := path.Select(selector)
		selectorKey := selector.CanonicalText()
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

		out.AddSet(itemSet)
	}

	return out.Build(path)
}
