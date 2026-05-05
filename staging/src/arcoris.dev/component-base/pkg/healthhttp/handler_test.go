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
	"net/http/httptest"
	"strings"
	"testing"

	"arcoris.dev/component-base/pkg/health"
	"arcoris.dev/component-base/pkg/healthtest"
)

func TestNewHandler(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Healthy("storage"),
	)

	handler, err := NewHandler(evaluator, health.TargetReady)
	if err != nil {
		t.Fatalf("NewHandler() = %v, want nil", err)
	}

	if handler.Target() != health.TargetReady {
		t.Fatalf("Target() = %s, want ready", handler.Target())
	}
	if handler.Format() != FormatText {
		t.Fatalf("Format() = %s, want text", handler.Format())
	}
	if handler.DetailLevel() != DetailNone {
		t.Fatalf("DetailLevel() = %s, want none", handler.DetailLevel())
	}
}

func TestNewHandlerRejectsNilEvaluator(t *testing.T) {
	t.Parallel()

	handler, err := NewHandler(nil, health.TargetReady)
	if handler != nil {
		t.Fatalf("handler = %+v, want nil", handler)
	}
	if !errors.Is(err, ErrNilEvaluator) {
		t.Fatalf("NewHandler(nil) = %v, want ErrNilEvaluator", err)
	}
}

func TestNewHandlerRejectsInvalidTarget(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Healthy("storage"),
	)

	handler, err := NewHandler(evaluator, health.TargetUnknown)
	if handler != nil {
		t.Fatalf("handler = %+v, want nil", handler)
	}
	if !errors.Is(err, health.ErrInvalidTarget) {
		t.Fatalf("NewHandler(invalid target) = %v, want health.ErrInvalidTarget", err)
	}
}

func TestNewHandlerRejectsInvalidOption(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Healthy("storage"),
	)

	handler, err := NewHandler(evaluator, health.TargetReady, WithFormat(Format(99)))
	if handler != nil {
		t.Fatalf("handler = %+v, want nil", handler)
	}
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("NewHandler(invalid option) = %v, want ErrInvalidFormat", err)
	}
}

func TestNewHandlerAppliesOptions(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Degraded("queue", health.ReasonOverloaded, "queue overloaded"),
	)

	handler, err := NewHandler(
		evaluator,
		health.TargetReady,
		WithFormat(FormatJSON),
		WithDetailLevel(DetailAll),
		WithPolicy(health.ReadyPolicy().WithDegraded(true)),
		WithFailedStatus(http.StatusTooManyRequests),
	)
	if err != nil {
		t.Fatalf("NewHandler() = %v, want nil", err)
	}

	if handler.Format() != FormatJSON {
		t.Fatalf("Format() = %s, want json", handler.Format())
	}
	if handler.DetailLevel() != DetailAll {
		t.Fatalf("DetailLevel() = %s, want all", handler.DetailLevel())
	}
	if !handler.config.policy.Passes(health.StatusDegraded) {
		t.Fatal("configured policy should pass degraded")
	}
	if handler.config.statusCodes.Failed != http.StatusTooManyRequests {
		t.Fatalf("failed status = %d, want %d", handler.config.statusCodes.Failed, http.StatusTooManyRequests)
	}
}

func TestHandlerServeHTTPReadyHealthy(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Healthy("storage"),
	)
	handler := mustNewHandler(t, evaluator, health.TargetReady)

	request := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultPassedStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultPassedStatus)
	}
	if got := recorder.Body.String(); got != textOK {
		t.Fatalf("body = %q, want %q", got, textOK)
	}
}

func TestHandlerServeHTTPReadyDegradedFailsByDefault(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Degraded("queue", health.ReasonOverloaded, "queue overloaded"),
	)
	handler := mustNewHandler(t, evaluator, health.TargetReady)

	request := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultFailedStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultFailedStatus)
	}
	if got := recorder.Body.String(); got != textUnhealthy {
		t.Fatalf("body = %q, want %q", got, textUnhealthy)
	}
}

func TestHandlerServeHTTPLiveDegradedPassesByDefault(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetLive,
		health.Degraded("queue", health.ReasonOverloaded, "queue overloaded"),
	)
	handler := mustNewHandler(t, evaluator, health.TargetLive)

	request := httptest.NewRequest(http.MethodGet, DefaultLivePath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultPassedStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultPassedStatus)
	}
	if got := recorder.Body.String(); got != textOK {
		t.Fatalf("body = %q, want %q", got, textOK)
	}
}

func TestHandlerServeHTTPWithCustomPolicy(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Degraded("queue", health.ReasonOverloaded, "queue overloaded"),
	)
	handler := mustNewHandler(
		t,
		evaluator,
		health.TargetReady,
		WithPolicy(health.ReadyPolicy().WithDegraded(true)),
	)

	request := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultPassedStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultPassedStatus)
	}
}

