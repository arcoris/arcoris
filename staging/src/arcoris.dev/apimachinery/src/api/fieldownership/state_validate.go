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

// ValidateStructure checks normalized ownership-state invariants.
//
// It does not reject overlapping ownership. Multiple owners may own exact or
// overlapping paths; conflict detection is a separate contextual operation.
func (s State) ValidateStructure() error {
	for i, entry := range s.entries {
		if err := validateStateEntry(i, entry); err != nil {
			return err
		}
		if i == 0 {
			continue
		}

		previous := s.entries[i-1].owner
		switch cmp := previous.Compare(entry.owner); {
		case cmp > 0:
			return errorfAt(
				entryOwnerPath(i),
				ErrInvalidState,
				ErrorReasonUnsortedStateEntries,
				"state entry owner %q sorts before previous owner %q",
				entry.owner,
				previous,
			)
		case cmp == 0:
			return errorfAt(
				entryOwnerPath(i),
				ErrInvalidState,
				ErrorReasonDuplicateStateOwner,
				"state contains duplicate owner %q",
				entry.owner,
			)
		}
	}

	return nil
}

// validateStateEntry checks one stored entry with state-level path context.
func validateStateEntry(index int, entry Entry) error {
	if entry.IsEmpty() {
		return errorfAt(
			entryFieldsPath(index),
			ErrInvalidState,
			ErrorReasonEmptyStateEntry,
			"state entry %d owns no fields",
			index,
		)
	}
	if err := entry.owner.ValidateLexical(); err != nil {
		return wrapAt(
			entryOwnerPath(index),
			ErrInvalidState,
			ErrorReasonInvalidEntryOwner,
			"state entry owner is invalid",
			err,
		)
	}
	if err := validateFieldsAt(
		entryFieldsPath(index),
		entry.fields,
		ErrorReasonInvalidEntryFields,
		"state entry field path is invalid",
	); err != nil {
		return wrapAt(
			entryFieldsPath(index),
			ErrInvalidState,
			ErrorReasonInvalidEntryFields,
			"state entry fields are invalid",
			err,
		)
	}

	return nil
}
