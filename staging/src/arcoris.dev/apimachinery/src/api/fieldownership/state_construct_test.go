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

package fieldownership

import (
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestEmptyState(t *testing.T) {
	state := EmptyState()

	requireEqual(t, state.IsEmpty(), true)
	requireEqual(t, len(state.Entries()), 0)
}

func TestNewStateSortsOwners(t *testing.T) {
	state, err := NewState(
		entry("user-cli", imagePath()),
		entry("autoscaler", replicasPath()),
	)

	requireNoError(t, err)
	requireOwners(t, state.Owners(), "autoscaler", "user-cli")
}

func TestNewStateMergesDuplicateOwners(t *testing.T) {
	state, err := NewState(
		entry("user-cli", imagePath()),
		entry("user-cli", replicasPath()),
	)

	requireNoError(t, err)
	requireEqual(t, len(state.Entries()), 1)
	requireSet(t, state.FieldsFor("user-cli"), "$.spec.image", "$.spec.replicas")
}

func TestNewStatePrunesEmptyEntries(t *testing.T) {
	state, err := NewState(emptyEntry("user-cli"))

	requireNoError(t, err)
	requireEqual(t, state.IsEmpty(), true)
}

func TestNewStateRejectsInvalidEntry(t *testing.T) {
	_, err := NewState(Entry{})

	requireErrorIs(t, err, ErrInvalidEntry)
}

func TestStateAllowsSharedPathOwnership(t *testing.T) {
	state, err := NewState(
		entry("user-cli", replicasPath()),
		entry("autoscaler", replicasPath()),
	)

	requireNoError(t, err)
	requireOwners(t, state.OwnersOf(replicasPath()), "autoscaler", "user-cli")
}

func TestMustStatePanicsOnInvalidEntry(t *testing.T) {
	requirePanic(t, func() {
		MustState(Entry{})
	})
}

func TestNewStateAllowsExplicitParentAndChildPaths(t *testing.T) {
	state, err := NewState(entry("user-cli", specPath(), replicasPath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("user-cli"), "$.spec", "$.spec.replicas")
}

func TestNewStateAllowsNoEntries(t *testing.T) {
	state, err := NewState()

	requireNoError(t, err)
	requireEqual(t, state.Fields().Equal(fieldpath.EmptySet()), true)
}
