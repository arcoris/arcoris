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

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

func TestErrorStringIncludesStructuredMembers(t *testing.T) {
	err := &Error{
		Record: diagnostic.NewRecord(
			"record.members[0].name",
			ErrEmptyMemberName,
			ErrorReasonEmptyMemberName,
			"record member name is empty",
		),
	}

	got := err.Error()

	if got == "" {
		t.Fatal("Error() returned empty string")
	}
	if !errors.Is(err, ErrEmptyMemberName) {
		t.Fatal("errors.Is did not preserve sentinel")
	}
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatal("errors.Is did not preserve broad value sentinel")
	}
	if !errors.Is(err, ErrInvalidRecord) {
		t.Fatal("errors.Is did not preserve broad record sentinel")
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
		Record: diagnostic.WrapRecord(
			"record.members[0].value",
			ErrInvalidRecordMember,
			ErrorReasonInvalidValue,
			"record member has invalid value",
			cause,
		),
	}

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorIs(t, err, ErrInvalidRecord)
	requireErrorIs(t, err, ErrInvalidRecordMember)
	requireErrorIs(t, err, cause)
}
