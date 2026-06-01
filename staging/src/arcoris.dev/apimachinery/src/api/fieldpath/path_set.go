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

// PathSet is a canonical sorted unique collection of semantic field paths.
//
// The set is immutable by API convention. Constructors normalize order and
// uniqueness so higher layers can rely on deterministic iteration, equality,
// and subtree-style operations.
//
// PathSet does not interpret descriptors or ownership semantics on its own. It
// simply provides a compact structured collection that future diff, validation,
// and managed-field layers can build upon.
type PathSet struct {
	paths []Path
}

// Len returns the number of canonical paths stored in s.
func (s PathSet) Len() int {
	return len(s.paths)
}

// IsZero reports whether s is the uninitialized zero set.
//
// The zero set behaves like an empty canonical set.
func (s PathSet) IsZero() bool {
	return len(s.paths) == 0
}

// Paths returns a detached slice of canonical paths.
//
// Callers may reorder or replace the returned slice without mutating the set.
// Individual paths are cloned so element slices also remain detached.
func (s PathSet) Paths() []Path {
	return clonePaths(s.paths)
}

// clonePaths returns a caller-owned slice of detached paths.
func clonePaths(paths []Path) []Path {
	if paths == nil {
		return nil
	}

	cloned := make([]Path, len(paths))
	for i, path := range paths {
		cloned[i] = Path{elements: cloneElements(path.elements)}
	}

	return cloned
}
