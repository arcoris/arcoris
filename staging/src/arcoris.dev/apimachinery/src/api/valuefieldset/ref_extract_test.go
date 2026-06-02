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

func TestExtractRefResolvesScalar(t *testing.T) {
	path := rootField("name")
	resolver := testResolver{
		"example.Name": types.Define("example.Name", types.String()),
	}

	got, err := ExtractAt(
		path,
		value.StringValue("api"),
		types.Ref("example.Name").Type(),
		Options{Resolver: resolver},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractRefResolvesObject(t *testing.T) {
	path := rootField("spec")
	resolver := testResolver{
		"example.Spec": types.Define(
			"example.Spec",
			types.Object(types.Field("image").String().Required()),
		),
	}
	val := value.MustObjectValue(
		value.ObjectMember("image", value.StringValue("api:v1")),
	)

	got, err := ExtractAt(
		path,
		val,
		types.Ref("example.Spec").Type(),
		Options{Resolver: resolver},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path.Field("image"))
}

func TestExtractRefRejectsMissingResolver(t *testing.T) {
	_, err := Extract(
		value.StringValue("api"),
		types.Ref("example.Name").Type(),
		Options{},
	)

	requireErrorIs(t, err, ErrUnresolvedRef)
	requireErrorReason(t, err, ErrorReasonUnresolvedRef)
}

func TestExtractRefRejectsUnresolvedReference(t *testing.T) {
	_, err := Extract(
		value.StringValue("api"),
		types.Ref("example.Missing").Type(),
		Options{Resolver: testResolver{}},
	)

	requireErrorIs(t, err, ErrUnresolvedRef)
	requireErrorReason(t, err, ErrorReasonUnresolvedRef)
}

func TestExtractRefRejectsReferenceCycle(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.B")),
		"example.B": types.Define("example.B", types.Ref("example.A")),
	}

	_, err := Extract(
		value.StringValue("api"),
		types.Ref("example.A").Type(),
		Options{Resolver: resolver},
	)

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
}

func TestExtractRefHonorsMaxDepth(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.B")),
		"example.B": types.Define("example.B", types.String()),
	}

	_, err := Extract(
		value.StringValue("api"),
		types.Ref("example.A").Type(),
		Options{Resolver: resolver, MaxDepth: 1},
	)

	requireErrorIs(t, err, ErrReferenceCycle)
	requireErrorReason(t, err, ErrorReasonReferenceCycle)
}
