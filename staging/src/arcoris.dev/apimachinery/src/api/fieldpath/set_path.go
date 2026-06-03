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

// cloneSetPath returns a detached path copy suitable for storing inside Set.
func cloneSetPath(path Path) Path {
	return Path{elements: cloneElements(path.elements)}
}

// appendSetPath appends one detached path copy to dst.
func appendSetPath(dst []Path, path Path) []Path {
	return append(dst, cloneSetPath(path))
}

// appendSetPaths appends detached path copies to dst in source order.
func appendSetPaths(dst []Path, paths []Path) []Path {
	for _, p := range paths {
		dst = appendSetPath(dst, p)
	}

	return dst
}

// compareSetPathSlices compares canonical path slices lexicographically.
func compareSetPathSlices(left []Path, right []Path) int {
	for i := 0; i < len(left) && i < len(right); i++ {
		if cmp := left[i].Compare(right[i]); cmp != 0 {
			return cmp
		}
	}

	switch {
	case len(left) < len(right):
		return -1
	case len(left) > len(right):
		return 1
	default:
		return 0
	}
}
