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

package valuefieldset

import "testing"

func TestErrorReasonStrings(t *testing.T) {
	tests := []struct {
		reason ErrorReason
		want   string
	}{
		{ErrorReasonInvalidZero, "invalid_zero"},
		{ErrorReasonInvalidDescriptor, "invalid_descriptor"},
		{ErrorReasonKindMismatch, "kind_mismatch"},
		{ErrorReasonUnknownField, "unknown_field"},
		{ErrorReasonUnresolvedRef, "unresolved_ref"},
		{ErrorReasonReferenceCycle, "reference_cycle"},
		{ErrorReasonMissingListKey, "missing_list_key"},
		{ErrorReasonInvalidListKey, "invalid_list_key"},
		{ErrorReasonDuplicateListKey, "duplicate_list_key"},
	}

	for _, tt := range tests {
		if string(tt.reason) != tt.want {
			t.Fatalf("reason = %q, want %q", tt.reason, tt.want)
		}
	}
}
