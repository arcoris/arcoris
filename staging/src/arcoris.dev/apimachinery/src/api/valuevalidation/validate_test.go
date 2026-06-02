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
		shape   types.Type
	}{
		{name: "bool", payload: value.BoolValue(true), shape: types.Bool().Type()},
		{name: "string", payload: value.StringValue("main"), shape: types.String().MinLen(1).Type()},
		{name: "bytes", payload: value.BytesValue([]byte("abc")), shape: types.Bytes().MinLen(3).Type()},
		{name: "int64", payload: value.Int64Value(-1), shape: types.Int64().Min(-2).Type()},
		{name: "uint64", payload: value.Uint64Value(3), shape: types.Uint64().Max(4).Type()},
		{name: "float64", payload: mustFloat(t, 1.5), shape: types.Float64().Range(1, 2).Type()},
		{name: "decimal", payload: mustDecimal(t, "12.30"), shape: types.Decimal().Precision(4).Scale(2).Type()},
		{name: "date", payload: mustDate(t, 2026, 6, 1), shape: types.Date().Type()},
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
		types.String().Type(),
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
		types.Type{},
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
		mustObject(t, value.ObjectMember("name", value.StringValue("main"))),
		types.MapOf(nil).Type(),
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
		types.ListOf(nil).Type(),
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
		types.Int64().Type(),
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
	basePath := fieldpath.RootPath().Field("desired").Field("replicas")
	err := valuevalidation.ValidateAt(
		basePath,
		value.StringValue("x"),
		types.Int32().Type(),
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

func TestValidateCollectsMultipleErrors(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
		types.Field("replicas").Int32().Required(),
	).Type()

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

func TestValidateHonorsMaxErrors(t *testing.T) {
	shape := types.Object(
		types.Field("a").String().Required(),
		types.Field("b").String().Required(),
		types.Field("c").String().Required(),
	).Type()

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
			types.Field("status").String().MinLen(1).Required(),
		),
	).Map("type").Type()
	payload := mustList(
		t,
		mustObject(t, value.ObjectMember("status", value.StringValue(""))),
		mustObject(t, value.ObjectMember("status", value.StringValue(""))),
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
		types.Field("labels").MapOf(types.String().MinLen(1)).Required(),
	).Type()
	payload := mustObject(
		t,
		value.ObjectMember(
			"labels",
			mustObject(
				t,
				value.ObjectMember("app", value.StringValue("")),
				value.ObjectMember("tier", value.StringValue("")),
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
