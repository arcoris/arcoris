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

package valuecompare

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestErrorAtBuildsStructuredError(t *testing.T) {
	err := errorAt(fieldpath.Root(), ErrInvalidValue, ErrorReasonInvalidZero, "bad")

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorReason(t, err, ErrorReasonInvalidZero)
	requireErrorPath(t, err, "$")
	requireErrorDetailContains(t, err, "bad")
}

func TestErrorfAtBuildsFormattedDetail(t *testing.T) {
	err := errorfAt(fieldpath.Root(), ErrUnknownField, ErrorReasonUnknownField, "field %q", "extra")

	requireErrorDetailContains(t, err, `field "extra"`)
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("cause")
	err := wrapAt(fieldpath.Root(), ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "bad", cause)

	requireErrorIs(t, err, ErrInvalidDescriptor)
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is(err, cause) = false")
	}
}
