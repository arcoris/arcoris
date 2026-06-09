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
	"errors"
	"regexp/syntax"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateStringConstraints(t *testing.T) {
	tests := []struct {
		name     string
		payload  value.Value
		shape    types.Descriptor
		sentinel error
		reason   valuevalidation.ErrorReason
	}{
		{
			name:     "too short",
			payload:  value.StringValue("a"),
			shape:    types.String().MinBytes(2).Descriptor(),
			sentinel: valuevalidation.ErrLengthOutOfRange,
			reason:   valuevalidation.ErrorReasonTooShort,
		},
		{
			name:     "too long",
			payload:  value.StringValue("abcd"),
			shape:    types.String().MaxBytes(3).Descriptor(),
			sentinel: valuevalidation.ErrLengthOutOfRange,
			reason:   valuevalidation.ErrorReasonTooLong,
		},
		{
			name:     "pattern mismatch",
			payload:  value.StringValue("abc"),
			shape:    types.String().Pattern(`^[0-9]+$`).Descriptor(),
			sentinel: valuevalidation.ErrPatternMismatch,
			reason:   valuevalidation.ErrorReasonPatternMismatch,
		},
		{
			name:     "enum mismatch",
			payload:  value.StringValue("blue"),
			shape:    types.String().Enum("red", "green").Descriptor(),
			sentinel: valuevalidation.ErrEnumMismatch,
			reason:   valuevalidation.ErrorReasonEnumMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := valuevalidation.Validate(
				tt.payload,
				tt.shape,
				valuevalidation.Options{},
			)

			requireError(t, err, tt.sentinel, tt.reason, "$")
		})
	}
}

func TestValidateStringPatternReusesPatternWithinRun(t *testing.T) {
	shape := types.ListOf(types.String().Pattern(`^[a-z]+$`)).Descriptor()
	payload := mustList(t, value.StringValue("ok"), value.StringValue("bad1"))

	err := valuevalidation.ValidateAt(
		fieldpath.RootPath().Field("names"),
		payload,
		shape,
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrPatternMismatch,
		valuevalidation.ErrorReasonPatternMismatch,
		"$.names[1]",
	)
}

func TestValidateStringInvalidPatternPreservesCompileError(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("anything"),
		types.String().Pattern(`[`).Descriptor(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrInvalidDescriptor,
		valuevalidation.ErrorReasonInvalidDescriptor,
		"$",
	)

	var syntaxError *syntax.Error
	if !errors.As(err, &syntaxError) {
		t.Fatalf("errors.As(*syntax.Error) = false: %v", err)
	}
}
