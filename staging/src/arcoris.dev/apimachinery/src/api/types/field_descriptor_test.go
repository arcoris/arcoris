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

func TestFieldDescriptorAccessorsDetachType(t *testing.T) {
	field := Field("name").String().Required().MinBytes(1).Description("name field").Field()
	desc := field.Descriptor()
	desc.string.minBytes = limit[int]{value: 99, set: true}

	requireEqual(t, field.Name(), FieldName("name"))
	requireEqual(t, field.Presence(), PresenceRequired)
	requireEqual(t, field.IsRequired(), true)
	requireEqual(t, field.IsOptional(), false)
	requireEqual(t, field.Description(), "name field")
	view, _ := field.Descriptor().AsString()
	min, _ := view.MinBytes()
	requireEqual(t, min, 1)
}

func TestFieldDescriptorShapeAndPresence(t *testing.T) {
	var zero FieldDescriptor
	requireEqual(t, zero.IsZero(), true)

	required := Field("name").String().Required().Nullable().Description("display name").Field()
	requireEqual(t, required.IsZero(), false)
	requireEqual(t, required.Name(), FieldName("name"))
	requireEqual(t, required.Presence(), PresenceRequired)
	requireEqual(t, required.IsRequired(), true)
	requireEqual(t, required.IsOptional(), false)
	requireEqual(t, required.Description(), "display name")
	requireCode(t, required.Descriptor(), DescriptorString)
	requireNullable(t, required.Descriptor(), true)

	optional := Field("enabled").Bool().Optional().Field()
	requireEqual(t, optional.Presence(), PresenceOptional)
	requireEqual(t, optional.IsRequired(), false)
	requireEqual(t, optional.IsOptional(), true)
	requireNullable(t, optional.Descriptor(), false)
}

func TestFieldDescriptorExpressionBoundaries(t *testing.T) {
	if _, ok := any(Field("name")).(FieldExpr); ok {
		t.Fatal("FieldBuilder must not implement FieldExpr")
	}
	if _, ok := any(Field("name").String()).(DescriptorExpr); ok {
		t.Fatal("field builders must not implement DescriptorExpr")
	}
}

func TestFieldDescriptorObjectValidationBoundaries(t *testing.T) {
	requireInvalidDescriptor(t, Object(FieldExpr(nil)).Descriptor(), nil, ErrInvalidField)

	missingPresence := Field("name").String().Field()
	requireInvalidDescriptor(t, objectTypeForField(missingPresence), nil, ErrInvalidField)

	duplicate := Object(
		Field("name").String().Required(),
		Field("name").Bool().Optional(),
	).Descriptor()
	requireInvalidDescriptor(t, duplicate, nil, ErrDuplicateField)
}

func TestFieldDescriptorObjectViewDetachAndOrder(t *testing.T) {
	desc := Object(
		Field("first").String().Required(),
		Field("second").Int64().Optional(),
	).Descriptor()

	fields := requireObjectView(t, desc).Fields()
	requireEqual(t, fields[0].Name(), FieldName("first"))
	requireEqual(t, fields[1].Name(), FieldName("second"))
	fields[0] = Field("changed").String().Required().Field()

	fields = requireObjectView(t, desc).Fields()
	requireEqual(t, fields[0].Name(), FieldName("first"))
}
