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

func TestExtractMapEntriesUseKeyPaths(t *testing.T) {
	path := rootField("metadata", "labels")
	val := value.MustObjectValue(
		value.ObjectMember("app", value.StringValue("api")),
	)

	got, err := ExtractAt(path, val, types.MapOf(types.String()).Descriptor(), Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Key("app"))
}

func TestExtractMapEmptyIncludesMapPath(t *testing.T) {
	path := rootField("metadata", "labels")
	val := value.MustObjectValue()

	got, err := ExtractAt(path, val, types.MapOf(types.String()).Descriptor(), Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractMapRejectsKindMismatch(t *testing.T) {
	path := rootField("metadata", "labels")
	val := value.MustListValue(value.StringValue("app"))

	_, err := ExtractAt(path, val, types.MapOf(types.String()).Descriptor(), Options{})

	requireErrorIs(t, err, ErrKindMismatch)
	requireErrorReason(t, err, ErrorReasonKindMismatch)
	requireErrorPath(t, err, `$.metadata.labels`)
}

func TestExtractMapNestedObjectValues(t *testing.T) {
	path := rootField("routes")
	descriptor := types.MapOf(
		types.Object(
			types.Field("backend").String().Required(),
		),
	).Descriptor()
	val := value.MustObjectValue(
		value.ObjectMember(
			"api",
			value.MustObjectValue(
				value.ObjectMember("backend", value.StringValue("svc")),
			),
		),
	)

	got, err := ExtractAt(path, val, descriptor, Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path.Key("api").Field("backend"))
}
