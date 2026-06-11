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
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// compareRecord compares declared object fields against record payload members.
//
// Declared fields use path.Field(name). Undeclared members are handled after
// declared fields so UnknownReject, UnknownPreserveOpaque, and UnknownPrune match the
// descriptor policy without affecting known-field traversal.
func (c *comparer) compareRecord(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Descriptor,
	depth int,
) (Result, error) {
	if err := requireKind(path, oldValue, value.KindRecord, descriptor.Code()); err != nil {
		return Result{}, err
	}
	if err := requireKind(path, newValue, value.KindRecord, descriptor.Code()); err != nil {
		return Result{}, err
	}

	objectView, ok := descriptor.AsObject()
	if !ok {
		return Result{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not an object",
		)
	}

	oldRecord, _ := oldValue.AsRecord()
	newRecord, _ := newValue.AsRecord()
	fields := recordFieldsByName(objectView.Fields())
	var result resultBuilder

	for _, field := range objectView.Fields() {
		name := string(field.Name())
		fieldPath, err := recordFieldPath(path, name)
		if err != nil {
			return Result{}, err
		}
		oldFieldValue, oldFound := oldRecord.Get(value.MemberName(name))
		newFieldValue, newFound := newRecord.Get(value.MemberName(name))

		child, err := c.compare(
			fieldPath,
			valuepresence.From(oldFieldValue, oldFound),
			valuepresence.From(newFieldValue, newFound),
			field.Descriptor(),
			depth+1,
		)
		if err != nil {
			return Result{}, err
		}

		result.AddResult(child)
	}

	unknown, err := c.compareUnknownRecordMembers(
		path,
		oldRecord,
		newRecord,
		fields,
		objectView.UnknownFields(),
	)
	if err != nil {
		return Result{}, err
	}

	result.AddResult(unknown)
	return result.Build()
}

// recordFieldsByName indexes declared fields so unknown-member passes can skip them.
func recordFieldsByName(fields []types.FieldDescriptor) map[string]types.FieldDescriptor {
	if len(fields) == 0 {
		return nil
	}

	declared := make(map[string]types.FieldDescriptor, len(fields))
	for _, field := range fields {
		declared[string(field.Name())] = field
	}

	return declared
}

// compareUnknownRecordMembers applies the descriptor's undeclared-member policy.
//
// UnknownReject fails fast, UnknownPreserveOpaque compares each unknown member as one
// opaque leaf, and UnknownPrune ignores unknown members completely.
func (c *comparer) compareUnknownRecordMembers(
	path fieldpath.Path,
	oldRecord value.RecordView,
	newRecord value.RecordView,
	declared map[string]types.FieldDescriptor,
	policy types.UnknownFieldPolicy,
) (Result, error) {
	switch policy {
	case types.UnknownReject:
		return c.rejectUnknownRecordMembers(path, oldRecord, newRecord, declared)
	case types.UnknownPreserveOpaque:
		return c.comparePreservedUnknownRecordMembers(path, oldRecord, newRecord, declared)
	case types.UnknownPrune:
		return EmptyResult(), nil
	default:
		return Result{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor unknown-field policy is invalid",
		)
	}
}

// rejectUnknownRecordMembers reports the first rejected unknown field deterministically.
func (c *comparer) rejectUnknownRecordMembers(
	path fieldpath.Path,
	oldRecord value.RecordView,
	newRecord value.RecordView,
	declared map[string]types.FieldDescriptor,
) (Result, error) {
	for _, name := range unknownMemberNames(oldRecord, newRecord, declared) {
		fieldPath, err := recordFieldPath(path, name)
		if err != nil {
			return Result{}, err
		}

		return Result{}, errorfAt(
			fieldPath,
			ErrUnknownField,
			ErrorReasonUnknownField,
			"record member %q is rejected by the object descriptor",
			name,
		)
	}

	return EmptyResult(), nil
}

// comparePreservedUnknownRecordMembers compares each preserved unknown as one leaf.
//
// Unknown values have no descriptor, so nested object or list changes must not
// produce nested semantic paths. Only the unknown member path is added, removed,
// or modified.
func (c *comparer) comparePreservedUnknownRecordMembers(
	path fieldpath.Path,
	oldRecord value.RecordView,
	newRecord value.RecordView,
	declared map[string]types.FieldDescriptor,
) (Result, error) {
	var result resultBuilder

	for _, name := range unknownMemberNames(oldRecord, newRecord, declared) {
		fieldPath, err := recordFieldPath(path, name)
		if err != nil {
			return Result{}, err
		}

		child, err := c.comparePreservedUnknownRecordMember(fieldPath, oldRecord, newRecord, name)
		if err != nil {
			return Result{}, err
		}

		result.AddResult(child)
	}

	return result.Build()
}

// comparePreservedUnknownRecordMember compares one preserved unknown member.
func (c *comparer) comparePreservedUnknownRecordMember(
	path fieldpath.Path,
	oldRecord value.RecordView,
	newRecord value.RecordView,
	name string,
) (Result, error) {
	oldMember, oldFound := oldRecord.Get(value.MemberName(name))
	newMember, newFound := newRecord.Get(value.MemberName(name))

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
	oldRecord value.RecordView,
	newRecord value.RecordView,
	declared map[string]types.FieldDescriptor,
) []string {
	seen := make(map[string]bool, oldRecord.Len()+newRecord.Len())
	addUnknownNames(seen, oldRecord, declared)
	addUnknownNames(seen, newRecord, declared)

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	slices.Sort(names)

	return names
}

// addUnknownNames records record member names that are not declared by descriptor.
func addUnknownNames(
	seen map[string]bool,
	record value.RecordView,
	declared map[string]types.FieldDescriptor,
) {
	record.ForEach(func(_ int, member value.RecordMember) bool {
		name := member.Name.String()
		if _, ok := declared[name]; !ok {
			seen[name] = true
		}
		return true
	})
}
