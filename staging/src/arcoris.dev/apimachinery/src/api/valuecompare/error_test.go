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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"errors"
	"testing"
)

func TestErrorIsSentinel(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidValue, ErrorReasonInvalidZero, "bad")

	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("errors.Is(err, ErrInvalidValue) = false")
	}
}

func TestErrorIsInvalidPath(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidPath, ErrorReasonInvalidPath, "bad")

	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("errors.Is(err, ErrInvalidPath) = false")
	}
}

func TestErrorIsInvalidResult(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidResult, ErrorReasonInvalidResult, "bad")

	if !errors.Is(err, ErrInvalidResult) {
		t.Fatalf("errors.Is(err, ErrInvalidResult) = false")
	}
}

func TestErrorAsValueCompareError(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidValue, ErrorReasonInvalidZero, "bad")

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(*Error) = false")
	}
	if got.Path != "$" {
		t.Fatalf("path = %q, want $", got.Path)
	}
}

func TestNilErrorStringAndUnwrap(t *testing.T) {
	var err *Error

	if err.Error() != "<nil>" {
		t.Fatalf("nil Error() = %q", err.Error())
	}
	if err.Unwrap() != nil {
		t.Fatalf("nil Unwrap() != nil")
	}
}
