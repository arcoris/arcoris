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

func TestExtractOwnershipFieldsRecordLeaves(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(
		types.Field("replicas").Int32().Required(),
		types.Field("image").String().Required(),
	).Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember("replicas", value.Int64Value(3)),
		value.MustRecordMember("image", value.StringValue("api:v1")),
	)

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		path.Field(testFieldName("replicas")),
		path.Field(testFieldName("image")),
	)
}

func TestExtractOwnershipFieldsRecordEmptyIncludesRecordPath(t *testing.T) {
	path := rootField("spec")
	val := value.MustRecordValue()

	got, err := ExtractOwnershipFieldsAt(path, val, types.Object().Descriptor(), Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractOwnershipFieldsRecordMissingFieldsNotIncluded(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(
		types.Field("replicas").Int32().Required(),
		types.Field("image").String().Required(),
	).Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember("replicas", value.Int64Value(3)),
	)

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field(testFieldName("replicas")))
}

func TestExtractOwnershipFieldsRecordUnknownRejectedReturnsError(t *testing.T) {
	path := rootField("spec")
	val := value.MustRecordValue(
		value.MustRecordMember("extra", value.StringValue("debug")),
	)

	_, err := ExtractOwnershipFieldsAt(path, val, types.Object().Descriptor(), Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, "$.spec.extra")
}

func TestExtractOwnershipFieldsRecordUnknownPreserveOpaqueIncludesOpaquePath(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object().
		UnknownFields(types.UnknownPreserveOpaque).
		Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember(
			"extra",
			value.MustRecordValue(
				value.MustRecordMember("nested", value.StringValue("debug")),
			),
		),
	)

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field(testFieldName("extra")))
}

func TestExtractOwnershipFieldsRecordUnknownPreserveOpaqueDoesNotTraverseNestedStructure(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object().
		UnknownFields(types.UnknownPreserveOpaque).
		Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember(
			"extra",
			value.MustRecordValue(
				value.MustRecordMember("nested", value.StringValue("debug")),
			),
		),
	)

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field(testFieldName("extra")))
}

func TestExtractOwnershipFieldsRecordUnknownPrunedSkipsPath(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(
		types.Field("name").String().Required(),
	).
		UnknownFields(types.UnknownPrune).
		Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember("name", value.StringValue("api")),
		value.MustRecordMember("extra", value.StringValue("debug")),
	)

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field(testFieldName("name")))
}

func TestExtractOwnershipFieldsRecordOnlyPrunedUnknownFieldsReturnsEmptySet(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object().
		UnknownFields(types.UnknownPrune).
		Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember("x-extra", value.StringValue("value")),
	)

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got)
}

func TestExtractOwnershipFieldsRecordNestedPaths(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(
		types.Field("template").Object(
			types.Field("image").String().Required(),
		).Required(),
	).Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember(
			"template",
			value.MustRecordValue(
				value.MustRecordMember("image", value.StringValue("api:v1")),
			),
		),
	)

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field(testFieldName("template")).Field(testFieldName("image")))
}
