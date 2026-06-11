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

package valuevalidation_test

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateAcceptsValidScalarValues(t *testing.T) {
	tests := []struct {
		name    string
		payload value.Value
		shape   types.Descriptor
	}{
		{name: "bool", payload: value.BoolValue(true), shape: types.Bool().Descriptor()},
		{name: "string", payload: value.StringValue("main"), shape: types.String().MinBytes(1).Descriptor()},
		{name: "bytes", payload: value.BytesValue([]byte("abc")), shape: types.Bytes().MinBytes(3).Descriptor()},
		{name: "int64", payload: value.Int64Value(-1), shape: types.Int64().Min(-2).Descriptor()},
		{name: "uint64", payload: value.Uint64Value(3), shape: types.Uint64().Max(4).Descriptor()},
		{name: "float64", payload: mustFloat(t, 1.5), shape: types.Float64().Range(1, 2).Descriptor()},
		{name: "decimal", payload: mustDecimal(t, "12.30"), shape: types.Decimal().Precision(4).Scale(2).Descriptor()},
		{name: "date", payload: mustDate(t, 2026, 6, 1), shape: types.Date().Descriptor()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireNoError(
				t,
				valuevalidation.Validate(
					tt.payload,
					tt.shape,
					valuevalidation.Options{},
				),
			)
		})
	}
}

func TestValidateRejectsInvalidZeroValue(t *testing.T) {
	err := valuevalidation.Validate(
		value.Value{},
		types.String().Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrInvalidValue,
		valuevalidation.ErrorReasonInvalidZero,
		"$",
	)
}

func TestValidateRejectsInvalidZeroDescriptorDefensively(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("x"),
		types.Descriptor{},
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrInvalidDescriptor,
		valuevalidation.ErrorReasonInvalidDescriptor,
		"$",
	)
}

func TestValidateRejectsMapWithInvalidValueDescriptorDefensively(t *testing.T) {
	err := valuevalidation.Validate(
		mustObject(t, value.MustRecordMember("name", value.StringValue("main"))),
		types.MapOf(nil).Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrInvalidDescriptor,
		valuevalidation.ErrorReasonInvalidDescriptor,
		"$",
	)
}

func TestValidateRejectsListWithInvalidElementDescriptorDefensively(t *testing.T) {
	err := valuevalidation.Validate(
		mustList(t, value.StringValue("main")),
		types.ListOf(nil).Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrInvalidDescriptor,
		valuevalidation.ErrorReasonInvalidDescriptor,
		"$",
	)
}

func TestValidateRejectsKindMismatch(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("x"),
		types.Int64().Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$",
	)
}

func TestValidateAtUsesBasePath(t *testing.T) {
	basePath := fieldpath.Root().Field(fieldpath.MustFieldName("desired")).Field(fieldpath.MustFieldName("replicas"))
	err := valuevalidation.ValidateAt(
		basePath,
		value.StringValue("x"),
		types.Int32().Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$.desired.replicas",
	)
}

func TestValidateAtRejectsInvalidBasePath(t *testing.T) {
	basePath := fieldpath.Root().Append(fieldpath.Element{})
	err := valuevalidation.ValidateAt(
		basePath,
		value.StringValue("x"),
		types.String().Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrInvalidPath,
		valuevalidation.ErrorReasonInvalidPath,
		"$.<invalid>",
	)
}

func TestValidateCollectsMultipleErrors(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
		types.Field("replicas").Int32().Required(),
	).Descriptor()

	payload := mustObject(t)
	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireErrorCount(t, err, 2)
	requireError(
		t,
		err,
		valuevalidation.ErrMissingField,
		valuevalidation.ErrorReasonMissingField,
		"$.name",
	)
	requireError(
		t,
		err,
		valuevalidation.ErrMissingField,
		valuevalidation.ErrorReasonMissingField,
		"$.replicas",
	)
}

func TestValidateMaxErrorsOneIsFailFast(t *testing.T) {
	shape := types.Object(
		types.Field("a").String().Required(),
		types.Field("b").String().Required(),
	).Descriptor()

	err := valuevalidation.Validate(
		mustObject(t),
		shape,
		valuevalidation.Options{MaxErrors: 1},
	)

	requireErrorCount(t, err, 1)
	requireError(
		t,
		err,
		valuevalidation.ErrMissingField,
		valuevalidation.ErrorReasonMissingField,
		"$.a",
	)
}

func TestValidateHonorsMaxErrors(t *testing.T) {
	shape := types.Object(
		types.Field("a").String().Required(),
		types.Field("b").String().Required(),
		types.Field("c").String().Required(),
	).Descriptor()

	payload := mustObject(t)
	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{MaxErrors: 2},
	)

	requireErrorCount(t, err, 2)
}

func TestValidateHonorsMaxErrorsWithListMapFallback(t *testing.T) {
	shape := types.ListOf(
		types.Object(
			types.Field("type").String().Required(),
			types.Field("status").String().MinBytes(1).Required(),
		),
	).Map("type").Descriptor()
	payload := mustList(
		t,
		mustObject(t, value.MustRecordMember("status", value.StringValue(""))),
		mustObject(t, value.MustRecordMember("status", value.StringValue(""))),
	)

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{MaxErrors: 2},
	)

	requireErrorCount(t, err, 2)
}

func TestValidateHonorsMaxErrorsWithObjectAndMapErrors(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
		types.Field("labels").MapOf(types.String().MinBytes(1)).Required(),
	).Descriptor()
	payload := mustObject(
		t,
		value.MustRecordMember(
			"labels",
			mustObject(
				t,
				value.MustRecordMember("app", value.StringValue("")),
				value.MustRecordMember("tier", value.StringValue("")),
			),
		),
	)

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{MaxErrors: 2},
	)

	requireErrorCount(t, err, 2)
}
