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
	"errors"
	"sync"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/health"
	"arcoris.dev/component-base/pkg/healthtest"
)

func TestRunnerRunPerformsInitialProbe(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := newTestRunner(t, newTestClock())
	done := make(chan error, 1)

	go func() {
		done <- runner.Run(ctx)
	}()

	snapshot := waitForSnapshot(t, runner, health.TargetReady)
	if !snapshot.IsFresh() {
		t.Fatalf("snapshot IsFresh() = false, want true: %#v", snapshot)
	}

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}

func TestRunnerRunWaitsForTickWhenInitialProbeDisabled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := time.Second
	clk := newTestClock()
	runner := newTestRunner(t, clk, WithInterval(interval), WithInitialProbe(false))
	done := make(chan error, 1)

	go func() {
		done <- runner.Run(ctx)
	}()

	waitForRunnerRunning(t, runner)
	if _, ok := runner.Snapshot(health.TargetReady); ok {
		t.Fatal("Snapshot before tick ok = true, want false")
	}

	snapshot := stepUntilGeneration(t, clk, runner, health.TargetReady, 1, interval)
	if snapshot.Generation != 1 {
		t.Fatalf("Generation = %d, want 1", snapshot.Generation)
	}

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}

func TestRunnerRunTickDrivenProbeIncrementsGeneration(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := time.Second
	clk := newTestClock()
	runner := newTestRunner(t, clk, WithInterval(interval), WithInitialProbe(false))
	done := make(chan error, 1)

	go func() {
		done <- runner.Run(ctx)
	}()

	waitForRunnerRunning(t, runner)
	first := stepUntilGeneration(t, clk, runner, health.TargetReady, 1, interval)
	second := stepUntilGeneration(t, clk, runner, health.TargetReady, 2, interval)

	if first.Generation != 1 {
		t.Fatalf("first Generation = %d, want 1", first.Generation)
	}
	if second.Generation != 2 {
		t.Fatalf("second Generation = %d, want 2", second.Generation)
	}
	if !second.Updated.After(first.Updated) {
		t.Fatalf("second Updated = %v, want after %v", second.Updated, first.Updated)
	}

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}

func TestRunnerRunRejectsConcurrentRun(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := newTestRunner(t, newTestClock())
	done := make(chan error, 1)

	go func() {
		done <- runner.Run(ctx)
	}()
	waitForRunnerRunning(t, runner)

	err := runner.Run(context.Background())
	if !errors.Is(err, ErrRunnerRunning) {
		t.Fatalf("concurrent Run() = %v, want ErrRunnerRunning", err)
	}

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}

func TestRunnerRunCanRestartAfterStop(t *testing.T) {
	t.Parallel()

	interval := time.Second
	clk := newTestClock()
	runner := newTestRunner(t, clk, WithInterval(interval), WithInitialProbe(false))

	firstCtx, firstCancel := context.WithCancel(context.Background())
	firstDone := make(chan error, 1)
	go func() {
		firstDone <- runner.Run(firstCtx)
	}()
	waitForRunnerRunning(t, runner)
	firstCancel()
	if err := waitForRunDone(t, firstDone); err != nil {
		t.Fatalf("first Run() = %v, want nil", err)
	}

	secondCtx, secondCancel := context.WithCancel(context.Background())
	defer secondCancel()
	secondDone := make(chan error, 1)
	go func() {
		secondDone <- runner.Run(secondCtx)
	}()
	waitForRunnerRunning(t, runner)
	snapshot := stepUntilGeneration(t, clk, runner, health.TargetReady, 1, interval)
	if snapshot.Generation != 1 {
		t.Fatalf("Generation = %d, want 1", snapshot.Generation)
	}

	secondCancel()
	if err := waitForRunDone(t, secondDone); err != nil {
		t.Fatalf("second Run() = %v, want nil", err)
	}
}

func TestRunnerRunWithAlreadyCanceledContextReturnsNilAndStoresNoSnapshot(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	runner := newTestRunner(t, newTestClock())

	if err := runner.Run(ctx); err != nil {
		t.Fatalf("Run(canceled ctx) = %v, want nil", err)
	}
	if _, ok := runner.Snapshot(health.TargetReady); ok {
		t.Fatal("Snapshot() ok = true, want false")
	}
}

func TestRunnerRunDoesNotStoreCancellationArtifacts(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	released := make(chan struct{})
	evaluator := healthtest.NewEvaluatorForTarget(
		t,
		health.TargetReady,
		healthtest.FuncChecker("ready_check", func(ctx context.Context) health.Result {
			close(started)
			<-ctx.Done()
			close(released)
			return health.Unknown("ready_check", health.ReasonCanceled, "canceled")
		}),
	)
	runner, err := NewRunner(evaluator, WithClock(newTestClock()), WithTargets(health.TargetReady))
	if err != nil {
		t.Fatalf("NewRunner() = %v, want nil", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- runner.Run(ctx)
	}()

	select {
	case <-started:
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for health check to start")
	}

	cancel()

	select {
	case <-released:
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for health check to observe cancellation")
	}
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
	if _, ok := runner.Snapshot(health.TargetReady); ok {
		t.Fatal("Snapshot() ok = true, want false")
	}
}

func TestRunnerRunNilReceiver(t *testing.T) {
	t.Parallel()

	var runner *Runner
	err := runner.Run(context.Background())

	if !errors.Is(err, ErrNilRunner) {
		t.Fatalf("Run() = %v, want ErrNilRunner", err)
	}
}

func TestRunnerRunPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(t, newTestClock())

	defer func() {
		recovered := recover()
		if recovered != "healthprobe: nil context" {
			t.Fatalf("Run(nil) panic = %v, want healthprobe: nil context", recovered)
		}
	}()

	_ = runner.Run(nil)
}

func TestRunnerConcurrentReadDuringRun(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := time.Second
	clk := newTestClock()
	runner := newTestRunner(t, clk, WithInterval(interval), WithInitialProbe(false))
	done := make(chan error, 1)
	go func() {
		done <- runner.Run(ctx)
	}()
	waitForRunnerRunning(t, runner)

	var readers sync.WaitGroup
	for i := 0; i < 4; i++ {
		readers.Add(1)
		go func() {
			defer readers.Done()
			for j := 0; j < 50; j++ {
				_, _ = runner.Snapshot(health.TargetReady)
				_ = runner.Snapshots()
			}
		}()
	}

	for i := 0; i < 50; i++ {
		clk.Step(interval)
	}
	readers.Wait()
	_ = stepUntilGeneration(t, clk, runner, health.TargetReady, 1, interval)

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}
