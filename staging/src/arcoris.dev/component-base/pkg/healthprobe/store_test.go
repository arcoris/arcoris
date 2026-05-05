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
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

func TestStoreSnapshotLifecycle(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady, health.TargetLive})

	if _, ok := s.snapshot(health.TargetReady); ok {
		t.Fatal("snapshot before update ok = true, want false")
	}

	updated := time.Unix(10, 0)
	report := health.Report{
		Target: health.TargetReady,
		Status: health.StatusHealthy,
		Checks: []health.Result{health.Healthy("database")},
	}

	s.update(health.TargetReady, report, updated)
	snapshot, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot after update ok = false, want true")
	}
	if snapshot.Generation != 1 {
		t.Fatalf("Generation = %d, want 1", snapshot.Generation)
	}
	if !snapshot.Updated.Equal(updated) {
		t.Fatalf("Updated = %v, want %v", snapshot.Updated, updated)
	}

	s.update(health.TargetReady, report, updated.Add(time.Second))
	snapshot, ok = s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot after second update ok = false, want true")
	}
	if snapshot.Generation != 2 {
		t.Fatalf("Generation = %d, want 2", snapshot.Generation)
	}
}

func TestStoreSnapshotsReturnConfiguredOrder(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady, health.TargetLive})
	updated := time.Unix(10, 0)

	s.update(health.TargetLive, health.Report{Target: health.TargetLive, Status: health.StatusHealthy}, updated)
	s.update(health.TargetReady, health.Report{Target: health.TargetReady, Status: health.StatusHealthy}, updated)

	snapshots := s.snapshots()
	if len(snapshots) != 2 {
		t.Fatalf("snapshots length = %d, want 2", len(snapshots))
	}
	if snapshots[0].Target != health.TargetReady || snapshots[1].Target != health.TargetLive {
		t.Fatalf("snapshot order = [%s %s], want [ready live]", snapshots[0].Target, snapshots[1].Target)
	}
}

func TestStoreCopiesReportChecks(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady})
	report := health.Report{
		Target: health.TargetReady,
		Status: health.StatusHealthy,
		Checks: []health.Result{health.Healthy("database")},
	}

	s.update(health.TargetReady, report, time.Unix(10, 0))
	report.Checks[0] = health.Unhealthy("mutated", health.ReasonFatal, "mutated")

	snapshot, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if snapshot.Report.Checks[0].Name != "database" {
		t.Fatalf("stored check name = %q, want database", snapshot.Report.Checks[0].Name)
	}

	snapshot.Report.Checks[0] = health.Unhealthy("mutated_again", health.ReasonFatal, "mutated")
	again, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if again.Report.Checks[0].Name != "database" {
		t.Fatalf("stored check name after read mutation = %q, want database", again.Report.Checks[0].Name)
	}
}
