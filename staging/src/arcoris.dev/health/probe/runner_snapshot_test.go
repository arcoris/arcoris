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

package probe

import (
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
)

func TestRunnerSnapshotNilRunner(t *testing.T) {
	t.Parallel()

	var runner *Runner

	if snap, ok := runner.Snapshot(health.TargetReady); ok || !snap.IsZero() {
		t.Fatalf("Snapshot() = %#v, %v; want zero false", snap, ok)
	}
	if snapshots := runner.Snapshots(); snapshots != nil {
		t.Fatalf("Snapshots() = %#v, want nil", snapshots)
	}
}

func TestRunnerSnapshotRejectsInvalidOrUnconfiguredTarget(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(t, newTestClock())

	tests := []health.Target{
		health.TargetUnknown,
		health.Target(255),
		health.TargetLive,
	}
	for _, target := range tests {
		if snap, ok := runner.Snapshot(target); ok || !snap.IsZero() {
			t.Fatalf("Snapshot(%s) = %#v, %v; want zero false", target, snap, ok)
		}
	}
}

func TestRunnerSnapshotComputesStaleAtReadTime(t *testing.T) {
	t.Parallel()

	clk := newTestClock()
	runner := newTestRunner(t, clk, WithStaleAfter(time.Second))
	if ok := runner.store.update(
		health.TargetReady,
		healthtest.HealthyReport(health.TargetReady),
	); !ok {
		t.Fatal("store.update() = false, want true")
	}

	snap, ok := runner.Snapshot(health.TargetReady)
	if !ok {
		t.Fatal("Snapshot() ok = false, want true")
	}
	if snap.Stale {
		t.Fatal("Stale = true, want false")
	}

	clk.Step(time.Second)
	snap, ok = runner.Snapshot(health.TargetReady)
	if !ok {
		t.Fatal("Snapshot() ok = false, want true")
	}
	if snap.Stale {
		t.Fatal("Stale at exact boundary = true, want false")
	}

	clk.Step(time.Nanosecond)
	snap, ok = runner.Snapshot(health.TargetReady)
	if !ok {
		t.Fatal("Snapshot() ok = false, want true")
	}
	if !snap.Stale {
		t.Fatal("Stale = false, want true")
	}
}

func TestRunnerSnapshotsComputeStaleAndOrder(t *testing.T) {
	t.Parallel()

	clk := newTestClock()
	evaluator := newTestEvaluator(t)
	runner, err := NewRunner(
		evaluator,
		WithClock(clk),
		WithTargets(health.TargetLive, health.TargetReady),
		WithStaleAfter(time.Second),
	)
	if err != nil {
		t.Fatalf("NewRunner() = %v, want nil", err)
	}
	if ok := runner.store.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady)); !ok {
		t.Fatal("store.update(ready) = false, want true")
	}
	if ok := runner.store.update(health.TargetLive, healthtest.HealthyReport(health.TargetLive)); !ok {
		t.Fatal("store.update(live) = false, want true")
	}

	clk.Step(time.Second + time.Nanosecond)
	snapshots := runner.Snapshots()
	if len(snapshots) != 2 {
		t.Fatalf("Snapshots length = %d, want 2", len(snapshots))
	}
	if snapshots[0].Target != health.TargetLive || snapshots[1].Target != health.TargetReady {
		t.Fatalf("snapshot order = [%s %s], want [live ready]", snapshots[0].Target, snapshots[1].Target)
	}
	if !snapshots[0].Stale || !snapshots[1].Stale {
		t.Fatalf("Snapshots stale = [%v %v], want [true true]", snapshots[0].Stale, snapshots[1].Stale)
	}
}

func TestRunnerSnapshotReadsAreDetached(t *testing.T) {
	t.Parallel()

	clk := newTestClock()
	runner := newTestRunner(t, clk)
	if ok := runner.store.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady)); !ok {
		t.Fatal("store.update() = false, want true")
	}

	snap, ok := runner.Snapshot(health.TargetReady)
	if !ok {
		t.Fatal("Snapshot() ok = false, want true")
	}
	snap.Report.Checks[0] = health.Unhealthy("mutated_snapshot", health.ReasonFatal, "mutated")

	snapshots := runner.Snapshots()
	if len(snapshots) != 1 {
		t.Fatalf("Snapshots length = %d, want 1", len(snapshots))
	}
	snapshots[0].Report.Checks[0] = health.Unhealthy("mutated_snapshots", health.ReasonFatal, "mutated")

	again, ok := runner.Snapshot(health.TargetReady)
	if !ok {
		t.Fatal("Snapshot() ok = false, want true")
	}
	if again.Report.Checks[0].Name != "ready_check" {
		t.Fatalf("stored check name = %q, want ready_check", again.Report.Checks[0].Name)
	}
}
