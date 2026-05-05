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

package healthprobe

import (
	"errors"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

func TestCloneReportCopiesChecks(t *testing.T) {
	t.Parallel()

	observed := time.Unix(10, 20)
	cause := errors.New("internal cause")

	source := health.Report{
		Target:   health.TargetReady,
		Status:   health.StatusUnhealthy,
		Observed: observed,
		Duration: 25 * time.Millisecond,
		Checks: []health.Result{
			health.Unhealthy(
				"database",
				health.ReasonDependencyUnavailable,
				"database is unavailable",
			).WithObserved(observed).WithDuration(10 * time.Millisecond).WithCause(cause),
			health.Healthy("cache").
				WithObserved(observed).
				WithDuration(5 * time.Millisecond),
		},
	}

	cloned := cloneReport(source)

	if cloned.Target != source.Target {
		t.Fatalf("Target = %s, want %s", cloned.Target, source.Target)
	}
	if cloned.Status != source.Status {
		t.Fatalf("Status = %s, want %s", cloned.Status, source.Status)
	}
	if !cloned.Observed.Equal(source.Observed) {
		t.Fatalf("Observed = %v, want %v", cloned.Observed, source.Observed)
	}
	if cloned.Duration != source.Duration {
		t.Fatalf("Duration = %v, want %v", cloned.Duration, source.Duration)
	}
	if len(cloned.Checks) != len(source.Checks) {
		t.Fatalf("Checks length = %d, want %d", len(cloned.Checks), len(source.Checks))
	}
	if len(cloned.Checks) > 0 && &cloned.Checks[0] == &source.Checks[0] {
		t.Fatal("cloneReport reused source Checks backing array")
	}

	cloned.Checks[0] = health.Healthy("clone_mutated")

	if source.Checks[0].Name != "database" {
		t.Fatalf("source check was mutated through clone: name=%q", source.Checks[0].Name)
	}
	if source.Checks[0].Cause != cause {
		t.Fatal("source check cause was unexpectedly changed")
	}

	source.Checks[1] = health.Unhealthy(
		"source_mutated",
		health.ReasonFatal,
		"source was mutated after clone",
	)

	if cloned.Checks[1].Name != "cache" {
		t.Fatalf("cloned check was mutated through source: name=%q", cloned.Checks[1].Name)
	}
}

func TestCloneReportHandlesEmptyChecks(t *testing.T) {
	t.Parallel()

	source := health.Report{
		Target:   health.TargetReady,
		Status:   health.StatusHealthy,
		Observed: time.Unix(10, 0),
	}

	cloned := cloneReport(source)

	if cloned.Target != source.Target {
		t.Fatalf("Target = %s, want %s", cloned.Target, source.Target)
	}
	if cloned.Status != source.Status {
		t.Fatalf("Status = %s, want %s", cloned.Status, source.Status)
	}
	if cloned.Checks != nil {
		t.Fatalf("Checks = %#v, want nil", cloned.Checks)
	}
}

func TestCloneSnapshotCopiesEmbeddedReportChecks(t *testing.T) {
	t.Parallel()

	observed := time.Unix(10, 20)
	updated := time.Unix(11, 0)

	source := Snapshot{
		Target: health.TargetReady,
		Report: health.Report{
			Target:   health.TargetReady,
			Status:   health.StatusDegraded,
			Observed: observed,
			Duration: time.Millisecond,
			Checks: []health.Result{
				health.Degraded(
					"queue",
					health.ReasonOverloaded,
					"queue is overloaded",
				).WithObserved(observed),
			},
		},
		Updated:    updated,
		Generation: 42,
		Stale:      true,
	}

	cloned := cloneSnapshot(source)

	if cloned.Target != source.Target {
		t.Fatalf("Target = %s, want %s", cloned.Target, source.Target)
	}
	if cloned.Report.Target != source.Report.Target {
		t.Fatalf("Report.Target = %s, want %s", cloned.Report.Target, source.Report.Target)
	}
	if cloned.Report.Status != source.Report.Status {
		t.Fatalf("Report.Status = %s, want %s", cloned.Report.Status, source.Report.Status)
	}
	if !cloned.Report.Observed.Equal(source.Report.Observed) {
		t.Fatalf("Report.Observed = %v, want %v", cloned.Report.Observed, source.Report.Observed)
	}
	if cloned.Report.Duration != source.Report.Duration {
		t.Fatalf("Report.Duration = %v, want %v", cloned.Report.Duration, source.Report.Duration)
	}
	if !cloned.Updated.Equal(source.Updated) {
		t.Fatalf("Updated = %v, want %v", cloned.Updated, source.Updated)
	}
	if cloned.Generation != source.Generation {
		t.Fatalf("Generation = %d, want %d", cloned.Generation, source.Generation)
	}
	if cloned.Stale != source.Stale {
		t.Fatalf("Stale = %v, want %v", cloned.Stale, source.Stale)
	}
	if len(cloned.Report.Checks) != len(source.Report.Checks) {
		t.Fatalf("Checks length = %d, want %d", len(cloned.Report.Checks), len(source.Report.Checks))
	}
	if len(cloned.Report.Checks) > 0 && &cloned.Report.Checks[0] == &source.Report.Checks[0] {
		t.Fatal("cloneSnapshot reused source Report.Checks backing array")
	}

	cloned.Report.Checks[0] = health.Healthy("clone_mutated")

	if source.Report.Checks[0].Name != "queue" {
		t.Fatalf("source snapshot check was mutated through clone: name=%q", source.Report.Checks[0].Name)
	}

	source.Report.Checks[0] = health.Unhealthy(
		"source_mutated",
		health.ReasonFatal,
		"source was mutated after clone",
	)

	if cloned.Report.Checks[0].Name != "clone_mutated" {
		t.Fatalf("cloned snapshot check was mutated through source: name=%q", cloned.Report.Checks[0].Name)
	}
}

func TestCloneSnapshotHandlesZeroSnapshot(t *testing.T) {
	t.Parallel()

	cloned := cloneSnapshot(Snapshot{})

	if !cloned.IsZero() {
		t.Fatalf("cloneSnapshot(zero).IsZero() = false, want true: %#v", cloned)
	}
	if !cloned.IsValid() {
		t.Fatalf("cloneSnapshot(zero).IsValid() = false, want true: %#v", cloned)
	}
}
