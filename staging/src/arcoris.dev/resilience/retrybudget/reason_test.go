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

package retrybudget

import "testing"

func TestReasonString(t *testing.T) {
	tests := []struct {
		reason Reason
		want   string
	}{
		{ReasonUnknown, "unknown"},
		{ReasonAllowed, "allowed"},
		{ReasonExhausted, "exhausted"},
		{Reason(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.reason.String(); got != tt.want {
			t.Fatalf("%v.String() = %q, want %q", uint8(tt.reason), got, tt.want)
		}
	}
}

func TestReasonClassification(t *testing.T) {
	tests := []struct {
		name    string
		reason  Reason
		valid   bool
		allowed bool
		denied  bool
	}{
		{name: "unknown", reason: ReasonUnknown},
		{name: "allowed", reason: ReasonAllowed, valid: true, allowed: true},
		{name: "exhausted", reason: ReasonExhausted, valid: true, denied: true},
		{name: "unknown numeric", reason: Reason(99)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.reason.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tt.valid)
			}
			if got := tt.reason.IsAllowed(); got != tt.allowed {
				t.Fatalf("IsAllowed() = %v, want %v", got, tt.allowed)
			}
			if got := tt.reason.IsDenied(); got != tt.denied {
				t.Fatalf("IsDenied() = %v, want %v", got, tt.denied)
			}
		})
	}
}
