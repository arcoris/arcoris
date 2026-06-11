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

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateRecordRequiredFields(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
	).Descriptor()

	err := valuevalidation.Validate(
		mustObject(t),
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrMissingField,
		valuevalidation.ErrorReasonMissingField,
		"$.name",
	)
}

func TestValidateRecordOptionalFields(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Optional(),
	).Descriptor()

	requireNoError(
		t,
		valuevalidation.Validate(
			mustObject(t),
			shape,
			valuevalidation.Options{},
		),
	)
}

func TestValidateRecordNullability(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
		types.Field("note").String().Nullable().Optional(),
	).Descriptor()

	payload := mustObject(
		t,
		value.MustRecordMember("name", value.StringValue("main")),
		value.MustRecordMember("note", value.NullValue()),
	)

	requireNoError(
		t,
		valuevalidation.Validate(
			payload,
			shape,
			valuevalidation.Options{},
		),
	)
}

func TestValidateRecordUnknownFieldRejected(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
	).UnknownFields(types.UnknownReject).Descriptor()

	payload := mustObject(
		t,
		value.MustRecordMember("name", value.StringValue("main")),
		value.MustRecordMember("extra", value.StringValue("x")),
	)

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrUnknownField,
		valuevalidation.ErrorReasonUnknownField,
		"$.extra",
	)
}

func TestValidateRecordUnknownFieldPruneAllowed(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
	).UnknownFields(types.UnknownPrune).Descriptor()

	payload := mustObject(
		t,
		value.MustRecordMember("name", value.StringValue("main")),
		value.MustRecordMember("extra", value.StringValue("x")),
	)

	requireNoError(
		t,
		valuevalidation.Validate(
			payload,
			shape,
			valuevalidation.Options{},
		),
	)
}

func TestValidateRecordUnknownFieldPreserveOpaqueAllowed(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
	).UnknownFields(types.UnknownPreserveOpaque).Descriptor()

	payload := mustObject(
		t,
		value.MustRecordMember("name", value.StringValue("main")),
		value.MustRecordMember("extra", mustObject(t, value.MustRecordMember("nested", value.StringValue("x")))),
	)

	requireNoError(
		t,
		valuevalidation.Validate(
			payload,
			shape,
			valuevalidation.Options{},
		),
	)
}

func TestValidateRecordInvalidUnknownFieldPolicy(t *testing.T) {
	shape := types.Object(
		types.Field("name").String().Required(),
	).UnknownFields(types.UnknownFieldPolicy(255)).Descriptor()

	err := valuevalidation.Validate(
		mustObject(t, value.MustRecordMember("name", value.StringValue("main"))),
		shape,
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

func TestValidateRecordNestedPath(t *testing.T) {
	shape := types.Object(
		types.Field("spec").Object(
			types.Field("replicas").Int32().Required(),
		).Required(),
	).Descriptor()

	payload := mustObject(
		t,
		value.MustRecordMember("spec", mustObject(
			t,
			value.MustRecordMember("replicas", value.StringValue("three")),
		)),
	)

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$.spec.replicas",
	)
}
