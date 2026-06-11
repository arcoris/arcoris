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

func TestExtractOwnershipFieldsOrderedListUsesIndexPaths(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue(
		value.StringValue("a"),
		value.StringValue("b"),
	)

	got, err := ExtractOwnershipFieldsAt(
		path,
		val,
		types.ListOf(types.String()).Ordered().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path.Index(0), path.Index(1))
}

func TestExtractOwnershipFieldsOrderedListObjectItemsUseIndexedChildPaths(t *testing.T) {
	path := rootField("containers")
	item := types.Object(
		types.Field("name").String().Required(),
		types.Field("image").String().Required(),
	)
	descriptor := types.ListOf(item).Ordered().Descriptor()
	val := value.MustListValue(
		value.MustRecordValue(
			value.MustRecordMember("name", value.StringValue("api")),
			value.MustRecordMember("image", value.StringValue("api:v1")),
		),
	)

	got, err := ExtractOwnershipFieldsAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		path.Index(0).Field(testFieldName("name")),
		path.Index(0).Field(testFieldName("image")),
	)
}

func TestExtractOwnershipFieldsOrderedListEmptyIncludesListPath(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue()

	got, err := ExtractOwnershipFieldsAt(
		path,
		val,
		types.ListOf(types.String()).Ordered().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractOwnershipFieldsAtomicListIncludesOnlyListPath(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue(
		value.StringValue("a"),
		value.StringValue("b"),
	)

	got, err := ExtractOwnershipFieldsAt(
		path,
		val,
		types.ListOf(types.String()).Atomic().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractOwnershipFieldsAtomicListEmptyIncludesListPath(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue()

	got, err := ExtractOwnershipFieldsAt(
		path,
		val,
		types.ListOf(types.String()).Atomic().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractOwnershipFieldsSetListIncludesOnlyListPath(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue(
		value.StringValue("a"),
		value.StringValue("b"),
	)

	got, err := ExtractOwnershipFieldsAt(
		path,
		val,
		types.ListOf(types.String()).Set().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractOwnershipFieldsSetListEmptyIncludesListPath(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue()

	got, err := ExtractOwnershipFieldsAt(
		path,
		val,
		types.ListOf(types.String()).Set().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}
