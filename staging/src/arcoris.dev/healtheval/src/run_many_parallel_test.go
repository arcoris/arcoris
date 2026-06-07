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

package eval

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/healthregistry"
)

func TestEvaluatorParallelPreservesResolverOrderWhenChecksFinishOutOfOrder(t *testing.T) {
	t.Parallel()

	registry := healthregistry.NewBuilder()
	releaseFirst := make(chan struct{})
	firstStarted := make(chan struct{})
	secondDone := make(chan struct{})
	thirdDone := make(chan struct{})

	mustRegisterExecutionCheck(t, registry, health.TargetReady, "first", func(context.Context) health.Result {
		close(firstStarted)
		<-releaseFirst
		return health.Healthy("first")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "second", func(context.Context) health.Result {
		close(secondDone)
		return health.Healthy("second")
	})
	mustRegisterExecutionCheck(t, registry, health.TargetReady, "third", func(context.Context) health.Result {
		close(thirdDone)
		return health.Healthy("third")
	})

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(health.TargetReady, 3),
	)

	done := make(chan health.Report, 1)
	go func() {
		report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
		if err != nil {
			t.Errorf("Evaluate() = %v, want nil", err)
		}
		done <- report
	}()

	<-firstStarted
	<-secondDone
	<-thirdDone
	close(releaseFirst)

	var report health.Report
	select {
	case report = <-done:
	case <-time.After(executionTestTimeout):
		t.Fatal("parallel evaluation did not finish")
	}

	names := executionResultNames(report.Checks)
	want := []string{"first", "second", "third"}

	if !sameStrings(names, want) {
		t.Fatalf("result names = %v, want %v", names, want)
	}
}

func TestEvaluatorParallelRespectsMaxConcurrency(t *testing.T) {
	t.Parallel()

	const checkCount = 8
	const limit = 3

	registry := healthregistry.NewBuilder()
	release := make(chan struct{})
	started := make(chan struct{}, checkCount)

	var active atomic.Int64
	var maxSeen atomic.Int64

	for i := 0; i < checkCount; i++ {
		name := executionCheckName(i)
		mustRegisterExecutionCheck(t, registry, health.TargetReady, name, func(context.Context) health.Result {
			cur := active.Add(1)
			updateMaxInt64(&maxSeen, cur)
			started <- struct{}{}

			<-release

			active.Add(-1)
			return health.Healthy(name)
		})
	}

	evaluator := mustExecutionEvaluator(
		t,
		registry,
		WithDefaultTimeout(0),
		WithTargetParallelChecks(health.TargetReady, limit),
	)

	done := make(chan health.Report, 1)
	go func() {
		report, err := evaluator.Evaluate(context.Background(), health.TargetReady)
		if err != nil {
			t.Errorf("Evaluate() = %v, want nil", err)
		}
		done <- report
	}()

	for i := 0; i < limit; i++ {
		<-started
	}

	if got := maxSeen.Load(); got != limit {
		t.Fatalf("max concurrency = %d, want exactly %d", got, limit)
	}

	close(release)

	var report health.Report
	select {
	case report = <-done:
	case <-time.After(executionTestTimeout):
		t.Fatal("parallel evaluation did not finish")
	}

	if got := maxSeen.Load(); got > limit {
		t.Fatalf("max concurrency after completion = %d, want <= %d", got, limit)
	}
	if len(report.Checks) != checkCount {
		t.Fatalf("checks = %d, want %d", len(report.Checks), checkCount)
	}
}
