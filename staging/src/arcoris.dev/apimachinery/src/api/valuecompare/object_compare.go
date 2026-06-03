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
	"slices"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// compareObject compares fixed object fields with field-path semantics.
//
// Declared fields use path.Field(name). Undeclared members are handled after
// declared fields so UnknownReject, UnknownPreserve, and UnknownPrune match the
// descriptor policy without affecting known-field traversal.
func (c *comparer) compareObject(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	depth int,
) (Result, error) {
	if err := requireKind(path, oldValue, value.KindObject, descriptor.Code()); err != nil {
		return Result{}, err
	}
	if err := requireKind(path, newValue, value.KindObject, descriptor.Code()); err != nil {
		return Result{}, err
	}

	objectView, ok := descriptor.Object()
	if !ok {
		return Result{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not an object",
		)
	}

	oldObject, _ := oldValue.Object()
	newObject, _ := newValue.Object()
	fields := objectFieldsByName(objectView.Fields())
	result := EmptyResult()

	for _, field := range objectView.Fields() {
		name := string(field.Name())
		fieldPath := path.Field(name)
		oldFieldValue, oldFound := oldObject.Get(name)
		newFieldValue, newFound := newObject.Get(name)

		child, err := c.compare(
			fieldPath,
			operand{value: oldFieldValue, present: oldFound},
			operand{value: newFieldValue, present: newFound},
			field.Type(),
			depth+1,
		)
		if err != nil {
			return Result{}, err
		}

		result = result.merge(child)
	}

	unknown, err := c.compareUnknownObjectMembers(
		path,
		oldObject,
		newObject,
		fields,
		objectView.UnknownFields(),
	)
	if err != nil {
		return Result{}, err
	}

	return result.merge(unknown), nil
}

// objectFieldsByName indexes declared fields so unknown-member passes can skip them.
func objectFieldsByName(fields []types.FieldDescriptor) map[string]types.FieldDescriptor {
	if len(fields) == 0 {
		return nil
	}

	declared := make(map[string]types.FieldDescriptor, len(fields))
	for _, field := range fields {
		declared[string(field.Name())] = field
	}

	return declared
}

// compareUnknownObjectMembers applies the descriptor's undeclared-member policy.
//
// UnknownReject fails fast, UnknownPreserve compares each unknown member as one
// opaque leaf, and UnknownPrune ignores unknown members completely.
func (c *comparer) compareUnknownObjectMembers(
	path fieldpath.Path,
	oldObject value.ObjectView,
	newObject value.ObjectView,
	declared map[string]types.FieldDescriptor,
	policy types.UnknownFieldPolicy,
) (Result, error) {
	switch policy {
	case types.UnknownReject:
		return c.rejectUnknownObjectMembers(path, oldObject, newObject, declared)
	case types.UnknownPreserve:
		return c.comparePreservedUnknownObjectMembers(path, oldObject, newObject, declared)
	case types.UnknownPrune:
		return EmptyResult(), nil
	default:
		return Result{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"object unknown-field policy is invalid",
		)
	}
}

// rejectUnknownObjectMembers reports the first rejected unknown field deterministically.
func (c *comparer) rejectUnknownObjectMembers(
	path fieldpath.Path,
	oldObject value.ObjectView,
	newObject value.ObjectView,
	declared map[string]types.FieldDescriptor,
) (Result, error) {
	for _, name := range unknownMemberNames(oldObject, newObject, declared) {
		return Result{}, errorfAt(
			path.Field(name),
			ErrUnknownField,
			ErrorReasonUnknownField,
			"field %q is rejected by the object descriptor",
			name,
		)
	}

	return EmptyResult(), nil
}

// comparePreservedUnknownObjectMembers compares each preserved unknown as one leaf.
//
// Unknown values have no descriptor, so nested object or list changes must not
// produce nested semantic paths. Only the unknown member path is added, removed,
// or modified.
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

// comparePreservedUnknownObjectMember compares one preserved unknown member.
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

// compareOpaqueLeaf marks only path when an opaque member's payload differs.
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

// unknownMemberNames returns deterministic undeclared names present on either side.
func unknownMemberNames(
	oldObject value.ObjectView,
	newObject value.ObjectView,
	declared map[string]types.FieldDescriptor,
) []string {
	seen := make(map[string]bool, oldObject.Len()+newObject.Len())
	addUnknownNames(seen, oldObject, declared)
	addUnknownNames(seen, newObject, declared)

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	slices.Sort(names)

	return names
}

// addUnknownNames records object member names that are not declared by descriptor.
func addUnknownNames(
	seen map[string]bool,
	object value.ObjectView,
	declared map[string]types.FieldDescriptor,
) {
	for _, member := range object.Members() {
		if _, ok := declared[member.Name]; !ok {
			seen[member.Name] = true
		}
	}
}
