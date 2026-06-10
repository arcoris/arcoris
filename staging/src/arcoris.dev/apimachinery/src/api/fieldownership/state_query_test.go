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
	"reflect"
	"testing"
)

func TestStateOwnersOfExact(t *testing.T) {
	owners, err := baseState().OwnersOf(replicasPath())

	requireNoError(t, err)
	requireOwners(t, owners, "autoscaler", "user-cli")
}

func TestStateOwnersOfDoesNotReturnOverlappingAncestor(t *testing.T) {
	state := MustState(entry("template-controller", specPath()))
	owners, err := state.OwnersOf(replicasPath())

	requireNoError(t, err)
	requireOwners(t, owners)
}

func TestStateOwnersOfUnknownPathIsEmpty(t *testing.T) {
	owners, err := baseState().OwnersOf(metadataPath())

	requireNoError(t, err)
	requireOwners(t, owners)
}

func TestStateOverlappingPathsExact(t *testing.T) {
	state := MustState(entry("autoscaler", replicasPath()))
	paths, err := state.OverlappingPaths(replicasPath())

	requireNoError(t, err)
	requireOwnedPaths(t, paths, "autoscaler:$.spec.replicas")
}

func TestStateOverlappingPathsAncestor(t *testing.T) {
	state := MustState(entry("template-controller", specPath()))
	paths, err := state.OverlappingPaths(replicasPath())

	requireNoError(t, err)
	requireOwnedPaths(t, paths, "template-controller:$.spec")
}

func TestStateOverlappingPathsDescendant(t *testing.T) {
	state := MustState(entry("label-controller", labelPath()))
	paths, err := state.OverlappingPaths(metadataPath())

	requireNoError(t, err)
	requireOwnedPaths(t, paths, `label-controller:$.metadata.labels["app"]`)
}

func TestStateOverlappingPathsNone(t *testing.T) {
	paths, err := baseState().OverlappingPaths(namePath())

	requireNoError(t, err)
	requireOwnedPaths(t, paths)
}

func TestStateOverlappingPathsDeterministicOrder(t *testing.T) {
	state := MustState(
		entry("user-cli", specPath()),
		entry("autoscaler", replicasPath()),
	)
	paths, err := state.OverlappingPaths(replicasPath())

	requireNoError(t, err)
	requireOwnedPaths(t,
		paths,
		"autoscaler:$.spec.replicas",
		"user-cli:$.spec",
	)
}

func requireOwnedPaths(t *testing.T, got OwnedPathSet, want ...string) {
	t.Helper()

	strings := make([]string, 0, got.Len())
	for _, path := range got.Paths() {
		strings = append(strings, path.Owner.String()+":"+path.Path.String())
	}

	if len(strings) == 0 && len(want) == 0 {
		return
	}

	if !reflect.DeepEqual(strings, want) {
		t.Fatalf("ownerships = %#v; want %#v", strings, want)
	}
}
