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
	"sync"
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
	"arcoris.dev/snapshot"
)

func TestStoreSnapshotLifecycle(t *testing.T) {
	t.Parallel()

	clk := newTestClock()
	s := newStore([]health.Target{health.TargetReady, health.TargetLive}, clk)

	if _, ok := s.snapshot(health.TargetReady); ok {
		t.Fatal("snapshot before update ok = true, want false")
	}

	report := healthtest.HealthyReport(health.TargetReady)
	if ok := s.update(health.TargetReady, report); !ok {
		t.Fatal("update() = false, want true")
	}
	snap, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot after update ok = false, want true")
	}
	if snap.Revision != snapshot.Revision(1) {
		t.Fatalf("Revision = %d, want 1", snap.Revision)
	}
	if !snap.Updated.Equal(testNow) {
		t.Fatalf("Updated = %v, want %v", snap.Updated, testNow)
	}

	clk.Step(time.Second)
	if ok := s.update(health.TargetReady, report); !ok {
		t.Fatal("second update() = false, want true")
	}
	snap, ok = s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot after second update ok = false, want true")
	}
	if snap.Revision != snapshot.Revision(2) {
		t.Fatalf("Revision = %d, want 2", snap.Revision)
	}
	if !snap.Updated.Equal(testNow.Add(time.Second)) {
		t.Fatalf("Updated = %v, want %v", snap.Updated, testNow.Add(time.Second))
	}
}

func TestStoreRevisionsAreIndependentPerTarget(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady, health.TargetLive}, newTestClock())

	if ok := s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady)); !ok {
		t.Fatal("update(ready) = false, want true")
	}
	if ok := s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady)); !ok {
		t.Fatal("second update(ready) = false, want true")
	}
	if ok := s.update(health.TargetLive, healthtest.HealthyReport(health.TargetLive)); !ok {
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
	if ready.Revision != snapshot.Revision(2) {
		t.Fatalf("ready Revision = %d, want 2", ready.Revision)
	}
	if live.Revision != snapshot.Revision(1) {
		t.Fatalf("live Revision = %d, want 1", live.Revision)
	}
}

func TestStoreUpdateUsesSnapshotClock(t *testing.T) {
	t.Parallel()

	clk := newTestClock()
	clk.Step(3 * time.Second)
	s := newStore([]health.Target{health.TargetReady}, clk)

	if ok := s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady)); !ok {
		t.Fatal("update() = false, want true")
	}

	snap, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if !snap.Updated.Equal(clk.Now()) {
		t.Fatalf("Updated = %v, want %v", snap.Updated, clk.Now())
	}
}

func TestStoreSnapshotsPreserveConfiguredOrder(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady, health.TargetLive}, newTestClock())

	if ok := s.update(health.TargetLive, healthtest.HealthyReport(health.TargetLive)); !ok {
		t.Fatal("update(live) = false, want true")
	}
	if ok := s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady)); !ok {
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

	s := newStore([]health.Target{health.TargetReady}, newTestClock())

	if ok := s.update(health.TargetLive, healthtest.HealthyReport(health.TargetLive)); ok {
		t.Fatal("update(unconfigured) = true, want false")
	}
	if _, ok := s.snapshot(health.TargetLive); ok {
		t.Fatal("snapshot(unconfigured) ok = true, want false")
	}
	if snapshots := s.snapshots(); len(snapshots) != 0 {
		t.Fatalf("snapshots length = %d, want 0", len(snapshots))
	}
}

func TestStoreSnapshotUnobservedTarget(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady}, newTestClock())

	if _, ok := s.snapshot(health.TargetReady); ok {
		t.Fatal("snapshot(unobserved) ok = true, want false")
	}
}

func TestStoreRejectsInvalidSnapshotInput(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	report := health.Report{
		Target:   health.TargetLive,
		Status:   health.StatusHealthy,
		Observed: testNow,
	}

	if ok := s.update(health.TargetReady, report); ok {
		t.Fatal("update(mismatched report) = true, want false")
	}
	if _, ok := s.snapshot(health.TargetReady); ok {
		t.Fatal("snapshot() ok = true, want false")
	}
}

func TestStoreUpdateClonesInputReport(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	report := healthtest.HealthyReport(health.TargetReady)

	if ok := s.update(health.TargetReady, report); !ok {
		t.Fatal("update() = false, want true")
	}
	report.Checks[0] = health.Unhealthy("mutated", health.ReasonFatal, "mutated")

	snap, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if snap.Report.Checks[0].Name != "ready_check" {
		t.Fatalf("stored check name = %q, want ready_check", snap.Report.Checks[0].Name)
	}
}

func TestStoreSnapshotReturnsDetachedReport(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	if ok := s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady)); !ok {
		t.Fatal("update() = false, want true")
	}

	snap, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	snap.Report.Checks[0] = health.Unhealthy("mutated_again", health.ReasonFatal, "mutated")

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

	s := newStore([]health.Target{health.TargetReady}, newTestClock())

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady))
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

	snap, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot ok = false, want true")
	}
	if snap.Revision == snapshot.ZeroRevision {
		t.Fatal("Revision = 0, want positive")
	}
}
