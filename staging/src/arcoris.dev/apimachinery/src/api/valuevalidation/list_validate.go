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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateList interprets value.KindList through a list descriptor.
func (v *validator) validateList(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Type,
	depth int,
) {
	if !v.requireKind(path, val, value.KindList, descriptor.Code()) {
		return
	}

	listView, ok := descriptor.List()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a list")
		return
	}

	element := listView.Element()
	if !element.IsValid() {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list element descriptor is invalid")
		return
	}

	valueView, _ := val.List()
	length := valueView.Len()
	if minItems, ok := listView.MinLen(); ok && length < minItems {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooShort,
			"list length %d is below minimum %d",
			length,
			minItems,
		)
	}
	if maxItems, ok := listView.MaxLen(); ok && length > maxItems {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooLong,
			"list length %d is above maximum %d",
			length,
			maxItems,
		)
	}

	switch listView.Semantics() {
	case types.ListAtomic, types.ListSet:
		v.validateIndexedList(path, valueView, element, depth)
	case types.ListMap:
		v.validateListMap(path, valueView, element, listView.MapKeys(), depth)
	default:
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list semantics are invalid")
	}
}

// validateIndexedList checks ordinary and set-like lists by physical index.
func (v *validator) validateIndexedList(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Type,
	depth int,
) {
	for index := 0; index < valueView.Len(); index++ {
		item, _ := valueView.At(index)
		v.validate(path.Index(index), item, element, depth+1)
	}
}

// validateListMap checks associative lists by stable selector identity.
func (v *validator) validateListMap(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Type,
	keys []types.FieldName,
	depth int,
) {
	if len(keys) == 0 {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list map has no key fields")
		return
	}

	seen := make(map[string]fieldpath.Path, valueView.Len())
	for index := 0; index < valueView.Len(); index++ {
		item, _ := valueView.At(index)
		indexPath := path.Index(index)

		selector, ok := v.tryListMapSelector(indexPath, item, element, keys)
		if !ok {
			v.validate(indexPath, item, element, depth+1)
			continue
		}

		selectorPath := path.Select(selector)
		selectorKey := selector.String()
		if previous, exists := seen[selectorKey]; exists {
			v.addf(
				selectorPath,
				ErrDuplicateListKey,
				ErrorReasonDuplicateListKey,
				"duplicate list map key; first occurrence at %s, duplicate at %s",
				previous,
				indexPath,
			)
		} else {
			seen[selectorKey] = indexPath
		}

		v.validate(selectorPath, item, element, depth+1)
	}
}
