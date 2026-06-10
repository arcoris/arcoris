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

package objectownership

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestInvalidPathErrorWrapsObjectOwnershipInvalidPath(t *testing.T) {
	err := Validate(document(documentEntry("user", "image")))

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, fieldpath.ErrInvalidPath)
}

func TestInvalidOwnerErrorWrapsFieldOwnershipInvalidOwner(t *testing.T) {
	err := Validate(document(Entry{
		Owner:  fieldownership.Owner{},
		Fields: []Path{"$.image"},
	}))

	requireErrorIs(t, err, ErrInvalidEntry)
	requireErrorIs(t, err, fieldownership.ErrInvalidOwner)
}

func TestUnsupportedVersionError(t *testing.T) {
	err := Validate(Document{Version: "v2"})

	requireErrorIs(t, err, ErrUnsupportedVersion)
}

func TestErrorAsObjectOwnershipError(t *testing.T) {
	err := Validate(Document{})

	var ownershipErr *Error
	if !errors.As(err, &ownershipErr) {
		t.Fatalf("errors.As(%T) = false", ownershipErr)
	}
}

func TestErrorDiagnosticPath(t *testing.T) {
	err := Validate(document(documentEntry("user", "")))

	requireObjectOwnershipError(
		t,
		err,
		"document.desired.entries[0].fields[0]",
		ErrorReasonInvalidPath,
	)
}

func TestNilError(t *testing.T) {
	var err *Error
	if err.Error() != "<nil>" {
		t.Fatalf("Error() = %q; want <nil>", err.Error())
	}
	if err.Unwrap() != nil {
		t.Fatalf("Unwrap() = %v; want nil", err.Unwrap())
	}
}
