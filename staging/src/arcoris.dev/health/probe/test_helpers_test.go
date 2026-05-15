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

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
	"arcoris.dev/snapshot"
)

const testTimeout = 5 * time.Second

var testNow = time.Unix(100, 0)

func newTestClock() *clock.FakeClock {
	return clock.NewFakeClock(testNow)
}

func newTestEvaluator(t *testing.T) health.Evaluator {
	t.Helper()

	return healthtest.NewEvaluatorForTarget(t, health.TargetReady, healthtest.HealthyChecker("ready_check"))
}

func newTestRunner(t *testing.T, clk clock.Clock, opts ...Option) *Runner {
	t.Helper()

	allOptions := []Option{WithClock(clk), WithTargets(health.TargetReady)}
	allOptions = append(allOptions, opts...)

	runner, err := NewRunner(newTestEvaluator(t), allOptions...)
	if err != nil {
		t.Fatalf("NewRunner() = %v, want nil", err)
	}

	return runner
}

func waitForSnapshot(t *testing.T, r *Runner, target health.Target) Snapshot {
	t.Helper()

	return waitForSnapshotWhere(t, r, target, func(Snapshot) bool {
		return true
	})
}

func waitForSnapshotWhere(
	t *testing.T,
	r *Runner,
	target health.Target,
	accept func(Snapshot) bool,
) Snapshot {
	t.Helper()

	deadline := time.After(testTimeout)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for snapshot target=%s", target)
		case <-ticker.C:
			snap, ok := r.Snapshot(target)
			if ok && accept(snap) {
				return snap
			}
		}
	}
}

func waitForRevision(t *testing.T, r *Runner, target health.Target, rev snapshot.Revision) Snapshot {
	t.Helper()

	return waitForSnapshotWhere(t, r, target, func(snap Snapshot) bool {
		return snap.Revision >= rev
	})
}

func waitForRunnerRunning(t *testing.T, r *Runner) {
	t.Helper()

	deadline := time.After(testTimeout)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for runner to start")
		case <-ticker.C:
			if r.running.Load() {
				return
			}
		}
	}
}

func stepUntilRevision(
	t *testing.T,
	clk *clock.FakeClock,
	r *Runner,
	target health.Target,
	rev snapshot.Revision,
	interval time.Duration,
) Snapshot {
	t.Helper()

	deadline := time.After(testTimeout)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for target=%s revision=%d", target, rev)
		case <-ticker.C:
			clk.Step(interval)
			if snap, ok := r.Snapshot(target); ok && snap.Revision >= rev {
				return snap
			}
		}
	}
}

func waitForRunDone(t *testing.T, done <-chan error) error {
	t.Helper()

	select {
	case err := <-done:
		return err
	case <-time.After(testTimeout):
		t.Fatal("Run did not stop")
		return nil
	}
}

func sameTargets(l []health.Target, r []health.Target) bool {
	if len(l) != len(r) {
		return false
	}
	for i := range l {
		if l[i] != r[i] {
			return false
		}
	}

	return true
}

func firstScheduleDelay(t *testing.T, sched delay.Schedule) time.Duration {
	t.Helper()

	seq := sched.NewSequence()
	if seq == nil {
		t.Fatal("schedule returned nil sequence")
	}
	d, ok := seq.Next()
	if !ok {
		t.Fatal("schedule sequence exhausted before first delay")
	}

	return d
}

type nilSequenceSchedule struct{}

func (nilSequenceSchedule) NewSequence() delay.Sequence {
	return nil
}

type negativeDelaySchedule struct{}

func (negativeDelaySchedule) NewSequence() delay.Sequence {
	return negativeDelaySequence{}
}

type negativeDelaySequence struct{}

func (negativeDelaySequence) Next() (time.Duration, bool) {
	return -time.Nanosecond, true
}
