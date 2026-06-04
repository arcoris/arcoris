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

package codecjson

import "testing"

// TestErrorReasonStrings keeps diagnostic reason strings stable.
func TestErrorReasonStrings(t *testing.T) {
	tests := []struct {
		name   string
		reason ErrorReason
		want   string
	}{
		{name: "invalid JSON", reason: ErrorReasonInvalidJSON, want: "invalid_json"},
		{name: "duplicate key", reason: ErrorReasonDuplicateKey, want: "duplicate_key"},
		{name: "trailing data", reason: ErrorReasonTrailingData, want: "trailing_data"},
		{name: "unsupported value", reason: ErrorReasonUnsupportedValue, want: "unsupported_value"},
		{name: "invalid number", reason: ErrorReasonInvalidNumber, want: "invalid_number"},
		{name: "invalid envelope", reason: ErrorReasonInvalidEnvelope, want: "invalid_envelope"},
		{name: "max depth", reason: ErrorReasonMaxDepthExceeded, want: "max_depth_exceeded"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := string(tc.reason); got != tc.want {
				t.Fatalf("reason = %q; want %q", got, tc.want)
			}
		})
	}
}
