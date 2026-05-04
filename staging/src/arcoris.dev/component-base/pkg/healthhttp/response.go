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
	"time"

	"arcoris.dev/component-base/pkg/health"
)

// Response is the safe HTTP representation of a health report.
//
// Response is an adapter DTO. It intentionally does not embed health.Report so
// future fields added to the core health package cannot accidentally become part
// of the public HTTP response contract.
//
// Response never contains Result.Cause, panic stacks, raw errors, context causes,
// connection strings, credentials, internal addresses, tenant identifiers, or
// other private diagnostic data. Renderers should encode this type, not
// health.Report directly.
type Response struct {
	// Target is the evaluated health target.
	//
	// The value is a stable diagnostic string such as "startup", "live", or
	// "ready".
	Target string `json:"target"`

	// Status is the aggregate health status of the evaluated target.
	//
	// The value is a stable diagnostic string such as "healthy", "degraded", or
	// "unhealthy".
	Status string `json:"status"`

	// Passed reports whether the aggregate report passed the handler policy.
	//
	// Passed is derived from health.TargetPolicy, not directly from Status. This
	// matters because a degraded status may pass liveness but fail readiness.
	Passed bool `json:"passed"`

	// Observed is the report observation timestamp encoded as UTC RFC3339Nano.
	//
	// The field is omitted when the report does not have an observation time.
	Observed string `json:"observed,omitempty"`

	// DurationMillis is the report evaluation duration in whole milliseconds.
	//
	// The field is omitted when the duration is zero or negative. Negative
	// durations should not be produced by Evaluator, but the DTO conversion stays
	// defensive because health.Report is a plain value type.
	DurationMillis int64 `json:"duration_ms,omitempty"`

	// Checks contains safe check-level details selected by DetailLevel.
	//
	// The slice is empty when DetailNone is used. It contains only failed checks
	// for DetailFailed and all checks for DetailAll.
	Checks []CheckResponse `json:"checks,omitempty"`
}

// CheckResponse is the safe HTTP representation of one health check result.
//
// CheckResponse intentionally exposes only stable, low-cardinality and
// caller-safe fields. It never contains Result.Cause or raw low-level diagnostic
// data.
type CheckResponse struct {
	// Name is the stable logical check name.
	Name string `json:"name"`

	// Status is the health status returned by the check.
	Status string `json:"status"`

	// Passed reports whether this check status passes the handler policy.
	//
	// Passed is policy-dependent. A degraded check may pass under a liveness
	// policy and fail under a readiness policy.
	Passed bool `json:"passed"`

	// Reason is the stable machine-readable reason for the check status.
	//
	// The field is omitted when the result has no reason.
	Reason string `json:"reason,omitempty"`

	// Message is the safe human-readable check message.
	//
	// The caller that created the original health.Result owns message safety.
	// Response conversion does not include Result.Cause and does not derive a
	// message from raw errors.
	Message string `json:"message,omitempty"`

	// Observed is the check observation timestamp encoded as UTC RFC3339Nano.
	//
	// The field is omitted when the check does not have an observation time.
	Observed string `json:"observed,omitempty"`

	// DurationMillis is the check duration in whole milliseconds.
	//
	// The field is omitted when the duration is zero or negative.
	DurationMillis int64 `json:"duration_ms,omitempty"`
}

// newResponse converts report into a safe HTTP response DTO.
//
// detail controls check-level exposure. DetailNone omits checks, DetailFailed
// includes only policy-failed checks, and DetailAll includes all checks.
//
// newResponse assumes policy is the same policy used to compute passed. It still
// receives passed explicitly because handler code already computes the aggregate
// decision when selecting the HTTP status code.
func newResponse(
	report health.Report,
	passed bool,
	policy health.TargetPolicy,
	detail DetailLevel,
) Response {
	response := Response{
		Target:         report.Target.String(),
		Status:         report.Status.String(),
		Passed:         passed,
		Observed:       formatTimestamp(report.Observed),
		DurationMillis: durationMillis(report.Duration),
	}

	checks := selectChecks(report, policy, detail)
	if len(checks) > 0 {
		response.Checks = newCheckResponses(checks, policy)
	}

	return response
}

// newCheckResponses converts check results into safe HTTP check DTOs.
//
// The returned slice preserves input order. It does not retain or expose
// Result.Cause.
func newCheckResponses(checks []health.Result, policy health.TargetPolicy) []CheckResponse {
	if len(checks) == 0 {
		return nil
	}

	responses := make([]CheckResponse, 0, len(checks))
	for _, check := range checks {
		responses = append(responses, newCheckResponse(check, policy))
	}

	return responses
}

// newCheckResponse converts one health result into a safe HTTP check DTO.
func newCheckResponse(result health.Result, policy health.TargetPolicy) CheckResponse {
	return CheckResponse{
		Name:           result.Name,
		Status:         result.Status.String(),
		Passed:         policy.Passes(result.Status),
		Reason:         formatReason(result.Reason),
		Message:        result.Message,
		Observed:       formatTimestamp(result.Observed),
		DurationMillis: durationMillis(result.Duration),
	}
}

// selectChecks returns the health results allowed by detail.
//
// The returned slice is caller-owned. DetailAll uses Report.ChecksCopy to avoid
// exposing the report's backing slice. DetailFailed uses Report.FailedChecks,
// which already returns a newly allocated slice.
func selectChecks(
	report health.Report,
	policy health.TargetPolicy,
	detail DetailLevel,
) []health.Result {
	switch detail {
	case DetailNone:
		return nil
	case DetailFailed:
		return report.FailedChecks(policy)
	case DetailAll:
		return report.ChecksCopy()
	default:
		return nil
	}
}

// formatTimestamp converts ts into the public timestamp representation.
//
// Zero timestamps are omitted by returning an empty string. Non-zero timestamps
// are normalized to UTC and encoded with RFC3339Nano so JSON responses are stable
// across hosts and local time zones.
func formatTimestamp(ts time.Time) string {
	if ts.IsZero() {
		return ""
	}

	return ts.UTC().Format(time.RFC3339Nano)
}

// durationMillis converts duration into whole milliseconds for public HTTP
// responses.
//
// Non-positive durations return zero so JSON omitempty suppresses the field.
// Sub-millisecond positive durations also become zero. healthhttp exposes
// durations as coarse diagnostics, not as a profiling or tracing API.
func durationMillis(duration time.Duration) int64 {
	if duration <= 0 {
		return 0
	}

	return duration.Milliseconds()
}

// formatReason converts reason into the public reason representation.
//
// ReasonNone is omitted by returning an empty string. Valid non-empty reasons are
// exposed as their underlying stable reason code. Invalid reasons are rendered as
// "invalid" defensively, although Evaluator and Result validation should prevent
// invalid reasons from reaching normal reports.
func formatReason(reason health.Reason) string {
	if reason == health.ReasonNone {
		return ""
	}
	if !reason.IsValid() {
		return reason.String()
	}

	return string(reason)
}
