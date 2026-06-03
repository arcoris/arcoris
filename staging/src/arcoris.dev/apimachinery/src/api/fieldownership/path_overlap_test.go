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

func TestPathsOverlapExact(t *testing.T) {
	requireEqual(t, pathsOverlap(specPath(), specPath()), true)
}

func TestPathsOverlapAncestor(t *testing.T) {
	requireEqual(t, pathsOverlap(specPath(), replicasPath()), true)
}

func TestPathsOverlapDescendant(t *testing.T) {
	requireEqual(t, pathsOverlap(replicasPath(), specPath()), true)
}

func TestPathsOverlapSibling(t *testing.T) {
	requireEqual(t, pathsOverlap(imagePath(), replicasPath()), false)
}

func TestPathsOverlapMapKey(t *testing.T) {
	requireEqual(t, pathsOverlap(path("$.metadata.labels"), labelPath()), true)
}

func TestPathsOverlapListIndex(t *testing.T) {
	requireEqual(t, pathsOverlap(argsPath(), argsIndexPath()), true)
}

func TestPathsOverlapListMapSelector(t *testing.T) {
	requireEqual(t, pathsOverlap(readyPath(), readyStatusPath()), true)
}
