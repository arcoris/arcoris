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

func TestOwnedPathStoresOwnerAndPath(t *testing.T) {
	ownedPath := OwnedPath{Owner: owner("user-cli"), Path: imagePath()}

	requireEqual(t, ownedPath.Owner, owner("user-cli"))
	requireEqual(t, ownedPath.Path.String(), "$.spec.image")
}

func TestOwnedPathSetSortsDeduplicatesAndDetaches(t *testing.T) {
	ownedPaths := NewOwnedPathSet(
		OwnedPath{Owner: owner("user-cli"), Path: imagePath()},
		OwnedPath{Owner: owner("autoscaler"), Path: replicasPath()},
		OwnedPath{Owner: owner("autoscaler"), Path: replicasPath()},
	)
	paths := ownedPaths.Paths()

	paths[0] = OwnedPath{Owner: owner("other"), Path: metadataPath()}

	requireEqual(t, ownedPaths.Len(), 2)
	requireOwners(t, ownedPaths.Owners(), "autoscaler", "user-cli")
	requireSet(t, ownedPaths.FieldPaths(), "$.spec.image", "$.spec.replicas")
}

func TestOwnedPathSetForEachStopsEarly(t *testing.T) {
	ownedPaths := NewOwnedPathSet(
		OwnedPath{Owner: owner("a"), Path: imagePath()},
		OwnedPath{Owner: owner("b"), Path: replicasPath()},
	)
	var visited []Owner

	ownedPaths.ForEach(func(index int, path OwnedPath) bool {
		visited = append(visited, path.Owner)
		return false
	})

	requireOwners(t, visited, "a")
}
