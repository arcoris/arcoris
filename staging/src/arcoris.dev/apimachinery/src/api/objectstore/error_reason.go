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
type ErrorReason uint8

const (
	// ErrorReasonNotFound reports a missing or tombstoned live object.
	ErrorReasonNotFound ErrorReason = iota + 1
	// ErrorReasonAlreadyExists reports a create request for an existing live object.
	ErrorReasonAlreadyExists
	// ErrorReasonConflict reports a compare-and-swap race after preconditions passed.
	ErrorReasonConflict
	// ErrorReasonStaleRevision reports an expected revision mismatch.
	ErrorReasonStaleRevision
	// ErrorReasonInvalidKey reports an invalid object store key.
	ErrorReasonInvalidKey
	// ErrorReasonInvalidState reports invalid object or ownership state.
	ErrorReasonInvalidState
	// ErrorReasonInvalidRevision reports a forbidden or missing revision.
	ErrorReasonInvalidRevision
	// ErrorReasonUninitializedStore reports use of a nil or zero implementation.
	ErrorReasonUninitializedStore
)

// IsValid reports whether r is a known objectstore error reason.
func (r ErrorReason) IsValid() bool {
	return r >= ErrorReasonNotFound && r <= ErrorReasonUninitializedStore
}

// String returns stable diagnostic text for r.
func (r ErrorReason) String() string {
	switch r {
	case ErrorReasonNotFound:
		return "not_found"
	case ErrorReasonAlreadyExists:
		return "already_exists"
	case ErrorReasonConflict:
		return "conflict"
	case ErrorReasonStaleRevision:
		return "stale_revision"
	case ErrorReasonInvalidKey:
		return "invalid_key"
	case ErrorReasonInvalidState:
		return "invalid_state"
	case ErrorReasonInvalidRevision:
		return "invalid_revision"
	case ErrorReasonUninitializedStore:
		return "uninitialized_store"
	default:
		return "unknown"
	}
}
