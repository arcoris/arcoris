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

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestExtractStartsAtRootPath(t *testing.T) {
	got, err := Extract(
		value.StringValue("api"),
		types.String().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, fieldpath.RootPath())
}

func TestExtractAtPreservesBasePath(t *testing.T) {
	path := rootField("spec", "image")

	got, err := ExtractAt(
		path,
		value.StringValue("api:v1"),
		types.String().Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractRejectsInvalidZeroValue(t *testing.T) {
	_, err := Extract(value.Value{}, types.String().Descriptor(), Options{})

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorReason(t, err, ErrorReasonInvalidZero)
	requireErrorPath(t, err, "$")
}

func TestExtractRejectsInvalidDescriptor(t *testing.T) {
	_, err := Extract(value.StringValue("api"), types.Descriptor{}, Options{})

	requireErrorIs(t, err, ErrInvalidDescriptor)
	requireErrorReason(t, err, ErrorReasonInvalidDescriptor)
	requireErrorPath(t, err, "$")
}

func TestExtractAtRejectsInvalidBasePath(t *testing.T) {
	_, err := ExtractAt(
		fieldpath.RootPath().Field(""),
		value.StringValue("api"),
		types.String().Descriptor(),
		Options{},
	)

	requireErrorIs(t, err, ErrInvalidDescriptor)
	requireErrorIs(t, err, fieldpath.ErrInvalidPath)
	requireErrorPath(t, err, `$.""`)
}

func TestExtractReturnsSortedUniqueSet(t *testing.T) {
	descriptor := types.Object(
		types.Field("zeta").String().Optional(),
		types.Field("alpha").String().Optional(),
	).Descriptor()
	val := value.MustRecordValue(
		value.MustRecordMember("zeta", value.StringValue("last")),
		value.MustRecordMember("alpha", value.StringValue("first")),
	)

	got, err := Extract(val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(
		t,
		got,
		fieldpath.RootPath().Field("alpha"),
		fieldpath.RootPath().Field("zeta"),
	)
}

func TestExtractDoesNotMutateValue(t *testing.T) {
	val := value.MustRecordValue(
		value.MustRecordMember("name", value.StringValue("api")),
	)

	_, err := Extract(
		val,
		types.Object(types.Field("name").String().Required()).Descriptor(),
		Options{},
	)
	requireNoError(t, err)

	objectView, ok := val.AsRecord()
	if !ok {
		t.Fatalf("Object() ok = false")
	}
	memberValue, ok := objectView.Get("name")
	if !ok {
		t.Fatalf("member name missing after extraction")
	}
	text, ok := memberValue.AsString()
	if !ok || text != "api" {
		t.Fatalf("member value = %q, %v; want api, true", text, ok)
	}
}
