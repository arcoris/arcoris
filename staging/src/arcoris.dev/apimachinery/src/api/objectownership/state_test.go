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

func TestStateDesired(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("user", "$.image")))

	requireOwnersOf(t, state.Desired(), path("$.image"), "user")
}

func TestStateObserved(t *testing.T) {
	state := NewStateWithSurfaces(
		ownershipState(ownershipEntry("desired", "$.image")),
		ownershipState(ownershipEntry("observer", "$.ready")),
		MetadataState{},
	)

	requireOwnersOf(t, state.Observed(), path("$.ready"), "observer")
}

func TestStateMetadata(t *testing.T) {
	state := NewStateWithSurfaces(
		ownershipState(ownershipEntry("desired", "$.image")),
		ownershipState(ownershipEntry("observer", "$.ready")),
		NewMetadataState(
			ownershipState(ownershipEntry("labeler", `$["app"]`)),
			ownershipState(ownershipEntry("annotator", `$["note"]`)),
		),
	)

	requireOwnersOf(t, state.Metadata().Labels(), path(`$["app"]`), "labeler")
	requireOwnersOf(t, state.Metadata().Annotations(), path(`$["note"]`), "annotator")
}

func TestStateIsEmpty(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("user", "$.image")))
	if state.IsEmpty() {
		t.Fatalf("non-empty state IsEmpty() = true")
	}
}
