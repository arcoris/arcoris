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
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
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

func TestRegistryRegisterEmptyBatchNoop(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	if err := registry.Register(TargetReady); err != nil {
		t.Fatalf("Register(empty) = %v, want nil", err)
	}
	if !registry.Empty() || registry.Len(TargetReady) != 0 {
		t.Fatal("empty registration mutated registry")
	}
}

func TestRegistryConcurrentRegisterAndReadRaceFree(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			name := fmt.Sprintf("check_%d", i)
			_ = registry.Register(TargetReady, mustCheck(t, name, Healthy(name)))
			_ = registry.Checks(TargetReady)
			_ = registry.Targets()
		}()
	}
	wg.Wait()
}

func TestRegistryBatchDuplicateCarriesPreviousIndex(t *testing.T) {
	t.Parallel()

	registry := NewRegistry()
	err := registry.Register(TargetReady, mustCheck(t, "dup", Healthy("dup")), mustCheck(t, "dup", Healthy("dup")))
	var duplicate DuplicateCheckError
	if !errors.As(err, &duplicate) {
		t.Fatalf("Register() = %v, want DuplicateCheckError", err)
	}
	if duplicate.Index != 1 || duplicate.PreviousIndex != 0 {
		t.Fatalf("duplicate = %+v, want index 1 previous 0", duplicate)
	}
}

func TestRegistryExistingDuplicateCarriesPreviousIndexMinusOne(t *testing.T) {
	t.Parallel()

	registry := mustRegistry(t, TargetReady, mustCheck(t, "dup", Healthy("dup")))
	err := registry.Register(TargetReady, mustCheck(t, "dup", Healthy("dup")))
	var duplicate DuplicateCheckError
	if !errors.As(err, &duplicate) {
		t.Fatalf("Register() = %v, want DuplicateCheckError", err)
	}
	if duplicate.Index != 0 || duplicate.PreviousIndex != -1 {
		t.Fatalf("duplicate = %+v, want index 0 previous -1", duplicate)
	}
}

func TestRegistryFailedExistingConflictDoesNotRegisterNonConflictingChecks(t *testing.T) {
	t.Parallel()

	registry := mustRegistry(t, TargetReady, mustCheck(t, "existing", Healthy("existing")))
	err := registry.Register(TargetReady, mustCheck(t, "existing", Healthy("existing")), mustCheck(t, "new_check", Healthy("new_check")))
	if !errors.Is(err, ErrDuplicateCheck) {
		t.Fatalf("Register() = %v, want ErrDuplicateCheck", err)
	}
	if registry.Has(TargetReady, "new_check") {
		t.Fatal("conflicting batch registered non-conflicting check")
	}
}

func TestRegistryRejectsTypedNilChecker(t *testing.T) {
	t.Parallel()

	var checker *typedNilChecker
	err := NewRegistry().Register(TargetReady, checker)
	if !errors.Is(err, ErrNilChecker) {
		t.Fatalf("Register(typed nil) = %v, want ErrNilChecker", err)
	}
}

func TestGateSetInvalidResultLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	gate, err := NewGate("ready_gate", Healthy("ready_gate"))
	if err != nil {
		t.Fatalf("NewGate() = %v, want nil", err)
	}
	if err := gate.Set(Result{Status: Status(99)}); !errors.Is(err, ErrInvalidGateResult) {
		t.Fatalf("Set(invalid) = %v, want ErrInvalidGateResult", err)
	}
	if got := gate.Check(context.Background()).Status; got != StatusHealthy {
		t.Fatalf("status after failed Set = %s, want healthy", got)
	}
}

func TestGateSetMismatchedNameLeavesStateUnchanged(t *testing.T) {
	t.Parallel()

	gate, err := NewGate("ready_gate", Healthy("ready_gate"))
	if err != nil {
		t.Fatalf("NewGate() = %v, want nil", err)
	}
	if err := gate.Set(Healthy("other_gate")); !errors.Is(err, ErrMismatchedGateResult) {
		t.Fatalf("Set(mismatch) = %v, want ErrMismatchedGateResult", err)
	}
	if got := gate.Check(context.Background()).Name; got != "ready_gate" {
		t.Fatalf("name after failed Set = %q, want ready_gate", got)
	}
}

func TestGateNilBehaviorMatrix(t *testing.T) {
	t.Parallel()

	var gate *Gate
	mutations := []func() error{
		func() error { return gate.Set(Healthy("ready_gate")) },
		gate.Healthy,
		func() error { return gate.Unknown(ReasonNotObserved, "unknown") },
		func() error { return gate.Starting(ReasonStarting, "starting") },
		func() error { return gate.Degraded(ReasonOverloaded, "degraded") },
		func() error { return gate.Unhealthy(ReasonFatal, "fatal") },
	}
	for i, mutate := range mutations {
		if err := mutate(); !errors.Is(err, ErrNilChecker) {
			t.Fatalf("mutation %d = %v, want ErrNilChecker", i, err)
		}
	}
	if res := gate.Check(context.Background()); !errors.Is(res.Cause, ErrNilChecker) || res.Status != StatusUnknown {
		t.Fatalf("nil Check() = %+v, want unknown ErrNilChecker", res)
	}
}

func TestGateConcurrentSetAndCheckRaceFree(t *testing.T) {
	t.Parallel()

	gate, err := NewUnknownGate("ready_gate")
	if err != nil {
		t.Fatalf("NewUnknownGate() = %v, want nil", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				_ = gate.Healthy()
			} else {
				_ = gate.Degraded(ReasonOverloaded, "overloaded")
			}
			_ = gate.Check(context.Background())
		}(i)
	}
	wg.Wait()
}

func TestGateStatusHelpersPublishExpectedResults(t *testing.T) {
	t.Parallel()

	gate, err := NewUnknownGate("ready_gate")
	if err != nil {
		t.Fatalf("NewUnknownGate() = %v, want nil", err)
	}
	tests := []struct {
		name   string
		set    func() error
		status Status
		reason Reason
	}{
		{"healthy", gate.Healthy, StatusHealthy, ReasonNone},
		{"unknown", func() error { return gate.Unknown(ReasonNotObserved, "unknown") }, StatusUnknown, ReasonNotObserved},
		{"starting", func() error { return gate.Starting(ReasonStarting, "starting") }, StatusStarting, ReasonStarting},
		{"degraded", func() error { return gate.Degraded(ReasonOverloaded, "degraded") }, StatusDegraded, ReasonOverloaded},
		{"unhealthy", func() error { return gate.Unhealthy(ReasonFatal, "fatal") }, StatusUnhealthy, ReasonFatal},
	}
	for _, tc := range tests {
		if err := tc.set(); err != nil {
			t.Fatalf("%s set = %v, want nil", tc.name, err)
		}
		res := gate.Check(context.Background())
		if res.Name != "ready_gate" || res.Status != tc.status || res.Reason != tc.reason {
			t.Fatalf("%s result = %+v, want status %s reason %s", tc.name, res, tc.status, tc.reason)
		}
	}
}

func namesOfResults(results []Result) []string {
	names := make([]string, len(results))
	for i, result := range results {
		names[i] = result.Name
	}
	return names
}
