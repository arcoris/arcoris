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
	"context"
	"errors"
	"sync"
	"testing"

	"arcoris.dev/component-base/pkg/health"
)

func TestSourceFunc(t *testing.T) {
	t.Parallel()

	source := SourceFunc(func(_ context.Context, target health.Target) (health.Report, error) {
		return HealthyReport(target), nil
	})

	report, err := source.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	AssertReportTarget(t, report, health.TargetReady)
	AssertReportStatus(t, report, health.StatusHealthy)
}

func TestStaticSourceReturnsDefensiveReports(t *testing.T) {
	t.Parallel()

	report := HealthyReport(health.TargetReady)
	source := NewStaticSource(report)
	report.Checks[0] = UnhealthyResult("mutated", health.ReasonFatal)

	first, err := source.Evaluate(context.Background(), health.TargetLive)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if first.Checks[0].Name != "ready_check" {
		t.Fatalf("first check name = %q, want ready_check", first.Checks[0].Name)
	}

	first.Checks[0] = UnhealthyResult("mutated_again", health.ReasonFatal)
	second, err := source.Evaluate(context.Background(), health.TargetLive)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if second.Checks[0].Name != "ready_check" {
		t.Fatalf("second check name = %q, want ready_check", second.Checks[0].Name)
	}
}

func TestTargetSource(t *testing.T) {
	t.Parallel()

	source := NewTargetSource(map[health.Target]health.Report{
		health.TargetReady: HealthyReport(health.TargetReady),
		health.TargetLive:  UnhealthyReport(health.TargetLive),
	})

	ready, err := source.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate(ready) = %v, want nil", err)
	}
	AssertReportStatus(t, ready, health.StatusHealthy)

	live, err := source.Evaluate(context.Background(), health.TargetLive)
	if err != nil {
		t.Fatalf("Evaluate(live) = %v, want nil", err)
	}
	AssertReportStatus(t, live, health.StatusUnhealthy)

	startup, err := source.Evaluate(context.Background(), health.TargetStartup)
	if err != nil {
		t.Fatalf("Evaluate(startup) = %v, want nil", err)
	}
	AssertReportStatus(t, startup, health.StatusUnknown)
	AssertReportTarget(t, startup, health.TargetStartup)
}

func TestErrorSource(t *testing.T) {
	t.Parallel()

	raw := errors.New("raw source failure")
	source := NewErrorSource(raw)

	_, err := source.Evaluate(context.Background(), health.TargetReady)
	if !errors.Is(err, raw) {
		t.Fatalf("Evaluate() = %v, want raw source failure", err)
	}
}

func TestCountingSource(t *testing.T) {
	t.Parallel()

	source := NewCountingSource(map[health.Target]health.Report{
		health.TargetReady: HealthyReport(health.TargetReady),
	})
	source.SetReport(health.TargetLive, DegradedReport(health.TargetLive))

	for i := 0; i < 2; i++ {
		report, err := source.Evaluate(context.Background(), health.TargetReady)
		if err != nil {
			t.Fatalf("Evaluate(ready) = %v, want nil", err)
		}
		AssertReportStatus(t, report, health.StatusHealthy)
	}
	live, err := source.Evaluate(context.Background(), health.TargetLive)
	if err != nil {
		t.Fatalf("Evaluate(live) = %v, want nil", err)
	}
	AssertReportStatus(t, live, health.StatusDegraded)

	if source.Calls(health.TargetReady) != 2 {
		t.Fatalf("ready calls = %d, want 2", source.Calls(health.TargetReady))
	}
	if source.Calls(health.TargetLive) != 1 {
		t.Fatalf("live calls = %d, want 1", source.Calls(health.TargetLive))
	}
}

func TestCountingSourceErrorAndUnknown(t *testing.T) {
	t.Parallel()

	raw := errors.New("private source failure")
	source := NewCountingSource(nil)
	source.SetError(health.TargetReady, raw)

	if _, err := source.Evaluate(context.Background(), health.TargetReady); !errors.Is(err, raw) {
		t.Fatalf("Evaluate(ready) = %v, want raw error", err)
	}
	source.SetError(health.TargetReady, nil)
	report, err := source.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate(ready after clear) = %v, want nil", err)
	}
	AssertReportStatus(t, report, health.StatusUnknown)
}

func TestCountingSourceConcurrentUse(t *testing.T) {
	t.Parallel()

	source := NewCountingSource(map[health.Target]health.Report{
		health.TargetReady: HealthyReport(health.TargetReady),
	})

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 25; j++ {
				_, _ = source.Evaluate(context.Background(), health.TargetReady)
			}
		}()
	}
	wg.Wait()

	if source.Calls(health.TargetReady) != 200 {
		t.Fatalf("ready calls = %d, want 200", source.Calls(health.TargetReady))
	}
}

func TestSequenceSource(t *testing.T) {
	t.Parallel()

	source := NewSequenceSource(
		health.TargetReady,
		HealthyReport(health.TargetReady),
		DegradedReport(health.TargetReady),
	)

	first, err := source.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("first Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, first, health.StatusHealthy)

	second, err := source.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("second Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, second, health.StatusDegraded)

	third, err := source.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("third Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, third, health.StatusDegraded)

	other, err := source.Evaluate(context.Background(), health.TargetLive)
	if err != nil {
		t.Fatalf("other Evaluate() = %v, want nil", err)
	}
	AssertReportStatus(t, other, health.StatusUnknown)
	AssertReportTarget(t, other, health.TargetLive)

	if source.Calls(health.TargetReady) != 3 {
		t.Fatalf("ready calls = %d, want 3", source.Calls(health.TargetReady))
	}
	if source.Calls(health.TargetLive) != 1 {
		t.Fatalf("live calls = %d, want 1", source.Calls(health.TargetLive))
	}
}

func TestSequenceSourceCopiesReports(t *testing.T) {
	t.Parallel()

	report := HealthyReport(health.TargetReady)
	source := NewSequenceSource(health.TargetReady, report)
	report.Checks[0] = UnhealthyResult("mutated", health.ReasonFatal)

	got, err := source.Evaluate(context.Background(), health.TargetReady)
	if err != nil {
		t.Fatalf("Evaluate() = %v, want nil", err)
	}
	if got.Checks[0].Name != "ready_check" {
		t.Fatalf("check name = %q, want ready_check", got.Checks[0].Name)
	}
}
