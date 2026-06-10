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

import "testing"

func TestConflictSetIsEmpty(t *testing.T) {
	requireEqual(t, ConflictSet{}.IsEmpty(), true)
	requireEqual(t, NewConflictSet(Conflict{Owner: owner("a")}).IsEmpty(), false)
}

func TestConflictSetLen(t *testing.T) {
	requireEqual(t, NewConflictSet(
		Conflict{Owner: owner("a")},
		Conflict{Owner: owner("b")},
	).Len(), 2)
}

func TestConflictSetOwnersSortedUnique(t *testing.T) {
	conflicts := NewConflictSet(
		Conflict{Owner: owner("user-cli"), OwnedPath: imagePath(), AttemptedPath: imagePath()},
		Conflict{Owner: owner("autoscaler"), OwnedPath: replicasPath(), AttemptedPath: replicasPath()},
		Conflict{Owner: owner("autoscaler"), OwnedPath: specPath(), AttemptedPath: replicasPath()},
		Conflict{Owner: owner("user-cli"), OwnedPath: specPath(), AttemptedPath: imagePath()},
	)

	requireOwners(t, conflicts.Owners(), "autoscaler", "user-cli")
}

func TestConflictSetOwnedPaths(t *testing.T) {
	conflicts := NewConflictSet(
		Conflict{Owner: owner("autoscaler"), OwnedPath: replicasPath(), AttemptedPath: specPath()},
		Conflict{Owner: owner("user-cli"), OwnedPath: imagePath(), AttemptedPath: specPath()},
	)

	requireSet(t, conflicts.OwnedPaths(), "$.spec.image", "$.spec.replicas")
}

func TestConflictSetAttemptedPaths(t *testing.T) {
	conflicts := NewConflictSet(
		Conflict{Owner: owner("autoscaler"), OwnedPath: replicasPath(), AttemptedPath: specPath()},
		Conflict{Owner: owner("user-cli"), OwnedPath: imagePath(), AttemptedPath: specPath()},
	)

	requireSet(t, conflicts.AttemptedPaths(), "$.spec")
}

func TestConflictSetErrorDeterministic(t *testing.T) {
	conflicts := NewConflictSet(
		Conflict{Owner: owner("user-cli"), OwnedPath: imagePath(), AttemptedPath: specPath()},
		Conflict{Owner: owner("autoscaler"), OwnedPath: replicasPath(), AttemptedPath: specPath()},
	)
	want := "field ownership conflicts: " +
		"autoscaler owns $.spec.replicas, attempted $.spec; " +
		"user-cli owns $.spec.image, attempted $.spec"

	requireEqual(t, conflicts.Error(), want)
}

func TestConflictSetConflictsReturnsDetachedSlice(t *testing.T) {
	conflicts := NewConflictSet(
		Conflict{Owner: owner("autoscaler"), OwnedPath: replicasPath(), AttemptedPath: specPath()},
	)
	got := conflicts.Conflicts()

	got[0].Owner = owner("other")

	requireOwners(t, conflicts.Owners(), "autoscaler")
}

func TestConflictSetForEachStopsEarly(t *testing.T) {
	conflicts := NewConflictSet(
		Conflict{Owner: owner("a"), OwnedPath: imagePath(), AttemptedPath: specPath()},
		Conflict{Owner: owner("b"), OwnedPath: replicasPath(), AttemptedPath: specPath()},
	)
	var visited []Owner

	conflicts.ForEach(func(index int, conflict Conflict) bool {
		visited = append(visited, conflict.Owner)
		return false
	})

	requireOwners(t, visited, "a")
}
