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
