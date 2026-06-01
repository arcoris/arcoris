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

// HasDescendant reports whether s contains a strict descendant of prefix.
//
// The prefix itself does not count. This distinction matters for higher layers
// that need to know whether a subtree contains deeper owned or changed fields.
func (s PathSet) HasDescendant(prefix Path) bool {
	for _, path := range s.paths {
		if path.IsDescendantOf(prefix) {
			return true
		}
	}

	return false
}

// IntersectsSubtree reports whether any stored path shares subtree ancestry
// with path.
//
// The operation is symmetric at the path level: a set intersects a subtree when
// it contains the exact path, an ancestor of the path, or a descendant below
// the path.
func (s PathSet) IntersectsSubtree(path Path) bool {
	for _, candidate := range s.paths {
		if candidate.IntersectsSubtree(path) {
			return true
		}
	}

	return false
}

// RemoveDescendants returns a detached set with strict descendants of prefix
// removed.
//
// The prefix path itself is preserved when present. Only deeper children are
// removed. This is useful when a higher layer wants one coarse subtree marker
// to replace more specific child ownership or diff entries.
func (s PathSet) RemoveDescendants(prefix Path) PathSet {
	if len(s.paths) == 0 {
		return PathSet{}
	}

	filtered := make([]Path, 0, len(s.paths))

	for _, path := range s.paths {
		if path.IsDescendantOf(prefix) {
			continue
		}

		filtered = append(filtered, Path{elements: cloneElements(path.elements)})
	}

	return PathSet{paths: filtered}
}
