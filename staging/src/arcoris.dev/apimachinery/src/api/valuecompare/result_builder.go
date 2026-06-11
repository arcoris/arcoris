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

// resultBuilder accumulates comparison buckets before canonical set construction.
type resultBuilder struct {
	added    []fieldpath.Path
	removed  []fieldpath.Path
	modified []fieldpath.Path
}

// AddAddedSet records paths that are present only in the new value.
func (b *resultBuilder) AddAddedSet(set fieldpath.Set) {
	set.ForEach(func(_ int, path fieldpath.Path) bool {
		b.added = append(b.added, path)
		return true
	})
}

// AddRemovedSet records paths that are present only in the old value.
func (b *resultBuilder) AddRemovedSet(set fieldpath.Set) {
	set.ForEach(func(_ int, path fieldpath.Path) bool {
		b.removed = append(b.removed, path)
		return true
	})
}

// AddModified records one changed semantic path.
func (b *resultBuilder) AddModified(path fieldpath.Path) {
	b.modified = append(b.modified, path)
}

// AddResult appends all buckets from result into the builder.
func (b *resultBuilder) AddResult(result Result) {
	b.AddAddedSet(result.added)
	b.AddRemovedSet(result.removed)
	b.AddModifiedSet(result.modified)
}

// AddModifiedSet records modified paths from a canonical set.
func (b *resultBuilder) AddModifiedSet(set fieldpath.Set) {
	set.ForEach(func(_ int, path fieldpath.Path) bool {
		b.modified = append(b.modified, path)
		return true
	})
}

// Build canonicalizes accumulated buckets into a validated Result.
func (b *resultBuilder) Build() (Result, error) {
	added, err := buildResultSet(b.added)
	if err != nil {
		return Result{}, err
	}
	removed, err := buildResultSet(b.removed)
	if err != nil {
		return Result{}, err
	}
	modified, err := buildResultSet(b.modified)
	if err != nil {
		return Result{}, err
	}

	return NewResult(added, removed, modified)
}

// buildResultSet canonicalizes a transient path slice into a fieldpath.Set.
func buildResultSet(paths []fieldpath.Path) (fieldpath.Set, error) {
	if len(paths) == 0 {
		return fieldpath.EmptySet(), nil
	}

	set, err := fieldpath.NewSet(paths...)
	if err != nil {
		return fieldpath.Set{}, wrapAt(
			fieldpath.Root(),
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"result field path cannot be stored in a set",
			err,
		)
	}

	return set, nil
}
