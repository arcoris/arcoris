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

	"arcoris.dev/health"
)

// AssertReportStatus fails if report has a different status.
//
// This is a health-domain assertion, not a generic equality helper. The failure
// message renders status names, which is more useful in health tests than raw
// numeric values.
func AssertReportStatus(t testing.TB, report health.Report, want health.Status) {
	t.Helper()

	if report.Status != want {
		t.Fatalf("report status = %s, want %s", report.Status, want)
	}
}

// AssertReportTarget fails if report has a different target.
//
// Target assertions are common in adapter tests because target/report mismatch
// bugs can otherwise be hidden behind status-only checks.
func AssertReportTarget(t testing.TB, report health.Report, want health.Target) {
	t.Helper()

	if report.Target != want {
		t.Fatalf("report target = %s, want %s", report.Target, want)
	}
}

// AssertCheckOrder fails if report checks do not have names in order.
//
// Registry order is part of package health's public contract and many adapters
// preserve that order in their DTOs. This assertion keeps that intent explicit.
func AssertCheckOrder(t testing.TB, report health.Report, names ...string) {
	t.Helper()

	if len(report.Checks) != len(names) {
		t.Fatalf("check count = %d, want %d", len(report.Checks), len(names))
	}
	for i, name := range names {
		if report.Checks[i].Name != name {
			t.Fatalf("check[%d] name = %q, want %q", i, report.Checks[i].Name, name)
		}
	}
}

// AssertReasons fails if report reasons do not match reasons in order.
//
// The assertion uses health.Report.Reasons, so ReasonNone is omitted just as the
// core report helper omits it.
func AssertReasons(t testing.TB, report health.Report, reasons ...health.Reason) {
	t.Helper()

	got := report.Reasons()
	if len(got) != len(reasons) {
		t.Fatalf("reason count = %d, want %d: got %v want %v", len(got), len(reasons), got, reasons)
	}
	for i, reason := range reasons {
		if got[i] != reason {
			t.Fatalf("reason[%d] = %s, want %s", i, got[i], reason)
		}
	}
}

// AssertValidReport fails if report is structurally invalid.
//
// Structural validity belongs to package health; this helper only gives tests a
// precise fixture-level failure message.
func AssertValidReport(t testing.TB, report health.Report) {
	t.Helper()

	if !report.IsValid() {
		t.Fatalf("report is invalid: %+v", report)
	}
}

// AssertInvalidReport fails if report is structurally valid.
//
// Use this only when invalid structure is the behavior under test. Normal
// fixtures in healthtest are intended to be valid.
func AssertInvalidReport(t testing.TB, report health.Report) {
	t.Helper()

	if report.IsValid() {
		t.Fatalf("report is valid, want invalid: %+v", report)
	}
}

// AssertResultStatus fails if result has a different status.
//
// Result status is asserted directly instead of through report aggregation when
// tests focus on checker or DTO conversion behavior.
func AssertResultStatus(t testing.TB, res health.Result, want health.Status) {
	t.Helper()

	if res.Status != want {
		t.Fatalf("result status = %s, want %s", res.Status, want)
	}
}

// AssertResultReason fails if result has a different reason.
//
// Reason is asserted separately from message and cause because only Reason is
// stable, machine-readable, and safe for public adapter logic.
func AssertResultReason(t testing.TB, res health.Result, want health.Reason) {
	t.Helper()

	if res.Reason != want {
		t.Fatalf("result reason = %s, want %s", res.Reason, want)
	}
}
