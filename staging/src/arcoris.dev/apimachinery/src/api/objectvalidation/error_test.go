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

package objectvalidation

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorFormatting(t *testing.T) {
	err := &Error{
		Path:   pathObjectDesired,
		Err:    ErrInvalidDesired,
		Reason: ErrorReasonInvalidDesired,
		Detail: "desired surface is invalid",
	}

	got := err.Error()
	for _, want := range []string{
		"objectvalidation",
		pathObjectDesired,
		ErrInvalidDesired.Error(),
		string(ErrorReasonInvalidDesired),
		"desired surface is invalid",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("Error() = %q, want segment %q", got, want)
		}
	}
}

func TestErrorUnwrapPreservesObjectSentinelAndCause(t *testing.T) {
	cause := errors.New("surface failed")
	err := nested(
		pathObjectDesired,
		ErrInvalidDesired,
		ErrorReasonInvalidDesired,
		"desired surface is invalid",
		cause,
	)

	requireErrorIs(t, err, ErrInvalidObject)
	requireErrorIs(t, err, ErrInvalidDesired)
	requireErrorIs(t, err, cause)
}

func TestErrorUnwrapPreservesInvalidPlanForMissingValidator(t *testing.T) {
	err := missingValidator(pathPlanDesiredValidator, "desired surface validator is required")

	requireErrorIs(t, err, ErrInvalidPlan)
	requireErrorIs(t, err, ErrMissingValidator)
}
