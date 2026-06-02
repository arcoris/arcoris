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

// Union returns all paths that appear in either set.
func (s Set) Union(other Set) Set {
	paths := make([]Path, 0, len(s.paths)+len(other.paths))

	left := 0
	right := 0
	for left < len(s.paths) && right < len(other.paths) {
		switch cmp := s.paths[left].Compare(other.paths[right]); {
		case cmp < 0:
			paths = appendSetPath(paths, s.paths[left])
			left++
		case cmp > 0:
			paths = appendSetPath(paths, other.paths[right])
			right++
		default:
			paths = appendSetPath(paths, s.paths[left])
			left++
			right++
		}
	}

	paths = appendSetPaths(paths, s.paths[left:])
	paths = appendSetPaths(paths, other.paths[right:])

	return Set{paths: paths}
}

// Intersection returns paths that appear in both sets.
func (s Set) Intersection(other Set) Set {
	intersectionCapacity := len(s.paths)
	if len(other.paths) < intersectionCapacity {
		intersectionCapacity = len(other.paths)
	}

	paths := make([]Path, 0, intersectionCapacity)

	left := 0
	right := 0
	for left < len(s.paths) && right < len(other.paths) {
		switch cmp := s.paths[left].Compare(other.paths[right]); {
		case cmp < 0:
			left++
		case cmp > 0:
			right++
		default:
			paths = appendSetPath(paths, s.paths[left])
			left++
			right++
		}
	}

	return Set{paths: paths}
}

// Difference returns paths that appear in s but not in other.
func (s Set) Difference(other Set) Set {
	paths := make([]Path, 0, len(s.paths))

	left := 0
	right := 0
	for left < len(s.paths) && right < len(other.paths) {
		switch cmp := s.paths[left].Compare(other.paths[right]); {
		case cmp < 0:
			paths = appendSetPath(paths, s.paths[left])
			left++
		case cmp > 0:
			right++
		default:
			left++
			right++
		}
	}

	paths = appendSetPaths(paths, s.paths[left:])

	return Set{paths: paths}
}
