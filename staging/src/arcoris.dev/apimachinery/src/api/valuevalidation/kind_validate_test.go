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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestValidateReportsKindMismatchDetail(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("three"),
		types.Int32().Type(),
		valuevalidation.Options{},
	)

	validationErr := findValidationError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$",
	)
	if validationErr == nil {
		t.Fatalf("kind mismatch diagnostic not found: %v", err)
	}

	for _, expected := range []string{"string", "int32", "integer"} {
		if !strings.Contains(validationErr.Detail, expected) {
			t.Fatalf("detail %q does not contain %q", validationErr.Detail, expected)
		}
	}
}

func TestValidateRejectsNonNullValueForNullDescriptor(t *testing.T) {
	err := valuevalidation.Validate(
		value.StringValue("not-null"),
		types.Null().Type(),
		valuevalidation.Options{},
	)

	requireError(
		t,
		err,
		valuevalidation.ErrKindMismatch,
		valuevalidation.ErrorReasonKindMismatch,
		"$",
	)
}
