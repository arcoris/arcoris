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

package types

import (
	"errors"
	"testing"
)

func TestValidateTypeBroadInvariantMatrix(t *testing.T) {
	validObject := Object(
		Field("spec").Object(
			Field("name").String().Required().MinBytes(1),
			Field("replicas").Int32().Optional().Range(1, 10),
			Field("labels").MapOf(String()).Optional(),
		).Required().UnknownFields(UnknownReject),
		Field("status").Object(
			Field("conditions").ListOf(Object(
				Field("type").String().Required(),
				Field("status").String().Required(),
			)).Optional().Map("type"),
		).Optional(),
	).UnknownFields(UnknownReject).Descriptor()

	requireNoError(t, ValidateLocal(validObject))

	invalidNested := validObject
	invalidNested.object.fields[0].descriptor.object.fields[0].descriptor.string.minBytes = limit[int]{value: 4, set: true}
	invalidNested.object.fields[0].descriptor.object.fields[0].descriptor.string.maxBytes = limit[int]{value: 1, set: true}
	requireErrorIs(t, ValidateLocal(invalidNested), ErrInvalidDescriptor)
}

func TestValidateDescriptorErrorWrappingAndPath(t *testing.T) {
	desc := Object(Field("name").String().Required()).Descriptor()
	desc.object.fields[0].descriptor.string.minBytes = limit[int]{value: 5, set: true}
	desc.object.fields[0].descriptor.string.maxBytes = limit[int]{value: 1, set: true}

	err := ValidateLocal(desc)
	requireErrorIs(t, err, ErrInvalidDescriptor)
	var descriptorErr *DescriptorError
	if !errors.As(err, &descriptorErr) {
		t.Fatalf("expected DescriptorError, got %T", err)
	}
	requireEqual(t, descriptorErr.Path, "descriptor.fields[name].type.bytes")
}

func TestValidateDefinitionRejectsInvalidNameAndCycles(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (Definition, bool) {
		switch name {
		case "example.Name":
			return Define("example.Name", String()), true
		case "example.Self":
			return Define("example.Self", Ref("example.Self")), true
		default:
			return Definition{}, false
		}
	})

	requireErrorIs(t, ValidateDefinitionResolved(Define("bad", String()), resolver), ErrInvalidDescriptorReference)
	requireErrorIs(t, ValidateDefinitionResolved(Define("example.Self", Ref("example.Self")), resolver), ErrInvalidDescriptorReference)
}
