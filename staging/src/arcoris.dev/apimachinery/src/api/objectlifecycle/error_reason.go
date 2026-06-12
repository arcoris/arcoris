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
	// ErrorReasonInvalidOwner reports malformed field owner input.
	ErrorReasonInvalidOwner ErrorReason = "invalid_owner"
	// ErrorReasonInvalidContext reports a nil lifecycle operation context.
	ErrorReasonInvalidContext ErrorReason = "invalid_context"
	// ErrorReasonInvalidExpectedRevision reports a missing delete/update revision.
	ErrorReasonInvalidExpectedRevision ErrorReason = "invalid_expected_revision"
	// ErrorReasonUnsupportedObservedApply reports attempted Observed apply intent.
	ErrorReasonUnsupportedObservedApply ErrorReason = "unsupported_observed_apply"
	// ErrorReasonInvalidExecutor reports missing executor dependencies.
	ErrorReasonInvalidExecutor ErrorReason = "invalid_executor"
	// ErrorReasonResourceNotFound reports a resource resolver miss.
	ErrorReasonResourceNotFound ErrorReason = "resource_not_found"
	// ErrorReasonInvalidResourceContract reports inconsistent resolved resource data.
	ErrorReasonInvalidResourceContract ErrorReason = "invalid_resource_contract"
	// ErrorReasonValidationFailed reports objectvalidation failure.
	ErrorReasonValidationFailed ErrorReason = "validation_failed"
	// ErrorReasonApplyFailed reports non-conflict objectapply or ownership-init failure.
	ErrorReasonApplyFailed ErrorReason = "apply_failed"
	// ErrorReasonOwnershipInitFailed reports create/apply-create ownership extraction failure.
	ErrorReasonOwnershipInitFailed ErrorReason = "ownership_init_failed"
	// ErrorReasonUnsupportedSurface reports an operation against an unsupported object surface.
	ErrorReasonUnsupportedSurface ErrorReason = "unsupported_surface"
	// ErrorReasonObservedNotDefined reports Observed updates for versions without Observed.
	ErrorReasonObservedNotDefined ErrorReason = "observed_not_defined"
	// ErrorReasonInvalidObserved reports invalid Observed update payload.
	ErrorReasonInvalidObserved ErrorReason = "invalid_observed"
	// ErrorReasonInvalidMetadataPatch reports malformed labels/annotations patch input.
	ErrorReasonInvalidMetadataPatch ErrorReason = "invalid_metadata_patch"
	// ErrorReasonEmptyMetadataPatch reports a metadata patch without changes.
	ErrorReasonEmptyMetadataPatch ErrorReason = "empty_metadata_patch"
	// ErrorReasonInvalidMetadataKey reports malformed label or annotation keys.
	ErrorReasonInvalidMetadataKey ErrorReason = "invalid_metadata_key"
	// ErrorReasonMetadataOwnershipFailed reports metadata ownership state construction failure.
	ErrorReasonMetadataOwnershipFailed ErrorReason = "metadata_ownership_failed"
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
	// ErrorReasonStoreInvalidState reports an invalid state rejected by objectstore.
	ErrorReasonStoreInvalidState ErrorReason = "store_invalid_state"
)

// IsValid reports whether r is a known lifecycle reason.
func (r ErrorReason) IsValid() bool {
	switch r {
	case ErrorReasonInvalidRequest,
		ErrorReasonInvalidOwner,
		ErrorReasonInvalidContext,
		ErrorReasonInvalidExpectedRevision,
		ErrorReasonUnsupportedObservedApply,
		ErrorReasonInvalidExecutor,
		ErrorReasonResourceNotFound,
		ErrorReasonInvalidResourceContract,
		ErrorReasonValidationFailed,
		ErrorReasonApplyFailed,
		ErrorReasonOwnershipInitFailed,
		ErrorReasonUnsupportedSurface,
		ErrorReasonObservedNotDefined,
		ErrorReasonInvalidObserved,
		ErrorReasonInvalidMetadataPatch,
		ErrorReasonEmptyMetadataPatch,
		ErrorReasonInvalidMetadataKey,
		ErrorReasonMetadataOwnershipFailed,
		ErrorReasonConflict,
		ErrorReasonNotFound,
		ErrorReasonAlreadyExists,
		ErrorReasonStaleRevision,
		ErrorReasonStoreFailed,
		ErrorReasonStoreInvalidState:
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
