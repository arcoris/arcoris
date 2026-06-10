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

// Conflicts returns ownership overlaps for owner attempting fields.
//
// The same owner never conflicts with itself. Empty attempted field sets return
// an empty ConflictSet. This method does not decide force behavior, request
// admission, or which fields should be attempted; callers provide the policy set.
func (s State) Conflicts(owner Owner, attempted fieldpath.Set) (ConflictSet, error) {
	if err := validateConflictInputs(owner, attempted); err != nil {
		return ConflictSet{}, err
	}
	if attempted.IsEmpty() || s.IsEmpty() {
		return ConflictSet{}, nil
	}

	conflicts := make([]Conflict, 0)
	attempted.ForEach(func(_ int, attemptedPath fieldpath.Path) bool {
		conflicts = append(conflicts, s.conflictsForPath(owner, attemptedPath)...)
		return true
	})

	return newConflictSetUnchecked(conflicts...), nil
}

// conflictsForPath collects overlaps for one already-validated attempted path.
func (s State) conflictsForPath(owner Owner, attemptedPath fieldpath.Path) []Conflict {
	conflicts := make([]Conflict, 0)
	for _, entry := range s.entries {
		if entry.owner == owner {
			continue
		}

		entry.fields.ForEach(func(_ int, ownedPath fieldpath.Path) bool {
			if ownedPath.Overlaps(attemptedPath) {
				conflicts = append(conflicts, Conflict{
					Owner:         entry.owner,
					OwnedPath:     ownedPath,
					AttemptedPath: attemptedPath,
				})
			}
			return true
		})
	}

	return conflicts
}
