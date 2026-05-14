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

package health

import "testing"

func TestStatusString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status Status
		want   string
	}{
		{StatusUnknown, "unknown"},
		{StatusStarting, "starting"},
		{StatusHealthy, "healthy"},
		{StatusDegraded, "degraded"},
		{StatusUnhealthy, "unhealthy"},
		{Status(99), "invalid"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()

			if got := tc.status.String(); got != tc.want {
				t.Fatalf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestStatusClassification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status      Status
		valid       bool
		affirmative bool
		negative    bool
		known       bool
		operational bool
	}{
		{StatusUnknown, true, false, false, false, false},
		{StatusStarting, true, false, false, true, true},
		{StatusHealthy, true, true, false, true, true},
		{StatusDegraded, true, false, false, true, true},
		{StatusUnhealthy, true, false, true, true, false},
		{Status(99), false, false, false, false, false},
	}

	for _, tc := range tests {
		t.Run(tc.status.String(), func(t *testing.T) {
			t.Parallel()

			if got := tc.status.IsValid(); got != tc.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tc.valid)
			}
			if got := tc.status.IsAffirmative(); got != tc.affirmative {
				t.Fatalf("IsAffirmative() = %v, want %v", got, tc.affirmative)
			}
			if got := tc.status.IsNegative(); got != tc.negative {
				t.Fatalf("IsNegative() = %v, want %v", got, tc.negative)
			}
			if got := tc.status.IsKnown(); got != tc.known {
				t.Fatalf("IsKnown() = %v, want %v", got, tc.known)
			}
			if got := tc.status.IsOperational(); got != tc.operational {
				t.Fatalf("IsOperational() = %v, want %v", got, tc.operational)
			}
		})
	}
}

func TestStatusMoreSevereThan(t *testing.T) {
	t.Parallel()

	if !StatusUnhealthy.MoreSevereThan(StatusUnknown) {
		t.Fatal("unhealthy should be more severe than unknown")
	}
	if !StatusUnknown.MoreSevereThan(StatusDegraded) {
		t.Fatal("unknown should be more severe than degraded")
	}
	if !Status(99).MoreSevereThan(StatusUnhealthy) {
		t.Fatal("invalid should be more severe than valid statuses")
	}
	if StatusHealthy.MoreSevereThan(StatusStarting) {
		t.Fatal("healthy should not be more severe than starting")
	}
}
