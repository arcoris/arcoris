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
	"testing"

	"arcoris.dev/health"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		target      health.Target
		status      health.Status
		shouldPass  bool
		format      Format
		detail      DetailLevel
		statusCodes HTTPStatusCodes
	}{
		{
			name:        "startup",
			target:      health.TargetStartup,
			status:      health.StatusStarting,
			shouldPass:  false,
			format:      FormatText,
			detail:      DetailNone,
			statusCodes: DefaultStatusCodes(),
		},
		{
			name:        "live",
			target:      health.TargetLive,
			status:      health.StatusDegraded,
			shouldPass:  true,
			format:      FormatText,
			detail:      DetailNone,
			statusCodes: DefaultStatusCodes(),
		},
		{
			name:        "ready",
			target:      health.TargetReady,
			status:      health.StatusDegraded,
			shouldPass:  false,
			format:      FormatText,
			detail:      DetailNone,
			statusCodes: DefaultStatusCodes(),
		},
		{
			name:        "unknown",
			target:      health.TargetUnknown,
			status:      health.StatusHealthy,
			shouldPass:  true,
			format:      FormatText,
			detail:      DetailNone,
			statusCodes: DefaultStatusCodes(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := defaultConfig(tc.target)

			if got := cfg.policy.Passes(tc.status); got != tc.shouldPass {
				t.Fatalf("policy.Passes(%s) = %v, want %v", tc.status, got, tc.shouldPass)
			}
			if cfg.format != tc.format {
				t.Fatalf("format = %s, want %s", cfg.format, tc.format)
			}
			if cfg.detailLevel != tc.detail {
				t.Fatalf("detailLevel = %s, want %s", cfg.detailLevel, tc.detail)
			}
			if cfg.statusCodes != tc.statusCodes {
				t.Fatalf("statusCodes = %+v, want %+v", cfg.statusCodes, tc.statusCodes)
			}
		})
	}
}

func TestDefaultConfigUsesHealthDefaultPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		target     health.Target
		status     health.Status
		shouldPass bool
	}{
		{name: "startup starting fails", target: health.TargetStartup, status: health.StatusStarting, shouldPass: false},
		{name: "startup healthy passes", target: health.TargetStartup, status: health.StatusHealthy, shouldPass: true},
		{name: "live starting passes", target: health.TargetLive, status: health.StatusStarting, shouldPass: true},
		{name: "live degraded passes", target: health.TargetLive, status: health.StatusDegraded, shouldPass: true},
		{name: "ready degraded fails", target: health.TargetReady, status: health.StatusDegraded, shouldPass: false},
		{name: "ready healthy passes", target: health.TargetReady, status: health.StatusHealthy, shouldPass: true},
		{name: "unknown zero policy healthy passes", target: health.TargetUnknown, status: health.StatusHealthy, shouldPass: true},
		{name: "unknown zero policy degraded fails", target: health.TargetUnknown, status: health.StatusDegraded, shouldPass: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			policy := defaultConfig(tc.target).policy
			if got := policy.Passes(tc.status); got != tc.shouldPass {
				t.Fatalf("defaultConfig(%s).policy.Passes(%s) = %v, want %v", tc.target, tc.status, got, tc.shouldPass)
			}
		})
	}
}

func TestApplyOptions(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(health.TargetReady)

	err := applyOptions(
		&cfg,
		WithFormat(FormatJSON),
		WithDetailLevel(DetailFailed),
		WithFailedStatus(http.StatusTooManyRequests),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}

	if cfg.format != FormatJSON {
		t.Fatalf("format = %s, want json", cfg.format)
	}
	if cfg.detailLevel != DetailFailed {
		t.Fatalf("detailLevel = %s, want failed", cfg.detailLevel)
	}
	if cfg.statusCodes.Failed != http.StatusTooManyRequests {
		t.Fatalf("failed status = %d, want %d", cfg.statusCodes.Failed, http.StatusTooManyRequests)
	}
}

func TestApplyOptionsRejectsNilOption(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(health.TargetReady)

	err := applyOptions(&cfg, nil)
	if !errors.Is(err, ErrNilOption) {
		t.Fatalf("applyOptions(nil) = %v, want ErrNilOption", err)
	}
}

func TestApplyOptionsAppliesInOrder(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(health.TargetReady)

	err := applyOptions(
		&cfg,
		WithFormat(FormatText),
		WithFormat(FormatJSON),
		WithDetailLevel(DetailNone),
		WithDetailLevel(DetailAll),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}

	if cfg.format != FormatJSON {
		t.Fatalf("format = %s, want json", cfg.format)
	}
	if cfg.detailLevel != DetailAll {
		t.Fatalf("detailLevel = %s, want all", cfg.detailLevel)
	}
}

