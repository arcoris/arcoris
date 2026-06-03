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

// compareList interprets value.KindList through a list descriptor.
func (c *comparer) compareList(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	depth int,
) (Result, error) {
	if err := requireKind(path, oldValue, value.KindList, descriptor.Code()); err != nil {
		return Result{}, err
	}
	if err := requireKind(path, newValue, value.KindList, descriptor.Code()); err != nil {
		return Result{}, err
	}

	listView, ok := descriptor.List()
	if !ok {
		return Result{}, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a list")
	}

	element := listView.Element()
	if !element.IsValid() {
		return Result{}, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "list element descriptor is invalid")
	}

	oldList, _ := oldValue.List()
	newList, _ := newValue.List()

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
// ListSet intentionally keeps whole-list equality for now. Stable item
// identity for arbitrary set elements is not defined yet, so emitting item
// paths would be misleading.
func (c *comparer) compareWholeList(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
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
func (c *comparer) compareOrderedList(
	path fieldpath.Path,
	oldList value.ListView,
	newList value.ListView,
	element types.Type,
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

// listOperand builds an operand from a list index.
func listOperand(list value.ListView, index int) operand {
	val, ok := list.At(index)
	return operand{value: val, present: ok}
}
