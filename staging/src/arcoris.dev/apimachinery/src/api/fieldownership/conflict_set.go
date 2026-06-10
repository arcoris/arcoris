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

import (
	"fmt"
	"slices"
	"strings"

	"arcoris.dev/apimachinery/api/fieldpath"
)

// ConflictSet is a deterministic immutable-by-convention collection of conflicts.
type ConflictSet struct {
	conflicts []Conflict
}

// NewConflictSet constructs a deterministic conflict set.
func NewConflictSet(conflicts ...Conflict) ConflictSet {
	if len(conflicts) == 0 {
		return ConflictSet{}
	}

	ordered := make([]Conflict, len(conflicts))
	copy(ordered, conflicts)
	slices.SortFunc(ordered, compareConflicts)

	return ConflictSet{conflicts: compactConflicts(ordered)}
}

// IsEmpty reports whether c contains no conflicts.
func (c ConflictSet) IsEmpty() bool {
	return len(c.conflicts) == 0
}

// Len returns the number of conflicts.
func (c ConflictSet) Len() int {
	return len(c.conflicts)
}

// Conflicts returns detached conflicts in deterministic order.
func (c ConflictSet) Conflicts() []Conflict {
	if len(c.conflicts) == 0 {
		return nil
	}

	conflicts := make([]Conflict, len(c.conflicts))
	copy(conflicts, c.conflicts)

	return conflicts
}

// ForEach visits conflicts in deterministic order until fn returns false.
func (c ConflictSet) ForEach(fn func(index int, conflict Conflict) bool) {
	for i, conflict := range c.conflicts {
		if !fn(i, conflict) {
			return
		}
	}
}

// Owners returns sorted unique conflicting owners.
func (c ConflictSet) Owners() []Owner {
	if len(c.conflicts) == 0 {
		return nil
	}

	owners := make([]Owner, 0, len(c.conflicts))
	for _, conflict := range c.conflicts {
		owners = append(owners, conflict.Owner)
	}

	sortOwners(owners)
	return compactSortedOwners(owners)
}

// OwnedPaths returns the set of paths already owned by conflicting owners.
func (c ConflictSet) OwnedPaths() fieldpath.Set {
	fields := fieldpath.EmptySet()
	for _, conflict := range c.conflicts {
		fields = fields.Insert(conflict.OwnedPath)
	}

	return fields
}

// AttemptedPaths returns the set of attempted paths involved in conflicts.
func (c ConflictSet) AttemptedPaths() fieldpath.Set {
	fields := fieldpath.EmptySet()
	for _, conflict := range c.conflicts {
		fields = fields.Insert(conflict.AttemptedPath)
	}

	return fields
}

// Error returns deterministic compact conflict text.
func (c ConflictSet) Error() string {
	if len(c.conflicts) == 0 {
		return "field ownership conflicts: none"
	}

	parts := make([]string, 0, len(c.conflicts))
	for _, conflict := range c.conflicts {
		parts = append(parts, fmt.Sprintf(
			"%s owns %s, attempted %s",
			conflict.Owner,
			conflict.OwnedPath,
			conflict.AttemptedPath,
		))
	}

	return "field ownership conflicts: " + strings.Join(parts, "; ")
}

// compareConflicts orders conflicts by owner, attempted path, then owned path.
func compareConflicts(left Conflict, right Conflict) int {
	if cmp := compareOwners(left.Owner, right.Owner); cmp != 0 {
		return cmp
	}
	if cmp := left.AttemptedPath.Compare(right.AttemptedPath); cmp != 0 {
		return cmp
	}

	return left.OwnedPath.Compare(right.OwnedPath)
}

// compactConflicts removes exact duplicate conflicts from a sorted slice.
func compactConflicts(conflicts []Conflict) []Conflict {
	if len(conflicts) == 0 {
		return nil
	}

	out := conflicts[:1]
	for _, conflict := range conflicts[1:] {
		last := out[len(out)-1]
		if conflict.Owner == last.Owner &&
			conflict.AttemptedPath.Equal(last.AttemptedPath) &&
			conflict.OwnedPath.Equal(last.OwnedPath) {
			continue
		}
		out = append(out, conflict)
	}

	return out
}
