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

package healthhttp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/healthtest"
)

func TestValidReportForTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		report health.Report
		want   bool
	}{
		{
			name:   "valid",
			report: healthtest.HealthyReport(health.TargetReady),
			want:   true,
		},
		{
			name:   "wrong target",
			report: healthtest.HealthyReport(health.TargetLive),
		},
		{
			name:   "invalid",
			report: httpInvalidReport(health.TargetReady),
		},
		{
			name:   "inconsistent",
			report: httpInconsistentReport(health.TargetReady),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := validReportForTarget(tc.report, health.TargetReady); got != tc.want {
				t.Fatalf("validReportForTarget() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHandlerServeHTTPRejectsMalformedEvaluatorReports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		report health.Report
	}{
		{
			name:   "wrong target",
			report: healthtest.HealthyReport(health.TargetLive),
		},
		{
			name:   "invalid",
			report: httpInvalidReport(health.TargetReady),
		},
		{
			name:   "inconsistent",
			report: httpInconsistentReport(health.TargetReady),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := mustNewHandler(
				t,
				healthtest.NewEvaluatorWithReport(tc.report),
				health.TargetReady,
				WithFormat(FormatJSON),
				WithDetailLevel(DetailAll),
				WithFailedStatus(http.StatusTeapot),
				WithErrorStatus(http.StatusBadGateway),
			)

			req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadGateway {
				t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadGateway)
			}

			body := recorder.Body.String()
			if !strings.Contains(body, "health handler error") {
				t.Fatalf("body = %q, want generic handler error", body)
			}
			if strings.Contains(body, "database") || strings.Contains(body, "wrong") {
				t.Fatalf("malformed report leaked through body: %q", body)
			}
		})
	}
}

func TestHandlerServeHTTPEvaluatorErrorIsGeneric(t *testing.T) {
	t.Parallel()

	raw := errors.New("database password=secret failed")
	handler := mustNewHandler(
		t,
		healthtest.NewErrorEvaluator(raw),
		health.TargetReady,
		WithErrorStatus(http.StatusBadGateway),
	)

	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadGateway)
	}
	if body := recorder.Body.String(); body != textHandlerError {
		t.Fatalf("body = %q, want %q", body, textHandlerError)
	}
	if strings.Contains(recorder.Body.String(), "password=secret") {
		t.Fatalf("body leaked raw evaluator error: %q", recorder.Body.String())
	}
}

func httpInvalidReport(target health.Target) health.Report {
	report := healthtest.HealthyReport(target)
	report.Status = health.Status(255)

	return report
}

func httpInconsistentReport(target health.Target) health.Report {
	return healthtest.Report(
		target,
		health.StatusHealthy,
		health.Unhealthy("database", health.ReasonFatal, "database failed"),
	)
}
