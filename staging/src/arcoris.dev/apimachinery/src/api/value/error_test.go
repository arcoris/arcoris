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

func TestErrorStringIncludesStructuredMembers(t *testing.T) {
	err := &Error{
		Path:   "object.members[0].name",
		Err:    ErrEmptyName,
		Reason: ErrorReasonEmptyName,
		Detail: "object member name is empty",
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

func TestErrorUnwrapPreservesMemberAndCause(t *testing.T) {
	cause := errors.New("nested")
	err := &Error{
		Path:   "object.members[0].value",
		Err:    ErrInvalidMember,
		Reason: ErrorReasonInvalidValue,
		Detail: "object member has invalid value",
		Cause:  cause,
	}

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorIs(t, err, ErrInvalidObject)
	requireErrorIs(t, err, ErrInvalidMember)
	requireErrorIs(t, err, cause)
}
