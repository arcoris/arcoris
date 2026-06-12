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

package objectstore

// ErrorReason is stable machine-readable detail for objectstore errors.
type ErrorReason string

const (
	// ErrorReasonNotFound reports a missing or tombstoned live object.
	ErrorReasonNotFound ErrorReason = "not_found"
	// ErrorReasonAlreadyExists reports a create request for an existing live object.
	ErrorReasonAlreadyExists ErrorReason = "already_exists"
	// ErrorReasonConflict reports a compare-and-swap race after preconditions passed.
	ErrorReasonConflict ErrorReason = "conflict"
	// ErrorReasonStaleRevision reports an expected revision mismatch.
	ErrorReasonStaleRevision ErrorReason = "stale_revision"
	// ErrorReasonInvalidKey reports an invalid object store key.
	ErrorReasonInvalidKey ErrorReason = "invalid_key"
	// ErrorReasonInvalidState reports otherwise invalid object store state.
	ErrorReasonInvalidState ErrorReason = "invalid_state"
	// ErrorReasonInvalidStateObject reports invalid object envelope metadata.
	ErrorReasonInvalidStateObject ErrorReason = "invalid_state_object"
	// ErrorReasonMissingDesired reports missing committed Desired data.
	ErrorReasonMissingDesired ErrorReason = "missing_desired"
	// ErrorReasonInvalidObserved reports invalid optional Observed data.
	ErrorReasonInvalidObserved ErrorReason = "invalid_observed"
	// ErrorReasonInvalidOwnership reports invalid or non-canonical ownership.
	ErrorReasonInvalidOwnership ErrorReason = "invalid_ownership"
	// ErrorReasonInvalidRevision reports a forbidden or missing revision.
	ErrorReasonInvalidRevision ErrorReason = "invalid_revision"
	// ErrorReasonNilContext reports a nil operation context.
	ErrorReasonNilContext ErrorReason = "nil_context"
	// ErrorReasonUninitializedStore reports use of a nil or zero implementation.
	ErrorReasonUninitializedStore ErrorReason = "uninitialized_store"
)

// IsValid reports whether r is a known objectstore error reason.
func (r ErrorReason) IsValid() bool {
	switch r {
	case ErrorReasonNotFound,
		ErrorReasonAlreadyExists,
		ErrorReasonConflict,
		ErrorReasonStaleRevision,
		ErrorReasonInvalidKey,
		ErrorReasonInvalidState,
		ErrorReasonInvalidStateObject,
		ErrorReasonMissingDesired,
		ErrorReasonInvalidObserved,
		ErrorReasonInvalidOwnership,
		ErrorReasonInvalidRevision,
		ErrorReasonNilContext,
		ErrorReasonUninitializedStore:
		return true
	default:
		return false
	}
}

// String returns stable diagnostic text for r.
func (r ErrorReason) String() string {
	if !r.IsValid() {
		return "unknown"
	}

	return string(r)
}
