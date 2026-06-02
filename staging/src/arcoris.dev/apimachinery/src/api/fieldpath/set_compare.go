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
func (s Set) Equal(other Set) bool {
	return s.Compare(other) == 0
}

// Compare imposes deterministic total ordering on sets.
//
// Ordering is lexicographic over canonical path order. This is a collection
// ordering contract only; it does not imply descriptor, diff, or ownership
// semantics.
func (s Set) Compare(other Set) int {
	return compareSetPathSlices(s.paths, other.paths)
}
