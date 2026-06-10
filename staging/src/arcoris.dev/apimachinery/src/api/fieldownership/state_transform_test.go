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

func TestStateSetFieldsAddsOwner(t *testing.T) {
	state, err := EmptyState().SetFields(owner("user-cli"), set(imagePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.spec.image")
}

func TestStateSetFieldsReplacesOwnerFields(t *testing.T) {
	state, err := baseState().SetFields(owner("user-cli"), set(namePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.metadata.name")
	requireSet(t, state.FieldsFor(owner("autoscaler")), "$.spec.replicas")
}

func TestStateSetFieldsEmptyRemovesOwner(t *testing.T) {
	state, err := baseState().SetFields(owner("user-cli"), fieldpath.EmptySet())

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")))
	requireOwners(t, state.Owners(), "autoscaler", "health-controller")
}

func TestStateAddFieldsAddsOwner(t *testing.T) {
	state, err := EmptyState().AddFields(owner("user-cli"), set(imagePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.spec.image")
}

func TestStateAddFieldsUnionsFields(t *testing.T) {
	state, err := MustState(entry("user-cli", imagePath())).AddFields(owner("user-cli"), set(replicasPath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.spec.image", "$.spec.replicas")
}

func TestStateAddFieldsEmptyNoop(t *testing.T) {
	state := baseState()
	got, err := state.AddFields(owner("user-cli"), fieldpath.EmptySet())

	requireNoError(t, err)
	requireSet(t, got.FieldsFor(owner("user-cli")), "$.spec.image", "$.spec.replicas")
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.spec.image", "$.spec.replicas")
}

func TestStateRemoveFieldsExact(t *testing.T) {
	state, err := baseState().RemoveFields(owner("user-cli"), set(imagePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.spec.replicas")
}

func TestStateRemoveFieldsDoesNotRemoveOverlappingPaths(t *testing.T) {
	state := MustState(entry("user-cli", specPath(), replicasPath()))
	got, err := state.RemoveFields(owner("user-cli"), set(specPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor(owner("user-cli")), "$.spec.replicas")
}

func TestStateRemoveOverlappingFieldsRemovesAncestor(t *testing.T) {
	state := MustState(entry("user-cli", specPath()))
	got, err := state.RemoveOverlappingFields(owner("user-cli"), set(replicasPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor(owner("user-cli")))
}

func TestStateRemoveOverlappingFieldsRemovesDescendant(t *testing.T) {
	state := MustState(entry("user-cli", replicasPath()))
	got, err := state.RemoveOverlappingFields(owner("user-cli"), set(specPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor(owner("user-cli")))
}

func TestStateRemoveFieldsFromOthers(t *testing.T) {
	state, err := baseState().RemoveFieldsFromOthers(owner("user-cli"), set(replicasPath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("autoscaler")))
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.spec.image", "$.spec.replicas")
}

func TestStateRemoveFieldsFromOthersKeepsCurrentOwner(t *testing.T) {
	state := MustState(entry("user-cli", replicasPath()), entry("autoscaler", replicasPath()))
	got, err := state.RemoveFieldsFromOthers(owner("user-cli"), set(replicasPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor(owner("user-cli")), "$.spec.replicas")
	requireSet(t, got.FieldsFor(owner("autoscaler")))
}

func TestStateRemoveOverlappingFieldsFromOthers(t *testing.T) {
	state := MustState(entry("user-cli", specPath()), entry("autoscaler", replicasPath()))
	got, err := state.RemoveOverlappingFieldsFromOthers(owner("user-cli"), set(specPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor(owner("user-cli")), "$.spec")
	requireSet(t, got.FieldsFor(owner("autoscaler")))
}

func TestStateRemoveOverlappingFieldsFromOthersKeepsCurrentOwner(t *testing.T) {
	state := MustState(entry("user-cli", replicasPath()), entry("autoscaler", specPath()))
	got, err := state.RemoveOverlappingFieldsFromOthers(owner("user-cli"), set(specPath()))

	requireNoError(t, err)
	requireSet(t, got.FieldsFor(owner("user-cli")), "$.spec.replicas")
	requireSet(t, got.FieldsFor(owner("autoscaler")))
}

func TestStateRemoveOwner(t *testing.T) {
	state, err := baseState().RemoveOwner(owner("autoscaler"))

	requireNoError(t, err)
	requireOwners(t, state.Owners(), "health-controller", "user-cli")
	requireSet(t, state.FieldsFor(owner("autoscaler")))
}

func TestStateTransformationsDoNotMutateReceiver(t *testing.T) {
	state := baseState()

	_, err := state.SetFields(owner("user-cli"), set(namePath()))

	requireNoError(t, err)
	requireSet(t, state.FieldsFor(owner("user-cli")), "$.spec.image", "$.spec.replicas")
	requireSet(t, state.FieldsFor(owner("autoscaler")), "$.spec.replicas")
}

func TestStateTransformationsKeepDeterministicOrdering(t *testing.T) {
	state := MustState(entry("z", imagePath()))
	got, err := state.AddFields(owner("a"), set(replicasPath()))

	requireNoError(t, err)
	requireOwners(t, got.Owners(), "a", "z")
}

func TestStateTransformationsRejectInvalidOwner(t *testing.T) {
	_, err := baseState().SetFields(Owner{}, set(imagePath()))

	requireErrorIs(t, err, ErrInvalidOwner)
}
