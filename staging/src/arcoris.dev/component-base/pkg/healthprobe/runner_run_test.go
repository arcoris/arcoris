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
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

func TestRunnerRunPerformsInitialProbe(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clk := newManualClock()
	runner := newTestRunner(t, clk)
	done := make(chan error, 1)

	go func() {
		done <- runner.Run(ctx)
	}()

	waitForSnapshot(t, runner, health.TargetReady)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run() = %v, want nil", err)
		}
	case <-time.After(testTimeout):
		t.Fatal("Run did not stop")
	}
}

func TestRunnerRunWaitsForTickWhenInitialProbeDisabled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clk := newManualClock()
	runner := newTestRunner(t, clk, WithInitialProbe(false))
	done := make(chan error, 1)

	go func() {
		done <- runner.Run(ctx)
	}()

	if _, ok := runner.Snapshot(health.TargetReady); ok {
		t.Fatal("Snapshot before tick ok = true, want false")
	}

	waitForTicker(t, clk)
	clk.TickAll()
	waitForSnapshot(t, runner, health.TargetReady)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run() = %v, want nil", err)
		}
	case <-time.After(testTimeout):
		t.Fatal("Run did not stop")
	}
}

func TestRunnerRunRejectsConcurrentRun(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clk := newManualClock()
	runner := newTestRunner(t, clk)
	started := make(chan struct{})
	done := make(chan error, 1)

	go func() {
		close(started)
		done <- runner.Run(ctx)
	}()
	<-started
	waitForRunnerRunning(t, runner)

	err := runner.Run(context.Background())
	if !errors.Is(err, ErrRunnerRunning) {
		t.Fatalf("concurrent Run() = %v, want ErrRunnerRunning", err)
	}

	cancel()
	<-done
}

func TestRunnerRunPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(t, newManualClock())

	defer func() {
		if recover() == nil {
			t.Fatal("Run(nil) did not panic")
		}
	}()

	_ = runner.Run(nil)
}
