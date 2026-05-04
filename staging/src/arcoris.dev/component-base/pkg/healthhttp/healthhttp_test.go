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
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/health"
)

func TestLowLevelPrimitives(t *testing.T) {
	t.Parallel()

	if path, ok := DefaultPath(health.TargetReady); !ok || path != DefaultReadyPath {
		t.Fatalf("DefaultPath(ready) = %q, %v", path, ok)
	}
	for _, path := range []string{DefaultStartupPath, DefaultLivePath, DefaultReadyPath, DefaultHealthPath, DefaultHealthPlainPath} {
		if err := ValidatePath(path); err != nil {
			t.Fatalf("ValidatePath(%q) = %v", path, err)
		}
	}
	for _, path := range []string{"", "readyz", "/", "/readyz?x=1", "/readyz#x", "http://host/readyz", "//host/readyz", "/readyz debug", "/readyz\\debug"} {
		if err := ValidatePath(path); !errors.Is(err, ErrInvalidPath) {
			t.Fatalf("ValidatePath(%q) = %v, want ErrInvalidPath", path, err)
		}
	}

	if !methodAllowed(http.MethodGet) || !methodAllowed(http.MethodHead) || methodAllowed(http.MethodPost) {
		t.Fatal("methodAllowed mismatch")
	}

	if FormatText.String() != "text" || FormatJSON.String() != "json" || Format(99).String() != "invalid" {
		t.Fatal("format String mismatch")
	}
	if !FormatText.IsValid() || !FormatJSON.IsValid() || Format(99).IsValid() {
		t.Fatal("format IsValid mismatch")
	}
	if FormatText.contentType() != contentTypeText || FormatJSON.contentType() != contentTypeJSON {
		t.Fatal("format content type mismatch")
	}

	if DetailNone.IncludesChecks() || !DetailFailed.IncludesFailedChecks() || !DetailAll.IncludesAllChecks() || DetailLevel(99).IsValid() {
		t.Fatal("detail predicates mismatch")
	}

	if err := DefaultStatusCodes().Validate(); err != nil {
		t.Fatalf("DefaultStatusCodes().Validate() = %v", err)
	}
	if err := (HTTPStatusCodes{Failed: http.StatusOK}).Validate(); !errors.Is(err, ErrInvalidHTTPStatusCode) {
		t.Fatalf("invalid status mapping = %v", err)
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	config := defaultConfig(health.TargetReady)
	err := applyOptions(
		&config,
		WithFormat(FormatJSON),
		WithDetailLevel(DetailAll),
		WithPolicy(health.ReadyPolicy().WithDegraded(true)),
		WithFailedStatus(http.StatusTooManyRequests),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v", err)
	}
	if config.format != FormatJSON || config.detailLevel != DetailAll || !config.policy.Passes(health.StatusDegraded) || config.statusCodes.Failed != http.StatusTooManyRequests {
		t.Fatalf("unexpected config: %+v", config)
	}
	if err := applyOptions(&config, nil); !errors.Is(err, ErrNilOption) {
		t.Fatalf("nil option = %v, want ErrNilOption", err)
	}
}

func TestResponseAndRenderingDoNotExposeCause(t *testing.T) {
	t.Parallel()

	report := testReport()
	policy := health.ReadyPolicy()

	response := newResponse(report, report.Passed(policy), policy, DetailAll)
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("json.Marshal() = %v", err)
	}
	if strings.Contains(string(data), "private cause") || strings.Contains(string(data), "Cause") || strings.Contains(string(data), "cause") {
		t.Fatalf("response exposes cause: %s", string(data))
	}

	config := defaultConfig(health.TargetReady)
	config.format = FormatJSON
	config.detailLevel = DetailAll
	recorder := httptest.NewRecorder()
	renderReport(recorder, httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil), config, report, false)
	if strings.Contains(recorder.Body.String(), "private cause") {
		t.Fatalf("render exposes cause: %q", recorder.Body.String())
	}
}

func TestHandler(t *testing.T) {
	t.Parallel()

	evaluator := mustEvaluator(t, health.TargetReady, health.Degraded("queue", health.ReasonOverloaded, "queue overloaded"))
	handler, err := NewHandler(evaluator, health.TargetReady)
	if err != nil {
		t.Fatalf("NewHandler() = %v", err)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil))
	if recorder.Code != DefaultFailedStatus || recorder.Body.String() != textUnhealthy {
		t.Fatalf("ready degraded response = %d %q", recorder.Code, recorder.Body.String())
	}

	handler, err = NewHandler(evaluator, health.TargetReady, WithPolicy(health.ReadyPolicy().WithDegraded(true)))
	if err != nil {
		t.Fatalf("NewHandler(custom) = %v", err)
	}
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil))
	if recorder.Code != DefaultPassedStatus {
		t.Fatalf("custom policy status = %d, want %d", recorder.Code, DefaultPassedStatus)
	}

	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodHead, DefaultReadyPath, nil))
	if recorder.Body.Len() != 0 {
		t.Fatalf("HEAD body length = %d", recorder.Body.Len())
	}

	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, DefaultReadyPath, nil))
	if recorder.Code != http.StatusMethodNotAllowed || recorder.Result().Header.Get(headerAllow) != allowedMethodsHeader {
		t.Fatalf("POST response = %d allow=%q", recorder.Code, recorder.Result().Header.Get(headerAllow))
	}
}