func TestWithPolicy(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(health.TargetReady)
	policy := health.ReadyPolicy().WithDegraded(true)

	err := WithPolicy(policy)(&cfg)
	if err != nil {
		t.Fatalf("WithPolicy() = %v, want nil", err)
	}

	if !cfg.policy.Passes(health.StatusDegraded) {
		t.Fatal("custom policy should pass degraded")
	}
}

func TestWithFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format Format
		want   bool
	}{
		{name: "text", format: FormatText, want: true},
		{name: "json", format: FormatJSON, want: true},
		{name: "invalid", format: Format(99), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := defaultConfig(health.TargetReady)
			err := WithFormat(tc.format)(&cfg)

			if got := err == nil; got != tc.want {
				t.Fatalf("WithFormat(%s) ok = %v, want %v; err=%v", tc.format, got, tc.want, err)
			}
			if tc.want && cfg.format != tc.format {
				t.Fatalf("format = %s, want %s", cfg.format, tc.format)
			}
			if !tc.want && !errors.Is(err, ErrInvalidFormat) {
				t.Fatalf("WithFormat(%s) = %v, want ErrInvalidFormat", tc.format, err)
			}
		})
	}
}

func TestWithDetailLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level DetailLevel
		want  bool
	}{
		{name: "none", level: DetailNone, want: true},
		{name: "failed", level: DetailFailed, want: true},
		{name: "all", level: DetailAll, want: true},
		{name: "invalid", level: DetailLevel(99), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := defaultConfig(health.TargetReady)
			err := WithDetailLevel(tc.level)(&cfg)

			if got := err == nil; got != tc.want {
				t.Fatalf("WithDetailLevel(%s) ok = %v, want %v; err=%v", tc.level, got, tc.want, err)
			}
			if tc.want && cfg.detailLevel != tc.level {
				t.Fatalf("detailLevel = %s, want %s", cfg.detailLevel, tc.level)
			}
			if !tc.want && !errors.Is(err, ErrInvalidDetailLevel) {
				t.Fatalf("WithDetailLevel(%s) = %v, want ErrInvalidDetailLevel", tc.level, err)
			}
		})
	}
}

func TestWithStatusCodes(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(health.TargetReady)

	err := WithStatusCodes(HTTPStatusCodes{
		Failed: http.StatusTooManyRequests,
	})(&cfg)
	if err != nil {
		t.Fatalf("WithStatusCodes() = %v, want nil", err)
	}

	want := HTTPStatusCodes{
		Passed: DefaultPassedStatus,
		Failed: http.StatusTooManyRequests,
		Error:  DefaultErrorStatus,
	}
	if cfg.statusCodes != want {
		t.Fatalf("statusCodes = %+v, want %+v", cfg.statusCodes, want)
	}
}

func TestWithStatusCodesRejectsInvalidMapping(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(health.TargetReady)

	err := WithStatusCodes(HTTPStatusCodes{
		Failed: http.StatusOK,
	})(&cfg)
	if !errors.Is(err, ErrInvalidHTTPStatusCode) {
		t.Fatalf("WithStatusCodes(invalid) = %v, want ErrInvalidHTTPStatusCode", err)
	}
}

func TestWithIndividualStatuses(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig(health.TargetReady)

	err := applyOptions(
		&cfg,
		WithPassedStatus(http.StatusNoContent),
		WithFailedStatus(http.StatusTooManyRequests),
		WithErrorStatus(http.StatusBadGateway),
	)
	if err != nil {
		t.Fatalf("applyOptions(status options) = %v, want nil", err)
	}

	want := HTTPStatusCodes{
		Passed: http.StatusNoContent,
		Failed: http.StatusTooManyRequests,
		Error:  http.StatusBadGateway,
	}
	if cfg.statusCodes != want {
		t.Fatalf("statusCodes = %+v, want %+v", cfg.statusCodes, want)
	}
}

func TestWithIndividualStatusesRejectInvalidCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opt  Option
	}{
		{name: "passed", opt: WithPassedStatus(http.StatusServiceUnavailable)},
		{name: "failed", opt: WithFailedStatus(http.StatusOK)},
		{name: "error", opt: WithErrorStatus(http.StatusBadRequest)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := defaultConfig(health.TargetReady)
			err := tc.opt(&cfg)

			if !errors.Is(err, ErrInvalidHTTPStatusCode) {
				t.Fatalf("%s option = %v, want ErrInvalidHTTPStatusCode", tc.name, err)
			}
		})
	}
}
