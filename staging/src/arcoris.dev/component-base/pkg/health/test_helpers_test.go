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

package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

const testTimeout = time.Second

var testObserved = time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)

type checkerFunc struct {
	name string
	fn   func(context.Context) Result
}

func (checker checkerFunc) Name() string {
	return checker.name
}

func (checker checkerFunc) Check(ctx context.Context) Result {
	return checker.fn(ctx)
}

type typedNilChecker struct{}

func (checker *typedNilChecker) Name() string {
	return "typed_nil"
}

func (checker *typedNilChecker) Check(context.Context) Result {
	return Healthy("typed_nil")
}

type stepClock struct {
	values  []time.Time
	next    int
	current time.Time
}

func newStepClock(values ...time.Time) *stepClock {
	return &stepClock{values: values}
}

func (clock *stepClock) Now() time.Time {
	if len(clock.values) == 0 {
		return time.Time{}
	}
	if clock.next >= len(clock.values) {
		clock.current = clock.values[len(clock.values)-1]
		return clock.current
	}

	value := clock.values[clock.next]
	clock.next++
	clock.current = value

	return value
}

func (clock *stepClock) Since(ts time.Time) time.Duration {
	return clock.current.Sub(ts)
}

func mustCheck(t *testing.T, name string, result Result) Checker {
	t.Helper()

	checker, err := NewCheck(name, func(context.Context) Result {
		return result
	})
	if err != nil {
		t.Fatalf("NewCheck(%q) = %v, want nil", name, err)
	}

	return checker
}

func mustRegistry(t *testing.T, target Target, checks ...Checker) *Registry {
	t.Helper()

	registry := NewRegistry()
	if err := registry.Register(target, checks...); err != nil {
		t.Fatalf("Register() = %v, want nil", err)
	}

	return registry
}

func mustEvaluator(t *testing.T, registry *Registry, options ...EvaluatorOption) *Evaluator {
	t.Helper()

	evaluator, err := NewEvaluator(registry, options...)
	if err != nil {
		t.Fatalf("NewEvaluator() = %v, want nil", err)
	}

	return evaluator
}

func mustErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false, want true", err, target)
	}
}

func mustPanicWith(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		if got != want {
			t.Fatalf("panic = %v, want %v", got, want)
		}
	}()

	fn()
}

func mustClose(t *testing.T, ch <-chan struct{}) {
	t.Helper()

	timer := time.NewTimer(testTimeout)
	defer timer.Stop()

	select {
	case <-ch:
	case <-timer.C:
		t.Fatal("channel did not close before timeout")
	}
}
