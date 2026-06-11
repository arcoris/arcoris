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
	"arcoris.dev/apimachinery/api/internal/listmapkey"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateList interprets value.KindList through a list descriptor.
func (v *validator) validateList(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) {
	if !v.requireKind(path, val, value.KindList, descriptor.Code()) {
		return
	}

	listView, ok := descriptor.AsList()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a list")
		return
	}

	element := listView.Element()
	if !element.IsValid() {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list element descriptor is invalid")
		return
	}

	valueView, _ := val.AsList()
	length := valueView.Len()
	if minItems, ok := listView.MinItems(); ok && length < minItems {
		v.addf(
			path,
			ErrLengthOutOfRange,
			ErrorReasonTooShort,
			"list length %d is below minimum %d",
			length,
			minItems,
		)
	}
	if maxItems, ok := listView.MaxItems(); ok && length > maxItems {
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
	case types.ListAtomic,
		types.ListOrdered:
		v.validateIndexedList(path, valueView, element, depth)
	case types.ListSet:
		v.validateListSet(path, valueView, element, depth)
	case types.ListMap:
		v.validateListMap(path, valueView, element, listView.MapKeys(), depth)
	default:
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list semantics are invalid")
	}
}

// validateIndexedList checks list item payloads by physical index.
//
// Validation uses indexes for atomic, ordered, and set-like lists because item
// locations produce better diagnostics. This intentionally differs from
// valuefieldset: ownership semantics treat atomic and set-like lists as one
// complete field, while validation still points users to the invalid item.
func (v *validator) validateIndexedList(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Descriptor,
	depth int,
) {
	for i := 0; i < valueView.Len(); i++ {
		if v.shouldStop() {
			return
		}

		item, _ := valueView.At(i)
		v.validate(path.Index(i), item, element, depth+1)
	}
}

// validateListMap checks ListMap items by stable selector identity.
func (v *validator) validateListMap(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Descriptor,
	keys []types.FieldName,
	depth int,
) {
	if len(keys) == 0 {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list map has no key fields")
		return
	}

	seen := make(map[string]fieldpath.Path, valueView.Len())
	for i := 0; i < valueView.Len(); i++ {
		if v.shouldStop() {
			return
		}

		item, _ := valueView.At(i)
		indexPath := path.Index(i)

		itemSelector, err := listmapkey.ExtractSelector(
			indexPath,
			item,
			element,
			keys,
			listmapkey.Options{
				Resolver: v.resolver,
				MaxDepth: v.maxDepth,
			},
		)
		if err != nil {
			if !v.reportListMapKeyFailure(err) {
				v.validate(indexPath, item, element, depth+1)
			}
			continue
		}

		selectorPath := path.Select(itemSelector)
		selectorIdentity := itemSelector.CanonicalText()
		if previous, exists := seen[selectorIdentity]; exists {
			v.addf(
				selectorPath,
				ErrDuplicateListKey,
				ErrorReasonDuplicateListKey,
				"duplicate list map key; first occurrence at %s, duplicate at %s",
				previous,
				indexPath,
			)
		} else {
			seen[selectorIdentity] = indexPath
		}

		v.validate(selectorPath, item, element, depth+1)
	}
}
