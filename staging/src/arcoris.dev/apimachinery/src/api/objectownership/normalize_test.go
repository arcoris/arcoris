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
