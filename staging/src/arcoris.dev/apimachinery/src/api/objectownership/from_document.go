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

package objectownership

import (
	"fmt"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
)

// StateFromDocument validates, normalizes, and converts doc into State.
//
// The conversion accepts valid raw document shape, canonicalizes it through
// Normalize, and then builds fieldownership.State for each object surface. It
// does not perform codec decoding or storage reads.
//
// StateFromDocument intentionally performs a second surface-to-state pass after
// Normalize. The extra pass keeps validation, canonical document output, and
// final State construction explicit until real performance targets justify
// merging those steps.
func StateFromDocument(doc Document) (State, error) {
	normalized, err := Normalize(doc)
	if err != nil {
		return State{}, err
	}

	desired, err := stateFromSurface(pathDocumentDesired, normalized.Desired)
	if err != nil {
		return State{}, err
	}

	return NewState(desired), nil
}

// stateFromSurface builds fieldownership.State for one normalized surface.
//
// The surface is expected to have valid canonical paths. fieldownership.NewState
// still performs the final invariant check and duplicate-owner merge.
func stateFromSurface(path string, surface Surface) (fieldownership.State, error) {
	entries := make([]fieldownership.Entry, 0, len(surface.Entries))
	for i, entry := range surface.Entries {
		stateEntry, err := entryFromDocument(entryPath(path, i), entry)
		if err != nil {
			return fieldownership.State{}, err
		}
		entries = append(entries, stateEntry)
	}

	state, err := fieldownership.NewState(entries...)
	if err != nil {
		return fieldownership.State{}, wrapAt(
			path,
			ErrInvalidSurface,
			ErrorReasonInvalidSurface,
			"surface ownership state is invalid",
			err,
		)
	}

	return state, nil
}

// entryFromDocument converts one document entry into fieldownership.Entry.
//
// Owner validation and empty-field allowances are delegated to fieldownership so
// objectownership does not duplicate field-level ownership rules.
func entryFromDocument(path string, entry Entry) (fieldownership.Entry, error) {
	fields, err := fieldsFromDocument(path, entry.Fields)
	if err != nil {
		return fieldownership.Entry{}, err
	}

	stateEntry, err := fieldownership.NewEntry(entry.Owner, fields)
	if err != nil {
		return fieldownership.Entry{}, wrapAt(
			path,
			ErrInvalidEntry,
			ErrorReasonInvalidEntry,
			"ownership entry is invalid",
			err,
		)
	}

	return stateEntry, nil
}

// fieldsFromDocument parses document path strings into a canonical field set.
//
// Duplicate fields are allowed in raw documents and are deduplicated by
// fieldpath.NewSet.
func fieldsFromDocument(path string, fields []Path) (fieldpath.Set, error) {
	paths := make([]fieldpath.Path, 0, len(fields))
	for i, field := range fields {
		parsed, err := parsePath(fieldPath(path, i), field)
		if err != nil {
			return fieldpath.Set{}, err
		}
		paths = append(paths, parsed)
	}

	set, err := fieldpath.NewSet(paths...)
	if err != nil {
		return fieldpath.Set{}, wrapAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"field set is invalid",
			err,
		)
	}

	return set, nil
}

// entryPath returns a stable diagnostic path for one surface entry.
func entryPath(path string, index int) string {
	return fmt.Sprintf("%s.entries[%d]", path, index)
}

// fieldPath returns a stable diagnostic path for one entry field.
func fieldPath(path string, index int) string {
	return fmt.Sprintf("%s.fields[%d]", path, index)
}
