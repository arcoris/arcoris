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

	for index := start; index < len(s.paths); index++ {
		path := s.paths[index]
		if !path.HasPrefix(prefix) {
			break
		}

		return true
	}

	return false
}

// Under returns paths equal to prefix or structurally below prefix.
func (s Set) Under(prefix Path) Set {
	if len(s.paths) == 0 {
		return Set{}
	}

	start, _ := findSetPath(s.paths, prefix)
	paths := make([]Path, 0, len(s.paths)-start)

	for index := start; index < len(s.paths); index++ {
		path := s.paths[index]
		if !path.HasPrefix(prefix) {
			break
		}

		paths = appendSetPath(paths, path)
	}

	return Set{paths: paths}
}
