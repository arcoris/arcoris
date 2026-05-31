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

func TestErrorfBuildsStructuredError(t *testing.T) {
	err := errorf(
		pathFloat,
		ErrInvalidFloat,
		ErrorReasonInvalidFloat,
		"float %q is not finite",
		"NaN",
	)

	valueErr := requireValueError(
		t,
		err,
		ErrInvalidFloat,
		pathFloat,
		ErrorReasonInvalidFloat,
	)

	requireEqual(t, valueErr.Detail, `float "NaN" is not finite`)
}

func TestNestedBuildsStructuredErrorWithCause(t *testing.T) {
	cause := errors.New("nested")
	err := nested(
		"object.members[0].value",
		ErrInvalidMember,
		ErrorReasonInvalidMember,
		"object member value is invalid",
		cause,
	)

	valueErr := requireValueError(
		t,
		err,
		ErrInvalidMember,
		"object.members[0].value",
		ErrorReasonInvalidMember,
	)

	requireEqual(t, valueErr.Cause, cause)
	requireErrorIs(t, err, cause)
}
