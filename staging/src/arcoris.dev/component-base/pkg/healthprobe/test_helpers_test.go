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
	"sync"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/clock"
	"arcoris.dev/component-base/pkg/health"
)

const testTimeout = 5 * time.Second

type manualClock struct {
	mu      sync.Mutex
	now     time.Time
	tickers []*manualTicker
}

func newManualClock() *manualClock {
	return &manualClock{now: time.Unix(100, 0)}
}

func (c *manualClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *manualClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

func (c *manualClock) After(d time.Duration) <-chan time.Time {
	timer := time.NewTimer(d)
	return timer.C
}

func (c *manualClock) NewTimer(d time.Duration) clock.Timer {
	return realTimer{timer: time.NewTimer(d)}
}

func (c *manualClock) NewTicker(d time.Duration) clock.Ticker {
	ticker := &manualTicker{ch: make(chan time.Time, 16)}
	c.mu.Lock()
	c.tickers = append(c.tickers, ticker)
	c.mu.Unlock()
	return ticker
}

func (c *manualClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (c *manualClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

func (c *manualClock) TickAll() {
	c.mu.Lock()
	now := c.now
	tickers := append([]*manualTicker(nil), c.tickers...)
	c.mu.Unlock()

	for _, ticker := range tickers {
		ticker.tick(now)
	}
}

type manualTicker struct {
	mu      sync.Mutex
	ch      chan time.Time
	stopped bool
}

func (t *manualTicker) C() <-chan time.Time { return t.ch }

func (t *manualTicker) Stop() {
	t.mu.Lock()
	t.stopped = true
	t.mu.Unlock()
}

func (t *manualTicker) Reset(d time.Duration) {}

func (t *manualTicker) tick(now time.Time) {
	t.mu.Lock()
	stopped := t.stopped
	t.mu.Unlock()
	if stopped {
		return
	}

	select {
	case t.ch <- now:
	default:
	}
}

type realTimer struct {
	timer *time.Timer
}

func (t realTimer) C() <-chan time.Time {
	return t.timer.C
}

func (t realTimer) Stop() bool {
	return t.timer.Stop()
}

func (t realTimer) Reset(d time.Duration) bool {
	return t.timer.Reset(d)
}

func newTestEvaluator(t *testing.T) *health.Evaluator {
	t.Helper()

	return newEvaluatorWithChecks(t, map[health.Target]health.CheckFunc{
		health.TargetReady: func(context.Context) health.Result {
			return health.Healthy("ready_check")
		},
	})
}

func newEvaluatorWithChecks(t *testing.T, checks map[health.Target]health.CheckFunc) *health.Evaluator {
	t.Helper()

	registry := health.NewRegistry()
	for target, fn := range checks {
		check, err := health.NewCheck(target.String()+"_check", fn)
		if err != nil {
			t.Fatalf("NewCheck() = %v, want nil", err)
		}
		if err := registry.Register(target, check); err != nil {
			t.Fatalf("Register() = %v, want nil", err)
		}
	}

	evaluator, err := health.NewEvaluator(registry, health.WithDefaultTimeout(0))
	if err != nil {
		t.Fatalf("NewEvaluator() = %v, want nil", err)
	}

	return evaluator
}

func newTestRunner(t *testing.T, clk *manualClock, options ...Option) *Runner {
	t.Helper()

	allOptions := []Option{WithClock(clk), WithTargets(health.TargetReady)}
	allOptions = append(allOptions, options...)

	runner, err := NewRunner(newTestEvaluator(t), allOptions...)
	if err != nil {
		t.Fatalf("NewRunner() = %v, want nil", err)
	}

	return runner
}

func waitForTicker(t *testing.T, clk *manualClock) {
	t.Helper()

	deadline := time.After(testTimeout)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for ticker creation")
		case <-ticker.C:
			clk.mu.Lock()
			count := len(clk.tickers)
			clk.mu.Unlock()
			if count > 0 {
				return
			}
		}
	}
}

func waitForSnapshot(t *testing.T, runner *Runner, target health.Target) {
	t.Helper()

	deadline := time.After(testTimeout)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for snapshot target=%s", target)
		case <-ticker.C:
			if _, ok := runner.Snapshot(target); ok {
				return
			}
		}
	}
}

func waitForRunnerRunning(t *testing.T, runner *Runner) {
	t.Helper()

	deadline := time.After(testTimeout)
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for runner to start")
		case <-ticker.C:
			if runner.running.Load() {
				return
			}
		}
	}
}

func sameHealthprobeTargets(left []health.Target, right []health.Target) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
