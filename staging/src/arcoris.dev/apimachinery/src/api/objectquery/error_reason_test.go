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

package objectquery

import "testing"

func TestErrorReasonStrings(t *testing.T) {
	tests := []struct {
		name   string
		reason ErrorReason
		want   string
	}{
		{name: "query", reason: ErrorReasonInvalidQuery, want: "invalid_query"},
		{name: "identity", reason: ErrorReasonInvalidIdentity, want: "invalid_identity"},
		{name: "selector", reason: ErrorReasonInvalidSelector, want: "invalid_selector"},
		{name: "requirement", reason: ErrorReasonInvalidRequirement, want: "invalid_requirement"},
		{name: "operator", reason: ErrorReasonInvalidOperator, want: "invalid_operator"},
		{name: "key", reason: ErrorReasonInvalidKey, want: "invalid_key"},
		{name: "value", reason: ErrorReasonInvalidValue, want: "invalid_value"},
		{name: "value count", reason: ErrorReasonInvalidValueCount, want: "invalid_value_count"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.reason); got != tt.want {
				t.Fatalf("reason = %q; want %q", got, tt.want)
			}
		})
	}
}
