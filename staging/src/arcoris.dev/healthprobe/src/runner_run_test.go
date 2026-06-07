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
	"context"
	"errors"
	"testing"
	"time"

	"arcoris.dev/chrono/delay"
	"arcoris.dev/health"
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

	snap := waitForSnapshot(t, runner, health.TargetReady)
	if !snap.IsFresh() {
		t.Fatalf("snapshot IsFresh() = false, want true: %#v", snap)
	}

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}

func TestRunnerRunWaitsForScheduleDelayWhenInitialProbeDisabled(t *testing.T) {
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
		t.Fatal("Snapshot before schedule delay ok = true, want false")
	}

	snap := stepUntilRevision(t, clk, runner, health.TargetReady, 1, interval)
	if snap.Revision != 1 {
		t.Fatalf("Revision = %d, want 1", snap.Revision)
	}

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}

func TestRunnerRunScheduleDrivenProbeIncrementsRevision(t *testing.T) {
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
	first := stepUntilRevision(t, clk, runner, health.TargetReady, 1, interval)
	second := stepUntilRevision(t, clk, runner, health.TargetReady, 2, interval)

	if first.Revision != 1 {
		t.Fatalf("first Revision = %d, want 1", first.Revision)
	}
	if second.Revision != 2 {
		t.Fatalf("second Revision = %d, want 2", second.Revision)
	}
	if !second.Updated.After(first.Updated) {
		t.Fatalf("second Updated = %v, want after %v", second.Updated, first.Updated)
	}

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}

func TestRunnerRunUsesScheduleDelays(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	sched := delay.Delays(time.Second, 3*time.Second)
	clk := newTestClock()
	runner := newTestRunner(t, clk, WithSchedule(sched), WithInitialProbe(false))
	done := make(chan error, 1)

	go func() {
		done <- runner.Run(ctx)
	}()

	waitForRunnerRunning(t, runner)
	first := stepUntilRevision(t, clk, runner, health.TargetReady, 1, time.Second)
	second := stepUntilRevision(t, clk, runner, health.TargetReady, 2, 3*time.Second)

	if first.Revision != 1 {
		t.Fatalf("first Revision = %d, want 1", first.Revision)
	}
	if second.Revision != 2 {
		t.Fatalf("second Revision = %d, want 2", second.Revision)
	}
	if !second.Updated.After(first.Updated) {
		t.Fatalf("second Updated = %v, want after %v", second.Updated, first.Updated)
	}
	if err := waitForRunDone(t, done); !errors.Is(err, ErrExhaustedSchedule) {
		t.Fatalf("Run() = %v, want ErrExhaustedSchedule", err)
	}
}

func TestRunnerRunAcceptsZeroScheduleDelay(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(
		t,
		newTestClock(),
		WithSchedule(delay.Delays(0)),
		WithInitialProbe(false),
	)

	err := runner.Run(context.Background())

	if !errors.Is(err, ErrExhaustedSchedule) {
		t.Fatalf("Run() = %v, want ErrExhaustedSchedule", err)
	}
	snap, ok := runner.Snapshot(health.TargetReady)
	if !ok {
		t.Fatal("Snapshot() ok = false, want true")
	}
	if snap.Revision != 1 {
		t.Fatalf("Revision = %d, want 1", snap.Revision)
	}
}

func TestRunnerRunInitialProbeBeforeScheduleDelay(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clk := newTestClock()
	runner := newTestRunner(
		t,
		clk,
		WithSchedule(delay.Delays(time.Hour)),
		WithInitialProbe(true),
	)
	done := make(chan error, 1)

	go func() {
		done <- runner.Run(ctx)
	}()

	snap := waitForRevision(t, runner, health.TargetReady, 1)
	if snap.Revision != 1 {
		t.Fatalf("Revision = %d, want 1", snap.Revision)
	}

	cancel()
	if err := waitForRunDone(t, done); err != nil {
		t.Fatalf("Run() = %v, want nil", err)
	}
}

func TestRunnerRunReturnsExhaustedSchedule(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(
		t,
		newTestClock(),
		WithSchedule(delay.Delays()),
		WithInitialProbe(false),
	)

	err := runner.Run(context.Background())

	if !errors.Is(err, ErrExhaustedSchedule) {
		t.Fatalf("Run() = %v, want ErrExhaustedSchedule", err)
	}
}

func TestRunnerRunReturnsNilWhenContextCanceledBeforeScheduleExhaustion(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	runner := newTestRunner(
		t,
		newTestClock(),
		WithSchedule(delay.Delays()),
		WithInitialProbe(false),
	)

	if err := runner.Run(ctx); err != nil {
		t.Fatalf("Run(canceled ctx) = %v, want nil", err)
	}
}
