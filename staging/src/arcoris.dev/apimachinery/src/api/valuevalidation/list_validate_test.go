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

func TestValidateListAcceptsValidItems(t *testing.T) {
	shape := types.ListOf(types.String().MinLen(1)).Type()
	payload := mustList(t, value.StringValue("a"), value.StringValue("b"))

	requireNoError(
		t,
		valuevalidation.Validate(
			payload,
			shape,
			valuevalidation.Options{},
		),
	)
}

func TestValidateListRejectsInvalidItem(t *testing.T) {
	shape := types.ListOf(types.String().MinLen(1)).Type()
	payload := mustList(t, value.StringValue("ok"), value.StringValue(""))

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$[1]",
	)
}

func TestValidateListUsesIndexPath(t *testing.T) {
	shape := types.ListOf(types.Int32()).Type()
	payload := mustList(t, value.StringValue("not-int"))

	err := valuevalidation.Validate(
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$[0]",
	)
}

func TestValidateOrderedListUsesIndexPaths(t *testing.T) {
	shape := types.ListOf(types.String().MinLen(1)).Ordered().Type()
	payload := mustList(t, value.StringValue("ok"), value.StringValue(""))

	err := valuevalidation.ValidateAt(
		rootField("args"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$.args[1]",
	)
}

func TestValidateOrderedListNestedObjectUsesIndexPaths(t *testing.T) {
	item := types.Object(
		types.Field("name").String().Required(),
		types.Field("image").String().MinLen(1).Required(),
	)
	shape := types.ListOf(item).Ordered().Type()
	payload := mustList(
		t,
		mustObject(
			t,
			value.ObjectMember("name", value.StringValue("api")),
			value.ObjectMember("image", value.StringValue("")),
		),
	)

	err := valuevalidation.ValidateAt(
		rootField("containers"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$.containers[0].image",
	)
}

func TestValidateAtomicListStillUsesIndexPathsForValidation(t *testing.T) {
	shape := types.ListOf(types.String().MinLen(1)).Atomic().Type()
	payload := mustList(t, value.StringValue(""))

	err := valuevalidation.ValidateAt(
		rootField("args"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$.args[0]",
	)
}

func TestValidateSetListStillUsesIndexPathsForValidation(t *testing.T) {
	shape := types.ListOf(types.String().MinLen(1)).Set().Type()
	payload := mustList(t, value.StringValue(""))

	err := valuevalidation.ValidateAt(
		rootField("args"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrLengthOutOfRange,
		valuevalidation.ErrorReasonTooShort,
		"$.args[0]",
	)
}
