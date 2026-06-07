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

	"arcoris.dev/health"
	"arcoris.dev/healthtest"
)

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

	evaluator := healthtest.NewEvaluatorWithReports(reports...)
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
