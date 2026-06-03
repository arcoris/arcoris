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

// equalObject reports descriptor-aware object equality without building diff sets.
func (c *comparer) equalObject(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	depth int,
) (bool, error) {
	if err := requireKind(path, oldValue, value.KindObject, descriptor.Code()); err != nil {
		return false, err
	}
	if err := requireKind(path, newValue, value.KindObject, descriptor.Code()); err != nil {
		return false, err
	}

	objectView, ok := descriptor.Object()
	if !ok {
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not an object")
	}

	oldObject, _ := oldValue.Object()
	newObject, _ := newValue.Object()
	fields := objectView.Fields()

	for _, field := range fields {
		name := string(field.Name())
		oldMember, oldFound := oldObject.Get(name)
		newMember, newFound := newObject.Get(name)
		if oldFound != newFound {
			return false, nil
		}
		if !oldFound {
			continue
		}

		equal, err := c.equalValue(path.Field(name), oldMember, newMember, field.Type(), depth+1)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return c.equalUnknownObjectMembers(
		path,
		oldObject,
		newObject,
		objectFieldsByName(fields),
		objectView.UnknownFields(),
	)
}

// equalUnknownObjectMembers applies unknown-field policy in equality-only mode.
func (c *comparer) equalUnknownObjectMembers(
	path fieldpath.Path,
	oldObject value.ObjectView,
	newObject value.ObjectView,
	declared map[string]types.FieldDescriptor,
	policy types.UnknownFieldPolicy,
) (bool, error) {
	switch policy {
	case types.UnknownReject:
		_, err := c.rejectUnknownObjectMembers(path, oldObject, newObject, declared)
		return err == nil, err
	case types.UnknownPreserve:
		return c.equalPreservedUnknownObjectMembers(path, oldObject, newObject, declared)
	case types.UnknownPrune:
		return true, nil
	default:
		return false, errorAt(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "object unknown-field policy is invalid")
	}
}

// equalPreservedUnknownObjectMembers compares preserved unknown fields as leaves.
func (c *comparer) equalPreservedUnknownObjectMembers(
	path fieldpath.Path,
	oldObject value.ObjectView,
	newObject value.ObjectView,
	declared map[string]types.FieldDescriptor,
) (bool, error) {
	for _, name := range unknownMemberNames(oldObject, newObject, declared) {
		oldMember, oldFound := oldObject.Get(name)
		newMember, newFound := newObject.Get(name)
		if oldFound != newFound {
			return false, nil
		}

		equal, err := c.equalOpaqueValue(path.Field(name), oldMember, newMember)
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}
