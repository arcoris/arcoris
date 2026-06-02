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

package valuefieldset

import (
	"errors"
	"testing"
)

func TestErrorAtBuildsStructuredDiagnostic(t *testing.T) {
	err := errorAt(
		rootField("spec"),
		ErrInvalidDescriptor,
		ErrorReasonInvalidDescriptor,
		"descriptor is invalid",
	)

	requireErrorIs(t, err, ErrInvalidDescriptor)
	requireErrorPath(t, err, "$.spec")
	requireErrorReason(t, err, ErrorReasonInvalidDescriptor)
	requireErrorDetailContains(t, err, "descriptor is invalid")
}

func TestErrorfAtFormatsDetail(t *testing.T) {
	err := errorfAt(
		rootField("items").Index(1),
		ErrInvalidListKey,
		ErrorReasonMissingListKey,
		"key %q is missing",
		"type",
	)

	requireErrorDetailContains(t, err, `key "type" is missing`)
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("cause")
	err := wrapAt(
		rootField("spec"),
		ErrInvalidDescriptor,
		ErrorReasonInvalidDescriptor,
		"wrapped",
		cause,
	)

	requireErrorIs(t, err, ErrInvalidDescriptor)
	requireErrorIs(t, err, cause)
}
