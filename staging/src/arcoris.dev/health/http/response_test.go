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
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"arcoris.dev/health"
	"arcoris.dev/health/healthtest"
)

func TestNewResponseWithDetailNone(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	policy := health.ReadyPolicy()

	response := newResponse(report, report.Passed(policy), policy, DetailNone)

	if response.Target != "ready" {
		t.Fatalf("Target = %q, want ready", response.Target)
	}
	if response.Status != "unhealthy" {
		t.Fatalf("Status = %q, want unhealthy", response.Status)
	}
	if response.Passed {
		t.Fatal("Passed = true, want false")
	}
	if response.Observed == "" {
		t.Fatal("Observed is empty, want timestamp")
	}
	if response.DurationMillis != 25 {
		t.Fatalf("DurationMillis = %d, want 25", response.DurationMillis)
	}
	if response.Checks != nil {
		t.Fatalf("Checks = %+v, want nil", response.Checks)
	}
}

func TestNewResponseWithDetailFailed(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	policy := health.ReadyPolicy()

	response := newResponse(report, report.Passed(policy), policy, DetailFailed)

	if len(response.Checks) != 3 {
		t.Fatalf("Checks length = %d, want 3", len(response.Checks))
	}

	wantNames := []string{"queue", "cache", "database"}
	for i, want := range wantNames {
		if response.Checks[i].Name != want {
			t.Fatalf("Checks[%d].Name = %q, want %q", i, response.Checks[i].Name, want)
		}
		if response.Checks[i].Passed {
			t.Fatalf("Checks[%d].Passed = true, want false", i)
		}
	}
}

func TestNewResponseWithDetailAll(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	policy := health.ReadyPolicy()

	response := newResponse(report, report.Passed(policy), policy, DetailAll)

	if len(response.Checks) != len(report.Checks) {
		t.Fatalf("Checks length = %d, want %d", len(response.Checks), len(report.Checks))
	}

	if response.Checks[0].Name != "storage" {
		t.Fatalf("Checks[0].Name = %q, want storage", response.Checks[0].Name)
	}
	if !response.Checks[0].Passed {
		t.Fatal("healthy storage check should pass ready policy")
	}

	if response.Checks[1].Name != "queue" {
		t.Fatalf("Checks[1].Name = %q, want queue", response.Checks[1].Name)
	}
	if response.Checks[1].Passed {
		t.Fatal("degraded queue check should fail ready policy")
	}
}

func TestNewResponseUsesPolicyForCheckPassed(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)

	readyResponse := newResponse(report, false, health.ReadyPolicy(), DetailFailed)
	if len(readyResponse.Checks) != 3 {
		t.Fatalf("ready failed checks = %d, want 3", len(readyResponse.Checks))
	}

	liveResponse := newResponse(report, false, health.LivePolicy(), DetailFailed)
	if len(liveResponse.Checks) != 2 {
		t.Fatalf("live failed checks = %d, want 2", len(liveResponse.Checks))
	}
	for _, check := range liveResponse.Checks {
		if check.Name == "queue" {
			t.Fatal("degraded queue should not fail live policy")
		}
	}
}

func TestNewCheckResponseOmitsReasonNone(t *testing.T) {
	t.Parallel()

	result := health.Healthy("storage")
	response := newCheckResponse(result, health.ReadyPolicy())

	if response.Reason != "" {
		t.Fatalf("Reason = %q, want empty", response.Reason)
	}
}

