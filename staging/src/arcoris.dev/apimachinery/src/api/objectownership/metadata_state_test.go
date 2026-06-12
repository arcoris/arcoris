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

func TestNewMetadataState(t *testing.T) {
	got := NewMetadataState(
		ownershipState(ownershipEntry("labels", `$["app"]`)),
		ownershipState(ownershipEntry("annotations", `$["scheduler.arcoris.dev/mode"]`)),
	)

	requireOwnersOf(t, got.Labels(), path(`$["app"]`), "labels")
	requireOwnersOf(t, got.Annotations(), path(`$["scheduler.arcoris.dev/mode"]`), "annotations")
}

func TestMetadataStateIsEmpty(t *testing.T) {
	if !NewMetadataState(ownershipState(), ownershipState()).IsEmpty() {
		t.Fatalf("empty metadata state IsEmpty() = false")
	}
	if NewMetadataState(ownershipState(ownershipEntry("labels", `$["app"]`)), ownershipState()).IsEmpty() {
		t.Fatalf("metadata state with labels IsEmpty() = true")
	}
	if NewMetadataState(ownershipState(), ownershipState(ownershipEntry("annotations", `$["note"]`))).IsEmpty() {
		t.Fatalf("metadata state with annotations IsEmpty() = true")
	}
}

func TestMetadataStateWithLabelsPreservesAnnotations(t *testing.T) {
	original := NewMetadataState(
		ownershipState(ownershipEntry("old", `$["app"]`)),
		ownershipState(ownershipEntry("annotations", `$["note"]`)),
	)

	got := original.WithLabels(ownershipState(ownershipEntry("new", `$["team"]`)))

	requireOwnersOf(t, got.Labels(), path(`$["team"]`), "new")
	requireOwnersOf(t, got.Annotations(), path(`$["note"]`), "annotations")
	requireOwnersOf(t, original.Labels(), path(`$["app"]`), "old")
}

func TestMetadataStateWithAnnotationsPreservesLabels(t *testing.T) {
	original := NewMetadataState(
		ownershipState(ownershipEntry("labels", `$["app"]`)),
		ownershipState(ownershipEntry("old", `$["note"]`)),
	)

	got := original.WithAnnotations(ownershipState(ownershipEntry("new", `$["scheduler.arcoris.dev/mode"]`)))

	requireOwnersOf(t, got.Labels(), path(`$["app"]`), "labels")
	requireOwnersOf(t, got.Annotations(), path(`$["scheduler.arcoris.dev/mode"]`), "new")
	requireOwnersOf(t, original.Annotations(), path(`$["note"]`), "old")
}
