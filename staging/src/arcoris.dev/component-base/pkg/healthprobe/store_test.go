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
	"sync"
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
	report := healthyReport(health.TargetReady, updated)

	if ok := s.update(health.TargetReady, report, updated); !ok {
		t.Fatal("update() = false, want true")
	}
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

	if ok := s.update(health.TargetReady, report, updated.Add(time.Second)); !ok {
		t.Fatal("second update() = false, want true")
	}
	snapshot, ok = s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot after second update ok = false, want true")
	}
	if snapshot.Generation != 2 {
		t.Fatalf("Generation = %d, want 2", snapshot.Generation)
	}
}

func TestStoreGenerationIsIndependentPerTarget(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady, health.TargetLive})
	updated := time.Unix(10, 0)

	if ok := s.update(health.TargetReady, healthyReport(health.TargetReady, updated), updated); !ok {
		t.Fatal("update(ready) = false, want true")
	}
	if ok := s.update(health.TargetReady, healthyReport(health.TargetReady, updated), updated); !ok {
		t.Fatal("second update(ready) = false, want true")
	}
	if ok := s.update(health.TargetLive, healthyReport(health.TargetLive, updated), updated); !ok {
		t.Fatal("update(live) = false, want true")
	}

	ready, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot(ready) ok = false, want true")
	}
	live, ok := s.snapshot(health.TargetLive)
	if !ok {
		t.Fatal("snapshot(live) ok = false, want true")
	}
	if ready.Generation != 2 {
		t.Fatalf("ready Generation = %d, want 2", ready.Generation)
	}
	if live.Generation != 1 {
		t.Fatalf("live Generation = %d, want 1", live.Generation)
	}
}

func TestStoreSnapshotsReturnConfiguredOrder(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady, health.TargetLive})
	updated := time.Unix(10, 0)

	if ok := s.update(health.TargetLive, healthyReport(health.TargetLive, updated), updated); !ok {
		t.Fatal("update(live) = false, want true")
	}
	if ok := s.update(health.TargetReady, healthyReport(health.TargetReady, updated), updated); !ok {
		t.Fatal("update(ready) = false, want true")
	}

	snapshots := s.snapshots()
	if len(snapshots) != 2 {
		t.Fatalf("snapshots length = %d, want 2", len(snapshots))
	}
	if snapshots[0].Target != health.TargetReady || snapshots[1].Target != health.TargetLive {
		t.Fatalf("snapshot order = [%s %s], want [ready live]", snapshots[0].Target, snapshots[1].Target)
	}
}

func TestStoreRejectsUnconfiguredTarget(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady})
	updated := time.Unix(10, 0)

	if ok := s.update(health.TargetLive, healthyReport(health.TargetLive, updated), updated); ok {
		t.Fatal("update(unconfigured) = true, want false")
	}
	if _, ok := s.snapshot(health.TargetLive); ok {
		t.Fatal("snapshot(unconfigured) ok = true, want false")
	}
	if snapshots := s.snapshots(); len(snapshots) != 0 {
		t.Fatalf("snapshots length = %d, want 0", len(snapshots))
	}
}

func TestStoreRejectsInvalidSnapshotInput(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady})
	updated := time.Unix(10, 0)
	report := health.Report{
		Target:   health.TargetLive,
		Status:   health.StatusHealthy,
		Observed: updated,
	}

	if ok := s.update(health.TargetReady, report, updated); ok {
		t.Fatal("update(mismatched report) = true, want false")
	}
	if _, ok := s.snapshot(health.TargetReady); ok {
		t.Fatal("snapshot() ok = true, want false")
	}
}

func TestStoreCopiesReportChecksOnWrite(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady})
	updated := time.Unix(10, 0)
	report := healthyReport(health.TargetReady, updated)

	if ok := s.update(health.TargetReady, report, updated); !ok {
		t.Fatal("update() = false, want true")
	}
	report.Checks[0] = health.Unhealthy("mutated", health.ReasonFatal, "mutated")

	snapshot, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if snapshot.Report.Checks[0].Name != "ready_check" {
		t.Fatalf("stored check name = %q, want ready_check", snapshot.Report.Checks[0].Name)
	}
}

func TestStoreCopiesReportChecksOnRead(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady})
	updated := time.Unix(10, 0)
	if ok := s.update(health.TargetReady, healthyReport(health.TargetReady, updated), updated); !ok {
		t.Fatal("update() = false, want true")
	}

	snapshot, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	snapshot.Report.Checks[0] = health.Unhealthy("mutated_again", health.ReasonFatal, "mutated")

	again, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if again.Report.Checks[0].Name != "ready_check" {
		t.Fatalf("stored check name after read mutation = %q, want ready_check", again.Report.Checks[0].Name)
	}

	snapshots := s.snapshots()
	snapshots[0].Report.Checks[0] = health.Unhealthy("mutated_snapshots", health.ReasonFatal, "mutated")
	again, ok = s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if again.Report.Checks[0].Name != "ready_check" {
		t.Fatalf("stored check name after snapshots mutation = %q, want ready_check", again.Report.Checks[0].Name)
	}
}

func TestStoreConcurrentReadUpdate(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady})
	updated := time.Unix(10, 0)

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = s.update(health.TargetReady, healthyReport(health.TargetReady, updated), updated)
			}
		}()
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = s.snapshot(health.TargetReady)
				_ = s.snapshots()
			}
		}()
	}
	wg.Wait()

	snapshot, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if snapshot.Generation == 0 {
		t.Fatal("Generation = 0, want positive")
	}
}
