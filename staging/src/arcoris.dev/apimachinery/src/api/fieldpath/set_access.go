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

import "slices"

// Paths returns detached paths in canonical order.
//
// Callers may reorder or replace the returned slice without mutating the set.
// Each Path is cloned so element slices also remain detached.
func (s Set) Paths() []Path {
	return clonePaths(s.paths)
}

// Has reports whether s contains path exactly.
//
// Membership uses binary search over canonical Path.Compare ordering.
func (s Set) Has(path Path) bool {
	_, found := findSetPath(s.paths, path)
	return found
}

// findSetPath returns the index of path or the canonical insertion position.
func findSetPath(paths []Path, path Path) (int, bool) {
	return slices.BinarySearchFunc(paths, path, func(candidate Path, target Path) int {
		return candidate.Compare(target)
	})
}
