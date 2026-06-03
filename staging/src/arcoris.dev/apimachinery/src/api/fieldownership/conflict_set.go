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

// ConflictSet is a deterministic collection of ownership conflicts.
type ConflictSet []Conflict

// IsEmpty reports whether c contains no conflicts.
func (c ConflictSet) IsEmpty() bool {
	return len(c) == 0
}

// Len returns the number of conflicts.
func (c ConflictSet) Len() int {
	return len(c)
}

// Owners returns sorted unique conflicting owners.
func (c ConflictSet) Owners() []Owner {
	if len(c) == 0 {
		return nil
	}

	owners := make([]Owner, 0, len(c))
	for _, conflict := range c {
		owners = append(owners, conflict.Owner)
	}

	slices.SortFunc(owners, compareOwners)
	return compactSortedOwners(owners)
}

// OwnedPaths returns the set of paths already owned by conflicting owners.
func (c ConflictSet) OwnedPaths() fieldpath.Set {
	fields := fieldpath.EmptySet()
	for _, conflict := range c {
		fields = fields.Insert(conflict.OwnedPath)
	}

	return fields
}

// AttemptedPaths returns the set of attempted paths involved in conflicts.
func (c ConflictSet) AttemptedPaths() fieldpath.Set {
	fields := fieldpath.EmptySet()
	for _, conflict := range c {
		fields = fields.Insert(conflict.AttemptedPath)
	}

	return fields
}

// Error returns deterministic compact conflict text.
func (c ConflictSet) Error() string {
	if len(c) == 0 {
		return "field ownership conflicts: none"
	}

	conflicts := sortedConflicts(c)
	parts := make([]string, 0, len(conflicts))
	for _, conflict := range conflicts {
		parts = append(parts, fmt.Sprintf(
			"%s owns %s, attempted %s",
			conflict.Owner,
			conflict.OwnedPath,
			conflict.AttemptedPath,
		))
	}

	return "field ownership conflicts: " + strings.Join(parts, "; ")
}

// sortedConflicts returns conflicts in Owner, AttemptedPath, OwnedPath order.
func sortedConflicts(conflicts ConflictSet) ConflictSet {
	sorted := make(ConflictSet, len(conflicts))
	copy(sorted, conflicts)
	slices.SortFunc(sorted, compareConflicts)

	return sorted
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
