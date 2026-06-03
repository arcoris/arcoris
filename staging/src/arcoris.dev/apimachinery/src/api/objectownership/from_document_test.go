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
)

func TestStateFromDocumentEmptyDocument(t *testing.T) {
	got, err := StateFromDocument(Document{Version: VersionV1})
	requireNoError(t, err)

	if !got.IsEmpty() {
		t.Fatalf("state from empty document IsEmpty() = false")
	}
}

func TestStateFromDocumentDesiredOwnership(t *testing.T) {
	got, err := StateFromDocument(document(documentEntry("user", "$.image")))
	requireNoError(t, err)

	requireOwners(t, got.Desired().OwnersOf(path("$.image")), "user")
}

func TestStateFromDocumentRejectsZeroVersion(t *testing.T) {
	_, err := StateFromDocument(Document{})

	requireErrorIs(t, err, ErrInvalidDocument)
}

func TestStateFromDocumentRejectsUnsupportedVersion(t *testing.T) {
	_, err := StateFromDocument(Document{Version: "v2"})

	requireErrorIs(t, err, ErrUnsupportedVersion)
}

func TestStateFromDocumentRejectsInvalidOwner(t *testing.T) {
	_, err := StateFromDocument(document(documentEntry(" ", "$.image")))

	requireErrorIs(t, err, ErrInvalidEntry)
	requireErrorIs(t, err, fieldownership.ErrInvalidOwner)
}

func TestStateFromDocumentRejectsInvalidPath(t *testing.T) {
	_, err := StateFromDocument(document(documentEntry("user", "image")))

	requireErrorIs(t, err, ErrInvalidPath)
}

func TestStateFromDocumentMergesDuplicateOwners(t *testing.T) {
	got, err := StateFromDocument(document(
		documentEntry("user", "$.image"),
		documentEntry("user", "$.replicas"),
	))
	requireNoError(t, err)

	if owners := got.Desired().Owners(); len(owners) != 1 {
		t.Fatalf("owners = %#v; want one owner", owners)
	}
	requirePaths(t, got.Desired().FieldsFor(owner("user")), "$.image", "$.replicas")
}

func TestStateFromDocumentDeduplicatesDuplicateFields(t *testing.T) {
	got, err := StateFromDocument(document(documentEntry("user", "$.image", "$.image")))
	requireNoError(t, err)

	requirePaths(t, got.Desired().FieldsFor(owner("user")), "$.image")
}

func TestStateFromDocumentPrunesEmptyEntries(t *testing.T) {
	got, err := StateFromDocument(document(documentEntry("user")))
	requireNoError(t, err)

	if !got.IsEmpty() {
		t.Fatalf("state from empty entry IsEmpty() = false")
	}
}

func TestStateFromDocumentPreservesSharedOwnership(t *testing.T) {
	got, err := StateFromDocument(document(documentEntry("a", "$.image"), documentEntry("b", "$.image")))
	requireNoError(t, err)

	requireOwners(t, got.Desired().OwnersOf(path("$.image")), "a", "b")
}

func TestStateFromDocumentPreservesAncestorDescendantFields(t *testing.T) {
	got, err := StateFromDocument(document(documentEntry("user", "$.spec", "$.spec.image")))
	requireNoError(t, err)

	requirePaths(t, got.Desired().FieldsFor(owner("user")), "$.spec", "$.spec.image")
}

func TestStateFromDocumentReturnsNormalizedState(t *testing.T) {
	got, err := StateFromDocument(document(documentEntry("z", "$.z"), documentEntry("a", "$.a")))
	requireNoError(t, err)

	requireOwners(t, got.Desired().Owners(), "a", "z")
}
