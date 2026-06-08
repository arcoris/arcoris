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

// Reason is stable machine-readable detail for objectstore errors.
type Reason uint8

const (
	// ReasonNotFound reports a missing or tombstoned live object.
	ReasonNotFound Reason = iota + 1
	// ReasonAlreadyExists reports a create request for an existing live object.
	ReasonAlreadyExists
	// ReasonConflict reports a compare-and-swap race after preconditions passed.
	ReasonConflict
	// ReasonStaleRevision reports an expected revision mismatch.
	ReasonStaleRevision
	// ReasonInvalidKey reports an invalid object store key.
	ReasonInvalidKey
	// ReasonInvalidState reports invalid object or ownership state.
	ReasonInvalidState
	// ReasonInvalidRevision reports a forbidden or missing revision.
	ReasonInvalidRevision
	// ReasonUninitializedStore reports use of a nil or zero implementation.
	ReasonUninitializedStore
)

// IsValid reports whether r is a known objectstore error reason.
func (r Reason) IsValid() bool {
	return r >= ReasonNotFound && r <= ReasonUninitializedStore
}

// String returns stable diagnostic text for r.
func (r Reason) String() string {
	switch r {
	case ReasonNotFound:
		return "not_found"
	case ReasonAlreadyExists:
		return "already_exists"
	case ReasonConflict:
		return "conflict"
	case ReasonStaleRevision:
		return "stale_revision"
	case ReasonInvalidKey:
		return "invalid_key"
	case ReasonInvalidState:
		return "invalid_state"
	case ReasonInvalidRevision:
		return "invalid_revision"
	case ReasonUninitializedStore:
		return "uninitialized_store"
	default:
		return "unknown"
	}
}
