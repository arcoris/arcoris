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

import "testing"

func TestErrorReasons(t *testing.T) {
	requireEqual(t, string(ErrorReasonEmptyOwner), "empty_owner")
	requireEqual(t, string(ErrorReasonInvalidOwnerUTF8), "invalid_owner_utf8")
	requireEqual(t, string(ErrorReasonWhitespaceOwner), "whitespace_owner")
	requireEqual(t, string(ErrorReasonOwnerBoundaryWhitespace), "owner_boundary_whitespace")
	requireEqual(t, string(ErrorReasonOwnerTooLong), "owner_too_long")
	requireEqual(t, string(ErrorReasonOwnerControlCharacter), "owner_control_character")
	requireEqual(t, string(ErrorReasonInvalidEntry), "invalid_entry")
	requireEqual(t, string(ErrorReasonInvalidEntryOwner), "invalid_entry_owner")
	requireEqual(t, string(ErrorReasonInvalidEntryFields), "invalid_entry_fields")
	requireEqual(t, string(ErrorReasonInvalidState), "invalid_state")
	requireEqual(t, string(ErrorReasonUnsortedStateEntries), "unsorted_state_entries")
	requireEqual(t, string(ErrorReasonDuplicateStateOwner), "duplicate_state_owner")
	requireEqual(t, string(ErrorReasonEmptyStateEntry), "empty_state_entry")
	requireEqual(t, string(ErrorReasonInvalidPath), "invalid_path")
	requireEqual(t, string(ErrorReasonInvalidOwnedPath), "invalid_owned_path")
	requireEqual(t, string(ErrorReasonInvalidOwnedPathOwner), "invalid_owned_path_owner")
	requireEqual(t, string(ErrorReasonInvalidOwnedPathPath), "invalid_owned_path_path")
	requireEqual(t, string(ErrorReasonInvalidAttemptedPath), "invalid_attempted_path")
	requireEqual(t, string(ErrorReasonInvalidConflict), "invalid_conflict")
	requireEqual(t, string(ErrorReasonInvalidConflictOwner), "invalid_conflict_owner")
	requireEqual(t, string(ErrorReasonInvalidConflictOwnedPath), "invalid_conflict_owned_path")
	requireEqual(t, string(ErrorReasonInvalidConflictAttemptedPath), "invalid_conflict_attempted_path")
	requireEqual(t, string(ErrorReasonConflict), "conflict")
}
