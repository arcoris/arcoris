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

package valuefieldset

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func conditionDescriptor() types.Descriptor {
	condition := types.Object(
		types.Field("type").String().Required(),
		types.Field("status").String().Required(),
	)

	return types.ListOf(condition).Map("type").Descriptor()
}

func readyConditionValue(status string) value.Value {
	return value.MustRecordValue(
		value.MustRecordMember("type", value.StringValue("Ready")),
		value.MustRecordMember("status", value.StringValue(status)),
	)
}

func TestExtractOwnershipFieldsListMapUsesSelectorPaths(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(readyConditionValue("True"))
	selectorPath := path.Select(readySelector())

	got, err := ExtractOwnershipFieldsAt(path, val, conditionDescriptor(), Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field(testFieldName("type")),
		selectorPath.Field(testFieldName("status")),
	)
}

func TestExtractOwnershipFieldsListMapEmptyIncludesListPath(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue()

	got, err := ExtractOwnershipFieldsAt(path, val, conditionDescriptor(), Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractOwnershipFieldsListMapDuplicateSelectorReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		readyConditionValue("True"),
		readyConditionValue("False"),
	)

	_, err := ExtractOwnershipFieldsAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrDuplicateListKey)
	requireErrorReason(t, err, ErrorReasonDuplicateListKey)
	requireErrorPath(t, err, `$.conditions[{"type":"Ready"}]`)
	requireErrorDetailContains(t, err, "first occurrence at $.conditions[0]")
	requireErrorDetailContains(t, err, "duplicate at $.conditions[1]")
}

func TestExtractOwnershipFieldsListMapMissingKeyReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		value.MustRecordValue(
			value.MustRecordMember("status", value.StringValue("True")),
		),
	)

	_, err := ExtractOwnershipFieldsAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestExtractOwnershipFieldsListMapNullKeyReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		value.MustRecordValue(
			value.MustRecordMember("type", value.NullValue()),
			value.MustRecordMember("status", value.StringValue("True")),
		),
	)

	_, err := ExtractOwnershipFieldsAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestExtractOwnershipFieldsListMapWrongKeyKindReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		value.MustRecordValue(
			value.MustRecordMember("type", value.BoolValue(true)),
			value.MustRecordMember("status", value.StringValue("True")),
		),
	)

	_, err := ExtractOwnershipFieldsAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestExtractOwnershipFieldsListMapNonObjectItemReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(value.StringValue("Ready"))

	_, err := ExtractOwnershipFieldsAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0]")
}

func TestExtractOwnershipFieldsListMapRefElementUsesSelectorPaths(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.dev.Condition": types.Define(
			"example.dev.Condition",
			types.Object(
				types.Field("type").String().Required(),
				types.Field("status").String().Required(),
			),
		),
	}
	descriptor := types.ListOf(types.Ref("example.dev.Condition")).
		Map("type").
		Descriptor()
	val := value.MustListValue(readyConditionValue("True"))
	selectorPath := path.Select(readySelector())

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{Resolver: resolver})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field(testFieldName("type")),
		selectorPath.Field(testFieldName("status")),
	)
}

func TestExtractOwnershipFieldsListMapUnresolvedRefReturnsUnresolvedRef(t *testing.T) {
	path := rootField("conditions")
	descriptor := types.ListOf(types.Ref("example.dev.Condition")).
		Map("type").
		Descriptor()
	val := value.MustListValue(readyConditionValue("True"))

	_, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})

	requireErrorIs(t, err, ErrUnresolvedRef)
	requireErrorReason(t, err, ErrorReasonUnresolvedRef)
	requireErrorPath(t, err, "$.conditions[0]")
}

func TestExtractOwnershipFieldsListMapReferenceCycleReturnsReferenceCycle(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.dev.Condition": types.Define(
			"example.dev.Condition",
			types.Ref("example.dev.Condition"),
		),
	}
	descriptor := types.ListOf(types.Ref("example.dev.Condition")).
		Map("type").
		Descriptor()
	val := value.MustListValue(readyConditionValue("True"))

	_, err := ExtractOwnershipFieldsAt(
		path,
		val,
		descriptor,
		Options{Resolver: resolver},
	)

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
	requireErrorPath(t, err, "$.conditions[0]")
}

func TestExtractOwnershipFieldsListMapRefKeyUsesSelectorLiteral(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.dev.ConditionType": types.Define(
			"example.dev.ConditionType",
			types.String(),
		),
	}
	condition := types.Object(
		types.Field("type").Ref("example.dev.ConditionType").Required(),
		types.Field("status").String().Required(),
	)
	descriptor := types.ListOf(condition).Map("type").Descriptor()
	val := value.MustListValue(readyConditionValue("True"))
	selectorPath := path.Select(readySelector())

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{Resolver: resolver})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field(testFieldName("type")),
		selectorPath.Field(testFieldName("status")),
	)
}

func TestExtractOwnershipFieldsListMapMultiKeySelector(t *testing.T) {
	path := rootField("routes")
	route := types.Object(
		types.Field("host").String().Required(),
		types.Field("port").Uint64().Required(),
		types.Field("backend").String().Required(),
	)
	descriptor := types.ListOf(route).Map("host", "port").Descriptor()
	val := value.MustListValue(
		value.MustRecordValue(
			value.MustRecordMember("host", value.StringValue("api.example.com")),
			value.MustRecordMember("port", value.Uint64Value(443)),
			value.MustRecordMember("backend", value.StringValue("svc")),
		),
	)
	selectorPath := path.Select(routeSelector())

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field(testFieldName("host")),
		selectorPath.Field(testFieldName("port")),
		selectorPath.Field(testFieldName("backend")),
	)
}

func TestExtractOwnershipFieldsListMapNestedObjectPaths(t *testing.T) {
	path := rootField("conditions")
	condition := types.Object(
		types.Field("type").String().Required(),
		types.Field("detail").Object(
			types.Field("message").String().Required(),
		).Required(),
	)
	descriptor := types.ListOf(condition).Map("type").Descriptor()
	val := value.MustListValue(
		value.MustRecordValue(
			value.MustRecordMember("type", value.StringValue("Ready")),
			value.MustRecordMember(
				"detail",
				value.MustRecordValue(
					value.MustRecordMember("message", value.StringValue("ok")),
				),
			),
		),
	)
	selectorPath := path.Select(readySelector())

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field(testFieldName("type")),
		selectorPath.Field(testFieldName("detail")).Field(testFieldName("message")),
	)
}
