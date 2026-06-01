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

// NewPathSet constructs a canonical sorted unique path set.
//
// Input order is accepted for ergonomics, but the stored representation is
// normalized into deterministic sort order with duplicates removed. Every input
// path is validated before canonicalization succeeds.
func NewPathSet(paths ...Path) (PathSet, error) {
	canonical := clonePaths(paths)

	for _, path := range canonical {
		if err := path.Validate(); err != nil {
			return PathSet{}, err
		}
	}

	sortPaths(canonical)
	canonical = compactPaths(canonical)

	return PathSet{paths: canonical}, nil
}

// MustPathSet constructs a canonical path set or panics on malformed input.
//
// It is intended for tests, fixtures, and package-level declarations where
// malformed paths should fail immediately.
func MustPathSet(paths ...Path) PathSet {
	set, err := NewPathSet(paths...)
	if err != nil {
		panic(err)
	}

	return set
}

// sortPaths orders paths in-place using their deterministic Compare contract.
func sortPaths(paths []Path) {
	slices.SortFunc(paths, func(left Path, right Path) int {
		return left.Compare(right)
	})
}

// compactPaths removes adjacent duplicates from an already sorted path slice.
func compactPaths(paths []Path) []Path {
	if len(paths) < 2 {
		return paths
	}

	write := 1

	for read := 1; read < len(paths); read++ {
		if paths[write-1].Equal(paths[read]) {
			continue
		}

		paths[write] = paths[read]
		write++
	}

	for i := write; i < len(paths); i++ {
		paths[i] = Path{}
	}

	return paths[:write:write]
}
