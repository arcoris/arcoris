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

package valuevalidation_test

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateRefResolvesScalar(t *testing.T) {
	resolver := testResolver{
		"example.Name": types.Define("example.Name", types.String().MinBytes(1)),
	}

	requireNoError(
		t,
		valuevalidation.Validate(
			value.StringValue("main"),
			types.Ref("example.Name").Descriptor(),
			valuevalidation.Options{Resolver: resolver},
		),
	)
}

func TestValidateRefResolvesObject(t *testing.T) {
	resolver := testResolver{
		"example.Spec": types.Define(
			"example.Spec",
			types.Object(types.Field("name").String().Required()),
		),
	}

	payload := mustObject(t, value.ObjectMember("name", value.StringValue("main")))

	requireNoError(
		t,
		valuevalidation.Validate(
			payload,
			types.Ref("example.Spec").Descriptor(),
			valuevalidation.Options{Resolver: resolver},
		),
	)
}

func TestValidateRefResolvesNullableTargetForNullValue(t *testing.T) {
	resolver := testResolver{
		"example.Note": types.Define(
			"example.Note",
			types.String().Nullable(),
		),
	}

	requireNoError(
		t,
		valuevalidation.Validate(
			value.NullValue(),
			types.Ref("example.Note").Descriptor(),
			valuevalidation.Options{Resolver: resolver},
		),
	)
}

func TestValidateRefRejectsMissingResolver(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("main"),
		types.Ref("example.Name").Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrUnresolvedRef,
		valuevalidation.ErrorReasonUnresolvedRef,
		"$",
	)
}

func TestValidateRefRejectsUnresolvedReference(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("main"),
		types.Ref("example.Name").Descriptor(),
		valuevalidation.Options{Resolver: testResolver{}},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrUnresolvedRef,
		valuevalidation.ErrorReasonUnresolvedRef,
		"$",
	)
}

func TestValidateRefRejectsReferenceCycle(t *testing.T) {
	resolver := testResolver{
		"example.A": types.Define("example.A", types.Ref("example.B")),
		"example.B": types.Define("example.B", types.Ref("example.A")),
	}

	err := valuevalidation.Validate(
		value.StringValue("main"),
		types.Ref("example.A").Descriptor(),
		valuevalidation.Options{Resolver: resolver},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrReferenceCycle,
		valuevalidation.ErrorReasonReferenceCycle,
		"$",
	)
}

func TestValidateRefRejectsMaxDepth(t *testing.T) {
	resolver := testResolver{
		"example.Name":       types.Define("example.Name", types.Ref("example.StringName")),
		"example.StringName": types.Define("example.StringName", types.String().MinBytes(1)),
	}

	err := valuevalidation.Validate(
		value.StringValue("main"),
		types.Ref("example.Name").Descriptor(),
		valuevalidation.Options{
			Resolver: resolver,
			MaxDepth: 0,
		},
	)

	requireNoError(t, err)

	err = valuevalidation.Validate(
		value.StringValue("main"),
		types.Ref("example.Name").Descriptor(),
		valuevalidation.Options{
			Resolver: resolver,
			MaxDepth: 1,
		},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrReferenceCycle,
		valuevalidation.ErrorReasonReferenceCycle,
		"$",
	)
}
