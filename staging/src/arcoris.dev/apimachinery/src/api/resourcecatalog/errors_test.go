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

package resourcecatalog

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorFormattingAndUnwrap(t *testing.T) {
	cause := errors.New("cause")
	err := &Error{
		Path:   "definitions[control.arcoris.dev:workers]",
		Err:    ErrDefinitionExists,
		Reason: ErrorReasonDefinitionExists,
		Detail: "resource already exists",
		Cause:  cause,
	}

	text := err.Error()
	for _, part := range []string{
		"resourcecatalog",
		"definitions[control.arcoris.dev:workers]",
		"API resource definition already exists",
		"definition_exists",
		"resource already exists",
	} {
		if !strings.Contains(text, part) {
			t.Fatalf("Error() = %q, missing %q", text, part)
		}
	}

	requireErrorIs(t, err, ErrDefinitionExists)
	requireErrorIs(t, err, cause)
}

func TestErrorNilReceiver(t *testing.T) {
	var err *Error
	requireEqual(t, err.Error(), "<nil>")
	if err.Unwrap() != nil {
		t.Fatalf("nil Error Unwrap() must return nil")
	}
}

func TestErrorUnwrapWithOnlyCause(t *testing.T) {
	cause := errors.New("cause")
	requireErrorIs(t, &Error{Cause: cause}, cause)
}
