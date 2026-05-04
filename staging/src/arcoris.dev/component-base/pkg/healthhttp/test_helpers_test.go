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

package healthhttp

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

// recordingMux is a small test mux that records Handle calls.
//
// It intentionally does not implement routing. Handler behavior is exercised
// separately through net/http/httptest against the registered handlers.
type recordingMux struct {
	handlers map[string]http.Handler
}

func newRecordingMux() *recordingMux {
	return &recordingMux{
		handlers: make(map[string]http.Handler),
	}
}

func (mux *recordingMux) Handle(pattern string, handler http.Handler) {
	mux.handlers[pattern] = handler
}

func (mux *recordingMux) handler(pattern string) (http.Handler, bool) {
	handler, ok := mux.handlers[pattern]
	return handler, ok
}

func mustTestEvaluator(t *testing.T, target health.Target, result health.Result) *health.Evaluator {
	t.Helper()

	checker, err := health.NewCheck(result.Name, func(context.Context) health.Result {
		return result
	})
	if err != nil {
		t.Fatalf("health.NewCheck() = %v, want nil", err)
	}

	registry := health.NewRegistry()
	if err := registry.Register(target, checker); err != nil {
		t.Fatalf("registry.Register() = %v, want nil", err)
	}

	evaluator, err := health.NewEvaluator(registry, health.WithDefaultTimeout(0))
	if err != nil {
		t.Fatalf("health.NewEvaluator() = %v, want nil", err)
	}

	return evaluator
}

func mustDefaultsEvaluator(t *testing.T) *health.Evaluator {
	t.Helper()

	registry := health.NewRegistry()

	for _, item := range []struct {
		target health.Target
		name   string
	}{
		{target: health.TargetStartup, name: "startup"},
		{target: health.TargetLive, name: "live"},
		{target: health.TargetReady, name: "ready"},
	} {
		checker, err := health.NewCheck(item.name, func(context.Context) health.Result {
			return health.Healthy(item.name)
		})
		if err != nil {
			t.Fatalf("health.NewCheck(%s) = %v, want nil", item.name, err)
		}
		if err := registry.Register(item.target, checker); err != nil {
			t.Fatalf("registry.Register(%s, %s) = %v, want nil", item.target, item.name, err)
		}
	}

	evaluator, err := health.NewEvaluator(registry, health.WithDefaultTimeout(0))
	if err != nil {
		t.Fatalf("health.NewEvaluator() = %v, want nil", err)
	}

	return evaluator
}

func mustNewHandler(t *testing.T, evaluator *health.Evaluator, target health.Target, options ...Option) *Handler {
	t.Helper()

	handler, err := NewHandler(evaluator, target, options...)
	if err != nil {
		t.Fatalf("NewHandler() = %v, want nil", err)
	}

	return handler
}

func testObservedTime() time.Time {
	return time.Date(2026, 5, 4, 12, 30, 15, 123456789, time.UTC)
}

func testReport() health.Report {
	observed := testObservedTime()

	return health.Report{
		Target:   health.TargetReady,
		Status:   health.StatusUnhealthy,
		Observed: observed,
		Duration: 25 * time.Millisecond,
		Checks: []health.Result{
			health.Healthy("storage").WithObserved(observed).WithDuration(2 * time.Millisecond),
			health.Degraded("queue", health.ReasonOverloaded, "queue is above soft capacity").
				WithObserved(observed).
				WithDuration(3 * time.Millisecond),
			health.Unknown("cache", health.ReasonNotObserved, "cache has not reported yet").
				WithObserved(observed).
				WithDuration(4 * time.Millisecond),
			health.Unhealthy("database", health.ReasonDependencyUnavailable, "database dependency is unavailable").
				WithObserved(observed).
				WithDuration(5 * time.Millisecond).
				WithCause(errors.New("private cause")),
		},
	}
}
