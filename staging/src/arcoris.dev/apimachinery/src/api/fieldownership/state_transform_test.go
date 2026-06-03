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

func TestStateWithFieldsAddsOwner(t *testing.T) {
	state, err := EmptyState().WithFields("user-cli", set(imagePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("user-cli"), "$.spec.image")
}

func TestStateWithFieldsReplacesOwnerFields(t *testing.T) {
	state, err := baseState().WithFields("user-cli", set(namePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("user-cli"), "$.metadata.name")
	requireSet(t, state.FieldsFor("autoscaler"), "$.spec.replicas")
}

func TestStateWithFieldsEmptyRemovesOwner(t *testing.T) {
	state, err := baseState().WithFields("user-cli", fieldpath.EmptySet())

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("user-cli"))
	requireOwners(t, state.Owners(), "autoscaler", "health-controller")
}

func TestStateAddFieldsAddsOwner(t *testing.T) {
	state, err := EmptyState().AddFields("user-cli", set(imagePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("user-cli"), "$.spec.image")
}

func TestStateAddFieldsUnionsFields(t *testing.T) {
	state, err := MustState(entry("user-cli", imagePath())).AddFields("user-cli", set(replicasPath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("user-cli"), "$.spec.image", "$.spec.replicas")
}

func TestStateAddFieldsEmptyNoop(t *testing.T) {
	state := baseState()
	got, err := state.AddFields("user-cli", fieldpath.EmptySet())

	requireNoError(t, err)
	requireSet(t, got.FieldsFor("user-cli"), "$.spec.image", "$.spec.replicas")
	requireSet(t, state.FieldsFor("user-cli"), "$.spec.image", "$.spec.replicas")
}

func TestStateRemoveFieldsExact(t *testing.T) {
	state, err := baseState().RemoveFields("user-cli", set(imagePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("user-cli"), "$.spec.replicas")
}

func TestStateRemoveFieldsDoesNotRemoveOverlaps(t *testing.T) {
	state := MustState(entry("user-cli", specPath(), replicasPath()))
	got, err := state.RemoveFields("user-cli", set(specPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor("user-cli"), "$.spec.replicas")
}

func TestStateRemoveOverlapsRemovesAncestor(t *testing.T) {
	state := MustState(entry("user-cli", specPath()))
	got, err := state.RemoveOverlaps("user-cli", set(replicasPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor("user-cli"))
}

func TestStateRemoveOverlapsRemovesDescendant(t *testing.T) {
	state := MustState(entry("user-cli", replicasPath()))
	got, err := state.RemoveOverlaps("user-cli", set(specPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor("user-cli"))
}

func TestStateRemoveFieldsFromOthers(t *testing.T) {
	state, err := baseState().RemoveFieldsFromOthers("user-cli", set(replicasPath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("autoscaler"))
	requireSet(t, state.FieldsFor("user-cli"), "$.spec.image", "$.spec.replicas")
}

func TestStateRemoveFieldsFromOthersKeepsCurrentOwner(t *testing.T) {
	state := MustState(entry("user-cli", replicasPath()), entry("autoscaler", replicasPath()))
	got, err := state.RemoveFieldsFromOthers("user-cli", set(replicasPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor("user-cli"), "$.spec.replicas")
	requireSet(t, got.FieldsFor("autoscaler"))
}

func TestStateRemoveOverlapsFromOthers(t *testing.T) {
	state := MustState(entry("user-cli", specPath()), entry("autoscaler", replicasPath()))
	got, err := state.RemoveOverlapsFromOthers("user-cli", set(specPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor("user-cli"), "$.spec")
	requireSet(t, got.FieldsFor("autoscaler"))
}

func TestStateRemoveOverlapsFromOthersKeepsCurrentOwner(t *testing.T) {
	state := MustState(entry("user-cli", replicasPath()), entry("autoscaler", specPath()))
	got, err := state.RemoveOverlapsFromOthers("user-cli", set(specPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor("user-cli"), "$.spec.replicas")
	requireSet(t, got.FieldsFor("autoscaler"))
}

func TestStateWithoutOwner(t *testing.T) {
	state, err := baseState().WithoutOwner("autoscaler")

	requireNoError(t, err)
	requireOwners(t, state.Owners(), "health-controller", "user-cli")
	requireSet(t, state.FieldsFor("autoscaler"))
}

func TestStateTransformationsDoNotMutateReceiver(t *testing.T) {
	state := baseState()

	_, err := state.WithFields("user-cli", set(namePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor("user-cli"), "$.spec.image", "$.spec.replicas")
	requireSet(t, state.FieldsFor("autoscaler"), "$.spec.replicas")
}

func TestStateTransformationsKeepDeterministicOrdering(t *testing.T) {
	state := MustState(entry("z", imagePath()))
	got, err := state.AddFields("a", set(replicasPath()))

	requireNoError(t, err)
	requireOwners(t, got.Owners(), "a", "z")
}

func TestStateTransformationsRejectInvalidOwner(t *testing.T) {
	_, err := baseState().WithFields("", set(imagePath()))

	requireErrorIs(t, err, ErrInvalidOwner)
}
