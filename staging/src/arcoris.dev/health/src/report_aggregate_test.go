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

package health

import "testing"

func TestReportAggregateHelpersRepairStaleStatus(t *testing.T) {
	t.Parallel()

	report := Report{
		Target: TargetReady,
		Status: StatusHealthy,
		Checks: []Result{
			Healthy("storage"),
			Unhealthy("database", ReasonFatal, "fatal"),
		},
	}

	if report.IsConsistent() {
		t.Fatal("stale report should be inconsistent")
	}
	if got := report.AggregateStatus(); got != StatusUnhealthy {
		t.Fatalf("AggregateStatus() = %s, want unhealthy", got)
	}
	if repaired := report.WithAggregateStatus(); !repaired.IsConsistent() {
		t.Fatalf("WithAggregateStatus() = %+v, want consistent", repaired)
	}
}

func TestAggregateStatusSeverityEdges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		checks []Result
		want   Status
	}{
		{name: "empty", checks: nil, want: StatusUnknown},
		{name: "healthy", checks: []Result{Healthy("storage")}, want: StatusHealthy},
		{name: "starting", checks: []Result{Healthy("storage"), Starting("boot", ReasonStarting, "starting")}, want: StatusStarting},
		{name: "degraded", checks: []Result{Starting("boot", ReasonStarting, "starting"), Degraded("queue", ReasonOverloaded, "overloaded")}, want: StatusDegraded},
		{name: "unknown", checks: []Result{Degraded("queue", ReasonOverloaded, "overloaded"), Unknown("cache", ReasonNotObserved, "unknown")}, want: StatusUnknown},
		{name: "unhealthy", checks: []Result{Unknown("cache", ReasonNotObserved, "unknown"), Unhealthy("db", ReasonFatal, "fatal")}, want: StatusUnhealthy},
		{name: "invalid", checks: []Result{Unhealthy("db", ReasonFatal, "fatal"), {Name: "bad", Status: Status(99)}}, want: Status(99)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := AggregateStatus(tc.checks); got != tc.want {
				t.Fatalf("AggregateStatus() = %s, want %s", got, tc.want)
			}
		})
	}
}

func TestReportConsistencyEdges(t *testing.T) {
	t.Parallel()

	if !(Report{}).IsConsistent() {
		t.Fatal("zero report should be consistent")
	}
	if !(Report{Target: TargetReady, Status: StatusUnknown}).IsConsistent() {
		t.Fatal("empty concrete unknown report should be consistent")
	}
	if (Report{Target: TargetReady, Status: StatusHealthy}).IsConsistent() {
		t.Fatal("empty concrete healthy report should be inconsistent")
	}
}
