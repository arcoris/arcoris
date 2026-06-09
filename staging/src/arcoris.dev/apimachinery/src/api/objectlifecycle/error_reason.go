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

package objectlifecycle

// ErrorReason is stable machine-readable detail for lifecycle errors.
type ErrorReason string

const (
	// ErrorReasonInvalidRequest reports malformed operation input.
	ErrorReasonInvalidRequest ErrorReason = "invalid_request"
	// ErrorReasonInvalidExecutor reports missing executor dependencies.
	ErrorReasonInvalidExecutor ErrorReason = "invalid_executor"
	// ErrorReasonResourceNotFound reports a resource resolver miss.
	ErrorReasonResourceNotFound ErrorReason = "resource_not_found"
	// ErrorReasonValidationFailed reports objectvalidation failure.
	ErrorReasonValidationFailed ErrorReason = "validation_failed"
	// ErrorReasonApplyFailed reports non-conflict objectapply or ownership-init failure.
	ErrorReasonApplyFailed ErrorReason = "apply_failed"
	// ErrorReasonConflict reports field ownership or store compare-and-swap conflict.
	ErrorReasonConflict ErrorReason = "conflict"
	// ErrorReasonNotFound reports a missing live object.
	ErrorReasonNotFound ErrorReason = "not_found"
	// ErrorReasonAlreadyExists reports create conflict with existing live state.
	ErrorReasonAlreadyExists ErrorReason = "already_exists"
	// ErrorReasonStaleRevision reports optimistic revision mismatch.
	ErrorReasonStaleRevision ErrorReason = "stale_revision"
	// ErrorReasonStoreFailed reports an unexpected objectstore error.
	ErrorReasonStoreFailed ErrorReason = "store_failed"
)

// IsValid reports whether r is a known lifecycle reason.
func (r ErrorReason) IsValid() bool {
	switch r {
	case ErrorReasonInvalidRequest,
		ErrorReasonInvalidExecutor,
		ErrorReasonResourceNotFound,
		ErrorReasonValidationFailed,
		ErrorReasonApplyFailed,
		ErrorReasonConflict,
		ErrorReasonNotFound,
		ErrorReasonAlreadyExists,
		ErrorReasonStaleRevision,
		ErrorReasonStoreFailed:
		return true
	default:
		return false
	}
}

// String returns stable diagnostic text for r.
func (r ErrorReason) String() string {
	if r.IsValid() {
		return string(r)
	}

	return "unknown"
}
