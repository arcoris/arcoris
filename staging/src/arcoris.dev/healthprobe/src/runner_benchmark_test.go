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
	"testing"
	"time"

	"arcoris.dev/health"
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
