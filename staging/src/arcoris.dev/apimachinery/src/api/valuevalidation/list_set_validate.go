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

// listSetElementKey is the concrete identity key for descriptor-approved set elements.
//
// The kind part prevents collisions between values such as string "1",
// integer 1, and bool true. The text part is produced only by exact scalar
// accessors, not by generic formatting.
type listSetElementKey struct {
	kind value.Kind
	text string
}

// validateListSet validates indexed items and then enforces concrete set uniqueness.
//
// api/types validates that ListSet descriptors are limited to stable scalar
// element descriptors. This package uses the same descriptor contract and owns
// the concrete duplicate check because only valuevalidation sees actual values.
func (v *validator) validateListSet(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Descriptor,
	depth int,
) {
	element, ok := v.resolveListSetElementDescriptor(path, element, depth)
	if !ok {
		return
	}

	v.validateIndexedList(path, valueView, element, depth)
	v.validateListSetDuplicates(path, valueView, element)
}

// resolveListSetElementDescriptor returns the terminal scalar descriptor for a set list.
func (v *validator) resolveListSetElementDescriptor(
	path fieldpath.Path,
	element types.Descriptor,
	depth int,
) (types.Descriptor, bool) {
	if element.Nullable() {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "ListSet element descriptor must be non-nullable")
		return types.Descriptor{}, false
	}

	resolved, err := v.refs.ResolveFinal(path, element, depth)
	if err != nil {
		v.addRefError(err)
		return types.Descriptor{}, false
	}
	if resolved.Nullable() {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "ListSet resolved element descriptor must be non-nullable")
		return types.Descriptor{}, false
	}
	if !isListSetElementDescriptor(resolved) {
		v.addf(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"ListSet element descriptor %s cannot be represented as a stable identity value",
			resolved.Code(),
		)
		return types.Descriptor{}, false
	}

	return resolved, true
}

// isListSetElementDescriptor reports whether descriptor kind has stable scalar identity.
func isListSetElementDescriptor(descriptor types.Descriptor) bool {
	switch descriptor.Code() {
	case types.DescriptorBool,
		types.DescriptorString,
		types.DescriptorInt8,
		types.DescriptorInt16,
		types.DescriptorInt32,
		types.DescriptorInt64,
		types.DescriptorUint8,
		types.DescriptorUint16,
		types.DescriptorUint32,
		types.DescriptorUint64:
		return true
	default:
		return false
	}
}

// validateListSetDuplicates reports repeated concrete set element identities.
func (v *validator) validateListSetDuplicates(
	path fieldpath.Path,
	valueView value.ListView,
	element types.Descriptor,
) {
	seen := make(map[listSetElementKey]int, valueView.Len())
	for i := 0; i < valueView.Len(); i++ {
		item, _ := valueView.At(i)
		key, ok := listSetElementKeyForValue(item, element)
		if !ok {
			continue
		}

		if first, exists := seen[key]; exists {
			v.addf(
				path.Index(i),
				ErrDuplicateListSetElement,
				ErrorReasonDuplicateListSetElement,
				"list set element duplicates value first seen at index %d",
				first,
			)
			continue
		}

		seen[key] = i
	}
}

// listSetElementKeyForValue builds a comparable identity key for supported scalars.
func listSetElementKeyForValue(item value.Value, descriptor types.Descriptor) (listSetElementKey, bool) {
	switch descriptor.Code() {
	case types.DescriptorBool:
		payload, ok := item.AsBool()
		if !ok {
			return listSetElementKey{}, false
		}
		if payload {
			return listSetElementKey{kind: value.KindBool, text: "true"}, true
		}
		return listSetElementKey{kind: value.KindBool, text: "false"}, true
	case types.DescriptorString:
		payload, ok := item.AsString()
		if !ok {
			return listSetElementKey{}, false
		}
		return listSetElementKey{kind: value.KindString, text: payload}, true
	case types.DescriptorInt8,
		types.DescriptorInt16,
		types.DescriptorInt32,
		types.DescriptorInt64,
		types.DescriptorUint8,
		types.DescriptorUint16,
		types.DescriptorUint32,
		types.DescriptorUint64:
		payload, ok := item.AsInteger()
		if !ok {
			return listSetElementKey{}, false
		}
		return listSetElementKey{kind: value.KindInteger, text: payload.String()}, true
	default:
		return listSetElementKey{}, false
	}
}
