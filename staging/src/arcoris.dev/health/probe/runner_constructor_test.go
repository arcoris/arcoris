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
	"errors"
	"testing"
	"time"

	"arcoris.dev/health"
)

func TestNewRunner(t *testing.T) {
	t.Parallel()

	evaluator := newTestEvaluator(t)
	clk := newTestClock()

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
	if d := firstScheduleDelay(t, runner.schedule); d != time.Second {
		t.Fatalf("schedule delay = %s, want 1s", d)
	}
	if runner.staleAfter != 2*time.Second {
		t.Fatalf("staleAfter = %s, want 2s", runner.staleAfter)
	}
	if runner.initialProbe {
		t.Fatal("initialProbe = true, want false")
	}
}

func TestNewRunnerDoesNotStartGoroutines(t *testing.T) {
	t.Parallel()

	runner := newTestRunner(t, newTestClock(), WithInitialProbe(false))

	if runner.running.Load() {
		t.Fatal("runner is running before Run")
	}
	if _, ok := runner.Snapshot(health.TargetReady); ok {
		t.Fatal("Snapshot() ok = true before Run, want false")
	}
}

func TestNewRunnerRejectsInvalidConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		eval    Evaluator
		opts    []Option
		wantErr error
	}{
		{
			name:    "nil evaluator",
			eval:    nil,
			opts:    []Option{WithTargets(health.TargetReady)},
			wantErr: ErrNilEvaluator,
		},
		{
			name:    "missing targets",
			eval:    newTestEvaluator(t),
			opts:    nil,
			wantErr: ErrNoTargets,
		},
		{
			name:    "invalid target",
			eval:    newTestEvaluator(t),
			opts:    []Option{WithTargets(health.TargetUnknown)},
			wantErr: health.ErrInvalidTarget,
		},
		{
			name:    "duplicate target",
			eval:    newTestEvaluator(t),
			opts:    []Option{WithTargets(health.TargetReady, health.TargetReady)},
			wantErr: ErrDuplicateTarget,
		},
		{
			name: "invalid interval",
			eval: newTestEvaluator(t),
			opts: []Option{
				WithTargets(health.TargetReady),
				WithInterval(0),
			},
			wantErr: ErrInvalidInterval,
		},
		{
			name: "invalid stale after",
			eval: newTestEvaluator(t),
			opts: []Option{
				WithTargets(health.TargetReady),
				WithStaleAfter(-time.Nanosecond),
			},
			wantErr: ErrInvalidStaleAfter,
		},
		{
			name: "nil option",
			eval: newTestEvaluator(t),
			opts: []Option{
				WithTargets(health.TargetReady),
				nil,
			},
			wantErr: ErrNilOption,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewRunner(tc.eval, tc.opts...)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewRunner() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestNewRunnerCopiesTargetList(t *testing.T) {
	t.Parallel()

	targets := []health.Target{health.TargetReady, health.TargetLive}
	runner, err := NewRunner(newTestEvaluator(t), WithClock(newTestClock()), WithTargets(targets...))
	if err != nil {
		t.Fatalf("NewRunner() = %v, want nil", err)
	}

	targets[0] = health.TargetStartup

	if !sameTargets(runner.targets, []health.Target{health.TargetReady, health.TargetLive}) {
		t.Fatalf("runner targets = %v, want [ready live]", runner.targets)
	}
	if !sameTargets(runner.store.targets, []health.Target{health.TargetReady, health.TargetLive}) {
		t.Fatalf("store targets = %v, want [ready live]", runner.store.targets)
	}
}
