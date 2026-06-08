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
	"net/http"
	"net/http/httptest"
	"testing"

	"arcoris.dev/health"
	"arcoris.dev/healthtest"
)

func BenchmarkHTTPHandlerServeTextHealthy(b *testing.B) {
	handler, err := NewHandler(
		healthtest.NewEvaluatorWithReport(healthtest.HealthyReport(health.TargetReady)),
		health.TargetReady,
	)
	if err != nil {
		b.Fatalf("NewHandler() = %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	b.ReportAllocs()
	for b.Loop() {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
	}
}

func BenchmarkHTTPHandlerServeJSONMixedDetailAll(b *testing.B) {
	handler, err := NewHandler(
		healthtest.NewEvaluatorWithReport(healthtest.MixedReport(health.TargetReady)),
		health.TargetReady,
		WithFormat(FormatJSON),
		WithDetailLevel(DetailAll),
	)
	if err != nil {
		b.Fatalf("NewHandler() = %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	b.ReportAllocs()
	for b.Loop() {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)
	}
}

func BenchmarkHTTPHandlerServeJSONMalformedReport(b *testing.B) {
	tests := []struct {
		name   string
		report health.Report
	}{
		{
			name:   "WrongTarget",
			report: healthtest.HealthyReport(health.TargetLive),
		},
		{
			name:   "Invalid",
			report: httpInvalidReport(health.TargetReady),
		},
		{
			name:   "Inconsistent",
			report: httpInconsistentReport(health.TargetReady),
		},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			handler, err := NewHandler(
				healthtest.NewEvaluatorWithReport(tc.report),
				health.TargetReady,
				WithFormat(FormatJSON),
			)
			if err != nil {
				b.Fatalf("NewHandler() = %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
			b.ReportAllocs()
			for b.Loop() {
				recorder := httptest.NewRecorder()
				handler.ServeHTTP(recorder, req)
			}
		})
	}
}

func BenchmarkHTTPRenderJSONMixedDetailAll(b *testing.B) {
	report := healthtest.MixedReport(health.TargetReady)
	cfg := defaultConfig(health.TargetReady)
	cfg.format = FormatJSON
	cfg.detailLevel = DetailAll
	req := httptest.NewRequest(http.MethodGet, DefaultReadyPath, nil)
	b.ReportAllocs()
	for b.Loop() {
		recorder := httptest.NewRecorder()
		renderReport(recorder, req, cfg, report, false)
	}
}
