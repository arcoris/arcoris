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

func TestDefaultPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		target Target
		want   TargetPolicy
	}{
		{TargetStartup, StartupPolicy()},
		{TargetLive, LivePolicy()},
		{TargetReady, ReadyPolicy()},
		{TargetUnknown, TargetPolicy{}},
		{Target(99), TargetPolicy{}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.target.String(), func(t *testing.T) {
			t.Parallel()

			if got := DefaultPolicy(test.target); got != test.want {
				t.Fatalf("DefaultPolicy() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestTargetPolicyPassesAndAllows(t *testing.T) {
	t.Parallel()

	policy := TargetPolicy{}.WithStarting(true).WithDegraded(true)

	if !policy.Passes(StatusHealthy) {
		t.Fatal("healthy should pass")
	}
	if !policy.Passes(StatusStarting) || !policy.Allows(StatusStarting) {
		t.Fatal("starting should be allowed")
	}
	if !policy.Passes(StatusDegraded) || !policy.Allows(StatusDegraded) {
		t.Fatal("degraded should be allowed")
	}
	if policy.Passes(StatusUnknown) || policy.Passes(StatusUnhealthy) || policy.Passes(Status(99)) {
		t.Fatal("unknown, unhealthy, and invalid statuses should fail")
	}
	if policy.Allows(StatusHealthy) {
		t.Fatal("healthy passes but is not explicitly allowed")
	}
	if !policy.Fails(StatusUnhealthy) {
		t.Fatal("unhealthy should fail")
	}
}
