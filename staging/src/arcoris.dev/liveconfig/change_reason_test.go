/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package liveconfig

import "testing"

func TestChangeReasonString(t *testing.T) {
	tests := []struct {
		name   string
		reason ChangeReason
		want   string
	}{
		{name: "unknown", reason: ChangeReasonUnknown, want: "unknown"},
		{name: "published", reason: ChangeReasonPublished, want: "published"},
		{name: "equal", reason: ChangeReasonEqual, want: "equal"},
		{name: "normalize failed", reason: ChangeReasonNormalizeFailed, want: "normalize_failed"},
		{name: "validate failed", reason: ChangeReasonValidateFailed, want: "validate_failed"},
		{name: "unknown numeric value", reason: ChangeReason(99), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reason.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestChangeReasonState(t *testing.T) {
	tests := []struct {
		name      string
		reason    ChangeReason
		valid     bool
		accepted  bool
		rejected  bool
		published bool
		equal     bool
	}{
		{
			name:   "unknown",
			reason: ChangeReasonUnknown,
		},
		{
			name:      "published",
			reason:    ChangeReasonPublished,
			valid:     true,
			accepted:  true,
			published: true,
		},
		{
			name:     "equal",
			reason:   ChangeReasonEqual,
			valid:    true,
			accepted: true,
			equal:    true,
		},
		{
			name:     "normalize failed",
			reason:   ChangeReasonNormalizeFailed,
			valid:    true,
			rejected: true,
		},
		{
			name:     "validate failed",
			reason:   ChangeReasonValidateFailed,
			valid:    true,
			rejected: true,
		},
		{
			name:   "unknown numeric value",
			reason: ChangeReason(99),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reason.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tt.valid)
			}
			if got := tt.reason.Accepted(); got != tt.accepted {
				t.Fatalf("Accepted() = %v, want %v", got, tt.accepted)
			}
			if got := tt.reason.Rejected(); got != tt.rejected {
				t.Fatalf("Rejected() = %v, want %v", got, tt.rejected)
			}
			if got := tt.reason.Published(); got != tt.published {
				t.Fatalf("Published() = %v, want %v", got, tt.published)
			}
			if got := tt.reason.Equal(); got != tt.equal {
				t.Fatalf("Equal() = %v, want %v", got, tt.equal)
			}
		})
	}
}
