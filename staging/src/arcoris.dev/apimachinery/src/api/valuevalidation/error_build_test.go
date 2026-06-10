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
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestErrorAtBuildsStructuredDiagnostic(t *testing.T) {
	err := errorAt(
		fieldpath.Root().Field(fieldpath.MustFieldName("name")),
		ErrMissingField,
		ErrorReasonMissingField,
		"name is required",
	)

	requireInternalError(
		t,
		err,
		ErrMissingField,
		ErrorReasonMissingField,
		"$.name",
	)
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("nested failure")
	err := wrapAt(
		fieldpath.Root().Field(fieldpath.MustFieldName("name")),
		ErrInvalidDescriptor,
		ErrorReasonInvalidDescriptor,
		"descriptor failed",
		cause,
	)

	requireInternalError(
		t,
		err,
		ErrInvalidDescriptor,
		ErrorReasonInvalidDescriptor,
		"$.name",
	)
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(cause) = false")
	}
}

func requireInternalError(
	t *testing.T,
	err error,
	target error,
	reason ErrorReason,
	path string,
) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v) = false", target)
	}

	var validationErr *Error
	if !errors.As(err, &validationErr) {
		t.Fatalf("errors.As(*Error) = false")
	}
	if validationErr.Path != path {
		t.Fatalf("Path = %q, want %q", validationErr.Path, path)
	}
	if validationErr.Reason != reason {
		t.Fatalf("Reason = %q, want %q", validationErr.Reason, reason)
	}
	if validationErr.Detail == "" {
		t.Fatalf("Detail is empty")
	}
}
