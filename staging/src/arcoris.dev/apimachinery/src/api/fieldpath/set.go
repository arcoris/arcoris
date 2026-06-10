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

// Set is a deterministic immutable-by-convention collection of semantic Paths.
//
// A Set stores paths in Path.Compare order and removes duplicates. It does not
// model ownership, apply semantics, managed fields, descriptor traversal, or
// value validation. Higher layers use Set as a stable collection primitive.
type Set struct {
	paths []Path
}

// Len returns the number of paths in s.
func (s Set) Len() int {
	return len(s.paths)
}

// IsZero reports whether s is the zero set value.
func (s Set) IsZero() bool {
	return s.paths == nil
}

// IsEmpty reports whether s contains no paths.
func (s Set) IsEmpty() bool {
	return len(s.paths) == 0
}

// Paths returns detached paths in canonical order.
//
// Callers may reorder or replace the returned slice without mutating the set.
// Each Path is cloned so element slices also remain detached.
func (s Set) Paths() []Path {
	return clonePaths(s.paths)
}

// ForEach visits canonical paths in order until fn returns false.
func (s Set) ForEach(fn func(index int, path Path) bool) {
	for i, path := range s.paths {
		if !fn(i, path) {
			return
		}
	}
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
