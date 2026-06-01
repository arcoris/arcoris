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

package resource

import (
	"errors"
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

func TestErrorFormattingAndUnwrap(t *testing.T) {
	cause := errors.New("cause")
	err := &Error{
		Record: diagnostic.WrapRecord(
			"definition.versions[v1].desired",
			ErrInvalidVersion,
			ErrorReasonDesiredNotObject,
			"desired root must be object",
			cause,
		),
	}
	text := err.Error()

	for _, part := range []string{
		"resource",
		"definition.versions[v1].desired",
		"invalid API resource version",
		"desired_not_object",
		"desired root must be object",
	} {
		if !strings.Contains(text, part) {
			t.Fatalf("Error() = %q, missing %q", text, part)
		}
	}

	if !errors.Is(err, ErrInvalidVersion) {
		t.Fatalf("errors.Is(err, ErrInvalidVersion) = false")
	}

	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(err, cause) = false")
	}
}

func TestErrorNilReceiverFormatting(t *testing.T) {
	var err *Error
	requireEqual(t, err.Error(), "<nil>")
	if err.Unwrap() != nil {
		t.Fatalf("nil Error Unwrap() must return nil")
	}
}

func TestErrorUnwrapWithOnlyCause(t *testing.T) {
	cause := errors.New("cause")
	err := &Error{
		Record: diagnostic.CauseRecord[ErrorReason](cause),
	}
	requireErrorIs(t, err, cause)
}
