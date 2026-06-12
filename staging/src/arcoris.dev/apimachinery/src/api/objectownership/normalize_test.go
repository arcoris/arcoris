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

func TestNormalizeWritesDocumentVersionV1(t *testing.T) {
	got, err := Normalize(document(documentEntry("user", "$.image")))
	requireNoError(t, err)

	if got.Version != DocumentVersionV1 {
		t.Fatalf("Version = %q; want %q", got.Version, DocumentVersionV1)
	}
}

func TestNormalizeSortsEntriesByOwner(t *testing.T) {
	got, err := Normalize(document(documentEntry("z", "$.z"), documentEntry("a", "$.a")))
	requireNoError(t, err)

	requireDocumentEntries(t, got.Desired, documentEntry("a", "$.a"), documentEntry("z", "$.z"))
}

func TestNormalizeMergesDuplicateOwners(t *testing.T) {
	got, err := Normalize(document(documentEntry("user", "$.image"), documentEntry("user", "$.replicas")))
	requireNoError(t, err)

	requireDocumentEntries(t, got.Desired, documentEntry("user", "$.image", "$.replicas"))
}

func TestNormalizeDeduplicatesFields(t *testing.T) {
	got, err := Normalize(document(documentEntry("user", "$.image", "$.image")))
	requireNoError(t, err)

	requireDocumentEntries(t, got.Desired, documentEntry("user", "$.image"))
}

func TestNormalizePrunesEmptyEntries(t *testing.T) {
	got, err := Normalize(document(documentEntry("user")))
	requireNoError(t, err)

	if !got.Desired.IsEmpty() {
		t.Fatalf("normalized desired = %#v; want empty", got.Desired)
	}
}

func TestNormalizeSortsFields(t *testing.T) {
	got, err := Normalize(document(documentEntry("user", "$.z", "$.a")))
	requireNoError(t, err)

	requireDocumentEntries(t, got.Desired, documentEntry("user", "$.a", "$.z"))
}

func TestNormalizePreservesSharedOwnership(t *testing.T) {
	got, err := Normalize(document(documentEntry("a", "$.image"), documentEntry("b", "$.image")))
	requireNoError(t, err)

	requireDocumentEntries(t, got.Desired, documentEntry("a", "$.image"), documentEntry("b", "$.image"))
}

func TestNormalizePreservesAncestorDescendantFields(t *testing.T) {
	got, err := Normalize(document(documentEntry("user", "$.spec.image", "$.spec")))
	requireNoError(t, err)

	requireDocumentEntries(t, got.Desired, documentEntry("user", "$.spec", "$.spec.image"))
}

func TestNormalizeDoesNotMutateInput(t *testing.T) {
	doc := document(documentEntry("user", "$.z", "$.a"))

	_, err := Normalize(doc)
	requireNoError(t, err)

	requireDocumentEntries(t, doc.Desired, documentEntry("user", "$.z", "$.a"))
}

func TestNormalizeRejectsInvalidRawDocument(t *testing.T) {
	_, err := Normalize(Document{})

	requireErrorIs(t, err, ErrInvalidDocument)
}

func TestNormalizeOutputPassesValidateNormalized(t *testing.T) {
	normalized, err := Normalize(document(
		documentEntry("b", "$.z", "$.z"),
		documentEntry("a", "$.spec.image", "$.spec"),
		documentEntry("b", "$.a"),
		documentEntry("empty"),
	))
	requireNoError(t, err)

	requireNoError(t, ValidateNormalized(normalized))
}
