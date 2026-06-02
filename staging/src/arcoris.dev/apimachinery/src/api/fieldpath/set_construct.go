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

// EmptySet returns an empty canonical path set.
func EmptySet() Set {
	return Set{}
}

// NewSet constructs a sorted duplicate-free path set.
//
// Every input path is validated before the set is normalized. The returned set
// owns its path slice and does not retain caller-owned slice aliases.
func NewSet(paths ...Path) (Set, error) {
	canonical, err := normalizeSetPaths(paths)
	if err != nil {
		return Set{}, err
	}

	return Set{paths: canonical}, nil
}

// MustSet constructs a path set or panics when any path is invalid.
//
// It is intended for tests, fixtures, and package-level declarations where
// malformed semantic paths should fail immediately.
func MustSet(paths ...Path) Set {
	set, err := NewSet(paths...)
	if err != nil {
		panic(err)
	}

	return set
}

// normalizeSetPaths validates, sorts, and deduplicates caller-provided paths.
func normalizeSetPaths(paths []Path) ([]Path, error) {
	canonical := clonePaths(paths)

	for _, path := range canonical {
		if err := path.Validate(); err != nil {
			return nil, err
		}
	}

	sortPaths(canonical)
	return compactPaths(canonical), nil
}
