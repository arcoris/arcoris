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

package fieldpath

import (
	"errors"
	"testing"
)

func TestNewErrorBuildsStructuredError(t *testing.T) {
	err := newError(
		ErrInvalidElement,
		ErrorReasonInvalidElement,
		"element kind is invalid",
	)

	var pathErr *Error
	if !errors.As(err, &pathErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	requireErrorIs(t, err, ErrInvalidElement)
	requireEqual(t, pathErr.Reason, ErrorReasonInvalidElement)
	requireEqual(t, pathErr.Detail, "element kind is invalid")
}

func TestNestedBuildsStructuredErrorWithCause(t *testing.T) {
	cause := errors.New("cause")
	err := nested(
		ErrInvalidPath,
		ErrorReasonInvalidElement,
		"path element is invalid",
		cause,
	)

	var pathErr *Error
	if !errors.As(err, &pathErr) {
		t.Fatalf("expected *Error, got %T", err)
	}

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, cause)
	requireEqual(t, pathErr.Cause, cause)
}
