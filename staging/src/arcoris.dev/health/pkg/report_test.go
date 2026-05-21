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
	"testing"
	"time"
)

func TestReportPredicates(t *testing.T) {
	t.Parallel()

	report := Report{
		Target:   TargetReady,
		Status:   StatusDegraded,
		Observed: testObserved,
		Duration: time.Second,
		Checks: []Result{
			Healthy("storage"),
			Degraded("queue", ReasonOverloaded, "queue overloaded"),
		},
	}

	if !report.IsValid() || !report.IsObserved() || report.Empty() {
		t.Fatalf("report predicates mismatch for %+v", report)
	}
	if !report.Passed(ReadyPolicy().WithDegraded(true)) {
		t.Fatal("degraded report should pass permissive ready policy")
	}
	if !report.Failed(ReadyPolicy()) {
		t.Fatal("degraded report should fail default ready policy")
	}
	if !report.HasReason(ReasonOverloaded) || report.HasReason(ReasonFatal) {
		t.Fatalf("reason predicates mismatch for %+v", report)
	}

	zero := Report{}
	if !zero.IsValid() || zero.IsObserved() || !zero.Empty() {
		t.Fatalf("zero report predicates mismatch for %+v", zero)
	}

	invalid := Report{Target: Target(99), Status: StatusHealthy, Duration: -time.Second}
	if invalid.IsValid() {
		t.Fatal("invalid report IsValid() = true, want false")
	}

	unknownWithChecks := Report{
		Target: TargetUnknown,
		Status: StatusHealthy,
		Checks: []Result{Healthy("storage")},
	}
	if unknownWithChecks.IsValid() {
		t.Fatal("TargetUnknown report with checks should be invalid")
	}

	invalidCheck := Report{
		Target: TargetReady,
		Status: StatusHealthy,
		Checks: []Result{{Status: StatusHealthy, Reason: Reason("bad-reason")}},
	}
	if invalidCheck.IsValid() {
		t.Fatal("report with invalid check reason should be invalid")
	}
}

func TestReportCheckAccessors(t *testing.T) {
	t.Parallel()

	report := Report{
		Target: TargetReady,
		Status: StatusUnknown,
		Checks: []Result{
			Healthy("storage"),
			Degraded("queue", ReasonOverloaded, "queue overloaded"),
			Unknown("cache", ReasonNotObserved, "cache unknown"),
			Unhealthy("database", ReasonFatal, "database failed"),
		},
	}

	res, ok := report.Check("queue")
	if !ok || res.Status != StatusDegraded {
		t.Fatalf("Check(queue) = %+v, %v; want degraded true", res, ok)
	}
	if _, ok := report.Check("missing"); ok {
		t.Fatal("Check(missing) ok = true, want false")
	}

	failed := report.FailedChecks(ReadyPolicy())
	if len(failed) != 3 {
		t.Fatalf("FailedChecks() = %d, want 3", len(failed))
	}
	if degraded := report.DegradedChecks(); len(degraded) != 1 || degraded[0].Name != "queue" {
		t.Fatalf("DegradedChecks() = %+v, want queue", degraded)
	}
	if unknown := report.UnknownChecks(); len(unknown) != 1 || unknown[0].Name != "cache" {
		t.Fatalf("UnknownChecks() = %+v, want cache", unknown)
	}
}

func TestReportReasonAccessors(t *testing.T) {
	t.Parallel()

	report := Report{
		Target: TargetReady,
		Status: StatusUnknown,
		Checks: []Result{
			Healthy("storage"),
			Degraded("queue", ReasonOverloaded, "queue overloaded"),
			Unhealthy("admission", ReasonAdmissionClosed, "admission closed"),
			Unknown("cache", ReasonNotObserved, "cache unknown"),
			Degraded("worker_pool", ReasonOverloaded, "worker pool overloaded"),
		},
	}

	overloaded := report.ChecksByReason(ReasonOverloaded)
	if len(overloaded) != 2 || overloaded[0].Name != "queue" || overloaded[1].Name != "worker_pool" {
		t.Fatalf("ChecksByReason(overloaded) = %+v, want queue and worker_pool", overloaded)
	}

	noReason := report.ChecksByReason(ReasonNone)
	if len(noReason) != 1 || noReason[0].Name != "storage" {
		t.Fatalf("ChecksByReason(none) = %+v, want storage", noReason)
	}

	reasons := report.Reasons()
	want := []Reason{ReasonOverloaded, ReasonAdmissionClosed, ReasonNotObserved}
	if len(reasons) != len(want) {
		t.Fatalf("Reasons() = %+v, want %+v", reasons, want)
	}
	for i := range want {
		if reasons[i] != want[i] {
			t.Fatalf("Reasons()[%d] = %s, want %s", i, reasons[i], want[i])
		}
	}

	if empty := (Report{}).Reasons(); empty != nil {
		t.Fatalf("empty Reasons() = %+v, want nil", empty)
	}
}

func TestReportChecksCopy(t *testing.T) {
	t.Parallel()

	report := Report{Target: TargetReady, Status: StatusHealthy, Checks: []Result{Healthy("storage"), Healthy("queue")}}

	copied := report.ChecksCopy()
	copied[0] = Unhealthy("storage", ReasonFatal, "fatal")

	if report.Checks[0].Status != StatusHealthy {
		t.Fatalf("report check mutated through copy = %s, want healthy", report.Checks[0].Status)
	}
	if empty := (Report{}).ChecksCopy(); empty != nil {
		t.Fatalf("empty ChecksCopy() = %+v, want nil", empty)
	}
}
