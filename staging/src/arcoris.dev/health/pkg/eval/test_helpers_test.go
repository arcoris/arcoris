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

package eval

import (
	"context"
	"testing"
	"time"

	"arcoris.dev/health"
)

const testTimeout = time.Second

var testObserved = time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)

type checkerFunc struct {
	name string
	fn   func(context.Context) health.Result
}

func (checker checkerFunc) Name() string {
	return checker.name
}

func (checker checkerFunc) Check(ctx context.Context) health.Result {
	return checker.fn(ctx)
}

type stepClock struct {
	values  []time.Time
	next    int
	current time.Time
}

func newStepClock(values ...time.Time) *stepClock {
	return &stepClock{values: values}
}

func (clk *stepClock) Now() time.Time {
	if len(clk.values) == 0 {
		return time.Time{}
	}
	if clk.next >= len(clk.values) {
		clk.current = clk.values[len(clk.values)-1]
		return clk.current
	}

	val := clk.values[clk.next]
	clk.next++
	clk.current = val

	return val
}

func (clk *stepClock) Since(ts time.Time) time.Duration {
	return clk.current.Sub(ts)
}

func mustCheck(t *testing.T, name string, res health.Result) health.Checker {
	t.Helper()

	checker, err := health.NewCheck(name, func(context.Context) health.Result {
		return res
	})
	if err != nil {
		t.Fatalf("health.NewCheck(%q) = %v, want nil", name, err)
	}

	return checker
}

func mustRegistry(t *testing.T, target health.Target, checks ...health.Checker) *health.Registry {
	t.Helper()

	registry := health.NewRegistry()
	if err := registry.Register(target, checks...); err != nil {
		t.Fatalf("Register() = %v, want nil", err)
	}

	return registry
}

func mustEvaluator(t *testing.T, r *health.Registry, opts ...EvaluatorOption) *Evaluator {
	t.Helper()

	evaluator, err := NewEvaluator(r, opts...)
	if err != nil {
		t.Fatalf("NewEvaluator() = %v, want nil", err)
	}

	return evaluator
}
