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

func TestAddedSubtreeErrorUsesValueCompareErrorModel(t *testing.T) {
	descriptor := types.Object(types.Field("name").String().Optional()).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember("name", value.BoolValue(true)))

	_, err := Compare(valueRecord(), newValue, descriptor, Options{})

	requireErrorIs(t, err, ErrKindMismatch)
	requireErrorReason(t, err, ErrorReasonKindMismatch)
	requireErrorPath(t, err, "$.name")
}

func TestCompareAddedSubtreeWrapsValueFieldSetUnknownField(t *testing.T) {
	descriptor := types.Object(types.Field("child").Object().Optional()).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember("child", valueRecord("extra", "new")))

	_, err := Compare(valueRecord(), newValue, descriptor, Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, "$.child.extra")
}

func TestCompareRemovedSubtreeWrapsValueFieldSetUnknownField(t *testing.T) {
	descriptor := types.Object(types.Field("child").Object().Optional()).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("child", valueRecord("extra", "old")))

	_, err := Compare(oldValue, valueRecord(), descriptor, Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, "$.child.extra")
}

func TestCompareAddedSubtreeWrapsInvalidListKey(t *testing.T) {
	descriptor := types.Object(
		types.Field("conditions").ListOf(conditionExpr()).Map("type").Optional(),
	).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember(
		"conditions",
		value.MustListValue(valueRecord("status", "True")),
	))

	_, err := Compare(valueRecord(), newValue, descriptor, Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareRemovedSubtreeWrapsInvalidListKey(t *testing.T) {
	descriptor := types.Object(
		types.Field("conditions").ListOf(conditionExpr()).Map("type").Optional(),
	).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember(
		"conditions",
		value.MustListValue(valueRecord("status", "True")),
	))

	_, err := Compare(oldValue, valueRecord(), descriptor, Options{})

	requireErrorIs(t, err, ErrInvalidListKey)
	requireErrorReason(t, err, ErrorReasonMissingListKey)
	requireErrorPath(t, err, "$.conditions[0].type")
}

func TestCompareAddedSubtreeWrapsUnresolvedRef(t *testing.T) {
	descriptor := types.Object(types.Field("name").Ref("example.Name").Optional()).Descriptor()
	newValue := value.MustRecordValue(value.MustRecordMember("name", value.StringValue("api")))

	_, err := Compare(valueRecord(), newValue, descriptor, Options{})

	requireErrorIs(t, err, ErrUnresolvedRef)
	requireErrorReason(t, err, ErrorReasonUnresolvedRef)
	requireErrorPath(t, err, "$.name")
}

func TestCompareRemovedSubtreeWrapsReferenceCycle(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.B")),
		"example.B": types.Define("example.B", types.Ref("example.A")),
	}
	descriptor := types.Object(types.Field("name").Ref("example.A").Optional()).Descriptor()
	oldValue := value.MustRecordValue(value.MustRecordMember("name", value.StringValue("api")))

	_, err := Compare(oldValue, valueRecord(), descriptor, Options{Resolver: resolver})

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
	requireErrorPath(t, err, "$.name")
}
