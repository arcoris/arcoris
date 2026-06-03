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

func TestStateConflictsEmptyAttemptedIsEmpty(t *testing.T) {
	conflicts, err := baseState().Conflicts("user-cli", set())

	requireNoError(t, err)
	requireEqual(t, conflicts.IsEmpty(), true)
}

func TestStateConflictsUnownedFieldsIsEmpty(t *testing.T) {
	conflicts, err := baseState().Conflicts("user-cli", set(namePath()))

	requireNoError(t, err)
	requireEqual(t, conflicts.IsEmpty(), true)
}

func TestStateConflictsIgnoresSameOwner(t *testing.T) {
	state := MustState(entry("user-cli", specPath()))

	conflicts, err := state.Conflicts("user-cli", set(replicasPath()))

	requireNoError(t, err)
	requireEqual(t, conflicts.IsEmpty(), true)
}

func TestStateConflictsExactOwnedPath(t *testing.T) {
	conflicts, err := baseState().Conflicts("other", set(replicasPath()))

	requireNoError(t, err)
	requireConflictStrings(t, conflicts,
		"autoscaler:$.spec.replicas->$.spec.replicas",
		"user-cli:$.spec.replicas->$.spec.replicas",
	)
}

func TestStateConflictsOwnedAncestor(t *testing.T) {
	state := MustState(entry("template-controller", specPath()))

	conflicts, err := state.Conflicts("user-cli", set(replicasPath()))

	requireNoError(t, err)
	requireConflictStrings(t, conflicts, "template-controller:$.spec->$.spec.replicas")
}

func TestStateConflictsOwnedDescendant(t *testing.T) {
	state := MustState(entry("label-controller", labelPath()))

	conflicts, err := state.Conflicts("user-cli", set(metadataPath()))

	requireNoError(t, err)
	requireConflictStrings(t, conflicts, `label-controller:$.metadata.labels["app"]->$.metadata`)
}

func TestStateConflictsMultipleOwners(t *testing.T) {
	state := MustState(entry("a", specPath()), entry("b", replicasPath()))

	conflicts, err := state.Conflicts("c", set(replicasPath()))

	requireNoError(t, err)
	requireConflictStrings(t, conflicts,
		"a:$.spec->$.spec.replicas",
		"b:$.spec.replicas->$.spec.replicas",
	)
}

func TestStateConflictsSharedPathOwnership(t *testing.T) {
	state := MustState(entry("a", replicasPath()), entry("b", replicasPath()))

	conflicts, err := state.Conflicts("c", set(replicasPath()))

	requireNoError(t, err)
	requireConflictStrings(t, conflicts,
		"a:$.spec.replicas->$.spec.replicas",
		"b:$.spec.replicas->$.spec.replicas",
	)
}

func TestStateConflictsDeterministicOrder(t *testing.T) {
	state := MustState(entry("b", specPath()), entry("a", replicasPath()))

	conflicts, err := state.Conflicts("c", set(specPath(), replicasPath()))

	requireNoError(t, err)
	requireConflictStrings(t, conflicts,
		"a:$.spec.replicas->$.spec",
		"a:$.spec.replicas->$.spec.replicas",
		"b:$.spec->$.spec",
		"b:$.spec->$.spec.replicas",
	)
}

func TestStateConflictsRejectsInvalidOwner(t *testing.T) {
	_, err := baseState().Conflicts("", set(replicasPath()))

	requireErrorIs(t, err, ErrInvalidOwner)
}

func requireConflictStrings(t *testing.T, conflicts ConflictSet, want ...string) {
	t.Helper()

	got := make([]string, 0, len(conflicts))
	for _, conflict := range conflicts {
		got = append(
			got,
			conflict.Owner.String()+":"+
				conflict.OwnedPath.String()+"->"+
				conflict.AttemptedPath.String(),
		)
	}

	if len(got) == 0 && len(want) == 0 {
		return
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("conflicts = %#v; want %#v", got, want)
	}
}
