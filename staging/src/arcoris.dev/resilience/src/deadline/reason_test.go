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

package deadline

import "testing"

func TestReasonString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		reason Reason
		want   string
	}{
		{ReasonAllowed, "allowed"},
		{ReasonContextDone, "context_done"},
		{ReasonNoDeadline, "no_deadline"},
		{ReasonExpired, "expired"},
		{ReasonInsufficientBudget, "insufficient_budget"},
		{Reason(255), "unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			if got := tt.reason.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReasonValidationHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reason  Reason
		valid   bool
		allowed bool
		denied  bool
	}{
		{name: "allowed", reason: ReasonAllowed, valid: true, allowed: true},
		{name: "context done", reason: ReasonContextDone, valid: true, denied: true},
		{name: "no deadline", reason: ReasonNoDeadline, valid: true, allowed: true},
		{name: "expired", reason: ReasonExpired, valid: true, denied: true},
		{name: "insufficient budget", reason: ReasonInsufficientBudget, valid: true, denied: true},
		{name: "unknown", reason: Reason(255)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.reason.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tt.valid)
			}
			if got := tt.reason.IsAllowedReason(); got != tt.allowed {
				t.Fatalf("IsAllowedReason() = %v, want %v", got, tt.allowed)
			}
			if got := tt.reason.IsDeniedReason(); got != tt.denied {
				t.Fatalf("IsDeniedReason() = %v, want %v", got, tt.denied)
			}
		})
	}
}
