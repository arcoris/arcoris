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

	"arcoris.dev/component-base/pkg/health"
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

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			config := defaultConfig(test.target)

			if got := config.policy.Passes(test.status); got != test.shouldPass {
				t.Fatalf("policy.Passes(%s) = %v, want %v", test.status, got, test.shouldPass)
			}
			if config.format != test.format {
				t.Fatalf("format = %s, want %s", config.format, test.format)
			}
			if config.detailLevel != test.detail {
				t.Fatalf("detailLevel = %s, want %s", config.detailLevel, test.detail)
			}
			if config.statusCodes != test.statusCodes {
				t.Fatalf("statusCodes = %+v, want %+v", config.statusCodes, test.statusCodes)
			}
		})
	}
}

func TestDefaultTargetPolicy(t *testing.T) {
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

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			policy := defaultTargetPolicy(test.target)
			if got := policy.Passes(test.status); got != test.shouldPass {
				t.Fatalf("defaultTargetPolicy(%s).Passes(%s) = %v, want %v", test.target, test.status, got, test.shouldPass)
			}
		})
	}
}

func TestApplyOptions(t *testing.T) {
	t.Parallel()

	config := defaultConfig(health.TargetReady)

	err := applyOptions(
		&config,
		WithFormat(FormatJSON),
		WithDetailLevel(DetailFailed),
		WithFailedStatus(http.StatusTooManyRequests),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}

	if config.format != FormatJSON {
		t.Fatalf("format = %s, want json", config.format)
	}
	if config.detailLevel != DetailFailed {
		t.Fatalf("detailLevel = %s, want failed", config.detailLevel)
	}
	if config.statusCodes.Failed != http.StatusTooManyRequests {
		t.Fatalf("failed status = %d, want %d", config.statusCodes.Failed, http.StatusTooManyRequests)
	}
}

func TestApplyOptionsRejectsNilOption(t *testing.T) {
	t.Parallel()

	config := defaultConfig(health.TargetReady)

	err := applyOptions(&config, nil)
	if !errors.Is(err, ErrNilOption) {
		t.Fatalf("applyOptions(nil) = %v, want ErrNilOption", err)
	}
}

func TestApplyOptionsAppliesInOrder(t *testing.T) {
	t.Parallel()

	config := defaultConfig(health.TargetReady)

	err := applyOptions(
		&config,
		WithFormat(FormatText),
		WithFormat(FormatJSON),
		WithDetailLevel(DetailNone),
		WithDetailLevel(DetailAll),
	)
	if err != nil {
		t.Fatalf("applyOptions() = %v, want nil", err)
	}

	if config.format != FormatJSON {
		t.Fatalf("format = %s, want json", config.format)
	}
	if config.detailLevel != DetailAll {
		t.Fatalf("detailLevel = %s, want all", config.detailLevel)
	}
}

func TestWithPolicy(t *testing.T) {
	t.Parallel()

	config := defaultConfig(health.TargetReady)
	policy := health.ReadyPolicy().WithDegraded(true)

	err := WithPolicy(policy)(&config)
	if err != nil {
		t.Fatalf("WithPolicy() = %v, want nil", err)
	}

	if !config.policy.Passes(health.StatusDegraded) {
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

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			config := defaultConfig(health.TargetReady)
			err := WithFormat(test.format)(&config)

			if got := err == nil; got != test.want {
				t.Fatalf("WithFormat(%s) ok = %v, want %v; err=%v", test.format, got, test.want, err)
			}
			if test.want && config.format != test.format {
				t.Fatalf("format = %s, want %s", config.format, test.format)
			}
			if !test.want && !errors.Is(err, ErrInvalidFormat) {
				t.Fatalf("WithFormat(%s) = %v, want ErrInvalidFormat", test.format, err)
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

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			config := defaultConfig(health.TargetReady)
			err := WithDetailLevel(test.level)(&config)

			if got := err == nil; got != test.want {
				t.Fatalf("WithDetailLevel(%s) ok = %v, want %v; err=%v", test.level, got, test.want, err)
			}
			if test.want && config.detailLevel != test.level {
				t.Fatalf("detailLevel = %s, want %s", config.detailLevel, test.level)
			}
			if !test.want && !errors.Is(err, ErrInvalidDetailLevel) {
				t.Fatalf("WithDetailLevel(%s) = %v, want ErrInvalidDetailLevel", test.level, err)
			}
		})
	}
}

func TestWithStatusCodes(t *testing.T) {
	t.Parallel()

	config := defaultConfig(health.TargetReady)

	err := WithStatusCodes(HTTPStatusCodes{
		Failed: http.StatusTooManyRequests,
	})(&config)
	if err != nil {
		t.Fatalf("WithStatusCodes() = %v, want nil", err)
	}

	want := HTTPStatusCodes{
		Passed: DefaultPassedStatus,
		Failed: http.StatusTooManyRequests,
		Error:  DefaultErrorStatus,
	}
	if config.statusCodes != want {
		t.Fatalf("statusCodes = %+v, want %+v", config.statusCodes, want)
	}
}

func TestWithStatusCodesRejectsInvalidMapping(t *testing.T) {
	t.Parallel()

	config := defaultConfig(health.TargetReady)

	err := WithStatusCodes(HTTPStatusCodes{
		Failed: http.StatusOK,
	})(&config)
	if !errors.Is(err, ErrInvalidHTTPStatusCode) {
		t.Fatalf("WithStatusCodes(invalid) = %v, want ErrInvalidHTTPStatusCode", err)
	}
}

func TestWithIndividualStatuses(t *testing.T) {
	t.Parallel()

	config := defaultConfig(health.TargetReady)

	err := applyOptions(
		&config,
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
	if config.statusCodes != want {
		t.Fatalf("statusCodes = %+v, want %+v", config.statusCodes, want)
	}
}

func TestWithIndividualStatusesRejectInvalidCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		option Option
	}{
		{name: "passed", option: WithPassedStatus(http.StatusServiceUnavailable)},
		{name: "failed", option: WithFailedStatus(http.StatusOK)},
		{name: "error", option: WithErrorStatus(http.StatusBadRequest)},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			config := defaultConfig(health.TargetReady)
			err := test.option(&config)

			if !errors.Is(err, ErrInvalidHTTPStatusCode) {
				t.Fatalf("%s option = %v, want ErrInvalidHTTPStatusCode", test.name, err)
			}
		})
	}
}
