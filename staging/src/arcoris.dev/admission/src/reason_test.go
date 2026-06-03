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

package admission

import "testing"

func TestReasonIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reason Reason
		want   bool
	}{
		{name: "admitted", reason: ReasonAdmitted, want: true},
		{name: "capacity exhausted", reason: Reason("capacity_exhausted"), want: true},
		{name: "budget exhausted", reason: Reason("budget_exhausted"), want: true},
		{name: "custom stable reason", reason: "tenant_quota_exhausted", want: true},
		{name: "custom with version", reason: "scheduler_unschedulable_v2", want: true},
		{name: "empty", reason: "", want: false},
		{name: "uppercase", reason: "CapacityExhausted", want: false},
		{name: "hyphen", reason: "capacity-exhausted", want: false},
		{name: "leading underscore", reason: "_capacity", want: false},
		{name: "trailing underscore", reason: "capacity_", want: false},
		{name: "repeated underscore", reason: "capacity__exhausted", want: false},
		{name: "leading digit", reason: "123capacity", want: false},
		{name: "too long", reason: Reason(string(make([]byte, maxReasonLength+1))), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.reason.IsValid(); got != tt.want {
				t.Fatalf("%q IsValid = %v, want %v", tt.reason, got, tt.want)
			}
		})
	}
}

func TestReasonString(t *testing.T) {
	t.Parallel()

	if got := Reason("capacity_exhausted").String(); got != "capacity_exhausted" {
		t.Fatalf("String = %q, want capacity_exhausted", got)
	}
}
