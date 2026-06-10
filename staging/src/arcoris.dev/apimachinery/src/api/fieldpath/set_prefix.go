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

// HasAnyUnder reports whether s contains prefix or any path below prefix.
//
// The prefix itself counts. This inclusive behavior is useful for subtree
// ownership and change filtering, where an exact subtree marker covers the
// subtree root as well as descendants.
func (s Set) HasAnyUnder(prefix Path) bool {
	start, _ := findSetPath(s.paths, prefix)
	return start < len(s.paths) && s.paths[start].HasPrefix(prefix)
}

// Under returns paths equal to prefix or structurally below prefix.
func (s Set) Under(prefix Path) Set {
	if len(s.paths) == 0 {
		return Set{}
	}

	start, _ := findSetPath(s.paths, prefix)
	if start >= len(s.paths) {
		return Set{}
	}

	paths := make([]Path, 0, len(s.paths)-start)

	for i := start; i < len(s.paths); i++ {
		p := s.paths[i]
		if !p.HasPrefix(prefix) {
			break
		}

		paths = appendSetPath(paths, p)
	}

	return Set{paths: paths}
}

// HasDescendant reports whether s contains a strict descendant of prefix.
func (s Set) HasDescendant(prefix Path) bool {
	for _, path := range s.paths {
		if path.IsDescendantOf(prefix) {
			return true
		}
	}

	return false
}

// RemoveDescendants returns a detached set without strict descendants of prefix.
//
// The prefix path itself is preserved when present.
func (s Set) RemoveDescendants(prefix Path) Set {
	if len(s.paths) == 0 {
		return Set{}
	}

	paths := make([]Path, 0, len(s.paths))
	for _, path := range s.paths {
		if path.IsDescendantOf(prefix) {
			continue
		}

		paths = appendSetPath(paths, path)
	}

	return Set{paths: paths}
}

// Overlaps reports whether any stored path overlaps path.
func (s Set) Overlaps(path Path) bool {
	for _, candidate := range s.paths {
		if candidate.Overlaps(path) {
			return true
		}
	}

	return false
}

// CompactSubtrees removes descendants whenever an ancestor is already present.
//
// Compaction is explicit so callers can choose between exact path-set algebra
// and subtree-marker semantics.
func (s Set) CompactSubtrees() Set {
	if len(s.paths) == 0 {
		return Set{}
	}

	paths := make([]Path, 0, len(s.paths))
	for _, path := range s.paths {
		hasAncestor := false
		for _, kept := range paths {
			if path.IsDescendantOf(kept) {
				hasAncestor = true
				break
			}
		}

		if !hasAncestor {
			paths = appendSetPath(paths, path)
		}
	}

	return Set{paths: paths}
}
