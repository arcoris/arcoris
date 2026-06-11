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
// Result is immutable to callers. Bucket accessors return fieldpath.Set values,
// whose contents are canonical and immutable by convention. A structurally valid
// result places an exact path in at most one bucket.
type Result struct {
	added    fieldpath.Set
	removed  fieldpath.Set
	modified fieldpath.Set
}

// EmptyResult returns a result whose buckets are canonical empty sets.
func EmptyResult() Result {
	return Result{
		added:    fieldpath.EmptySet(),
		removed:  fieldpath.EmptySet(),
		modified: fieldpath.EmptySet(),
	}
}

// NewResult constructs a structurally valid comparison result.
//
// Added, removed, and modified buckets must be valid fieldpath sets. An exact
// path may appear in only one bucket.
func NewResult(added fieldpath.Set, removed fieldpath.Set, modified fieldpath.Set) (Result, error) {
	result := Result{
		added:    added,
		removed:  removed,
		modified: modified,
	}
	if err := result.ValidateStructure(); err != nil {
		return Result{}, err
	}

	return result, nil
}

// MustResult is the panic-on-error form of NewResult for tests and static fixtures.
func MustResult(added fieldpath.Set, removed fieldpath.Set, modified fieldpath.Set) Result {
	result, err := NewResult(added, removed, modified)
	if err != nil {
		panic(err)
	}

	return result
}

// Added returns paths mentioned by the new value but not the old value.
func (r Result) Added() fieldpath.Set {
	return r.added
}

// Removed returns paths mentioned by the old value but not the new value.
func (r Result) Removed() fieldpath.Set {
	return r.removed
}

// Modified returns paths present in both values whose semantic payload changed.
func (r Result) Modified() fieldpath.Set {
	return r.modified
}

// IsEmpty reports whether no semantic changes were found.
func (r Result) IsEmpty() bool {
	return r.added.IsEmpty() && r.removed.IsEmpty() && r.modified.IsEmpty()
}

// Changed returns every path mentioned by Added, Removed, or Modified.
func (r Result) Changed() fieldpath.Set {
	return unionSets(unionSets(r.added, r.removed), r.modified)
}

// merge combines child comparison output into r.
func (r Result) merge(other Result) Result {
	if other.IsEmpty() {
		return r
	}

	r.added = unionSets(r.added, other.added)
	r.removed = unionSets(r.removed, other.removed)
	r.modified = unionSets(r.modified, other.modified)

	return r
}

// withAdded returns a copy of r with added paths merged into Added.
func (r Result) withAdded(set fieldpath.Set) Result {
	r.added = unionSets(r.added, set)
	return r
}

// withRemoved returns a copy of r with removed paths merged into Removed.
func (r Result) withRemoved(set fieldpath.Set) Result {
	r.removed = unionSets(r.removed, set)
	return r
}

// withModified returns a copy of r with path recorded as one modified leaf.
func (r Result) withModified(path fieldpath.Path) (Result, error) {
	set, err := setAt(path)
	if err != nil {
		return Result{}, err
	}

	r.modified = unionSets(r.modified, set)
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
