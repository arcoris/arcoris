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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
)

func TestInstall(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewEvaluatorWithResults(t, health.TargetReady, health.Healthy("storage"))

	err := Install(mux, "/internal/readyz", evaluator, health.TargetReady)
	if err != nil {
		t.Fatalf("Install() = %v, want nil", err)
	}

	handler, ok := mux.handler("/internal/readyz")
	if !ok {
		t.Fatal("registered handler not found")
	}

	req := httptest.NewRequest(http.MethodGet, "/internal/readyz", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != DefaultPassedStatus {
		t.Fatalf("status = %d, want %d", resp.StatusCode, DefaultPassedStatus)
	}
}

func TestInstallAppliesOptions(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewEvaluatorWithResults(
		t,
		health.TargetReady,
		health.Degraded("queue", health.ReasonOverloaded, "queue overloaded"),
	)

	err := Install(
		mux,
		DefaultReadyPath,
		evaluator,
		health.TargetReady,
		WithPolicy(health.ReadyPolicy().WithDegraded(true)),
		WithFormat(FormatJSON),
		WithDetailLevel(DetailAll),
	)
	if err != nil {
		t.Fatalf("Install() = %v, want nil", err)
	}

	handler, ok := mux.handler(DefaultReadyPath)
	if !ok {
		t.Fatal("registered handler not found")
	}

	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != DefaultPassedStatus {
		t.Fatalf("status = %d, want %d", resp.StatusCode, DefaultPassedStatus)
	}
	if got := resp.Header.Get(headerContentType); got != contentTypeJSON {
		t.Fatalf("Content-Type = %q, want %q", got, contentTypeJSON)
	}
}

func TestInstallRejectsNilMux(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t, health.TargetReady, health.Healthy("storage"))

	err := Install(nil, DefaultReadyPath, evaluator, health.TargetReady)
	if !errors.Is(err, ErrNilMux) {
		t.Fatalf("Install(nil mux) = %v, want ErrNilMux", err)
	}
}

func TestInstallRejectsTypedNilMux(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewEvaluatorWithResults(t, health.TargetReady, health.Healthy("storage"))

	var mux *recordingMux
	err := Install(mux, DefaultReadyPath, evaluator, health.TargetReady)

	if !errors.Is(err, ErrNilMux) {
		t.Fatalf("Install(typed nil mux) = %v, want ErrNilMux", err)
	}
}

func TestInstallRejectsInvalidPath(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewEvaluatorWithResults(t, health.TargetReady, health.Healthy("storage"))

	err := Install(mux, "readyz", evaluator, health.TargetReady)
	if !errors.Is(err, ErrInvalidPath) {
		t.Fatalf("Install(invalid path) = %v, want ErrInvalidPath", err)
	}
	if len(mux.handlers) != 0 {
		t.Fatalf("mux handlers = %d, want 0", len(mux.handlers))
	}
}

func TestInstallRejectsNilEvaluator(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()

	err := Install(mux, DefaultReadyPath, nil, health.TargetReady)
	if !errors.Is(err, ErrNilEvaluator) {
		t.Fatalf("Install(nil evaluator) = %v, want ErrNilEvaluator", err)
	}
	if len(mux.handlers) != 0 {
		t.Fatalf("mux handlers = %d, want 0", len(mux.handlers))
	}
}

func TestInstallRejectsInvalidTarget(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewEvaluatorWithResults(t, health.TargetReady, health.Healthy("storage"))

	err := Install(mux, DefaultReadyPath, evaluator, health.TargetUnknown)
	if !errors.Is(err, health.ErrInvalidTarget) {
		t.Fatalf("Install(invalid target) = %v, want health.ErrInvalidTarget", err)
	}
	if len(mux.handlers) != 0 {
		t.Fatalf("mux handlers = %d, want 0", len(mux.handlers))
	}
}

func TestInstallRejectsInvalidOption(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewEvaluatorWithResults(t, health.TargetReady, health.Healthy("storage"))

	err := Install(mux, DefaultReadyPath, evaluator, health.TargetReady, WithFormat(Format(99)))
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("Install(invalid option) = %v, want ErrInvalidFormat", err)
	}
	if len(mux.handlers) != 0 {
		t.Fatalf("mux handlers = %d, want 0", len(mux.handlers))
	}
}

