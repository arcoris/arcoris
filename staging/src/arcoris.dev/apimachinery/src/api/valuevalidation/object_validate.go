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

// validateObject interprets value.KindObject as a fixed-field object descriptor.
func (v *validator) validateObject(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) {
	if !v.requireKind(path, val, value.KindObject, descriptor.Code()) {
		return
	}

	objectView, ok := descriptor.AsObject()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not an object")
		return
	}

	valueView, _ := val.Object()
	fields := objectView.Fields()
	declared := make(map[string]types.FieldDescriptor, len(fields))

	for _, fieldDescriptor := range fields {
		name := string(fieldDescriptor.Name())
		declared[name] = fieldDescriptor

		memberValue, found := valueView.Get(name)
		fieldPath := path.Field(name)
		if !found {
			if fieldDescriptor.IsRequired() {
				v.addf(fieldPath, ErrMissingField, ErrorReasonMissingField, "required field %q is missing", name)
			}
			continue
		}

		v.validate(fieldPath, memberValue, fieldDescriptor.Descriptor(), depth+1)
	}

	if objectView.UnknownFields() == types.UnknownReject {
		v.validateUnknownObjectMembers(path, valueView, declared)
	}
}

// validateUnknownObjectMembers reports undeclared members under reject policy.
func (v *validator) validateUnknownObjectMembers(
	path fieldpath.Path,
	valueView value.ObjectView,
	declared map[string]types.FieldDescriptor,
) {
	for _, objectMember := range valueView.Members() {
		if _, ok := declared[objectMember.Name]; ok {
			continue
		}

		v.addf(
			path.Field(objectMember.Name),
			ErrUnknownField,
			ErrorReasonUnknownField,
			"field %q is not declared by the object descriptor",
			objectMember.Name,
		)
	}
}
