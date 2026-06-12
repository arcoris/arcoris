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

import "testing"

func TestErrorReasonStringAndValidity(t *testing.T) {
	tests := []struct {
		name   string
		reason ErrorReason
		text   string
		valid  bool
	}{
		{name: "invalid request", reason: ErrorReasonInvalidRequest, text: "invalid_request", valid: true},
		{name: "invalid owner", reason: ErrorReasonInvalidOwner, text: "invalid_owner", valid: true},
		{name: "invalid context", reason: ErrorReasonInvalidContext, text: "invalid_context", valid: true},
		{name: "invalid expected revision", reason: ErrorReasonInvalidExpectedRevision, text: "invalid_expected_revision", valid: true},
		{name: "unsupported observed apply", reason: ErrorReasonUnsupportedObservedApply, text: "unsupported_observed_apply", valid: true},
		{name: "invalid executor", reason: ErrorReasonInvalidExecutor, text: "invalid_executor", valid: true},
		{name: "resource not found", reason: ErrorReasonResourceNotFound, text: "resource_not_found", valid: true},
		{name: "invalid resource contract", reason: ErrorReasonInvalidResourceContract, text: "invalid_resource_contract", valid: true},
		{name: "validation failed", reason: ErrorReasonValidationFailed, text: "validation_failed", valid: true},
		{name: "apply failed", reason: ErrorReasonApplyFailed, text: "apply_failed", valid: true},
		{name: "ownership init failed", reason: ErrorReasonOwnershipInitFailed, text: "ownership_init_failed", valid: true},
		{name: "unsupported surface", reason: ErrorReasonUnsupportedSurface, text: "unsupported_surface", valid: true},
		{name: "observed not defined", reason: ErrorReasonObservedNotDefined, text: "observed_not_defined", valid: true},
		{name: "invalid observed", reason: ErrorReasonInvalidObserved, text: "invalid_observed", valid: true},
		{name: "invalid metadata patch", reason: ErrorReasonInvalidMetadataPatch, text: "invalid_metadata_patch", valid: true},
		{name: "empty metadata patch", reason: ErrorReasonEmptyMetadataPatch, text: "empty_metadata_patch", valid: true},
		{name: "invalid metadata key", reason: ErrorReasonInvalidMetadataKey, text: "invalid_metadata_key", valid: true},
		{name: "metadata ownership failed", reason: ErrorReasonMetadataOwnershipFailed, text: "metadata_ownership_failed", valid: true},
		{name: "conflict", reason: ErrorReasonConflict, text: "conflict", valid: true},
		{name: "not found", reason: ErrorReasonNotFound, text: "not_found", valid: true},
		{name: "already exists", reason: ErrorReasonAlreadyExists, text: "already_exists", valid: true},
		{name: "stale revision", reason: ErrorReasonStaleRevision, text: "stale_revision", valid: true},
		{name: "store failed", reason: ErrorReasonStoreFailed, text: "store_failed", valid: true},
		{name: "store invalid state", reason: ErrorReasonStoreInvalidState, text: "store_invalid_state", valid: true},
		{name: "unknown", reason: ErrorReason("other"), text: "unknown", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reason.String(); got != tt.text {
				t.Fatalf("String() = %q; want %q", got, tt.text)
			}
			if got := tt.reason.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
		})
	}
}