func TestNewCheckResponseIncludesSafeFields(t *testing.T) {
	t.Parallel()

	result := health.Degraded(
		"queue",
		health.ReasonOverloaded,
		"queue is above soft capacity",
	).WithObserved(healthtest.ObservedTime).
		WithDuration(1500 * time.Millisecond).
		WithCause(errors.New("private cause"))

	response := newCheckResponse(result, health.ReadyPolicy())

	if response.Name != "queue" {
		t.Fatalf("Name = %q, want queue", response.Name)
	}
	if response.Status != "degraded" {
		t.Fatalf("Status = %q, want degraded", response.Status)
	}
	if response.Passed {
		t.Fatal("Passed = true, want false")
	}
	if response.Reason != "overloaded" {
		t.Fatalf("Reason = %q, want overloaded", response.Reason)
	}
	if response.Message != "queue is above soft capacity" {
		t.Fatalf("Message = %q, want safe message", response.Message)
	}
	if response.Observed == "" {
		t.Fatal("Observed is empty, want timestamp")
	}
	if response.DurationMillis != 1500 {
		t.Fatalf("DurationMillis = %d, want 1500", response.DurationMillis)
	}
}

func TestResponseJSONDoesNotExposeCause(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	policy := health.ReadyPolicy()
	response := newResponse(report, false, policy, DetailAll)

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("json.Marshal() = %v, want nil", err)
	}

	encoded := string(data)
	if strings.Contains(encoded, "private cause") {
		t.Fatalf("encoded response exposes private cause: %s", encoded)
	}
	if strings.Contains(encoded, "Cause") || strings.Contains(encoded, "cause") {
		t.Fatalf("encoded response contains cause field: %s", encoded)
	}
}

func TestSelectChecks(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	policy := health.ReadyPolicy()

	tests := []struct {
		name   string
		detail DetailLevel
		want   int
	}{
		{name: "none", detail: DetailNone, want: 0},
		{name: "failed", detail: DetailFailed, want: 3},
		{name: "all", detail: DetailAll, want: len(report.Checks)},
		{name: "invalid", detail: DetailLevel(99), want: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			checks := selectChecks(report, policy, tc.detail)
			if len(checks) != tc.want {
				t.Fatalf("selectChecks(%s) length = %d, want %d", tc.detail, len(checks), tc.want)
			}
		})
	}
}

func TestSelectChecksReturnsDefensiveCopyForAll(t *testing.T) {
	t.Parallel()

	report := healthtest.MixedReport(health.TargetReady)
	checks := selectChecks(report, health.ReadyPolicy(), DetailAll)

	checks[0] = health.Unhealthy("mutated", health.ReasonFatal, "mutated")

	if report.Checks[0].Name != "storage" {
		t.Fatalf("report mutated through selected checks: %+v", report.Checks[0])
	}
}

func TestFormatTimestamp(t *testing.T) {
	t.Parallel()

	if got := formatTimestamp(time.Time{}); got != "" {
		t.Fatalf("zero timestamp = %q, want empty", got)
	}

	ts := time.Date(2026, 5, 4, 12, 30, 15, 123456789, time.FixedZone("UTC+3", 3*60*60))
	got := formatTimestamp(ts)
	want := "2026-05-04T09:30:15.123456789Z"

	if got != want {
		t.Fatalf("formatTimestamp() = %q, want %q", got, want)
	}
}

func TestDurationMillis(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration time.Duration
		want     int64
	}{
		{name: "zero", duration: 0, want: 0},
		{name: "negative", duration: -time.Second, want: 0},
		{name: "sub millisecond", duration: time.Microsecond, want: 0},
		{name: "one millisecond", duration: time.Millisecond, want: 1},
		{name: "second", duration: time.Second, want: 1000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := durationMillis(tc.duration); got != tc.want {
				t.Fatalf("durationMillis(%s) = %d, want %d", tc.duration, got, tc.want)
			}
		})
	}
}

func TestFormatReason(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		reason health.Reason
		want   string
	}{
		{name: "none", reason: health.ReasonNone, want: ""},
		{name: "builtin", reason: health.ReasonOverloaded, want: "overloaded"},
		{name: "custom", reason: health.Reason("custom_reason"), want: "custom_reason"},
		{name: "invalid", reason: health.Reason("bad-reason"), want: "invalid"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := formatReason(tc.reason); got != tc.want {
				t.Fatalf("formatReason(%s) = %q, want %q", tc.reason, got, tc.want)
			}
		})
	}
}
