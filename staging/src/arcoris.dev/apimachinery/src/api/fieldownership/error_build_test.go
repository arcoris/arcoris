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

package fieldownership

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestErrorAtBuildsStructuredError(t *testing.T) {
	err := errorAt("owner", ErrInvalidOwner, ErrorReasonEmptyOwner, "owner is empty")
	want := "fieldownership: owner: invalid field owner: empty_owner: owner is empty"

	requireErrorIs(t, err, ErrInvalidOwner)
	requireEqual(t, err.Error(), want)
}

func TestErrorfAtFormatsDetail(t *testing.T) {
	err := errorfAt(
		"",
		ErrInvalidOwner,
		ErrorReasonOwnerTooLong,
		"owner exceeds %d bytes",
		MaxOwnerLength,
	)
	want := "fieldownership: " +
		"invalid field owner: owner_too_long: owner exceeds 128 bytes"

	requireErrorIs(t, err, ErrInvalidOwner)
	requireEqual(t, err.Error(), want)
}

func TestWrapAtPreservesCause(t *testing.T) {
	cause := errors.New("cause")
	err := wrapAt("", ErrInvalidEntry, ErrorReasonInvalidEntry, "entry is invalid", cause)

	requireErrorIs(t, err, ErrInvalidEntry)
	requireErrorIs(t, err, cause)
}

func TestWrapPathErrorPreservesInvalidPathCause(t *testing.T) {
	invalidPath := fieldpath.Root()
	_, cause := fieldpath.NewPath(fieldpath.Element{})
	err := wrapPathError(invalidPath, "path is invalid", cause)

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, cause)
}
