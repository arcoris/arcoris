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

// clonePaths returns detached paths in source order.
func clonePaths(paths []Path) []Path {
	if paths == nil {
		return nil
	}

	cloned := make([]Path, 0, len(paths))
	return appendSetPaths(cloned, paths)
}

// sortPaths orders paths with Path.Compare.
func sortPaths(paths []Path) {
	slices.SortFunc(paths, func(left Path, right Path) int {
		return left.Compare(right)
	})
}

// compactPaths removes exact duplicate paths from an already sorted path slice.
func compactPaths(paths []Path) []Path {
	if len(paths) < 2 {
		return paths
	}

	write := 1
	for read := 1; read < len(paths); read++ {
		if paths[read].Equal(paths[write-1]) {
			continue
		}

		paths[write] = paths[read]
		write++
	}

	return paths[:write]
}
