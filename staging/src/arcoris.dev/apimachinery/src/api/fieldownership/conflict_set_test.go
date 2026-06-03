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
	requireEqual(t, ConflictSet{{Owner: "a"}}.IsEmpty(), false)
}

func TestConflictSetLen(t *testing.T) {
	requireEqual(t, ConflictSet{{Owner: "a"}, {Owner: "b"}}.Len(), 2)
}

func TestConflictSetOwnersSortedUnique(t *testing.T) {
	conflicts := ConflictSet{
		{Owner: "user-cli", OwnedPath: imagePath(), AttemptedPath: imagePath()},
		{Owner: "autoscaler", OwnedPath: replicasPath(), AttemptedPath: replicasPath()},
		{Owner: "autoscaler", OwnedPath: specPath(), AttemptedPath: replicasPath()},
		{Owner: "user-cli", OwnedPath: specPath(), AttemptedPath: imagePath()},
	}

	requireOwners(t, conflicts.Owners(), "autoscaler", "user-cli")
}

func TestConflictSetOwnedPaths(t *testing.T) {
	conflicts := ConflictSet{
		{Owner: "autoscaler", OwnedPath: replicasPath(), AttemptedPath: specPath()},
		{Owner: "user-cli", OwnedPath: imagePath(), AttemptedPath: specPath()},
	}

	requireSet(t, conflicts.OwnedPaths(), "$.spec.image", "$.spec.replicas")
}

func TestConflictSetAttemptedPaths(t *testing.T) {
	conflicts := ConflictSet{
		{Owner: "autoscaler", OwnedPath: replicasPath(), AttemptedPath: specPath()},
		{Owner: "user-cli", OwnedPath: imagePath(), AttemptedPath: specPath()},
	}

	requireSet(t, conflicts.AttemptedPaths(), "$.spec")
}

func TestConflictSetErrorDeterministic(t *testing.T) {
	conflicts := ConflictSet{
		{Owner: "user-cli", OwnedPath: imagePath(), AttemptedPath: specPath()},
		{Owner: "autoscaler", OwnedPath: replicasPath(), AttemptedPath: specPath()},
	}
	want := "field ownership conflicts: " +
		"autoscaler owns $.spec.replicas, attempted $.spec; " +
		"user-cli owns $.spec.image, attempted $.spec"

	requireEqual(t, conflicts.Error(), want)
}
