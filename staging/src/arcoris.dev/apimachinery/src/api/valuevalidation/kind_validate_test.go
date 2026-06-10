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

package valuevalidation

import (
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestValidateReportsKindMismatchDetail(t *testing.T) {
	err := Validate(
		value.StringValue("three"),
		types.Int32().Descriptor(),
		Options{},
	)

	requireInternalError(
		t,
		err,
		ErrKindMismatch,
		ErrorReasonKindMismatch,
		"$",
	)

	var validationErr *Error
	if !errors.As(err, &validationErr) {
		t.Fatalf("errors.As(*Error) = false: %v", err)
	}

	for _, expected := range []string{"string", "int32", "integer"} {
		if !strings.Contains(validationErr.Detail, expected) {
			t.Fatalf("detail %q does not contain %q", validationErr.Detail, expected)
		}
	}
}

func TestValidateRejectsNonNullValueForNullDescriptor(t *testing.T) {
	err := Validate(
		value.StringValue("not-null"),
		types.Null().Descriptor(),
		Options{},
	)

	requireInternalError(
		t,
		err,
		ErrKindMismatch,
		ErrorReasonKindMismatch,
		"$",
	)
}

func TestRequireKindAcceptsExpectedKind(t *testing.T) {
	run := newValidator(Options{})

	ok := run.requireKind(
		fieldpath.Root(),
		value.StringValue("main"),
		value.KindString,
		types.DescriptorString,
	)

	if !ok {
		t.Fatalf("requireKind() = false")
	}
	if result := run.result(); result != nil {
		t.Fatalf("result() = %v, want nil", result)
	}
}

func TestRequireKindRecordsMismatch(t *testing.T) {
	run := newValidator(Options{})

	ok := run.requireKind(
		fieldpath.Root(),
		value.StringValue("main"),
		value.KindInteger,
		types.DescriptorInt32,
	)

	if ok {
		t.Fatalf("requireKind() = true")
	}

	result := run.result()
	if !errors.Is(result, ErrKindMismatch) {
		t.Fatalf("errors.Is(ErrKindMismatch) = false: %v", result)
	}
}
