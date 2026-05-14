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

import (
	"strings"
	"testing"
)

func TestReasonString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reason Reason
		want   string
	}{
		{name: "none", reason: ReasonNone, want: "none"},
		{name: "builtin", reason: ReasonTimeout, want: "timeout"},
		{name: "custom", reason: Reason("custom_reason"), want: "custom_reason"},
		{name: "invalid", reason: Reason("Bad"), want: "invalid"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.reason.String(); got != tc.want {
				t.Fatalf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestReasonValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reason Reason
		want   bool
	}{
		{name: "none", reason: ReasonNone, want: true},
		{name: "custom", reason: Reason("custom_reason"), want: true},
		{name: "custom with digit", reason: Reason("custom_reason_1"), want: true},
		{name: "max length", reason: Reason(strings.Repeat("a", maxReasonLength)), want: true},
		{name: "starts with digit", reason: Reason("1custom"), want: false},
		{name: "starts with underscore", reason: Reason("_custom"), want: false},
		{name: "ends with underscore", reason: Reason("custom_"), want: false},
		{name: "repeated underscore", reason: Reason("custom__reason"), want: false},
		{name: "hyphen", reason: Reason("custom-reason"), want: false},
		{name: "dot", reason: Reason("custom.reason"), want: false},
		{name: "slash", reason: Reason("custom/reason"), want: false},
		{name: "space", reason: Reason("custom reason"), want: false},
		{name: "upper case", reason: Reason("CustomReason"), want: false},
		{name: "too long", reason: Reason(strings.Repeat("a", maxReasonLength+1)), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.reason.IsValid(); got != tc.want {
				t.Fatalf("IsValid() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestReasonBuiltins(t *testing.T) {
	t.Parallel()

	for _, reason := range builtinReasonsForTest() {
		t.Run(reason.String(), func(t *testing.T) {
			t.Parallel()

			if !reason.IsBuiltin() {
				t.Fatalf("%s.IsBuiltin() = false, want true", reason.String())
			}
			if !reason.IsValid() {
				t.Fatalf("%s.IsValid() = false, want true", reason.String())
			}
			if reason != ReasonNone && reason.String() == "invalid" {
				t.Fatalf("%s.String() = invalid, want builtin diagnostic string", string(reason))
			}
		})
	}

	if Reason("custom_reason").IsBuiltin() {
		t.Fatal("custom reason IsBuiltin() = true, want false")
	}
	if Reason("invalid-reason").IsBuiltin() {
		t.Fatal("invalid custom reason IsBuiltin() = true, want false")
	}
}

func TestReasonNone(t *testing.T) {
	t.Parallel()

	if !ReasonNone.IsNone() {
		t.Fatal("ReasonNone.IsNone() = false, want true")
	}

	for _, reason := range builtinReasonsForTest() {
		t.Run(reason.String(), func(t *testing.T) {
			t.Parallel()

			if reason != ReasonNone && reason.IsNone() {
				t.Fatalf("%s.IsNone() = true, want false", reason.String())
			}
		})
	}
}

func TestReasonObservationClassification(t *testing.T) {
	t.Parallel()
	assertReasonSet(t, "observation", Reason.IsObservationReason, []Reason{ReasonNotObserved, ReasonTimeout, ReasonCanceled, ReasonPanic})
}

func TestReasonExecutionClassification(t *testing.T) {
	t.Parallel()
	assertReasonSet(t, "execution", Reason.IsExecutionReason, []Reason{ReasonTimeout, ReasonCanceled, ReasonPanic})
}

func TestReasonLifecycleClassification(t *testing.T) {
	t.Parallel()
	assertReasonSet(t, "lifecycle", Reason.IsLifecycleReason, []Reason{ReasonStarting, ReasonDraining, ReasonShuttingDown})
}

func TestReasonDependencyClassification(t *testing.T) {
	t.Parallel()
	assertReasonSet(t, "dependency", Reason.IsDependencyReason, []Reason{ReasonDependencyUnavailable, ReasonDependencyDegraded})
}

func TestReasonControlClassification(t *testing.T) {
	t.Parallel()
	assertReasonSet(t, "control", Reason.IsControlReason, []Reason{ReasonOverloaded, ReasonBackpressured, ReasonRateLimited, ReasonAdmissionClosed, ReasonCapacityExhausted, ReasonResourceExhausted})
}

func TestReasonFreshnessClassification(t *testing.T) {
	t.Parallel()
	assertReasonSet(t, "freshness", Reason.IsFreshnessReason, []Reason{ReasonStale, ReasonNotSynced, ReasonSyncFailed, ReasonLagging})
}

func TestReasonConnectivityClassification(t *testing.T) {
	t.Parallel()
	assertReasonSet(t, "connectivity", Reason.IsConnectivityReason, []Reason{ReasonPartitioned})
}

func TestReasonConfigurationClassification(t *testing.T) {
	t.Parallel()
	assertReasonSet(t, "configuration", Reason.IsConfigurationReason, []Reason{ReasonMisconfigured})
}

func builtinReasonsForTest() []Reason {
	return []Reason{
		ReasonNone,
		ReasonNotObserved,
		ReasonTimeout,
		ReasonCanceled,
		ReasonPanic,
		ReasonStarting,
		ReasonDraining,
		ReasonShuttingDown,
		ReasonDependencyUnavailable,
		ReasonDependencyDegraded,
		ReasonOverloaded,
		ReasonBackpressured,
		ReasonRateLimited,
		ReasonAdmissionClosed,
		ReasonCapacityExhausted,
		ReasonResourceExhausted,
		ReasonStale,
		ReasonNotSynced,
		ReasonSyncFailed,
		ReasonLagging,
		ReasonPartitioned,
		ReasonMisconfigured,
		ReasonFatal,
	}
}

func assertReasonSet(t *testing.T, name string, predicate func(Reason) bool, want []Reason) {
	t.Helper()

	wanted := make(map[Reason]struct{}, len(want))
	for _, reason := range want {
		wanted[reason] = struct{}{}
		if !predicate(reason) {
			t.Fatalf("%s predicate for %s = false, want true", name, reason.String())
		}
	}

	for _, reason := range builtinReasonsForTest() {
		_, shouldMatch := wanted[reason]
		if got := predicate(reason); got != shouldMatch {
			t.Fatalf("%s predicate for %s = %v, want %v", name, reason.String(), got, shouldMatch)
		}
	}

	custom := Reason(name + "_custom")
	if predicate(custom) {
		t.Fatalf("%s predicate for custom reason %s = true, want false", name, custom)
	}

	invalid := Reason(name + "-invalid")
	if predicate(invalid) {
		t.Fatalf("%s predicate for invalid reason %s = true, want false", name, invalid)
	}
}
