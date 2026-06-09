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

func TestToDocumentEmptyState(t *testing.T) {
	got := ToDocument(EmptyState())

	if got.Version != DocumentVersionV1 {
		t.Fatalf("Version = %q; want %q", got.Version, DocumentVersionV1)
	}
	if !got.Desired.IsEmpty() {
		t.Fatalf("Desired = %#v; want empty", got.Desired)
	}
}

func TestToDocumentDesiredOwnership(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("user", "$.image")))

	got := ToDocument(state)

	requireDocumentEntries(t, got.Desired, documentEntry("user", "$.image"))
}

func TestToDocumentWritesDocumentVersionV1(t *testing.T) {
	if got := ToDocument(EmptyState()).Version; got != DocumentVersionV1 {
		t.Fatalf("ToDocument().Version = %q; want %q", got, DocumentVersionV1)
	}
}

func TestToDocumentDeterministicOwnerOrder(t *testing.T) {
	state := NewState(ownershipState(
		ownershipEntry("z", "$.z"),
		ownershipEntry("a", "$.a"),
	))

	got := ToDocument(state)

	requireDocumentEntries(t, got.Desired, documentEntry("a", "$.a"), documentEntry("z", "$.z"))
}

func TestToDocumentDeterministicFieldOrder(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("user", "$.z", "$.a")))

	got := ToDocument(state)

	requireDocumentEntries(t, got.Desired, documentEntry("user", "$.a", "$.z"))
}

func TestToDocumentPreservesSharedOwnership(t *testing.T) {
	state := NewState(ownershipState(
		ownershipEntry("a", "$.image"),
		ownershipEntry("b", "$.image"),
	))

	got := ToDocument(state)

	requireDocumentEntries(t, got.Desired, documentEntry("a", "$.image"), documentEntry("b", "$.image"))
}

func TestToDocumentPreservesAncestorDescendantFields(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("user", "$.spec", "$.spec.image")))

	got := ToDocument(state)

	requireDocumentEntries(t, got.Desired, documentEntry("user", "$.spec", "$.spec.image"))
}

func TestToDocumentDoesNotMutateState(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("user", "$.image")))

	got := ToDocument(state)
	got.Desired.Entries[0].Owner = owner("other")

	requireOwners(t, state.Desired().OwnersOf(path("$.image")), "user")
}