func TestInstall(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := mustDefaultsEvaluator(t)

	if err := InstallDefaults(mux, evaluator, WithFormat(FormatJSON)); err != nil {
		t.Fatalf("InstallDefaults() = %v", err)
	}

	for _, path := range []string{DefaultStartupPath, DefaultLivePath, DefaultReadyPath} {
		handler, ok := mux.handler(path)
		if !ok {
			t.Fatalf("handler for %s not registered", path)
		}
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != DefaultPassedStatus {
			t.Fatalf("%s status = %d", path, recorder.Code)
		}
	}
	if _, ok := mux.handler(DefaultHealthPath); ok {
		t.Fatalf("DefaultHealthPath must not be registered by InstallDefaults")
	}
	if err := Install(nil, DefaultReadyPath, evaluator, health.TargetReady); !errors.Is(err, ErrNilMux) {
		t.Fatalf("Install(nil mux) = %v", err)
	}
	if err := Install(newRecordingMux(), "readyz", evaluator, health.TargetReady); !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("Install(invalid path) = %v", err)
	}
}

func TestNilAndTypedNilMux(t *testing.T) {
	t.Parallel()

	var typedNil *recordingMux
	if !nilMux(nil) || !nilMux(typedNil) || nilMux(newRecordingMux()) || nilMux(http.NewServeMux()) {
		t.Fatal("nilMux mismatch")
	}
}

type recordingMux struct {
	handlers map[string]http.Handler
}

func newRecordingMux() *recordingMux {
	return &recordingMux{handlers: make(map[string]http.Handler)}
}

func (mux *recordingMux) Handle(pattern string, handler http.Handler) {
	mux.handlers[pattern] = handler
}

func (mux *recordingMux) handler(pattern string) (http.Handler, bool) {
	handler, ok := mux.handlers[pattern]
	return handler, ok
}

func testReport() health.Report {
	observed := time.Date(2026, 5, 4, 12, 30, 15, 123456789, time.UTC)
	return health.Report{
		Target:   health.TargetReady,
		Status:   health.StatusUnhealthy,
		Observed: observed,
		Duration: 25 * time.Millisecond,
		Checks: []health.Result{
			health.Healthy("storage").WithObserved(observed).WithDuration(2 * time.Millisecond),
			health.Degraded("queue", health.ReasonOverloaded, "queue overloaded").WithObserved(observed).WithDuration(3 * time.Millisecond),
			health.Unhealthy("database", health.ReasonDependencyUnavailable, "database unavailable").
				WithObserved(observed).
				WithDuration(5 * time.Millisecond).
				WithCause(errors.New("private cause")),
		},
	}
}

func mustEvaluator(t *testing.T, target health.Target, result health.Result) *health.Evaluator {
	t.Helper()

	checker, err := health.NewCheck(result.Name, func(context.Context) health.Result { return result })
	if err != nil {
		t.Fatalf("health.NewCheck() = %v", err)
	}

	registry := health.NewRegistry()
	if err := registry.Register(target, checker); err != nil {
		t.Fatalf("registry.Register() = %v", err)
	}

	evaluator, err := health.NewEvaluator(registry, health.WithDefaultTimeout(0))
	if err != nil {
		t.Fatalf("health.NewEvaluator() = %v", err)
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
		{health.TargetStartup, "startup"},
		{health.TargetLive, "live"},
		{health.TargetReady, "ready"},
	} {
		item := item
		checker, err := health.NewCheck(item.name, func(context.Context) health.Result {
			return health.Healthy(item.name)
		})
		if err != nil {
			t.Fatalf("health.NewCheck(%s) = %v", item.name, err)
		}
		if err := registry.Register(item.target, checker); err != nil {
			t.Fatalf("registry.Register(%s) = %v", item.target, err)
		}
	}

	evaluator, err := health.NewEvaluator(registry, health.WithDefaultTimeout(0))
	if err != nil {
		t.Fatalf("health.NewEvaluator() = %v", err)
	}

	return evaluator
}
