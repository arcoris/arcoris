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

// Equal reports whether s and other contain the same canonical paths.
//
// Because PathSet constructors canonicalize sort order and uniqueness, equality
// is purely structural and independent from caller input order.
func (s PathSet) Equal(other PathSet) bool {
	return s.Compare(other) == 0
}

// Compare imposes deterministic total ordering on canonical path sets.
//
// Ordering is lexicographic over canonical path order. This is a deterministic
// collection-ordering contract only; it does not imply descriptor or ownership
// semantics.
func (s PathSet) Compare(other PathSet) int {
	for i := 0; i < len(s.paths) && i < len(other.paths); i++ {
		if cmp := s.paths[i].Compare(other.paths[i]); cmp != 0 {
			return cmp
		}
	}

	switch {
	case len(s.paths) < len(other.paths):
		return -1
	case len(s.paths) > len(other.paths):
		return 1
	default:
		return 0
	}
}
