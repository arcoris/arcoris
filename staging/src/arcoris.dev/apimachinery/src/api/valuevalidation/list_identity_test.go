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

func TestValidateListMapAcceptsValidConditions(t *testing.T) {
	shape := conditionListShape()
	payload := mustList(
		t,
		conditionValue(t, "Ready", "True"),
		conditionValue(t, "Healthy", "False"),
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

func TestValidateListMapMissingKeyStillValidatesItemAtIndexPath(t *testing.T) {
	shape := conditionListShape()
	payload := mustList(
		t,
		mustObject(t, value.ObjectMember("status", value.StringValue(""))),
	)

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("conditions"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrMissingField,
		valuevalidation.ErrorReasonMissingField,
		"$.conditions[0].type",
	)
	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$.conditions[0].status",
	)
}

func TestValidateListMapWrongKeyKindStillValidatesItemAtIndexPath(t *testing.T) {
	shape := conditionListShape()
	payload := mustList(
		t,
		mustObject(
			t,
			value.ObjectMember("type", value.BoolValue(true)),
			value.ObjectMember("status", value.StringValue("")),
		),
	)

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("conditions"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$.conditions[0].type",
	)
	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$.conditions[0].status",
	)
}

func TestValidateListMapNullKeyStillValidatesItemAtIndexPath(t *testing.T) {
	shape := conditionListShape()
	payload := mustList(
		t,
		mustObject(
			t,
			value.ObjectMember("type", value.NullValue()),
			value.ObjectMember("status", value.StringValue("")),
		),
	)

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("conditions"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrNullNotAllowed,
		valuevalidation.ErrorReasonNullNotAllowed,
		"$.conditions[0].type",
	)
	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$.conditions[0].status",
	)
}

func TestValidateListMapNonObjectItemReportsIndexPathKindMismatch(t *testing.T) {
	shape := conditionListShape()
	payload := mustList(t, value.StringValue("not-object"))

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("conditions"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$.conditions[0]",
	)
}

func TestValidateListMapDescriptorFailureReportsIdentityError(t *testing.T) {
	shape := types.ListOf(types.Ref("example.dev.Condition")).Map("type").Descriptor()
	payload := mustList(t, conditionValue(t, "Ready", "True"))

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("conditions"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrUnresolvedRef,
		valuevalidation.ErrorReasonUnresolvedRef,
		"$.conditions[0]",
	)
}

func TestValidateListMapSelectorSuccessStillUsesSelectorPath(t *testing.T) {
	shape := conditionListShape()
	payload := mustList(t, conditionValue(t, "Ready", ""))

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("conditions"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		`$.conditions[{"type":"Ready"}].status`,
	)
}

func TestValidateListMapDuplicateKeyUsesSelectorPath(t *testing.T) {
	shape := conditionListShape()
	payload := mustList(
		t,
		conditionValue(t, "Ready", "True"),
		conditionValue(t, "Ready", "False"),
	)

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("conditions"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrDuplicateListKey,
		valuevalidation.ErrorReasonDuplicateListKey,
		`$.conditions[{"type":"Ready"}]`,
	)
}

func TestValidateListMapMultiKeySelectorPath(t *testing.T) {
	routeShape := types.Object(
		types.Field("host").String().Required(),
		types.Field("port").Uint64().Required(),
		types.Field("backend").String().MinBytes(1).Required(),
	)
	shape := types.ListOf(routeShape).Map("host", "port").Descriptor()
	payload := mustList(
		t,
		mustObject(
			t,
			value.ObjectMember("host", value.StringValue("api.example.com")),
			value.ObjectMember("port", value.Uint64Value(443)),
			value.ObjectMember("backend", value.StringValue("")),
		),
	)

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("routes"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		`$.routes[{"host":"api.example.com","port":443}].backend`,
	)
}

func conditionListShape() types.Descriptor {
	conditionShape := types.Object(
		types.Field("type").String().Required(),
		types.Field("status").String().MinBytes(1).Required(),
	)

	return types.ListOf(conditionShape).Map("type").Descriptor()
}

func conditionValue(t *testing.T, conditionType string, status string) value.Value {
	t.Helper()

	return mustObject(
		t,
		value.ObjectMember("type", value.StringValue(conditionType)),
		value.ObjectMember("status", value.StringValue(status)),
	)
}
