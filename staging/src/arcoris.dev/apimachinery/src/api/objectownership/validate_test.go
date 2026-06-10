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
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestValidateAcceptsDocumentVersionV1(t *testing.T) {
	requireNoError(t, Validate(Document{Version: DocumentVersionV1}))
}

func TestValidateRejectsZeroVersion(t *testing.T) {
	err := Validate(Document{})

	requireErrorIs(t, err, ErrInvalidDocument)
	requireObjectOwnershipError(t, err, pathDocumentVersion, ErrorReasonInvalidDocument)
}

func TestValidateRejectsUnsupportedVersion(t *testing.T) {
	err := Validate(Document{Version: "v2"})

	requireErrorIs(t, err, ErrUnsupportedVersion)
	requireObjectOwnershipError(t, err, pathDocumentVersion, ErrorReasonUnsupportedVersion)
}

func TestValidateAcceptsUnsortedEntries(t *testing.T) {
	err := Validate(document(documentEntry("z", "$.z"), documentEntry("a", "$.a")))

	requireNoError(t, err)
}

func TestValidateAcceptsDuplicateOwners(t *testing.T) {
	err := Validate(document(documentEntry("user", "$.image"), documentEntry("user", "$.replicas")))

	requireNoError(t, err)
}

func TestValidateAcceptsEmptyFieldEntryWithValidOwner(t *testing.T) {
	err := Validate(document(documentEntry("user")))

	requireNoError(t, err)
}

func TestValidateRejectsInvalidOwner(t *testing.T) {
	err := Validate(document(Entry{
		Owner:  fieldownership.Owner{},
		Fields: []Path{"$.image"},
	}))

	requireErrorIs(t, err, ErrInvalidEntry)
	requireErrorIs(t, err, fieldownership.ErrInvalidOwner)
	requireObjectOwnershipError(
		t,
		err,
		"document.desired.entries[0].owner",
		ErrorReasonInvalidOwner,
	)
}

func TestValidateRejectsEmptyPath(t *testing.T) {
	err := Validate(document(documentEntry("user", "")))

	requireErrorIs(t, err, ErrInvalidPath)
	requireObjectOwnershipError(
		t,
		err,
		"document.desired.entries[0].fields[0]",
		ErrorReasonInvalidPath,
	)
}

func TestValidateRejectsInvalidPathGrammar(t *testing.T) {
	err := Validate(document(documentEntry("user", "image")))

	requireErrorIs(t, err, ErrInvalidPath)
	requireErrorIs(t, err, fieldpath.ErrInvalidPath)
}

func TestValidateRejectsNonCanonicalPathIfApplicable(t *testing.T) {
	err := Validate(document(documentEntry("user", `$."image"`)))

	requireErrorIs(t, err, ErrInvalidPath)
}

func TestValidateRejectsNonCanonicalPathWithObjectOwnershipError(t *testing.T) {
	err := Validate(document(documentEntry("user", `$."image"`)))

	requireErrorIs(t, err, ErrInvalidPath)
	requireObjectOwnershipError(
		t,
		err,
		"document.desired.entries[0].fields[0]",
		ErrorReasonInvalidPath,
	)
}
