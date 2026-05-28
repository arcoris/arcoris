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
	"arcoris.dev/healthtest"
)

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

func BenchmarkStaleCalculation(b *testing.B) {
	age := 2 * time.Second
	staleAfter := time.Second
	b.ReportAllocs()
	for b.Loop() {
		_ = isStale(age, staleAfter)
	}
}

func newBenchmarkRunner(b *testing.B) *Runner {
	b.Helper()

	evaluator := healthtest.NewEvaluatorWithReports(
		healthtest.HealthyReport(health.TargetStartup),
		healthtest.HealthyReport(health.TargetLive),
		healthtest.HealthyReport(health.TargetReady),
	)
	runner, err := NewRunner(evaluator, WithTargets(health.TargetStartup, health.TargetLive, health.TargetReady))
	if err != nil {
		b.Fatalf("NewRunner() = %v", err)
	}
	runner.runCycle(context.Background())
	return runner
}
