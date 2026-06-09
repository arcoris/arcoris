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

// Reason is stable machine-readable detail for lifecycle errors.
type Reason string

const (
	// ReasonInvalidRequest reports malformed operation input.
	ReasonInvalidRequest Reason = "invalid_request"
	// ReasonInvalidExecutor reports missing executor dependencies.
	ReasonInvalidExecutor Reason = "invalid_executor"
	// ReasonResourceNotFound reports a resource resolver miss.
	ReasonResourceNotFound Reason = "resource_not_found"
	// ReasonValidationFailed reports objectvalidation failure.
	ReasonValidationFailed Reason = "validation_failed"
	// ReasonApplyFailed reports non-conflict objectapply or ownership-init failure.
	ReasonApplyFailed Reason = "apply_failed"
	// ReasonConflict reports field ownership or store compare-and-swap conflict.
	ReasonConflict Reason = "conflict"
	// ReasonNotFound reports a missing live object.
	ReasonNotFound Reason = "not_found"
	// ReasonAlreadyExists reports create conflict with existing live state.
	ReasonAlreadyExists Reason = "already_exists"
	// ReasonStaleRevision reports optimistic revision mismatch.
	ReasonStaleRevision Reason = "stale_revision"
	// ReasonStoreFailed reports an unexpected objectstore error.
	ReasonStoreFailed Reason = "store_failed"
)

// IsValid reports whether r is a known lifecycle reason.
func (r Reason) IsValid() bool {
	switch r {
	case ReasonInvalidRequest,
		ReasonInvalidExecutor,
		ReasonResourceNotFound,
		ReasonValidationFailed,
		ReasonApplyFailed,
		ReasonConflict,
		ReasonNotFound,
		ReasonAlreadyExists,
		ReasonStaleRevision,
		ReasonStoreFailed:
		return true
	default:
		return false
	}
}

// String returns stable diagnostic text for r.
func (r Reason) String() string {
	if r.IsValid() {
		return string(r)
	}

	return "unknown"
}
