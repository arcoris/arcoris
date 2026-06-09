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
	desc := Object().Descriptor()

	requireNoError(t, ValidateLocal(desc))
	view, ok := desc.AsObject()
	requireEqual(t, ok, true)
	requireEqual(t, len(view.Fields()), 0)
	requireEqual(t, view.UnknownFields(), UnknownReject)
}

func TestObjectFieldOrderAndDetach(t *testing.T) {
	desc := Object(
		Field("first").String().Required(),
		Field("second").Int64().Optional(),
	).UnknownFields(UnknownPrune).Descriptor()

	view, ok := desc.AsObject()
	requireEqual(t, ok, true)
	fields := view.Fields()
	requireEqual(t, fields[0].Name(), FieldName("first"))
	requireEqual(t, fields[1].Name(), FieldName("second"))
	requireEqual(t, view.UnknownFields(), UnknownPrune)

	fields[0].name = "changed"
	requireEqual(t, view.Fields()[0].Name(), FieldName("first"))
}

func TestObjectDuplicateFieldsRejected(t *testing.T) {
	desc := Object(
		Field("name").String().Required(),
		Field("name").String().Optional(),
	).Descriptor()

	requireErrorIs(t, ValidateLocal(desc), ErrDuplicateField)
}

func TestObjectInvalidUnknownPolicyRejected(t *testing.T) {
	desc := Object().Descriptor()
	desc.object.unknown = UnknownFieldPolicy(99)

	requireErrorIs(t, ValidateLocal(desc), ErrInvalidDescriptor)
}

func TestObjectNestedValidation(t *testing.T) {
	desc := Object(
		Field("spec").Object(
			Field("maxConcurrency").Int64().Required().Min(1),
			Field("image").String().Optional().MinBytes(1),
		).Required().UnknownFields(UnknownReject),
	).Descriptor()

	requireNoError(t, ValidateLocal(desc))
}

func TestObjectDescriptorExprMarker(t *testing.T) {
	Object().descriptorExpr()
}
