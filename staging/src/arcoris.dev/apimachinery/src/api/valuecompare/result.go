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

package valuecompare

import "arcoris.dev/apimachinery/api/fieldpath"

// Result stores the semantic paths changed by a comparison.
//
// Sets are canonical fieldpath.Set values and are immutable by convention. A
// well-formed comparison result places a path in at most one bucket.
type Result struct {
	// Added contains paths mentioned by the new value but not the old value.
	Added fieldpath.Set
	// Removed contains paths mentioned by the old value but not the new value.
	Removed fieldpath.Set
	// Modified contains paths present in both values whose payload changed.
	Modified fieldpath.Set
}

// EmptyResult returns a result whose buckets are canonical empty sets.
func EmptyResult() Result {
	return Result{
		Added:    fieldpath.EmptySet(),
		Removed:  fieldpath.EmptySet(),
		Modified: fieldpath.EmptySet(),
	}
}

// IsEmpty reports whether no semantic changes were found.
func (r Result) IsEmpty() bool {
	return r.Added.IsEmpty() && r.Removed.IsEmpty() && r.Modified.IsEmpty()
}

// Changed returns every path mentioned by Added, Removed, or Modified.
func (r Result) Changed() fieldpath.Set {
	return unionSets(unionSets(r.Added, r.Removed), r.Modified)
}

// merge combines child comparison output into r.
func (r Result) merge(other Result) Result {
	if other.IsEmpty() {
		return r
	}

	r.Added = unionSets(r.Added, other.Added)
	r.Removed = unionSets(r.Removed, other.Removed)
	r.Modified = unionSets(r.Modified, other.Modified)

	return r
}

// withAdded returns a copy of r with added paths merged into Added.
func (r Result) withAdded(set fieldpath.Set) Result {
	r.Added = unionSets(r.Added, set)
	return r
}

// withRemoved returns a copy of r with removed paths merged into Removed.
func (r Result) withRemoved(set fieldpath.Set) Result {
	r.Removed = unionSets(r.Removed, set)
	return r
}

// withModified returns a copy of r with path recorded as one modified leaf.
func (r Result) withModified(path fieldpath.Path) (Result, error) {
	set, err := setAt(path)
	if err != nil {
		return Result{}, err
	}

	r.Modified = unionSets(r.Modified, set)
	return r, nil
}

// unionSets keeps common empty-set merges allocation-free.
func unionSets(left fieldpath.Set, right fieldpath.Set) fieldpath.Set {
	switch {
	case left.IsEmpty():
		return right
	case right.IsEmpty():
		return left
	default:
		return left.Union(right)
	}
}
