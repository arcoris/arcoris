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

func TestEmptyState(t *testing.T) {
	if !EmptyState().IsEmpty() {
		t.Fatalf("EmptyState().IsEmpty() = false")
	}
}

func TestNewState(t *testing.T) {
	desired := ownershipState(ownershipEntry("user", "$.image"))

	got := NewState(desired)

	requireOwners(t, got.Desired().OwnersOf(path("$.image")), "user")
}

func TestStateDesired(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("user", "$.image")))

	requireOwners(t, state.Desired().OwnersOf(path("$.image")), "user")
}

func TestStateWithDesired(t *testing.T) {
	original := NewState(ownershipState(ownershipEntry("old", "$.image")))
	replacement := ownershipState(ownershipEntry("new", "$.replicas"))

	got := original.WithDesired(replacement)

	requireOwners(t, got.Desired().OwnersOf(path("$.replicas")), "new")
	requireOwners(t, original.Desired().OwnersOf(path("$.image")), "old")
}

func TestStateIsEmpty(t *testing.T) {
	state := NewState(ownershipState(ownershipEntry("user", "$.image")))
	if state.IsEmpty() {
		t.Fatalf("non-empty state IsEmpty() = true")
	}
}

func TestStateImmutability(t *testing.T) {
	original := NewState(ownershipState(ownershipEntry("user", "$.image")))

	_ = original.WithDesired(fieldownership.EmptyState())

	requireOwners(t, original.Desired().OwnersOf(path("$.image")), "user")
}
