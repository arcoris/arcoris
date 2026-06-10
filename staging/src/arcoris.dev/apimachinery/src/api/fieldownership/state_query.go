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

import "arcoris.dev/apimachinery/api/fieldpath"

// OwnersOf returns owners that own path exactly.
func (s State) OwnersOf(path fieldpath.Path) ([]Owner, error) {
	if err := validateOwnedPathQuery(path, "owned path is invalid"); err != nil {
		return nil, err
	}

	owners := make([]Owner, 0)
	for _, entry := range s.entries {
		if entry.fields.Has(path) {
			owners = append(owners, entry.owner)
		}
	}

	return owners, nil
}

// OverlappingPaths returns owned paths that structurally overlap path.
//
// Overlap is inclusive: exact matches, ancestors, and descendants all match.
func (s State) OverlappingPaths(path fieldpath.Path) (OwnedPathSet, error) {
	if err := validateOwnedPathQuery(path, "overlap query path is invalid"); err != nil {
		return OwnedPathSet{}, err
	}

	ownedPaths := make([]OwnedPath, 0)
	for _, entry := range s.entries {
		for _, owned := range entry.fields.Paths() {
			if owned.Overlaps(path) {
				ownedPaths = append(ownedPaths, OwnedPath{Owner: entry.owner, Path: owned})
			}
		}
	}

	return NewOwnedPathSet(ownedPaths...), nil
}
