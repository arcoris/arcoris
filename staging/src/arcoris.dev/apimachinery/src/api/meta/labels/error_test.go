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

package labels

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorPreservesStructuredLabelDiagnostics(t *testing.T) {
	cause := errors.New("nested")
	err := &Error{
		Path:   "labels[key]",
		Err:    ErrInvalidSet,
		Reason: ErrorReasonInvalidForm,
		Detail: "nested value is invalid",
		Cause:  cause,
	}

	requireErrorIs(t, err, ErrInvalidSet)
	requireErrorIs(t, err, cause)

	var got *Error
	if !errors.As(err, &got) {
		t.Fatalf("errors.As(%T) = false", got)
	}
	if got.Path != "labels[key]" {
		t.Fatalf("Path = %q, want %q", got.Path, "labels[key]")
	}
	if got.Reason != ErrorReasonInvalidForm {
		t.Fatalf("Reason = %q, want %q", got.Reason, ErrorReasonInvalidForm)
	}
	if !strings.Contains(err.Error(), "nested value is invalid") {
		t.Fatalf("Error() = %q, want detail", err.Error())
	}
}

func TestErrorNilLabelDiagnostic(t *testing.T) {
	var err *Error
	if got := err.Error(); got != "<nil>" {
		t.Fatalf("Error() = %q, want <nil>", got)
	}
	if got := err.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v, want nil", got)
	}
}
