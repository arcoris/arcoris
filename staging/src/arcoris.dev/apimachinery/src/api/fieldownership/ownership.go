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

// OwnedPath reports one owner/path record returned by ownership queries.
type OwnedPath struct {
	// Owner is the field ownership identity that owns Path.
	Owner Owner

	// Path is the owned semantic field path.
	Path fieldpath.Path
}

// ValidateStructure checks whether p is a well-formed owner/path record.
//
// It validates only the ownership record shape. Root paths are allowed, and the
// method does not check object surfaces, descriptors, admission, storage, or
// policy.
func (p OwnedPath) ValidateStructure() error {
	return validateOwnedPath(p)
}

// OwnedPathSet is a deterministic immutable-by-convention owner/path collection.
type OwnedPathSet struct {
	paths []OwnedPath
}

// NewOwnedPathSet constructs a validated deterministic owner/path set.
func NewOwnedPathSet(paths ...OwnedPath) (OwnedPathSet, error) {
	for _, path := range paths {
		if err := path.ValidateStructure(); err != nil {
			return OwnedPathSet{}, err
		}
	}

	return newOwnedPathSetUnchecked(paths...), nil
}

// MustOwnedPathSet constructs an owner/path set or panics when any record is invalid.
func MustOwnedPathSet(paths ...OwnedPath) OwnedPathSet {
	set, err := NewOwnedPathSet(paths...)
	if err != nil {
		panic(err)
	}

	return set
}

// newOwnedPathSetUnchecked constructs a deterministic owner/path set from trusted records.
func newOwnedPathSetUnchecked(paths ...OwnedPath) OwnedPathSet {
	if len(paths) == 0 {
		return OwnedPathSet{}
	}

	ordered := make([]OwnedPath, len(paths))
	copy(ordered, paths)
	sortOwnedPaths(ordered)

	return OwnedPathSet{paths: compactOwnedPaths(ordered)}
}

// Len returns the number of owner/path records.
func (s OwnedPathSet) Len() int {
	return len(s.paths)
}

// IsEmpty reports whether s contains no owner/path records.
func (s OwnedPathSet) IsEmpty() bool {
	return len(s.paths) == 0
}

// Paths returns detached owner/path records in deterministic order.
func (s OwnedPathSet) Paths() []OwnedPath {
	if len(s.paths) == 0 {
		return nil
	}

	paths := make([]OwnedPath, len(s.paths))
	copy(paths, s.paths)

	return paths
}

// ForEach visits owner/path records in deterministic order until fn returns false.
//
// fn must be non-nil. Passing nil is programmer error and may panic.
func (s OwnedPathSet) ForEach(fn func(index int, path OwnedPath) bool) {
	for i, path := range s.paths {
		if !fn(i, path) {
			return
		}
	}
}

// Owners returns sorted unique owners represented in s.
func (s OwnedPathSet) Owners() []Owner {
	if len(s.paths) == 0 {
		return nil
	}

	owners := make([]Owner, 0, len(s.paths))
	for _, path := range s.paths {
		owners = append(owners, path.Owner)
	}
	sortOwners(owners)

	return compactSortedOwners(owners)
}

// FieldPaths returns the canonical set of field paths represented in s.
func (s OwnedPathSet) FieldPaths() fieldpath.Set {
	fields := fieldpath.EmptySet()
	for _, path := range s.paths {
		fields = fields.Insert(path.Path)
	}

	return fields
}
