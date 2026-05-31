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

package value

import (
	"errors"
	"testing"
)

func TestErrorStringIncludesStructuredFields(t *testing.T) {
	err := &Error{
		Path:   "object.fields[0].name",
		Err:    ErrEmptyName,
		Reason: ErrorReasonEmptyName,
		Detail: "object field name is empty",
	}

	got := err.Error()

	if got == "" {
		t.Fatal("Error() returned empty string")
	}
	if !errors.Is(err, ErrEmptyName) {
		t.Fatal("errors.Is did not preserve sentinel")
	}
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatal("errors.Is did not preserve broad value sentinel")
	}
	if !errors.Is(err, ErrInvalidObject) {
		t.Fatal("errors.Is did not preserve broad object sentinel")
	}
}

func TestNilErrorMethodsAreSafe(t *testing.T) {
	var err *Error

	requireEqual(t, err.Error(), "<nil>")
	requireEqual(t, err.Unwrap() == nil, true)
}

func TestErrorUnwrapPreservesMapAndCause(t *testing.T) {
	cause := errors.New("nested")
	err := &Error{
		Path:   "map.entries[0].value",
		Err:    ErrInvalidEntry,
		Reason: ErrorReasonInvalidValue,
		Detail: "map entry has invalid value",
		Cause:  cause,
	}

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorIs(t, err, ErrInvalidMap)
	requireErrorIs(t, err, ErrInvalidEntry)
	requireErrorIs(t, err, cause)
}
