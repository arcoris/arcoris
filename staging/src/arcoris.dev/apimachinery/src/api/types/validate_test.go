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
			Field("name").String().Required().MinLen(1),
			Field("replicas").Int32().Optional().Range(1, 10),
			Field("labels").MapOf(String()).Optional(),
		).Required().UnknownFields(UnknownReject),
		Field("status").Object(
			Field("conditions").ListOf(Object(
				Field("type").String().Required(),
				Field("status").String().Required(),
			)).Optional().Map("type"),
		).Optional(),
	).UnknownFields(UnknownReject).Type()

	requireNoError(t, ValidateType(validObject, nil))

	invalidNested := validObject
	invalidNested.object.fields[0].typ.object.fields[0].typ.string.minLen = intLimit{value: 4, set: true}
	invalidNested.object.fields[0].typ.object.fields[0].typ.string.maxLen = intLimit{value: 1, set: true}
	requireErrorIs(t, ValidateType(invalidNested, nil), ErrInvalidType)
}

func TestValidateTypeErrorWrappingAndPath(t *testing.T) {
	typ := Object(Field("name").String().Required()).Type()
	typ.object.fields[0].typ.string.minLen = intLimit{value: 5, set: true}
	typ.object.fields[0].typ.string.maxLen = intLimit{value: 1, set: true}

	err := ValidateType(typ, nil)
	requireErrorIs(t, err, ErrInvalidType)
	var typeErr *TypeError
	if !errors.As(err, &typeErr) {
		t.Fatalf("expected TypeError, got %T", err)
	}
	requireEqual(t, typeErr.Path, "type.fields[name].type.len")
}

func TestValidateDefinitionRejectsInvalidNameAndCycles(t *testing.T) {
	resolver := resolverFunc(func(name TypeName) (TypeDefinition, bool) {
		switch name {
		case "example.Name":
			return Define("example.Name", String()), true
		case "example.Self":
			return Define("example.Self", Ref("example.Self")), true
		default:
			return TypeDefinition{}, false
		}
	})

	requireErrorIs(t, ValidateDefinition(Define("bad", String()), resolver), ErrInvalidTypeReference)
	requireErrorIs(t, ValidateDefinition(Define("example.Self", Ref("example.Self")), resolver), ErrInvalidTypeReference)
}
