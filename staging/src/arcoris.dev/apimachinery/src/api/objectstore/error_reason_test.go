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

import "testing"

func TestErrorReasonStringAndValidity(t *testing.T) {
	tests := []struct {
		name  string
		in    ErrorReason
		text  string
		valid bool
	}{
		{name: "not found", in: ErrorReasonNotFound, text: "not_found", valid: true},
		{name: "already exists", in: ErrorReasonAlreadyExists, text: "already_exists", valid: true},
		{name: "conflict", in: ErrorReasonConflict, text: "conflict", valid: true},
		{name: "stale revision", in: ErrorReasonStaleRevision, text: "stale_revision", valid: true},
		{name: "invalid key", in: ErrorReasonInvalidKey, text: "invalid_key", valid: true},
		{name: "invalid state", in: ErrorReasonInvalidState, text: "invalid_state", valid: true},
		{name: "invalid state object", in: ErrorReasonInvalidStateObject, text: "invalid_state_object", valid: true},
		{name: "missing desired", in: ErrorReasonMissingDesired, text: "missing_desired", valid: true},
		{name: "invalid observed", in: ErrorReasonInvalidObserved, text: "invalid_observed", valid: true},
		{name: "invalid ownership", in: ErrorReasonInvalidOwnership, text: "invalid_ownership", valid: true},
		{name: "invalid revision", in: ErrorReasonInvalidRevision, text: "invalid_revision", valid: true},
		{name: "nil context", in: ErrorReasonNilContext, text: "nil_context", valid: true},
		{name: "uninitialized store", in: ErrorReasonUninitializedStore, text: "uninitialized_store", valid: true},
		{name: "unknown", in: "", text: "unknown", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.String(); got != tt.text {
				t.Fatalf("String() = %q; want %q", got, tt.text)
			}
			if got := tt.in.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
		})
	}
}
