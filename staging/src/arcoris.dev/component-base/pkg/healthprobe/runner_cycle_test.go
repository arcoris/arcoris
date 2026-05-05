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
	"context"
	"testing"

	"arcoris.dev/component-base/pkg/health"
	"arcoris.dev/component-base/pkg/healthtest"
)

func TestRunnerRunCycleEvaluatesTargetsInOrder(t *testing.T) {
	t.Parallel()

	clk := newTestClock()
	order := make(chan health.Target, 2)
	evaluator := healthtest.NewEvaluator(t, healthtest.NewRegistry(
		t,
		healthtest.ForTarget(health.TargetReady, healthtest.FuncChecker("ready_check", func(context.Context) health.Result {
			order <- health.TargetReady
			return health.Healthy("ready_check")
		})),
		healthtest.ForTarget(health.TargetLive, healthtest.FuncChecker("live_check", func(context.Context) health.Result {
			order <- health.TargetLive
			return health.Healthy("live_check")
		})),
	))
	runner, err := NewRunner(evaluator, WithClock(clk), WithTargets(health.TargetReady, health.TargetLive))
	if err != nil {
		t.Fatalf("NewRunner() = %v, want nil", err)
	}

	runner.runCycle(context.Background())

	if got := <-order; got != health.TargetReady {
		t.Fatalf("first target = %s, want ready", got)
	}
	if got := <-order; got != health.TargetLive {
		t.Fatalf("second target = %s, want live", got)
	}
	if snapshots := runner.Snapshots(); len(snapshots) != 2 {
		t.Fatalf("Snapshots length = %d, want 2", len(snapshots))
	}
}

func TestRunnerRunCycleDoesNotStoreAfterContextCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	runner := newTestRunner(t, newTestClock())
	runner.runCycle(ctx)

	if _, ok := runner.Snapshot(health.TargetReady); ok {
		t.Fatal("Snapshot() ok = true, want false")
	}
}

func TestUnknownReport(t *testing.T) {
	t.Parallel()

	observed := newTestClock().Now()
	report := unknownReport(health.TargetReady, observed)
	if report.Target != health.TargetReady {
		t.Fatalf("Target = %s, want ready", report.Target)
	}
	if report.Status != health.StatusUnknown {
		t.Fatalf("Status = %s, want unknown", report.Status)
	}
	if !report.Observed.Equal(observed) {
		t.Fatalf("Observed = %v, want %v", report.Observed, observed)
	}
	if len(report.Checks) != 0 {
		t.Fatalf("Checks length = %d, want 0", len(report.Checks))
	}
	if !report.IsValid() {
		t.Fatalf("unknownReport().IsValid() = false, want true: %#v", report)
	}
}
