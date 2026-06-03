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

package valuecompare

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestCompareListMapSameReorderedIsEmpty(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(
		conditionValue("Ready", "True"),
		conditionValue("Degraded", "False"),
	)
	newValue := value.MustListValue(
		conditionValue("Degraded", "False"),
		conditionValue("Ready", "True"),
	)

	got, err := CompareAt(path, oldValue, newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, nil)
}

func TestCompareListMapModifiedItemField(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(readySelector()).Field("status")))
}

func TestCompareListMapAddedItem(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, paths(path.Select(readySelector()).Field("type"), path.Select(readySelector()).Field("status")), nil, nil)
}

func TestCompareListMapRemovedItem(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, value.MustListValue(), conditionsDescriptor(), Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, paths(path.Select(readySelector()).Field("type"), path.Select(readySelector()).Field("status")), nil)
}

func TestCompareListMapMultiKeySelector(t *testing.T) {
	path := rootField("routes")
	descriptor := types.ListOf(
		types.Object(
			types.Field("host").String().Required(),
			types.Field("port").Uint64().Required(),
			types.Field("backend").String().Required(),
		),
	).Map("host", "port").Type()
	oldValue := value.MustListValue(routeValue("old"))
	newValue := value.MustListValue(routeValue("new"))

	got, err := CompareAt(path, oldValue, newValue, descriptor, Options{})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(routeSelector()).Field("backend")))
}

func TestCompareListMapDuplicateOldSelectorReturnsError(t *testing.T) {
	path := rootField("conditions")
	oldValue := value.MustListValue(conditionValue("Ready", "True"), conditionValue("Ready", "False"))

	_, err := CompareAt(path, oldValue, value.MustListValue(), conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrDuplicateListKey)
	requireErrorReason(t, err, ErrorReasonDuplicateListKey)
	requireErrorPath(t, err, `$.conditions[{"type":"Ready"}]`)
	requireErrorDetailContains(t, err, "first occurrence at $.conditions[0]")
	requireErrorDetailContains(t, err, "duplicate at $.conditions[1]")
}

func TestCompareListMapDuplicateNewSelectorReturnsError(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(conditionValue("Ready", "True"), conditionValue("Ready", "False"))

	_, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrDuplicateListKey)
	requireErrorReason(t, err, ErrorReasonDuplicateListKey)
	requireErrorPath(t, err, `$.conditions[{"type":"Ready"}]`)
}

func TestCompareListMapMissingKeyReturnsError(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(valueObject("status", "True"))

	_, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareListMapWrongKeyKindReturnsError(t *testing.T) {
	path := rootField("conditions")
	newValue := value.MustListValue(value.MustObjectValue(
		value.ObjectMember("type", value.BoolValue(true)),
		value.ObjectMember("status", value.StringValue("True")),
	))

	_, err := CompareAt(path, value.MustListValue(), newValue, conditionsDescriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonInvalidListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareListMapRefElement(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.Condition": types.Define("example.Condition", conditionExpr()),
	}
	descriptor := types.ListOf(types.Ref("example.Condition")).Map("type").Type()
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, newValue, descriptor, Options{Resolver: resolver})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(readySelector()).Field("status")))
}

func TestCompareListMapRefKeyType(t *testing.T) {
	path := rootField("conditions")
	resolver := testResolver{
		"example.ConditionType": types.Define("example.ConditionType", types.String()),
	}
	descriptor := types.ListOf(
		types.Object(
			types.Field("type").Ref("example.ConditionType").Required(),
			types.Field("status").String().Required(),
		),
	).Map("type").Type()
	oldValue := value.MustListValue(conditionValue("Ready", "False"))
	newValue := value.MustListValue(conditionValue("Ready", "True"))

	got, err := CompareAt(path, oldValue, newValue, descriptor, Options{Resolver: resolver})
	requireNoError(t, err)
	requireResult(t, got, nil, nil, paths(path.Select(readySelector()).Field("status")))
}

func routeValue(backend string) value.Value {
	return value.MustObjectValue(
		value.ObjectMember("host", value.StringValue("api.example.com")),
		value.ObjectMember("port", value.Uint64Value(443)),
		value.ObjectMember("backend", value.StringValue(backend)),
	)
}
