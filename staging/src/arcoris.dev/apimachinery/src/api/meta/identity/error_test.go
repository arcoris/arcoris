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

package identity

import (
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

func TestErrorPreservesStructuredIdentityDiagnostics(t *testing.T) {
	cause := errors.New("nested")
	err := &Error{
		Value: "bad",
		Record: diagnostic.WrapRecord(
			"objectReference.name",
			ErrInvalidObjectReference,
			ErrorReasonInvalidForm,
			"nested value is invalid",
			cause,
		),
	}

	requireErrorIs(t, err, ErrInvalidObjectReference)
	requireErrorIs(t, err, cause)

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(%T) = false", got)
	}
	if got.Path != "objectReference.name" {
		t.Fatalf("Path = %q, want %q", got.Path, "objectReference.name")
	}
	if got.Reason != ErrorReasonInvalidForm {
		t.Fatalf("Reason = %q, want %q", got.Reason, ErrorReasonInvalidForm)
	}
	if got.Value != "bad" {
		t.Fatalf("Value = %q, want bad", got.Value)
	}
	if !strings.Contains(err.Error(), "nested value is invalid") {
		t.Fatalf("Error() = %q, want detail", err.Error())
	}
	if !strings.Contains(err.Error(), `value "bad"`) {
		t.Fatalf("Error() = %q, want rejected value", err.Error())
	}
}

func TestErrorRecordsRejectedValue(t *testing.T) {
	err := Name("Bad").ValidateLexical()
	requireErrorIs(t, err, ErrInvalidName)

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(%T) = false", got)
	}
	if got.Path != "name" {
		t.Fatalf("Path = %q, want name", got.Path)
	}
	if got.Value != "Bad" {
		t.Fatalf("Value = %q, want Bad", got.Value)
	}
	if got.Reason != ErrorReasonInvalidEdge {
		t.Fatalf("Reason = %q, want %q", got.Reason, ErrorReasonInvalidEdge)
	}
}

func TestErrorNilIdentityDiagnostic(t *testing.T) {
	var err *Error
	if got := err.Error(); got != "<nil>" {
		t.Fatalf("Error() = %q, want <nil>", got)
	}
	if got := err.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v, want nil", got)
	}
}
