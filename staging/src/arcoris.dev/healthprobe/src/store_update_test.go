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

package probe

import (
	"sync"
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/healthtest"
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

func TestStoreRejectsInconsistentReportWithoutAdvancingRevision(t *testing.T) {
	t.Parallel()

	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	if ok := s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady)); !ok {
		t.Fatal("initial update() = false, want true")
	}

	before, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("initial snapshot ok = false, want true")
	}

	inconsistent := health.Report{
		Target:   health.TargetReady,
		Status:   health.StatusHealthy,
		Observed: testNow,
	}
	if ok := s.update(health.TargetReady, inconsistent); ok {
		t.Fatal("update(inconsistent report) = true, want false")
	}

	after, ok := s.snapshot(health.TargetReady)
	if !ok {
		t.Fatal("snapshot after rejected update ok = false, want true")
	}
	if after.Revision != before.Revision {
		t.Fatalf("Revision after rejected update = %d, want %d", after.Revision, before.Revision)
	}
	if after.Report.Status != before.Report.Status || len(after.Report.Checks) != len(before.Report.Checks) {
		t.Fatalf("snapshot after rejected update = %+v, want previous report %+v", after.Report, before.Report)
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
