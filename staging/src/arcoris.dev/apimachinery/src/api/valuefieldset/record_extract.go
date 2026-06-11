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

package valuefieldset

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractRecord interprets value.KindRecord under a fixed-field object descriptor.
func (e *extractor) extractRecord(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) (fieldpath.Set, error) {
	if err := requireKind(path, val, value.KindRecord, descriptor.Code()); err != nil {
		return fieldpath.Set{}, err
	}

	objectView, ok := descriptor.AsObject()
	if !ok {
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not an object",
		)
	}

	valueView, _ := val.AsRecord()
	if valueView.IsEmpty() {
		return setAt(path)
	}

	fields := recordFieldsByName(objectView.Fields())
	var out setBuilder

	var extractErr error
	valueView.ForEach(func(_ int, recordMember value.RecordMember) bool {
		name := recordMember.Name.String()
		memberPath, err := recordMemberPath(path, name)
		if err != nil {
			extractErr = err
			return false
		}

		fieldDescriptor, declared := fields[name]

		if !declared {
			memberSet, err := e.extractUnknownRecordMember(
				memberPath,
				name,
				objectView.UnknownFields(),
			)
			if err != nil {
				extractErr = err
				return false
			}

			out.AddSet(memberSet)
			return true
		}

		memberSet, err := e.extract(
			memberPath,
			recordMember.Value,
			fieldDescriptor.Descriptor(),
			depth+1,
		)
		if err != nil {
			extractErr = err
			return false
		}

		out.AddSet(memberSet)
		return true
	})
	if extractErr != nil {
		return fieldpath.Set{}, extractErr
	}

	return out.Build(path)
}

// recordFieldsByName builds a declaration lookup for actual record members.
func recordFieldsByName(fields []types.FieldDescriptor) map[string]types.FieldDescriptor {
	declared := make(map[string]types.FieldDescriptor, len(fields))
	for _, fieldDescriptor := range fields {
		declared[string(fieldDescriptor.Name())] = fieldDescriptor
	}

	return declared
}

// extractUnknownRecordMember handles a record member not declared by an object descriptor.
func (e *extractor) extractUnknownRecordMember(
	path fieldpath.Path,
	name string,
	policy types.UnknownFieldPolicy,
) (fieldpath.Set, error) {
	switch policy {
	case types.UnknownReject:
		return fieldpath.Set{}, errorfAt(
			path,
			ErrUnknownField,
			ErrorReasonUnknownField,
			"record member %q is rejected by the object descriptor",
			name,
		)
	case types.UnknownPreserveOpaque:
		return setAt(path)
	case types.UnknownPrune:
		return fieldpath.EmptySet(), nil
	default:
		return fieldpath.Set{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"object unknown-field policy is invalid",
		)
	}
}
