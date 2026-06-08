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

package codecselection

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorStringNil(t *testing.T) {
	var selectionError *Error

	if got := selectionError.Error(); got != "<nil>" {
		t.Fatalf("Error() = %q", got)
	}
}

func TestErrorStringIncludesPackageName(t *testing.T) {
	err := errorAt(
		"codecselection.decodeBindings[0]",
		ErrInvalidBinding,
		ErrorReasonInvalidBinding,
		"binding is invalid",
	)

	if got := err.Error(); got == "" || !strings.Contains(got, "codecselection") {
		t.Fatalf("Error() = %q; want package name", got)
	}
}

func TestErrorUnwrapNil(t *testing.T) {
	var selectionError *Error

	if got := selectionError.Unwrap(); got != nil {
		t.Fatalf("Unwrap() = %v; want nil", got)
	}
}

func TestErrorUnwrapExposesSentinel(t *testing.T) {
	err := errorAt(
		"codecselection.decodeBindings[0]",
		ErrInvalidBinding,
		ErrorReasonInvalidBinding,
		"binding is invalid",
	)

	requireErrorIs(t, err, ErrInvalidBinding)
}

func TestErrorAsSelectionError(t *testing.T) {
	err := errorAt(
		"codecselection.decodeBindings[0]",
		ErrInvalidBinding,
		ErrorReasonInvalidBinding,
		"binding is invalid",
	)

	var selectionError *Error
	if !errors.As(err, &selectionError) {
		t.Fatalf("errors.As(..., *Error) = false")
	}
}

func TestErrorDiagnosticPath(t *testing.T) {
	err := errorAt(
		"codecselection.decodeBindings[0].entryID",
		ErrUnknownEntryID,
		ErrorReasonUnknownEntryID,
		"entry is not registered",
	)

	requireSelectionError(t, err, "codecselection.decodeBindings[0].entryID", ErrorReasonUnknownEntryID)
}
