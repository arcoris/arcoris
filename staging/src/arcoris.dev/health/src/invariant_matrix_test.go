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

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestValidCheckNameBoundaryMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"empty", "", false},
		{"one letter", "a", true},
		{"lower snake", "database_pool", true},
		{"digits after first", "database_pool_1", true},
		{"max length", strings.Repeat("a", maxCheckNameLength), true},
		{"too long", strings.Repeat("a", maxCheckNameLength+1), false},
		{"starts digit", "1database", false},
		{"starts underscore", "_database", false},
		{"trailing underscore", "database_", false},
		{"double underscore", "database__pool", false},
		{"uppercase", "Database", false},
		{"hyphen", "database-pool", false},
		{"dot", "database.pool", false},
		{"slash", "database/pool", false},
		{"space", "database pool", false},
		{"tab", "database\tpool", false},
		{"newline", "database\npool", false},
		{"non ascii", "databasé", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := ValidCheckName(tc.in); got != tc.want {
				t.Fatalf("ValidCheckName(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestReasonValidationBoundaryMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reason Reason
		want   bool
	}{
		{"none", ReasonNone, true},
		{"one letter", Reason("a"), true},
		{"lower snake", Reason("custom_reason"), true},
		{"digits after first", Reason("custom_reason_1"), true},
		{"max length", Reason(strings.Repeat("a", maxReasonLength)), true},
		{"too long", Reason(strings.Repeat("a", maxReasonLength+1)), false},
		{"starts digit", Reason("1custom"), false},
		{"starts underscore", Reason("_custom"), false},
		{"trailing underscore", Reason("custom_"), false},
		{"double underscore", Reason("custom__reason"), false},
		{"uppercase", Reason("Custom"), false},
		{"hyphen", Reason("custom-reason"), false},
		{"dot", Reason("custom.reason"), false},
		{"slash", Reason("custom/reason"), false},
		{"space", Reason("custom reason"), false},
		{"tab", Reason("custom\treason"), false},
		{"newline", Reason("custom\nreason"), false},
		{"non ascii", Reason("customé"), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.reason.IsValid(); got != tc.want {
				t.Fatalf("IsValid() = %v, want %v for %q", got, tc.want, string(tc.reason))
			}
		})
	}
}

func TestReasonCustomValidButNotBuiltin(t *testing.T) {
	t.Parallel()

	reason := Reason("tenant_123")
	if !reason.IsValid() {
		t.Fatal("custom stable reason should be syntactically valid")
	}
	if reason.IsBuiltin() {
		t.Fatal("custom reason IsBuiltin() = true, want false")
	}
}

func TestConcreteTargetsReturnsFreshDeterministicSlice(t *testing.T) {
	t.Parallel()

	first := ConcreteTargets()
	second := ConcreteTargets()
	want := []Target{TargetStartup, TargetLive, TargetReady}
	if len(first) != len(want) || len(second) != len(want) {
		t.Fatalf("ConcreteTargets lengths = %d/%d, want %d", len(first), len(second), len(want))
	}
	for i := range want {
		if first[i] != want[i] || second[i] != want[i] {
			t.Fatalf("ConcreteTargets order = %v and %v, want %v", first, second, want)
		}
	}
	first[0] = TargetUnknown
	if got := ConcreteTargets()[0]; got != TargetStartup {
		t.Fatalf("ConcreteTargets()[0] after mutation = %s, want startup", got)
	}
}

func TestResultNormalizeMatrix(t *testing.T) {
	t.Parallel()

	cause := errors.New("private cause")
	existingObserved := testObserved.Add(-time.Minute)
	res := Result{
		Name:     "custom",
		Status:   Status(99),
		Reason:   ReasonFatal,
		Message:  "message",
		Cause:    cause,
		Observed: existingObserved,
		Duration: -time.Second,
	}.Normalize("default", testObserved)

	if res.Name != "custom" || res.Status != StatusUnknown || res.Observed != existingObserved || res.Duration != 0 {
		t.Fatalf("Normalize() = %+v, want preserved name/observed, unknown status, zero duration", res)
	}
	if res.Reason != ReasonFatal || res.Message != "message" || res.Cause != cause {
		t.Fatalf("Normalize() did not preserve reason/message/cause: %+v", res)
	}

	filled := Result{Status: StatusHealthy}.Normalize("default", testObserved)
	if filled.Name != "default" || filled.Observed != testObserved || filled.Duration != 0 {
		t.Fatalf("Normalize(empty fields) = %+v, want filled name/observed", filled)
	}
}

func TestResultValidityMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		res  Result
		want bool
	}{
		{"zero", Result{}, true},
		{"invalid status", Result{Status: Status(99)}, false},
		{"invalid reason", Result{Status: StatusHealthy, Reason: Reason("bad-reason")}, false},
		{"negative duration", Result{Status: StatusHealthy, Duration: -time.Nanosecond}, false},
		{"named valid", Healthy("storage"), true},
		{"named invalid", Healthy("bad-name"), false},
		{"unnamed valid", Result{Status: StatusHealthy}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.res.IsValid(); got != tc.want {
				t.Fatalf("IsValid() = %v, want %v for %+v", got, tc.want, tc.res)
			}
		})
	}
}

func TestReportValidityMatrix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		report Report
		want   bool
	}{
		{"zero", Report{}, true},
		{"observed concrete", Report{Target: TargetReady, Status: StatusHealthy, Observed: testObserved}, true},
		{"unknown with checks", Report{Target: TargetUnknown, Status: StatusUnknown, Checks: []Result{Healthy("storage")}}, false},
		{"invalid status", Report{Target: TargetReady, Status: Status(99)}, false},
		{"invalid check", Report{Target: TargetReady, Status: StatusHealthy, Checks: []Result{{Status: StatusHealthy, Reason: Reason("bad-reason")}}}, false},
		{"negative duration", Report{Target: TargetReady, Status: StatusHealthy, Duration: -time.Nanosecond}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.report.IsValid(); got != tc.want {
				t.Fatalf("IsValid() = %v, want %v for %+v", got, tc.want, tc.report)
			}
		})
	}
}

func TestReportChecksCopyDetached(t *testing.T) {
	t.Parallel()

	report := Report{Target: TargetReady, Status: StatusHealthy, Checks: []Result{Healthy("first")}}
	copied := report.ChecksCopy()
	copied[0] = Unhealthy("first", ReasonFatal, "fatal")
	if report.Checks[0].Status != StatusHealthy {
		t.Fatalf("report mutated through ChecksCopy: %+v", report.Checks[0])
	}
}

func TestReportReasonsUniqueFirstSeenOrder(t *testing.T) {
	t.Parallel()

	report := Report{Target: TargetReady, Status: StatusUnhealthy, Checks: []Result{
		Healthy("storage"),
		Degraded("queue", ReasonOverloaded, "overloaded"),
		Unknown("cache", ReasonNotObserved, "unknown"),
		Unhealthy("database", ReasonOverloaded, "overloaded again"),
	}}
	want := []Reason{ReasonOverloaded, ReasonNotObserved}
	got := report.Reasons()
	if len(got) != len(want) {
		t.Fatalf("Reasons() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Reasons()[%d] = %s, want %s", i, got[i], want[i])
		}
	}
}

func TestReportFilterHelpersPreserveReportOrder(t *testing.T) {
	t.Parallel()

	report := Report{Target: TargetReady, Status: StatusUnhealthy, Checks: []Result{
		Unknown("first", ReasonNotObserved, "unknown"),
		Degraded("second", ReasonOverloaded, "degraded"),
		Unhealthy("third", ReasonFatal, "fatal"),
		Degraded("fourth", ReasonBackpressured, "degraded"),
	}}
	if got := namesOfResults(report.DegradedChecks()); strings.Join(got, ",") != "second,fourth" {
		t.Fatalf("DegradedChecks order = %v", got)
	}
	if got := namesOfResults(report.FailedChecks(ReadyPolicy())); strings.Join(got, ",") != "first,second,third,fourth" {
		t.Fatalf("FailedChecks order = %v", got)
	}
	if got := namesOfResults(report.UnknownChecks()); strings.Join(got, ",") != "first" {
		t.Fatalf("UnknownChecks order = %v", got)
	}
}

func TestReportStatusAggregationPolicy(t *testing.T) {
	t.Parallel()

	results := []Result{
		Healthy("healthy"),
		Starting("starting", ReasonStarting, "starting"),
		Degraded("degraded", ReasonOverloaded, "degraded"),
		Unknown("unknown", ReasonNotObserved, "unknown"),
		Unhealthy("unhealthy", ReasonFatal, "fatal"),
	}
	status := StatusHealthy
	for _, result := range results {
		if result.Status.MoreSevereThan(status) {
			status = result.Status
		}
	}
	if status != StatusUnhealthy {
		t.Fatalf("aggregated status = %s, want unhealthy", status)
	}
}

func namesOfResults(results []Result) []string {
	names := make([]string, len(results))
	for i, result := range results {
		names[i] = result.Name
	}
	return names
}
