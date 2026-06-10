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
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// compareList dispatches list comparison by descriptor list semantics.
//
// Atomic and set lists are one semantic field. Ordered lists use physical index
// paths. ListMap values use selector identity and must not fall back to index
// comparison when selector extraction fails.
func (c *comparer) compareList(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Descriptor,
	depth int,
) (Result, error) {
	if err := requireKind(path, oldValue, value.KindList, descriptor.Code()); err != nil {
		return Result{}, err
	}
	if err := requireKind(path, newValue, value.KindList, descriptor.Code()); err != nil {
		return Result{}, err
	}

	listView, ok := descriptor.AsList()
	if !ok {
		return Result{}, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a list")
	}

	element := listView.Element()
	if !element.IsValid() {
		return Result{}, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list element descriptor is invalid")
	}

	oldList, _ := oldValue.AsList()
	newList, _ := newValue.AsList()

	switch listView.Semantics() {
	case types.ListAtomic,
		types.ListSet:
		return c.compareWholeList(path, oldValue, newValue, descriptor, depth)
	case types.ListOrdered:
		return c.compareOrderedList(path, oldList, newList, element, depth)
	case types.ListMap:
		return c.compareListMap(path, oldList, newList, element, listView.MapKeys(), depth)
	default:
		return Result{}, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list semantics are invalid")
	}
}

// compareWholeList treats atomic and set-like lists as one semantic field.
//
// ListSet intentionally keeps whole-list comparison for now. valuevalidation
// enforces concrete scalar-set uniqueness, but compare/apply do not yet expose
// unordered set item paths or move semantics.
func (c *comparer) compareWholeList(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Descriptor,
	depth int,
) (Result, error) {
	equal, err := c.equalValue(path, oldValue, newValue, descriptor, depth)
	if err != nil {
		return Result{}, err
	}
	if equal {
		return EmptyResult(), nil
	}

	return EmptyResult().withModified(path)
}

// compareOrderedList compares list items by physical index.
//
// Ordered list semantics make indexes part of the API contract, so added,
// removed, and modified item paths are allowed to use path.Index(i).
func (c *comparer) compareOrderedList(
	path fieldpath.Path,
	oldList value.ListView,
	newList value.ListView,
	element types.Descriptor,
	depth int,
) (Result, error) {
	result := EmptyResult()
	oldLen := oldList.Len()
	newLen := newList.Len()
	maxLen := max(oldLen, newLen)

	for i := 0; i < maxLen; i++ {
		child, err := c.compare(
			path.Index(i),
			listOperand(oldList, i),
			listOperand(newList, i),
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

// listOperand converts an index lookup into presence-aware compare input.
func listOperand(list value.ListView, index int) valuepresence.Operand {
	val, ok := list.At(index)
	return valuepresence.From(val, ok)
}

// equalList compares list payloads according to descriptor list semantics.
//
// It is an equality-only companion to compareList. For ListSet, equality
// remains whole-list and order-sensitive even though valuevalidation enforces
// duplicate-free scalar set values.
func (c *comparer) equalList(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Descriptor,
	depth int,
) (bool, error) {
	if err := requireKind(path, oldValue, value.KindList, descriptor.Code()); err != nil {
		return false, err
	}
	if err := requireKind(path, newValue, value.KindList, descriptor.Code()); err != nil {
		return false, err
	}

	listView, ok := descriptor.AsList()
	if !ok {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a list")
	}

	oldList, _ := oldValue.AsList()
	newList, _ := newValue.AsList()
	element := listView.Element()

	switch listView.Semantics() {
	case types.ListAtomic,
		types.ListSet,
		types.ListOrdered:
		return c.equalListByIndex(path, oldList, newList, element, depth)
	case types.ListMap:
		return c.equalListMap(path, oldList, newList, element, listView.MapKeys(), depth)
	default:
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list semantics are invalid")
	}
}

// equalListByIndex compares list items in physical order.
//
// Atomic, set, and ordered whole-list equality all use this exact sequence
// comparison today. Only ListMap gets selector-based equality.
func (c *comparer) equalListByIndex(
	path fieldpath.Path,
	oldList value.ListView,
	newList value.ListView,
	element types.Descriptor,
	depth int,
) (bool, error) {
	n := oldList.Len()
	if n != newList.Len() {
		return false, nil
	}
	if !element.IsValid() {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list element descriptor is invalid")
	}

	for i := 0; i < n; i++ {
		oldItem, _ := oldList.At(i)
		newItem, _ := newList.At(i)
		equal, err := c.equalValue(path.Index(i), oldItem, newItem, element, depth+1)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}
