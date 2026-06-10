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

// ErrorReason gives stable machine-readable detail inside a broad error category.
type ErrorReason string

const (
	// ErrorReasonEmptyOwner reports empty owner identity text.
	ErrorReasonEmptyOwner ErrorReason = "empty_owner"

	// ErrorReasonInvalidOwnerUTF8 reports owner text that is not valid UTF-8.
	ErrorReasonInvalidOwnerUTF8 ErrorReason = "invalid_owner_utf8"

	// ErrorReasonWhitespaceOwner reports owner text made only of whitespace.
	ErrorReasonWhitespaceOwner ErrorReason = "whitespace_owner"

	// ErrorReasonOwnerBoundaryWhitespace reports leading or trailing whitespace.
	ErrorReasonOwnerBoundaryWhitespace ErrorReason = "owner_boundary_whitespace"

	// ErrorReasonOwnerTooLong reports owner text exceeding MaxOwnerLength.
	ErrorReasonOwnerTooLong ErrorReason = "owner_too_long"

	// ErrorReasonOwnerControlCharacter reports control characters in owner text.
	ErrorReasonOwnerControlCharacter ErrorReason = "owner_control_character"

	// ErrorReasonInvalidEntry reports malformed ownership entry state.
	ErrorReasonInvalidEntry ErrorReason = "invalid_entry"

	// ErrorReasonInvalidEntryOwner reports a malformed entry owner.
	ErrorReasonInvalidEntryOwner ErrorReason = "invalid_entry_owner"

	// ErrorReasonInvalidEntryFields reports malformed entry field paths.
	ErrorReasonInvalidEntryFields ErrorReason = "invalid_entry_fields"

	// ErrorReasonInvalidState reports malformed ownership state.
	ErrorReasonInvalidState ErrorReason = "invalid_state"

	// ErrorReasonUnsortedStateEntries reports entries not sorted by owner.
	ErrorReasonUnsortedStateEntries ErrorReason = "unsorted_state_entries"

	// ErrorReasonDuplicateStateOwner reports repeated owner entries in State.
	ErrorReasonDuplicateStateOwner ErrorReason = "duplicate_state_owner"

	// ErrorReasonEmptyStateEntry reports an empty entry stored in State.
	ErrorReasonEmptyStateEntry ErrorReason = "empty_state_entry"

	// ErrorReasonInvalidPath reports malformed semantic field paths.
	ErrorReasonInvalidPath ErrorReason = "invalid_path"

	// ErrorReasonInvalidOwnedPath reports a malformed stored owned path.
	ErrorReasonInvalidOwnedPath ErrorReason = "invalid_owned_path"

	// ErrorReasonInvalidOwnedPathOwner reports a malformed OwnedPath owner.
	ErrorReasonInvalidOwnedPathOwner ErrorReason = "invalid_owned_path_owner"

	// ErrorReasonInvalidOwnedPathPath reports a malformed OwnedPath path.
	ErrorReasonInvalidOwnedPathPath ErrorReason = "invalid_owned_path_path"

	// ErrorReasonInvalidAttemptedPath reports a malformed attempted path.
	ErrorReasonInvalidAttemptedPath ErrorReason = "invalid_attempted_path"

	// ErrorReasonInvalidConflict reports malformed conflict state.
	ErrorReasonInvalidConflict ErrorReason = "invalid_conflict"

	// ErrorReasonInvalidConflictOwner reports a malformed conflict owner.
	ErrorReasonInvalidConflictOwner ErrorReason = "invalid_conflict_owner"

	// ErrorReasonInvalidConflictOwnedPath reports a malformed conflict owned path.
	ErrorReasonInvalidConflictOwnedPath ErrorReason = "invalid_conflict_owned_path"

	// ErrorReasonInvalidConflictAttemptedPath reports a malformed conflict attempted path.
	ErrorReasonInvalidConflictAttemptedPath ErrorReason = "invalid_conflict_attempted_path"

	// ErrorReasonConflict reports overlapping ownership for an attempted field set.
	ErrorReasonConflict ErrorReason = "conflict"
)