func TestHandlerServeHTTPJSON(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Unhealthy(
			"database",
			health.ReasonDependencyUnavailable,
			"database dependency is unavailable",
		).WithCause(errors.New("private cause")),
	)
	handler := mustNewHandler(
		t,
		evaluator,
		health.TargetReady,
		WithFormat(FormatJSON),
		WithDetailLevel(DetailAll),
	)

	request := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultFailedStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultFailedStatus)
	}
	if got := response.Header.Get(headerContentType); got != contentTypeJSON {
		t.Fatalf("Content-Type = %q, want %q", got, contentTypeJSON)
	}
	if !strings.Contains(recorder.Body.String(), `"target":"ready"`) {
		t.Fatalf("JSON body does not contain target: %q", recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "private cause") {
		t.Fatalf("JSON body exposes private cause: %q", recorder.Body.String())
	}
}

func TestHandlerServeHTTPHeadSuppressesBody(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Healthy("storage"),
	)
	handler := mustNewHandler(
		t,
		evaluator,
		health.TargetReady,
		WithFormat(FormatJSON),
		WithDetailLevel(DetailAll),
	)

	request := httptest.NewRequest(http.MethodHead, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultPassedStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultPassedStatus)
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("HEAD body length = %d, want 0; body=%q", recorder.Body.Len(), recorder.Body.String())
	}
}

func TestHandlerServeHTTPRejectsUnsupportedMethod(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t,
		health.TargetReady,
		health.Healthy("storage"),
	)
	handler := mustNewHandler(t, evaluator, health.TargetReady)

	request := httptest.NewRequest(http.MethodPost, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusMethodNotAllowed)
	}
	if got := response.Header.Get(headerAllow); got != allowedMethodsHeader {
		t.Fatalf("Allow header = %q, want %q", got, allowedMethodsHeader)
	}
}

func TestHandlerServeHTTPNoChecksFails(t *testing.T) {
	t.Parallel()

	registry := health.NewRegistry()
	evaluator, err := health.NewEvaluator(registry)
	if err != nil {
		t.Fatalf("health.NewEvaluator() = %v, want nil", err)
	}

	handler := mustNewHandler(t, evaluator, health.TargetReady)

	request := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultFailedStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultFailedStatus)
	}
}

func TestHandlerServeHTTPPassesRequestContext(t *testing.T) {
	t.Parallel()

	type contextKey struct{}

	key := contextKey{}
	value := "request-value"

	checker, err := health.NewCheck("context_check", func(ctx context.Context) health.Result {
		if ctx.Value(key) != value {
			return health.Unhealthy(
				"context_check",
				health.ReasonMisconfigured,
				"request context value missing",
			)
		}

		return health.Healthy("context_check")
	})
	if err != nil {
		t.Fatalf("health.NewCheck() = %v, want nil", err)
	}

	registry := health.NewRegistry()
	if err := registry.Register(health.TargetReady, checker); err != nil {
		t.Fatalf("registry.Register() = %v, want nil", err)
	}

	evaluator, err := health.NewEvaluator(registry, health.WithDefaultTimeout(0))
	if err != nil {
		t.Fatalf("health.NewEvaluator() = %v, want nil", err)
	}

	handler := mustNewHandler(t, evaluator, health.TargetReady)

	request := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	request = request.WithContext(context.WithValue(request.Context(), key, value))

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultPassedStatus {
		t.Fatalf("status = %d, want %d; body=%q", response.StatusCode, DefaultPassedStatus, recorder.Body.String())
	}
}

func TestHandlerServeHTTPRejectsNilRequest(t *testing.T) {
	t.Parallel()

	handler := mustNewHandler(
		t,
		healthtest.NewEvaluatorWithResults(t, health.TargetReady, health.Healthy("storage")),
		health.TargetReady,
	)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, nil)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultErrorStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultErrorStatus)
	}
	if got := recorder.Body.String(); got != textHandlerError {
		t.Fatalf("body = %q, want %q", got, textHandlerError)
	}
}

func TestNilHandlerServeHTTP(t *testing.T) {
	t.Parallel()

	var handler *Handler

	request := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != DefaultErrorStatus {
		t.Fatalf("status = %d, want %d", response.StatusCode, DefaultErrorStatus)
	}
	if got := recorder.Body.String(); got != textHandlerError {
		t.Fatalf("body = %q, want %q", got, textHandlerError)
	}
}

func TestNilHandlerAccessors(t *testing.T) {
	t.Parallel()

	var handler *Handler

	if handler.Target() != health.TargetUnknown {
		t.Fatalf("Target() = %s, want unknown", handler.Target())
	}
	if handler.Format() != FormatText {
		t.Fatalf("Format() = %s, want text", handler.Format())
	}
	if handler.DetailLevel() != DetailNone {
		t.Fatalf("DetailLevel() = %s, want none", handler.DetailLevel())
	}
}