func TestInstallDefaults(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewDefaultTargetsEvaluator(t)

	err := InstallDefaults(mux, evaluator)
	if err != nil {
		t.Fatalf("InstallDefaults() = %v, want nil", err)
	}

	wantPaths := []string{
		DefaultStartupPath,
		DefaultLivePath,
		DefaultReadyPath,
	}

	for _, path := range wantPaths {
		path := path
		t.Run(path, func(t *testing.T) {
			handler, ok := mux.handler(path)
			if !ok {
				t.Fatalf("handler for %s not registered", path)
			}

			req := httptest.NewRequest(http.MethodGet, path, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if resp.StatusCode != DefaultPassedStatus {
				t.Fatalf("status for %s = %d, want %d", path, resp.StatusCode, DefaultPassedStatus)
			}
		})
	}
}

func TestInstallDefaultsDoesNotInstallGeneralHealthPath(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewDefaultTargetsEvaluator(t)

	err := InstallDefaults(mux, evaluator)
	if err != nil {
		t.Fatalf("InstallDefaults() = %v, want nil", err)
	}

	if _, ok := mux.handler(DefaultHealthPath); ok {
		t.Fatalf("general health path %s was registered by InstallDefaults", DefaultHealthPath)
	}
}

func TestInstallDefaultsAppliesOptionsToAllHandlers(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewDefaultTargetsEvaluator(t)

	err := InstallDefaults(
		mux,
		evaluator,
		WithFormat(FormatJSON),
		WithDetailLevel(DetailAll),
	)
	if err != nil {
		t.Fatalf("InstallDefaults() = %v, want nil", err)
	}

	for _, path := range []string{
		DefaultStartupPath,
		DefaultLivePath,
		DefaultReadyPath,
	} {
		handler, ok := mux.handler(path)
		if !ok {
			t.Fatalf("handler for %s not registered", path)
		}

		req := httptest.NewRequest(http.MethodGet, path, nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		resp := recorder.Result()
		defer resp.Body.Close()

		if got := resp.Header.Get(headerContentType); got != contentTypeJSON {
			t.Fatalf("Content-Type for %s = %q, want %q", path, got, contentTypeJSON)
		}
	}
}

func TestInstallDefaultsRejectsNilMux(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewDefaultTargetsEvaluator(t)

	err := InstallDefaults(nil, evaluator)
	if !errors.Is(err, ErrNilMux) {
		t.Fatalf("InstallDefaults(nil mux) = %v, want ErrNilMux", err)
	}
}

func TestInstallDefaultsRejectsTypedNilMux(t *testing.T) {
	t.Parallel()

	evaluator := healthtest.NewDefaultTargetsEvaluator(t)

	var mux *recordingMux
	err := InstallDefaults(mux, evaluator)

	if !errors.Is(err, ErrNilMux) {
		t.Fatalf("InstallDefaults(typed nil mux) = %v, want ErrNilMux", err)
	}
}

func TestInstallDefaultsRejectsNilEvaluatorWithoutMutatingMux(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()

	err := InstallDefaults(mux, nil)
	if !errors.Is(err, ErrNilEvaluator) {
		t.Fatalf("InstallDefaults(nil evaluator) = %v, want ErrNilEvaluator", err)
	}
	if len(mux.handlers) != 0 {
		t.Fatalf("mux handlers = %d, want 0", len(mux.handlers))
	}
}

func TestInstallDefaultsRejectsInvalidOptionsWithoutMutatingMux(t *testing.T) {
	t.Parallel()

	mux := newRecordingMux()
	evaluator := healthtest.NewDefaultTargetsEvaluator(t)

	err := InstallDefaults(mux, evaluator, WithFormat(Format(99)))
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("InstallDefaults(invalid option) = %v, want ErrInvalidFormat", err)
	}
	if len(mux.handlers) != 0 {
		t.Fatalf("mux handlers = %d, want 0", len(mux.handlers))
	}
}

func TestDefaultHandlersUsePrimaryDefaultPathsOnly(t *testing.T) {
	t.Parallel()

	if len(defaultHandlers) != 3 {
		t.Fatalf("defaultHandlers length = %d, want 3", len(defaultHandlers))
	}

	seen := make(map[string]health.Target, len(defaultHandlers))
	for _, item := range defaultHandlers {
		if err := ValidatePath(item.path); err != nil {
			t.Fatalf("default path %q invalid: %v", item.path, err)
		}
		if !item.target.IsConcrete() {
			t.Fatalf("default target for path %q is not concrete: %s", item.path, item.target)
		}
		if _, exists := seen[item.path]; exists {
			t.Fatalf("duplicate default path: %s", item.path)
		}

		seen[item.path] = item.target
	}

	if seen[DefaultStartupPath] != health.TargetStartup {
		t.Fatalf("startup path target = %s, want startup", seen[DefaultStartupPath])
	}
	if seen[DefaultLivePath] != health.TargetLive {
		t.Fatalf("live path target = %s, want live", seen[DefaultLivePath])
	}
	if seen[DefaultReadyPath] != health.TargetReady {
		t.Fatalf("ready path target = %s, want ready", seen[DefaultReadyPath])
	}

	if _, ok := seen[DefaultHealthPath]; ok {
		t.Fatalf("DefaultHealthPath must not be installed by default")
	}
}
