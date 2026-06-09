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
		{name: "invalid executor", reason: ErrorReasonInvalidExecutor, text: "invalid_executor", valid: true},
		{name: "resource not found", reason: ErrorReasonResourceNotFound, text: "resource_not_found", valid: true},
		{name: "validation failed", reason: ErrorReasonValidationFailed, text: "validation_failed", valid: true},
		{name: "apply failed", reason: ErrorReasonApplyFailed, text: "apply_failed", valid: true},
		{name: "conflict", reason: ErrorReasonConflict, text: "conflict", valid: true},
		{name: "not found", reason: ErrorReasonNotFound, text: "not_found", valid: true},
		{name: "already exists", reason: ErrorReasonAlreadyExists, text: "already_exists", valid: true},
		{name: "stale revision", reason: ErrorReasonStaleRevision, text: "stale_revision", valid: true},
		{name: "store failed", reason: ErrorReasonStoreFailed, text: "store_failed", valid: true},
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
