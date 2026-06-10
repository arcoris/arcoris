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

package fieldownership

import "slices"

// sortOwnedPaths orders owned paths by owner and path.
func sortOwnedPaths(paths []OwnedPath) {
	slices.SortFunc(paths, compareOwnedPaths)
}

// compareOwnedPaths orders owned paths deterministically.
func compareOwnedPaths(left OwnedPath, right OwnedPath) int {
	if cmp := compareOwners(left.Owner, right.Owner); cmp != 0 {
		return cmp
	}

	return left.Path.Compare(right.Path)
}

// compactOwnedPaths removes exact duplicate owner/path records.
func compactOwnedPaths(paths []OwnedPath) []OwnedPath {
	if len(paths) == 0 {
		return nil
	}

	out := paths[:1]
	for _, path := range paths[1:] {
		last := out[len(out)-1]
		if path.Owner == last.Owner && path.Path.Equal(last.Path) {
			continue
		}
		out = append(out, path)
	}

	return out
}
