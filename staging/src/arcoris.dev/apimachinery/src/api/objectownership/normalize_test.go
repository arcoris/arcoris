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
	"reflect"
	"testing"
)

func TestNormalizeEmptyState(t *testing.T) {
	got := Normalize(EmptyState())

	if !got.IsEmpty() {
		t.Fatalf("Normalize(EmptyState()).IsEmpty() = false")
	}
}

func TestNormalizePreservesEverySurface(t *testing.T) {
	state := NewStateWithSurfaces(
		ownershipState(ownershipEntry("desired", "$.image")),
		ownershipState(ownershipEntry("observed", "$.ready")),
		NewMetadataState(
			ownershipState(ownershipEntry("labels", `$["app"]`)),
			ownershipState(ownershipEntry("annotations", `$["scheduler.arcoris.dev/mode"]`)),
		),
	)

	got := Normalize(state)

	requireOwnersOf(t, got.Desired(), path("$.image"), "desired")
	requireOwnersOf(t, got.Observed(), path("$.ready"), "observed")
	requireOwnersOf(t, got.Metadata().Labels(), path(`$["app"]`), "labels")
	requireOwnersOf(t, got.Metadata().Annotations(), path(`$["scheduler.arcoris.dev/mode"]`), "annotations")
}

func TestNormalizeDoesNotMutateInput(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("desired", "$.image")))

	_ = Normalize(state)

	requireOwnersOf(t, state.Desired(), path("$.image"), "desired")
}

func TestNormalizeIsIdempotent(t *testing.T) {
	state := NewStateWithSurfaces(
		ownershipState(
			ownershipEntry("zeta", "$.replicas"),
			ownershipEntry("alpha", "$.image"),
			ownershipEntry("alpha", "$.replicas", "$.image"),
		),
		ownershipState(
			ownershipEntry("observer", "$.ready"),
			ownershipEntry("observer", "$.phase"),
		),
		NewMetadataState(
			ownershipState(
				ownershipEntry("labeler", `$["team"]`),
				ownershipEntry("labeler", `$["app"]`),
			),
			ownershipState(
				ownershipEntry("annotator", `$["scheduler.arcoris.dev/mode"]`),
				ownershipEntry("annotator", `$["note"]`),
			),
		),
	)

	once := Normalize(state)
	twice := Normalize(once)

	if !reflect.DeepEqual(twice, once) {
		t.Fatalf("Normalize(Normalize(state)) != Normalize(state)")
	}
}

func TestNormalizeIsDeterministicAcrossEquivalentInputOrder(t *testing.T) {
	first := NewStateWithSurfaces(
		ownershipState(
			ownershipEntry("zeta", "$.replicas"),
			ownershipEntry("alpha", "$.image"),
			ownershipEntry("alpha", "$.replicas"),
		),
		ownershipState(ownershipEntry("observer", "$.ready")),
		NewMetadataState(
			ownershipState(ownershipEntry("labeler", `$["team"]`), ownershipEntry("labeler", `$["app"]`)),
			ownershipState(ownershipEntry("annotator", `$["note"]`)),
		),
	)
	second := NewStateWithSurfaces(
		ownershipState(
			ownershipEntry("alpha", "$.replicas"),
			ownershipEntry("zeta", "$.replicas"),
			ownershipEntry("alpha", "$.image"),
		),
		ownershipState(ownershipEntry("observer", "$.ready")),
		NewMetadataState(
			ownershipState(ownershipEntry("labeler", `$["app"]`), ownershipEntry("labeler", `$["team"]`)),
			ownershipState(ownershipEntry("annotator", `$["note"]`)),
		),
	)

	if !reflect.DeepEqual(Normalize(first), Normalize(second)) {
		t.Fatalf("equivalent ownership states normalized differently")
	}
}

func TestNormalizeDoesNotMergeSurfaces(t *testing.T) {
	state := NewStateWithSurfaces(
		ownershipState(ownershipEntry("desired-owner", "$.ready")),
		ownershipState(ownershipEntry("observed-owner", "$.ready")),
		NewMetadataState(
			ownershipState(ownershipEntry("label-owner", `$["ready"]`)),
			ownershipState(ownershipEntry("annotation-owner", `$["ready"]`)),
		),
	)

	got := Normalize(state)

	requireOwnersOf(t, got.Desired(), path("$.ready"), "desired-owner")
	requireOwnersOf(t, got.Observed(), path("$.ready"), "observed-owner")
	requireOwnersOf(t, got.Metadata().Labels(), path(`$["ready"]`), "label-owner")
	requireOwnersOf(t, got.Metadata().Annotations(), path(`$["ready"]`), "annotation-owner")
}
