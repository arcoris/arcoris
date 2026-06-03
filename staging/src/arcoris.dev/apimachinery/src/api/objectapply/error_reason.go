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

package objectapply

// ErrorReason is the stable machine-readable reason for an objectapply error.
type ErrorReason string

const (
	// ErrorReasonInvalidOwner reports malformed field owner input.
	ErrorReasonInvalidOwner ErrorReason = "invalid_owner"

	// ErrorReasonInvalidResource reports missing or unusable resource input.
	ErrorReasonInvalidResource ErrorReason = "invalid_resource"

	// ErrorReasonInvalidLiveObject reports live object validation failure.
	ErrorReasonInvalidLiveObject ErrorReason = "invalid_live_object"

	// ErrorReasonInvalidAppliedObject reports applied object validation failure.
	ErrorReasonInvalidAppliedObject ErrorReason = "invalid_applied_object"

	// ErrorReasonIdentityMismatch reports incompatible object identity.
	ErrorReasonIdentityMismatch ErrorReason = "identity_mismatch"

	// ErrorReasonVersionMismatch reports unsupported cross-version apply.
	ErrorReasonVersionMismatch ErrorReason = "version_mismatch"

	// ErrorReasonUnsupportedObservedApply reports attempted Observed apply.
	ErrorReasonUnsupportedObservedApply ErrorReason = "unsupported_observed_apply"

	// ErrorReasonUnsupportedMetadataChange reports attempted metadata apply.
	ErrorReasonUnsupportedMetadataChange ErrorReason = "unsupported_metadata_change"

	// ErrorReasonDesiredApplyFailed reports non-conflict valueapply failure.
	ErrorReasonDesiredApplyFailed ErrorReason = "desired_apply_failed"

	// ErrorReasonConflict reports Desired ownership conflict.
	ErrorReasonConflict ErrorReason = "conflict"
)
