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

package codecregistry

// checkEntryConflicts rejects duplicate entry identities.
//
// Formats and media types are grouping attributes and may repeat. EntryID is the
// only unique registry identity.
func checkEntryConflicts(
	index int,
	entry Entry,
	seenIDs map[EntryID]int,
) error {
	if previous, ok := seenIDs[entry.id]; ok {
		return duplicateEntryIDError(index, previous, entry.id)
	}
	seenIDs[entry.id] = index

	return nil
}

// duplicateEntryIDError creates a stable duplicate entry ID diagnostic.
func duplicateEntryIDError(
	index int,
	previous int,
	id EntryID,
) error {
	return errorfAt(
		registrationIDPath(index),
		ErrDuplicateEntryID,
		ErrorReasonDuplicateEntryID,
		"entry ID %q duplicates registrations[%d]",
		id,
		previous,
	)
}
