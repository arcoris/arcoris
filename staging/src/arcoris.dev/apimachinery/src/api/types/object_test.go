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

import "testing"

func TestObjectEmptyObjectValidAndUnknownDefault(t *testing.T) {
	typ := Object().Type()

	requireNoError(t, ValidateType(typ, nil))
	view, ok := typ.Object()
	requireEqual(t, ok, true)
	requireEqual(t, len(view.Fields()), 0)
	requireEqual(t, view.UnknownFields(), UnknownReject)
}

func TestObjectFieldOrderAndDetach(t *testing.T) {
	typ := Object(
		Field("first").String().Required(),
		Field("second").Int64().Optional(),
	).UnknownFields(UnknownPrune).Type()

	view, ok := typ.Object()
	requireEqual(t, ok, true)
	fields := view.Fields()
	requireEqual(t, fields[0].Name(), FieldName("first"))
	requireEqual(t, fields[1].Name(), FieldName("second"))
	requireEqual(t, view.UnknownFields(), UnknownPrune)

	fields[0].name = "changed"
	requireEqual(t, view.Fields()[0].Name(), FieldName("first"))
}

func TestObjectDuplicateFieldsRejected(t *testing.T) {
	typ := Object(
		Field("name").String().Required(),
		Field("name").String().Optional(),
	).Type()

	requireErrorIs(t, ValidateType(typ, nil), ErrDuplicateField)
}

func TestObjectInvalidUnknownPolicyRejected(t *testing.T) {
	typ := Object().Type()
	typ.object.unknown = UnknownFieldPolicy(99)

	requireErrorIs(t, ValidateType(typ, nil), ErrInvalidType)
}

func TestObjectNestedValidation(t *testing.T) {
	typ := Object(
		Field("spec").Object(
			Field("maxConcurrency").Int64().Required().Min(1),
			Field("image").String().Optional().MinLen(1),
		).Required().UnknownFields(UnknownReject),
	).Type()

	requireNoError(t, ValidateType(typ, nil))
}

func TestObjectTypeExprMarker(t *testing.T) {
	Object().typeExpr()
}
