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
//
// Invalid paths are treated as no match because this query API intentionally has
// no error return. Error-returning operations, such as Conflicts, validate
// caller-supplied paths.
func (s State) OwnersOf(path fieldpath.Path) []Owner {
	if path.Validate() != nil {
		return nil
	}

	owners := make([]Owner, 0)
	for _, entry := range s.entries {
		if entry.fields.Has(path) {
			owners = append(owners, entry.owner)
		}
	}

	return owners
}

// OverlappingOwners returns owned paths that structurally overlap path.
//
// Overlap is inclusive: exact matches, ancestors, and descendants all match.
// Invalid paths are treated as no match because this query API has no error
// return.
func (s State) OverlappingOwners(path fieldpath.Path) []Ownership {
	if path.Validate() != nil {
		return nil
	}

	ownerships := make([]Ownership, 0)
	for _, entry := range s.entries {
		for _, owned := range entry.fields.Paths() {
			if pathsOverlap(owned, path) {
				ownerships = append(ownerships, Ownership{Owner: entry.owner, Path: owned})
			}
		}
	}

	return ownerships
}
