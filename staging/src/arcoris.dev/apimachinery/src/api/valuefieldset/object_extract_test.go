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

func TestExtractObjectLeaves(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(
		types.Field("replicas").Int32().Required(),
		types.Field("image").String().Required(),
	).Descriptor()
	val := value.MustObjectValue(
		value.ObjectMember("replicas", value.Int64Value(3)),
		value.ObjectMember("image", value.StringValue("api:v1")),
	)

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		path.Field("replicas"),
		path.Field("image"),
	)
}

func TestExtractObjectEmptyIncludesObjectPath(t *testing.T) {
	path := rootField("spec")
	val := value.MustObjectValue()

	got, err := ExtractAt(path, val, types.Object().Descriptor(), Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractObjectMissingFieldsNotIncluded(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(
		types.Field("replicas").Int32().Required(),
		types.Field("image").String().Required(),
	).Descriptor()
	val := value.MustObjectValue(
		value.ObjectMember("replicas", value.Int64Value(3)),
	)

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field("replicas"))
}

func TestExtractObjectUnknownRejectedReturnsError(t *testing.T) {
	path := rootField("spec")
	val := value.MustObjectValue(
		value.ObjectMember("extra", value.StringValue("debug")),
	)

	_, err := ExtractAt(path, val, types.Object().Descriptor(), Options{})

	requireErrorIs(t, err, ErrUnknownField)
	requireErrorReason(t, err, ErrorReasonUnknownField)
	requireErrorPath(t, err, "$.spec.extra")
}

func TestExtractObjectUnknownPreservedIncludesOpaquePath(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object().
		UnknownFields(types.UnknownPreserve).
		Descriptor()
	val := value.MustObjectValue(
		value.ObjectMember(
			"extra",
			value.MustObjectValue(
				value.ObjectMember("nested", value.StringValue("debug")),
			),
		),
	)

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field("extra"))
}

func TestExtractObjectUnknownPreservedDoesNotTraverseNestedStructure(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object().
		UnknownFields(types.UnknownPreserve).
		Descriptor()
	val := value.MustObjectValue(
		value.ObjectMember(
			"extra",
			value.MustObjectValue(
				value.ObjectMember("nested", value.StringValue("debug")),
			),
		),
	)

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field("extra"))
}

func TestExtractObjectUnknownPrunedSkipsPath(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(
		types.Field("name").String().Required(),
	).
		UnknownFields(types.UnknownPrune).
		Descriptor()
	val := value.MustObjectValue(
		value.ObjectMember("name", value.StringValue("api")),
		value.ObjectMember("extra", value.StringValue("debug")),
	)

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field("name"))
}

func TestExtractObjectOnlyPrunedUnknownFieldsReturnsEmptySet(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object().
		UnknownFields(types.UnknownPrune).
		Descriptor()
	val := value.MustObjectValue(
		value.ObjectMember("x-extra", value.StringValue("value")),
	)

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got)
}

func TestExtractObjectNestedPaths(t *testing.T) {
	path := rootField("spec")
	descriptor := types.Object(
		types.Field("template").Object(
			types.Field("image").String().Required(),
		).Required(),
	).Descriptor()
	val := value.MustObjectValue(
		value.ObjectMember(
			"template",
			value.MustObjectValue(
				value.ObjectMember("image", value.StringValue("api:v1")),
			),
		),
	)

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field("template").Field("image"))
}
