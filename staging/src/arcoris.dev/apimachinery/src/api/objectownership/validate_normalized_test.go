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

import "testing"

func TestValidateNormalizedAcceptsNormalizeOutput(t *testing.T) {
	normalized, err := Normalize(document(
		documentEntry("z", "$.z", "$.z"),
		documentEntry("a", "$.spec.image", "$.spec"),
		documentEntry("z", "$.a"),
		documentEntry("empty"),
	))
	requireNoError(t, err)

	requireNoError(t, ValidateNormalized(normalized))
}

func TestValidateNormalizedRejectsUnsortedEntries(t *testing.T) {
	err := ValidateNormalized(document(documentEntry("z", "$.z"), documentEntry("a", "$.a")))

	requireErrorIs(t, err, ErrNotNormalized)
	requireObjectOwnershipError(t, err, pathDocument, ErrorReasonNotNormalized)
}

func TestValidateNormalizedRejectsDuplicateOwnerEntries(t *testing.T) {
	err := ValidateNormalized(document(documentEntry("user", "$.image"), documentEntry("user", "$.replicas")))

	requireErrorIs(t, err, ErrNotNormalized)
	requireObjectOwnershipError(t, err, pathDocument, ErrorReasonNotNormalized)
}

func TestValidateNormalizedRejectsDuplicateFields(t *testing.T) {
	err := ValidateNormalized(document(documentEntry("user", "$.image", "$.image")))

	requireErrorIs(t, err, ErrNotNormalized)
	requireObjectOwnershipError(t, err, pathDocument, ErrorReasonNotNormalized)
}

func TestValidateNormalizedRejectsEmptyEntries(t *testing.T) {
	err := ValidateNormalized(document(documentEntry("user")))

	requireErrorIs(t, err, ErrNotNormalized)
	requireObjectOwnershipError(t, err, pathDocument, ErrorReasonNotNormalized)
}

func TestValidateNormalizedRejectsNonNilEmptySurface(t *testing.T) {
	err := ValidateNormalized(Document{
		Version: DocumentVersionV1,
		Desired: Surface{Entries: []Entry{}},
	})

	requireErrorIs(t, err, ErrNotNormalized)
	requireObjectOwnershipError(t, err, pathDocument, ErrorReasonNotNormalized)
}

func TestValidateNormalizedPreservesValidationErrors(t *testing.T) {
	err := ValidateNormalized(Document{})

	requireErrorIs(t, err, ErrInvalidDocument)
	requireObjectOwnershipError(t, err, pathDocumentVersion, ErrorReasonMissingVersion)
}

func TestValidateNormalizedRejectsUnsupportedVersion(t *testing.T) {
	err := ValidateNormalized(Document{Version: "v2"})

	requireErrorIs(t, err, ErrUnsupportedVersion)
	requireObjectOwnershipError(t, err, pathDocumentVersion, ErrorReasonUnsupportedVersion)
}
