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
	"fmt"
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/healthregistry"
)

func newBenchmarkEvaluator(b *testing.B, opts ...EvaluatorOption) *Evaluator {
	b.Helper()

	return newBenchmarkEvaluatorWithChecks(b, 4, opts...)
}

func newBenchmarkEvaluatorWithChecks(b *testing.B, count int, opts ...EvaluatorOption) *Evaluator {
	b.Helper()

	builder := healthregistry.NewBuilder()
	for i := 0; i < count; i++ {
		checkName := fmt.Sprintf("check_%d", i)
		checker, err := health.NewCheck(checkName, func(context.Context) health.Result {
			return health.Healthy(checkName)
		})
		if err != nil {
			b.Fatalf("NewCheck() = %v", err)
		}
		if err := builder.Register(health.TargetReady, checker); err != nil {
			b.Fatalf("Register() = %v", err)
		}
	}

	return newBenchmarkEvaluatorFromRegistry(b, builder, opts...)
}

func newBenchmarkEvaluatorWithFunc(
	b *testing.B,
	name string,
	fn health.CheckFunc,
	opts ...EvaluatorOption,
) *Evaluator {
	b.Helper()

	builder := healthregistry.NewBuilder()
	checker, err := health.NewCheck(name, fn)
	if err != nil {
		b.Fatalf("NewCheck() = %v", err)
	}
	if err := builder.Register(health.TargetReady, checker); err != nil {
		b.Fatalf("Register() = %v", err)
	}

	return newBenchmarkEvaluatorFromRegistry(b, builder, opts...)
}

func newBenchmarkEvaluatorFromRegistry(
	b *testing.B,
	builder *healthregistry.Builder,
	opts ...EvaluatorOption,
) *Evaluator {
	b.Helper()

	registry, err := builder.Build()
	if err != nil {
		b.Fatalf("Build() = %v", err)
	}

	evaluator, err := NewEvaluator(registry, opts...)
	if err != nil {
		b.Fatalf("NewEvaluator() = %v", err)
	}
	return evaluator
}
