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

func conditionDescriptor() types.Type {
	condition := types.Object(
		types.Field("type").String().Required(),
		types.Field("status").String().Required(),
	)

	return types.ListOf(condition).Map("type").Type()
}

func readyConditionValue(status string) value.Value {
	return value.MustObjectValue(
		value.ObjectMember("type", value.StringValue("Ready")),
		value.ObjectMember("status", value.StringValue(status)),
	)
}

func TestExtractListMapUsesSelectorPaths(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(readyConditionValue("True"))
	selectorPath := path.Select(readySelector())

	got, err := ExtractAt(path, val, conditionDescriptor(), Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field("type"),
		selectorPath.Field("status"),
	)
}

func TestExtractListMapEmptyIncludesListPath(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue()

	got, err := ExtractAt(path, val, conditionDescriptor(), Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractListMapDuplicateSelectorReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		readyConditionValue("True"),
		readyConditionValue("False"),
	)

	_, err := ExtractAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrDuplicateListKey)
	requireErrorReason(t, err, ErrorReasonDuplicateListKey)
	requireErrorPath(t, err, `$.conditions[{"type":"Ready"}]`)
	requireErrorDetailContains(t, err, "first occurrence at $.conditions[0]")
	requireErrorDetailContains(t, err, "duplicate at $.conditions[1]")
}

func TestExtractListMapMissingKeyReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		value.MustObjectValue(
			value.ObjectMember("status", value.StringValue("True")),
		),
	)

	_, err := ExtractAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestExtractListMapNullKeyReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		value.MustObjectValue(
			value.ObjectMember("type", value.NullValue()),
			value.ObjectMember("status", value.StringValue("True")),
		),
	)

	_, err := ExtractAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestExtractListMapWrongKeyKindReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(
		value.MustObjectValue(
			value.ObjectMember("type", value.BoolValue(true)),
			value.ObjectMember("status", value.StringValue("True")),
		),
	)

	_, err := ExtractAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestExtractListMapNonObjectItemReturnsError(t *testing.T) {
	path := rootField("conditions")
	val := value.MustListValue(value.StringValue("Ready"))

	_, err := ExtractAt(path, val, conditionDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0]")
}

func TestExtractListMapRefElementUsesSelectorPaths(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.Condition": types.Define(
			"example.Condition",
			types.Object(
				types.Field("type").String().Required(),
				types.Field("status").String().Required(),
			),
		),
	}
	descriptor := types.ListOf(types.Ref("example.Condition")).
		Map("type").
		Type()
	val := value.MustListValue(readyConditionValue("True"))
	selectorPath := path.Select(readySelector())

	got, err := ExtractAt(path, val, descriptor, Options{Resolver: resolver})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field("type"),
		selectorPath.Field("status"),
	)
}

func TestExtractListMapRefKeyUsesSelectorLiteral(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.ConditionType": types.Define(
			"example.ConditionType",
			types.String(),
		),
	}
	condition := types.Object(
		types.Field("type").Ref("example.ConditionType").Required(),
		types.Field("status").String().Required(),
	)
	descriptor := types.ListOf(condition).Map("type").Type()
	val := value.MustListValue(readyConditionValue("True"))
	selectorPath := path.Select(readySelector())

	got, err := ExtractAt(path, val, descriptor, Options{Resolver: resolver})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field("type"),
		selectorPath.Field("status"),
	)
}

func TestExtractListMapMultiKeySelector(t *testing.T) {
	path := rootField("routes")
	route := types.Object(
		types.Field("host").String().Required(),
		types.Field("port").Uint64().Required(),
		types.Field("backend").String().Required(),
	)
	descriptor := types.ListOf(route).Map("host", "port").Type()
	val := value.MustListValue(
		value.MustObjectValue(
			value.ObjectMember("host", value.StringValue("api.example.com")),
			value.ObjectMember("port", value.Uint64Value(443)),
			value.ObjectMember("backend", value.StringValue("svc")),
		),
	)
	selectorPath := path.Select(routeSelector())

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field("host"),
		selectorPath.Field("port"),
		selectorPath.Field("backend"),
	)
}

func TestExtractListMapNestedObjectPaths(t *testing.T) {
	path := rootField("conditions")
	condition := types.Object(
		types.Field("type").String().Required(),
		types.Field("detail").Object(
			types.Field("message").String().Required(),
		).Required(),
	)
	descriptor := types.ListOf(condition).Map("type").Type()
	val := value.MustListValue(
		value.MustObjectValue(
			value.ObjectMember("type", value.StringValue("Ready")),
			value.ObjectMember(
				"detail",
				value.MustObjectValue(
					value.ObjectMember("message", value.StringValue("ok")),
				),
			),
		),
	)
	selectorPath := path.Select(readySelector())

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		selectorPath.Field("type"),
		selectorPath.Field("detail").Field("message"),
	)
}
