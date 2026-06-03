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

// Result stores semantic paths added, removed, or modified by a comparison.
//
// Sets are canonical fieldpath.Set values. A path appears in at most one bucket.
type Result struct {
	// Added contains paths present only in the new value.
	Added fieldpath.Set
	// Removed contains paths present only in the old value.
	Removed fieldpath.Set
	// Modified contains paths present in both values whose semantic payload changed.
	Modified fieldpath.Set
}

// EmptyResult returns a result with canonical empty sets.
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

// Changed returns the union of Added, Removed, and Modified.
func (r Result) Changed() fieldpath.Set {
	return unionSets(unionSets(r.Added, r.Removed), r.Modified)
}

// merge folds another comparison result into r.
func (r Result) merge(other Result) Result {
	if other.IsEmpty() {
		return r
	}

	r.Added = unionSets(r.Added, other.Added)
	r.Removed = unionSets(r.Removed, other.Removed)
	r.Modified = unionSets(r.Modified, other.Modified)

	return r
}

// withAdded returns r with added paths unioned in.
func (r Result) withAdded(set fieldpath.Set) Result {
	r.Added = unionSets(r.Added, set)
	return r
}

// withRemoved returns r with removed paths unioned in.
func (r Result) withRemoved(set fieldpath.Set) Result {
	r.Removed = unionSets(r.Removed, set)
	return r
}

// withModified returns r with path recorded as modified.
func (r Result) withModified(path fieldpath.Path) (Result, error) {
	set, err := setAt(path)
	if err != nil {
		return Result{}, err
	}

	r.Modified = unionSets(r.Modified, set)
	return r, nil
}

// unionSets avoids fieldpath.Set.Union allocation when either side is empty.
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
