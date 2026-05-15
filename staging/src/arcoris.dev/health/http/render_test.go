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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
)

func TestRenderReportText(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	cfg := defaultConfig(health.TargetReady)
	cfg.detailLevel = DetailAll

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)

	renderReport(recorder, req, cfg, report, false)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != DefaultFailedStatus {
		t.Fatalf("status = %d, want %d", resp.StatusCode, DefaultFailedStatus)
	}
	if got := resp.Header.Get(headerCacheControl); got != headerValueNoStore {
		t.Fatalf("Cache-Control = %q, want %q", got, headerValueNoStore)
	}
	if got := resp.Header.Get(headerXContentTypeOptions); got != headerValueNoSniff {
		t.Fatalf("X-Content-Type-Options = %q, want %q", got, headerValueNoSniff)
	}
	if strings.Contains(recorder.Body.String(), "private cause") {
		t.Fatalf("text body exposes private cause: %q", recorder.Body.String())
	}
}

func TestRenderReportJSON(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	cfg := defaultConfig(health.TargetReady)
	cfg.format = FormatJSON
	cfg.detailLevel = DetailAll

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)

	renderReport(recorder, req, cfg, report, false)

	resp := recorder.Result()
	defer resp.Body.Close()

	if got := resp.Header.Get(headerContentType); got != contentTypeJSON {
		t.Fatalf("Content-Type = %q, want %q", got, contentTypeJSON)
	}
	if strings.Contains(recorder.Body.String(), "private cause") {
		t.Fatalf("json body exposes private cause: %q", recorder.Body.String())
	}
}

func TestRenderHandlerErrorText(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)

	renderHandlerError(recorder, req, defaultConfig(health.TargetReady))

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != DefaultErrorStatus {
		t.Fatalf("status = %d, want %d", resp.StatusCode, DefaultErrorStatus)
	}
	if got := recorder.Body.String(); got != textHandlerError {
		t.Fatalf("body = %q, want %q", got, textHandlerError)
	}
}

func TestRenderHandlerErrorJSON(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(health.TargetReady)
	cfg.format = FormatJSON

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)

	renderHandlerError(recorder, req, cfg)

	resp := recorder.Result()
	defer resp.Body.Close()

	if got := resp.Header.Get(headerContentType); got != contentTypeJSON {
		t.Fatalf("Content-Type = %q, want %q", got, contentTypeJSON)
	}
	if strings.Contains(recorder.Body.String(), "private") {
		t.Fatalf("error body must stay generic: %q", recorder.Body.String())
	}
}

func TestRenderSuppressesHeadBody(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	cfg := defaultConfig(health.TargetReady)
	cfg.format = FormatJSON
	cfg.detailLevel = DetailAll

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodHead, DefaultReadyPath, nil)

	renderReport(recorder, req, cfg, report, false)

	if recorder.Body.Len() != 0 {
		t.Fatalf("HEAD body length = %d, want 0", recorder.Body.Len())
	}
}
