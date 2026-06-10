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
	"strings"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestErrorIsInvalidOwner(t *testing.T) {
	requireErrorIs(t, Owner{}.ValidateLexical(), ErrInvalidOwner)
}

func TestErrorIsInvalidEntry(t *testing.T) {
	requireErrorIs(t, Entry{}.ValidateStructure(), ErrInvalidEntry)
}

func TestErrorIsInvalidState(t *testing.T) {
	state := State{entries: []Entry{emptyEntry("user")}}

	requireErrorIs(t, state.ValidateStructure(), ErrInvalidState)
}

func TestErrorIsInvalidOwnedPath(t *testing.T) {
	err := OwnedPath{Owner: Owner{}, Path: imagePath()}.ValidateStructure()

	requireErrorIs(t, err, ErrInvalidOwnedPath)
}

func TestErrorIsInvalidConflict(t *testing.T) {
	err := Conflict{Owner: Owner{}, OwnedPath: imagePath(), AttemptedPath: imagePath()}.ValidateStructure()

	requireErrorIs(t, err, ErrInvalidConflict)
}

func TestErrorIsInvalidPath(t *testing.T) {
	err := wrapPathError(
		fieldpath.Root(),
		"bad path",
		errors.New("cause"),
	)

	requireErrorIs(t, err, ErrInvalidPath)
}

func TestErrorMessageIncludesReason(t *testing.T) {
	err := Owner{}.ValidateLexical()

	if !strings.Contains(err.Error(), string(ErrorReasonEmptyOwner)) {
		t.Fatalf("error %q does not contain reason", err.Error())
	}
}

func TestNilErrorString(t *testing.T) {
	var err *Error

	requireEqual(t, err.Error(), "<nil>")
}

func TestNilErrorUnwrap(t *testing.T) {
	var err *Error

	requireEqual(t, err.Unwrap() == nil, true)
}
