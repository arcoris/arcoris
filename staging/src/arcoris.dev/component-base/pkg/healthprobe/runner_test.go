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

func TestNewRunner(t *testing.T) {
	t.Parallel()

	evaluator := newTestEvaluator(t)
	clk := newManualClock()

	runner, err := NewRunner(
		evaluator,
		WithClock(clk),
		WithTargets(health.TargetReady),
		WithInterval(time.Second),
		WithStaleAfter(2*time.Second),
		WithInitialProbe(false),
	)
	if err != nil {
		t.Fatalf("NewRunner() = %v, want nil", err)
	}
	if runner.evaluator != evaluator {
		t.Fatal("evaluator not stored")
	}
	if runner.clock != clk {
		t.Fatal("clock not stored")
	}
	if runner.interval != time.Second {
		t.Fatalf("interval = %s, want 1s", runner.interval)
	}
	if runner.staleAfter != 2*time.Second {
		t.Fatalf("staleAfter = %s, want 2s", runner.staleAfter)
	}
	if runner.initialProbe {
		t.Fatal("initialProbe = true, want false")
	}
}

func TestNewRunnerRejectsNilEvaluator(t *testing.T) {
	t.Parallel()

	_, err := NewRunner(nil, WithTargets(health.TargetReady))
	if !errors.Is(err, ErrNilEvaluator) {
		t.Fatalf("NewRunner(nil) = %v, want ErrNilEvaluator", err)
	}
}

func TestRunnerSnapshotComputesStale(t *testing.T) {
	t.Parallel()

	clk := newManualClock()
	runner := newTestRunner(t, clk, WithStaleAfter(time.Second))
	runner.store.update(
		health.TargetReady,
		health.Report{Target: health.TargetReady, Status: health.StatusHealthy},
		clk.Now(),
	)

	snapshot, ok := runner.Snapshot(health.TargetReady)
	if !ok {
		t.Fatal("Snapshot() ok = false, want true")
	}
	if snapshot.Stale {
		t.Fatal("Stale = true, want false")
	}

	clk.Advance(time.Second + time.Nanosecond)
	snapshot, ok = runner.Snapshot(health.TargetReady)
	if !ok {
		t.Fatal("Snapshot() ok = false, want true")
	}
	if !snapshot.Stale {
		t.Fatal("Stale = false, want true")
	}
}

func TestRunnerSnapshotsComputeStale(t *testing.T) {
	t.Parallel()

	clk := newManualClock()
	runner := newTestRunner(t, clk, WithStaleAfter(time.Second))
	runner.store.update(health.TargetReady, health.Report{Target: health.TargetReady, Status: health.StatusHealthy}, clk.Now())

	clk.Advance(time.Second + time.Nanosecond)
	snapshots := runner.Snapshots()
	if len(snapshots) != 1 {
		t.Fatalf("Snapshots length = %d, want 1", len(snapshots))
	}
	if !snapshots[0].Stale {
		t.Fatal("Stale = false, want true")
	}
}
