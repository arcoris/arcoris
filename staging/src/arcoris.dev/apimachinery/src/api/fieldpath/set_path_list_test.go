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

package fieldpath

import "testing"

func TestClonePathsDetachesPathElements(t *testing.T) {
	paths := []Path{setReplicasPath()}
	cloned := clonePaths(paths)

	cloned[0].elements[0] = testFieldElement("status")

	requireEqual(t, paths[0].String(), "$.spec.replicas")
	requireEqual(t, cloned[0].String(), "$.status.replicas")
}

func TestSortPaths(t *testing.T) {
	paths := []Path{setReplicasPath(), setImagePath()}

	sortPaths(paths)

	requireEqual(t, paths[0].String(), "$.spec.image")
	requireEqual(t, paths[1].String(), "$.spec.replicas")
}

func TestCompactPathsRemovesExactDuplicates(t *testing.T) {
	paths := []Path{setImagePath(), setImagePath(), setReplicasPath()}

	compacted := compactPaths(paths)

	requireEqual(t, len(compacted), 2)
	requireEqual(t, compacted[0].String(), "$.spec.image")
	requireEqual(t, compacted[1].String(), "$.spec.replicas")
}
