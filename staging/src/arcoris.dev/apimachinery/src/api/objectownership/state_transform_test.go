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

func TestStateWithDesired(t *testing.T) {
	original := NewStateWithSurfaces(
		ownershipState(ownershipEntry("old", "$.image")),
		ownershipState(ownershipEntry("observed", "$.ready")),
		NewMetadataState(ownershipState(ownershipEntry("labels", `$["app"]`)), ownershipState()),
	)
	replacement := ownershipState(ownershipEntry("new", "$.replicas"))

	got := original.WithDesired(replacement)

	requireOwnersOf(t, got.Desired(), path("$.replicas"), "new")
	requireOwnersOf(t, got.Observed(), path("$.ready"), "observed")
	requireOwnersOf(t, got.Metadata().Labels(), path(`$["app"]`), "labels")
	requireOwnersOf(t, original.Desired(), path("$.image"), "old")
}

func TestStateWithObservedPreservesDesiredAndMetadata(t *testing.T) {
	original := NewStateWithSurfaces(
		ownershipState(ownershipEntry("desired", "$.image")),
		ownershipState(ownershipEntry("old", "$.ready")),
		NewMetadataState(ownershipState(ownershipEntry("labels", `$["app"]`)), ownershipState()),
	)

	got := original.WithObserved(ownershipState(ownershipEntry("new", "$.phase")))

	requireOwnersOf(t, got.Desired(), path("$.image"), "desired")
	requireOwnersOf(t, got.Observed(), path("$.phase"), "new")
	requireOwnersOf(t, got.Metadata().Labels(), path(`$["app"]`), "labels")
	requireOwnersOf(t, original.Observed(), path("$.ready"), "old")
}

func TestStateWithMetadataPreservesDesiredAndObserved(t *testing.T) {
	original := NewStateWithSurfaces(
		ownershipState(ownershipEntry("desired", "$.image")),
		ownershipState(ownershipEntry("observed", "$.ready")),
		NewMetadataState(ownershipState(ownershipEntry("old-label", `$["app"]`)), ownershipState()),
	)
	replacement := NewMetadataState(
		ownershipState(ownershipEntry("new-label", `$["team"]`)),
		ownershipState(ownershipEntry("new-annotation", `$["note"]`)),
	)

	got := original.WithMetadata(replacement)

	requireOwnersOf(t, got.Desired(), path("$.image"), "desired")
	requireOwnersOf(t, got.Observed(), path("$.ready"), "observed")
	requireOwnersOf(t, got.Metadata().Labels(), path(`$["team"]`), "new-label")
	requireOwnersOf(t, got.Metadata().Annotations(), path(`$["note"]`), "new-annotation")
	requireOwnersOf(t, original.Metadata().Labels(), path(`$["app"]`), "old-label")
}

func TestStateImmutability(t *testing.T) {
	original := NewState(ownershipState(ownershipEntry("user", "$.image")))

	_ = original.WithDesired(ownershipState())

	requireOwnersOf(t, original.Desired(), path("$.image"), "user")
}
