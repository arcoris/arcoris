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
	requireOwners(t, baseState().OwnersOf(replicasPath()), "autoscaler", "user-cli")
}

func TestStateOwnersOfDoesNotReturnOverlappingAncestor(t *testing.T) {
	state := MustState(entry("template-controller", specPath()))

	requireOwners(t, state.OwnersOf(replicasPath()))
}

func TestStateOwnersOfUnknownPathIsEmpty(t *testing.T) {
	requireOwners(t, baseState().OwnersOf(metadataPath()))
}

func TestStateOverlappingOwnersExact(t *testing.T) {
	state := MustState(entry("autoscaler", replicasPath()))

	requireOwnerships(t, state.OverlappingOwners(replicasPath()), "autoscaler:$.spec.replicas")
}

func TestStateOverlappingOwnersAncestor(t *testing.T) {
	state := MustState(entry("template-controller", specPath()))

	requireOwnerships(t, state.OverlappingOwners(replicasPath()), "template-controller:$.spec")
}

func TestStateOverlappingOwnersDescendant(t *testing.T) {
	state := MustState(entry("label-controller", labelPath()))

	requireOwnerships(t,
		state.OverlappingOwners(metadataPath()),
		`label-controller:$.metadata.labels["app"]`,
	)
}

func TestStateOverlappingOwnersNone(t *testing.T) {
	requireOwnerships(t, baseState().OverlappingOwners(namePath()))
}

func TestStateOverlappingOwnersDeterministicOrder(t *testing.T) {
	state := MustState(
		entry("user-cli", specPath()),
		entry("autoscaler", replicasPath()),
	)

	requireOwnerships(t,
		state.OverlappingOwners(replicasPath()),
		"autoscaler:$.spec.replicas",
		"user-cli:$.spec",
	)
}

func requireOwnerships(t *testing.T, got []Ownership, want ...string) {
	t.Helper()

	strings := make([]string, 0, len(got))
	for _, ownership := range got {
		strings = append(strings, ownership.Owner.String()+":"+ownership.Path.String())
	}

	if len(strings) == 0 && len(want) == 0 {
		return
	}

	if !reflect.DeepEqual(strings, want) {
		t.Fatalf("ownerships = %#v; want %#v", strings, want)
	}
}
