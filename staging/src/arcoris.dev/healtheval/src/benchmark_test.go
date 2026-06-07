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
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/healthregistry"
)

func BenchmarkEvaluateSequentialOneCheck(b *testing.B) {
	evaluator := newBenchmarkEvaluatorWithChecks(b, 1, WithDefaultTimeout(0))
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluateSequentialManyChecks(b *testing.B) {
	evaluator := newBenchmarkEvaluatorWithChecks(b, 16, WithDefaultTimeout(0))
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluateParallelManyChecks(b *testing.B) {
	evaluator := newBenchmarkEvaluatorWithChecks(
		b,
		16,
		WithDefaultTimeout(0),
		WithExecutionPolicy(ParallelExecutionPolicy(4)),
	)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluateTimeoutDisabled(b *testing.B) {
	evaluator := newBenchmarkEvaluatorWithChecks(b, 4, WithDefaultTimeout(0))
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluateTimeoutEnabled(b *testing.B) {
	evaluator := newBenchmarkEvaluatorWithFunc(
		b,
		"timeout_check",
		func(ctx context.Context) health.Result {
			<-ctx.Done()
			return health.Unknown("timeout_check", health.ReasonTimeout, "timed out")
		},
		WithDefaultTimeout(time.Nanosecond),
	)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluateTimeoutEnabledFastCheck(b *testing.B) {
	evaluator := newBenchmarkEvaluatorWithFunc(
		b,
		"fast_timeout_check",
		func(context.Context) health.Result {
			return health.Healthy("fast_timeout_check")
		},
		WithDefaultTimeout(time.Second),
	)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluatePanicRecovery(b *testing.B) {
	evaluator := newBenchmarkEvaluatorWithFunc(
		b,
		"panic_check",
		func(context.Context) health.Result {
			panic("benchmark panic")
		},
		WithDefaultTimeout(0),
	)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluateResolverMismatch(b *testing.B) {
	checker := health.MustCheck("wrong_target_check", func(context.Context) health.Result {
		return health.Healthy("wrong_target_check")
	})
	wrongSet := health.MustCheckSet(health.TargetLive, checker)
	resolver := health.CheckResolverFunc(func(health.Target) (health.CheckSet, error) {
		return wrongSet, nil
	})
	evaluator, err := NewEvaluator(resolver, WithDefaultTimeout(0))
	if err != nil {
		b.Fatalf("NewEvaluator() = %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluateResolverEmpty(b *testing.B) {
	evaluator := newBenchmarkEvaluatorFromRegistry(b, healthregistry.NewBuilder())
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluatorSequential(b *testing.B) {
	evaluator := newBenchmarkEvaluator(b)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}

func BenchmarkEvaluatorParallel(b *testing.B) {
	evaluator := newBenchmarkEvaluator(b, WithExecutionPolicy(ParallelExecutionPolicy(4)))
	b.ReportAllocs()
	for b.Loop() {
		_, _ = evaluator.Evaluate(context.Background(), health.TargetReady)
	}
}
