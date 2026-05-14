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

package healthtest

import (
	"testing"
	"time"

	"arcoris.dev/health"
)

func TestReportFixturesAreValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		report     health.Report
		wantStatus health.Status
	}{
		{"healthy", HealthyReport(health.TargetReady), health.StatusHealthy},
		{"starting", StartingReport(health.TargetStartup), health.StatusStarting},
		{"degraded", DegradedReport(health.TargetReady), health.StatusDegraded},
		{"unhealthy", UnhealthyReport(health.TargetLive), health.StatusUnhealthy},
		{"unknown", UnknownReport(health.TargetReady), health.StatusUnknown},
		{"mixed", MixedReport(health.TargetReady), health.StatusUnhealthy},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			AssertValidReport(t, tc.report)
			AssertReportStatus(t, tc.report, tc.wantStatus)
			if !tc.report.Observed.Equal(ObservedTime) {
				t.Fatalf("Observed = %v, want %v", tc.report.Observed, ObservedTime)
			}
		})
	}
}

func TestMixedReportContents(t *testing.T) {
	t.Parallel()

	report := MixedReport(health.TargetReady)

	AssertCheckOrder(t, report, "storage", "queue", "cache", "database")
	AssertReasons(
		t,
		report,
		health.ReasonOverloaded,
		health.ReasonNotObserved,
		health.ReasonDependencyUnavailable,
	)
	if report.Duration != 25*time.Millisecond {
		t.Fatalf("Duration = %s, want 25ms", report.Duration)
	}
	if report.Checks[3].Cause == nil || report.Checks[3].Cause.Error() != "private cause" {
		t.Fatalf("database cause = %v, want private cause", report.Checks[3].Cause)
	}
}

func TestUnknownReportForTargetUnknownUsesValidZeroShape(t *testing.T) {
	t.Parallel()

	report := UnknownReport(health.TargetUnknown)

	AssertValidReport(t, report)
	AssertReportTarget(t, report, health.TargetUnknown)
	AssertReportStatus(t, report, health.StatusUnknown)
	if report.IsObserved() {
		t.Fatalf("Observed = %v, want zero", report.Observed)
	}
	if len(report.Checks) != 0 {
		t.Fatalf("Checks length = %d, want 0", len(report.Checks))
	}
}

func TestReportReturnsDefensiveChecksCopy(t *testing.T) {
	t.Parallel()

	checks := []health.Result{HealthyResult("storage")}
	report := Report(health.TargetReady, health.StatusHealthy, checks...)
	checks[0] = UnhealthyResult("mutated", health.ReasonFatal)

	if report.Checks[0].Name != "storage" {
		t.Fatalf("report check mutated through input slice: %+v", report.Checks[0])
	}

	report.Checks[0] = UnhealthyResult("caller_mutated", health.ReasonFatal)
	again := MixedReport(health.TargetReady)
	if again.Checks[0].Name != "storage" {
		t.Fatalf("new fixture was affected by caller mutation: %+v", again.Checks[0])
	}
}
