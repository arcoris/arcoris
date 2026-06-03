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

// equalList compares list payloads according to list semantic identity.
func (c *comparer) equalList(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	depth int,
) (bool, error) {
	if err := requireKind(path, oldValue, value.KindList, descriptor.Code()); err != nil {
		return false, err
	}
	if err := requireKind(path, newValue, value.KindList, descriptor.Code()); err != nil {
		return false, err
	}

	listView, ok := descriptor.List()
	if !ok {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a list")
	}

	oldList, _ := oldValue.List()
	newList, _ := newValue.List()
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
