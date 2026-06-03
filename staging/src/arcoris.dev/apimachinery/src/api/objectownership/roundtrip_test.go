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

func TestDocumentRoundTripEmptyState(t *testing.T) {
	state, err := StateFromDocument(ToDocument(EmptyState()))
	requireNoError(t, err)

	if !state.IsEmpty() {
		t.Fatalf("round-trip empty state IsEmpty() = false")
	}
}

func TestDocumentRoundTripDesiredState(t *testing.T) {
	state := NewState(ownershipState(
		ownershipEntry("a", "$.spec", "$.spec.image"),
		ownershipEntry("b", "$.spec.image"),
	))

	roundTripped, err := StateFromDocument(ToDocument(state))
	requireNoError(t, err)

	requirePaths(t, roundTripped.Desired().FieldsFor(owner("a")), "$.spec", "$.spec.image")
	requirePaths(t, roundTripped.Desired().FieldsFor(owner("b")), "$.spec.image")
}

func TestNormalizeDocumentRoundTrip(t *testing.T) {
	doc := document(
		documentEntry("b", "$.z", "$.z"),
		documentEntry("a", "$.spec.image", "$.spec"),
		documentEntry("b", "$.a"),
	)

	normalized, err := Normalize(doc)
	requireNoError(t, err)

	state, err := StateFromDocument(normalized)
	requireNoError(t, err)

	requireDocumentEntries(t, ToDocument(state).Desired, normalized.Desired.Entries...)
}
