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

func TestEmptyState(t *testing.T) {
	if !EmptyState().IsEmpty() {
		t.Fatalf("EmptyState().IsEmpty() = false")
	}
}

func TestNewState(t *testing.T) {
	desired := ownershipState(ownershipEntry("user", "$.image"))

	got := NewState(desired)

	requireOwnersOf(t, got.Desired(), path("$.image"), "user")
	if !got.Observed().IsEmpty() {
		t.Fatalf("NewState observed surface is not empty")
	}
	if !got.Metadata().IsEmpty() {
		t.Fatalf("NewState metadata surface is not empty")
	}
}

func TestNewStateWithSurfaces(t *testing.T) {
	desired := ownershipState(ownershipEntry("desired", "$.image"))
	observed := ownershipState(ownershipEntry("observed", "$.ready"))
	metadata := NewMetadataState(
		ownershipState(ownershipEntry("labels", `$["app"]`)),
		ownershipState(ownershipEntry("annotations", `$["scheduler.arcoris.dev/mode"]`)),
	)

	got := NewStateWithSurfaces(desired, observed, metadata)

	requireOwnersOf(t, got.Desired(), path("$.image"), "desired")
	requireOwnersOf(t, got.Observed(), path("$.ready"), "observed")
	requireOwnersOf(t, got.Metadata().Labels(), path(`$["app"]`), "labels")
	requireOwnersOf(t, got.Metadata().Annotations(), path(`$["scheduler.arcoris.dev/mode"]`), "annotations")
}
