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

// comparePreservedUnknownObjectMembers treats each unknown member as an opaque leaf.
func (c *comparer) comparePreservedUnknownObjectMembers(
	path fieldpath.Path,
	oldObject value.ObjectView,
	newObject value.ObjectView,
	declared map[string]types.FieldDescriptor,
) (Result, error) {
	result := EmptyResult()

	for _, name := range unknownMemberNames(oldObject, newObject, declared) {
		child, err := c.comparePreservedUnknownObjectMember(path.Field(name), oldObject, newObject, name)
		if err != nil {
			return Result{}, err
		}

		result = result.merge(child)
	}

	return result, nil
}

// comparePreservedUnknownObjectMember compares one opaque unknown member.
func (c *comparer) comparePreservedUnknownObjectMember(
	path fieldpath.Path,
	oldObject value.ObjectView,
	newObject value.ObjectView,
	name string,
) (Result, error) {
	oldMember, oldFound := oldObject.Get(name)
	newMember, newFound := newObject.Get(name)

	switch {
	case !oldFound && newFound:
		set, err := setAt(path)
		if err != nil {
			return Result{}, err
		}
		return EmptyResult().withAdded(set), nil
	case oldFound && !newFound:
		set, err := setAt(path)
		if err != nil {
			return Result{}, err
		}
		return EmptyResult().withRemoved(set), nil
	case oldFound && newFound:
		return c.compareOpaqueLeaf(path, oldMember, newMember)
	default:
		return EmptyResult(), nil
	}
}

// compareOpaqueLeaf marks path modified when same unknown member payload differs.
func (c *comparer) compareOpaqueLeaf(path fieldpath.Path, oldMember value.Value, newMember value.Value) (Result, error) {
	equal, err := c.equalOpaqueValue(path, oldMember, newMember)
	if err != nil {
		return Result{}, err
	}
	if equal {
		return EmptyResult(), nil
	}

	return EmptyResult().withModified(path)
}
