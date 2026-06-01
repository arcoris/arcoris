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

// Has reports whether s contains path.
//
// Membership uses the set's canonical sorted order, so callers get predictable
// O(log n) lookup without exposing internal indexing structures.
func (s PathSet) Has(path Path) bool {
	_, ok := s.search(path)
	return ok
}

// search locates path in the canonical sorted slice.
//
// It returns the index where path was found or, when absent, the insertion
// point that would preserve canonical ordering.
func (s PathSet) search(path Path) (int, bool) {
	return slices.BinarySearchFunc(s.paths, path, func(candidate Path, target Path) int {
		return candidate.Compare(target)
	})
}
