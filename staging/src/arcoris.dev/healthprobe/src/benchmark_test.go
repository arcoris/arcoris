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
	"fmt"
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/healthtest"
)

func BenchmarkRunnerRunCycleOneTarget(b *testing.B) {
	runner := newBenchmarkRunnerForTargets(b, []health.Target{health.TargetReady})
	b.ReportAllocs()
	for b.Loop() {
		runner.runCycle(context.Background())
	}
}

func BenchmarkRunnerRunCycleThreeTargets(b *testing.B) {
	runner := newBenchmarkRunnerForTargets(
		b,
		[]health.Target{health.TargetStartup, health.TargetLive, health.TargetReady},
	)
	b.ReportAllocs()
	for b.Loop() {
		runner.runCycle(context.Background())
	}
}

func BenchmarkRunnerSnapshot(b *testing.B) {
	runner := newBenchmarkRunner(b)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = runner.Snapshot(health.TargetReady)
	}
}

func BenchmarkRunnerSnapshots(b *testing.B) {
	runner := newBenchmarkRunner(b)
	b.ReportAllocs()
	for b.Loop() {
		_ = runner.Snapshots()
	}
}

func BenchmarkRunnerSnapshotFresh(b *testing.B) {
	runner := newBenchmarkRunner(b, WithStaleAfter(time.Hour))
	b.ReportAllocs()
	for b.Loop() {
		_, _ = runner.Snapshot(health.TargetReady)
	}
}

func BenchmarkRunnerSnapshotStale(b *testing.B) {
	clk := newTestClock()
	runner := newBenchmarkRunner(b, WithClock(clk), WithStaleAfter(time.Nanosecond))
	clk.Step(time.Second)

	b.ReportAllocs()
	for b.Loop() {
		_, _ = runner.Snapshot(health.TargetReady)
	}
}

func BenchmarkStoreUpdate(b *testing.B) {
	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	report := healthtest.HealthyReport(health.TargetReady)

	b.ReportAllocs()
	for b.Loop() {
		_ = s.update(health.TargetReady, report)
	}
}

func BenchmarkStoreSnapshot(b *testing.B) {
	s := newStore([]health.Target{health.TargetReady}, newTestClock())
	_ = s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady))

	b.ReportAllocs()
	for b.Loop() {
		_, _ = s.snapshot(health.TargetReady)
	}
}

func BenchmarkStoreSnapshots(b *testing.B) {
	s := newStore(
		[]health.Target{health.TargetStartup, health.TargetLive, health.TargetReady},
		newTestClock(),
	)
	_ = s.update(health.TargetStartup, healthtest.HealthyReport(health.TargetStartup))
	_ = s.update(health.TargetLive, healthtest.HealthyReport(health.TargetLive))
	_ = s.update(health.TargetReady, healthtest.HealthyReport(health.TargetReady))

	b.ReportAllocs()
	for b.Loop() {
		_ = s.snapshots()
	}
}

func BenchmarkSnapshotCloneSmallReport(b *testing.B) {
	snap := benchmarkSnapshotWithChecks(1)
	b.ReportAllocs()
	for b.Loop() {
		_ = cloneSnapshot(snap)
	}
}

func BenchmarkSnapshotCloneLargeReport(b *testing.B) {
	snap := benchmarkSnapshotWithChecks(64)
	b.ReportAllocs()
	for b.Loop() {
		_ = cloneSnapshot(snap)
	}
}

func BenchmarkStaleCalculation(b *testing.B) {
	age := 2 * time.Second
	staleAfter := time.Second
	b.ReportAllocs()
	for b.Loop() {
		_ = isStale(age, staleAfter)
	}
}

func newBenchmarkRunner(b *testing.B, opts ...Option) *Runner {
	b.Helper()

	runner := newBenchmarkRunnerForTargets(
		b,
		[]health.Target{health.TargetStartup, health.TargetLive, health.TargetReady},
		opts...,
	)
	runner.runCycle(context.Background())
	return runner
}

func newBenchmarkRunnerForTargets(
	b *testing.B,
	targets []health.Target,
	opts ...Option,
) *Runner {
	b.Helper()

	reports := make([]health.Report, 0, len(targets))
	for _, target := range targets {
		reports = append(reports, healthtest.HealthyReport(target))
	}

	evaluator := healthtest.NewEvaluatorWithReports(
		reports...,
	)

	options := append([]Option{WithTargets(targets...)}, opts...)
	runner, err := NewRunner(evaluator, options...)
	if err != nil {
		b.Fatalf("NewRunner() = %v", err)
	}

	return runner
}

func benchmarkSnapshotWithChecks(count int) Snapshot {
	observed := testNow
	checks := make([]health.Result, 0, count)
	for i := 0; i < count; i++ {
		checks = append(checks, health.Healthy(fmt.Sprintf("check_%d", i)).WithObserved(observed))
	}

	return Snapshot{
		Target: health.TargetReady,
		Report: health.Report{
			Target:   health.TargetReady,
			Status:   health.StatusHealthy,
			Observed: observed,
			Checks:   checks,
		},
		Revision: 1,
		Updated:  observed,
	}
}
