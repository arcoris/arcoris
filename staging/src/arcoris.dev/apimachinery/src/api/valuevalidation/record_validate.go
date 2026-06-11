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

// validateRecord interprets value.KindRecord under a fixed-field object descriptor.
func (v *validator) validateRecord(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) {
	if !v.requireKind(path, val, value.KindRecord, descriptor.Code()) {
		return
	}

	objectView, ok := descriptor.AsObject()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not an object")
		return
	}

	valueView, _ := val.AsRecord()
	fields := objectView.Fields()
	declared := make(map[string]types.FieldDescriptor, len(fields))

	for _, fieldDescriptor := range fields {
		if v.shouldStop() {
			return
		}

		name := string(fieldDescriptor.Name())
		fieldName, err := fieldpath.NewFieldName(name)
		if err != nil {
			v.addf(
				path,
				ErrInvalidDescriptor,
				ErrorReasonInvalidDescriptor,
				"object descriptor field name %q cannot become a field path element",
				name,
			)
			if v.shouldStop() {
				return
			}
			continue
		}
		declared[name] = fieldDescriptor

		memberValue, found := valueView.Get(value.MemberName(name))
		fieldPath := path.Field(fieldName)
		if !found {
			if fieldDescriptor.IsRequired() {
				v.addf(fieldPath, ErrMissingField, ErrorReasonMissingField, "required field %q is missing", name)
			}
			continue
		}

		v.validate(fieldPath, memberValue, fieldDescriptor.Descriptor(), depth+1)
	}

	switch objectView.UnknownFields() {
	case types.UnknownReject:
		v.validateUnknownRecordMembers(path, valueView, declared)
	case types.UnknownPreserveOpaque, types.UnknownPrune:
		return
	default:
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "object unknown-field policy is invalid")
	}
}

// validateUnknownRecordMembers reports undeclared record members under reject policy.
func (v *validator) validateUnknownRecordMembers(
	path fieldpath.Path,
	valueView value.RecordView,
	declared map[string]types.FieldDescriptor,
) {
	valueView.ForEach(func(_ int, recordMember value.RecordMember) bool {
		name := recordMember.Name.String()
		if _, ok := declared[name]; ok {
			return true
		}

		fieldName, err := fieldpath.NewFieldName(name)
		if err != nil {
			v.addf(
				path,
				ErrInvalidValue,
				ErrorReasonInvalidFieldName,
				"record member name %q cannot become a field path element",
				name,
			)
			return !v.shouldStop()
		}

		v.addf(
			path.Field(fieldName),
			ErrUnknownField,
			ErrorReasonUnknownField,
			"record member %q is not declared by the object descriptor",
			name,
		)
		return !v.shouldStop()
	})
}
