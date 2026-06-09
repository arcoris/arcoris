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
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestErrorIsSentinel(t *testing.T) {
	err := valuevalidation.Validate(value.StringValue("x"), types.Int64().Descriptor(), valuevalidation.Options{})

	if !errors.Is(err, valuevalidation.ErrKindMismatch) {
		t.Fatalf("errors.Is(ErrKindMismatch) = false")
	}
}

func TestErrorAsValueValidationError(t *testing.T) {
	err := valuevalidation.Validate(value.StringValue("x"), types.Int64().Descriptor(), valuevalidation.Options{})

	var validationError *valuevalidation.Error
	if !errors.As(err, &validationError) {
		t.Fatalf("errors.As(*Error) = false")
	}
	if got, want := validationError.Path, "$"; got != want {
		t.Fatalf("Path = %q, want %q", got, want)
	}
	if got, want := validationError.Reason, valuevalidation.ErrorReasonKindMismatch; got != want {
		t.Fatalf("Reason = %q, want %q", got, want)
	}
}
