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
// Response intentionally does not embed health.Report, so future fields added to
// the core health package cannot accidentally become part of the public HTTP
// response contract.
type Response struct {
	Target         string          `json:"target"`
	Status         string          `json:"status"`
	Passed         bool            `json:"passed"`
	Observed       string          `json:"observed,omitempty"`
	DurationMillis int64           `json:"duration_ms,omitempty"`
	Checks         []CheckResponse `json:"checks,omitempty"`
}

// CheckResponse is the safe HTTP representation of one health check result.
type CheckResponse struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	Passed         bool   `json:"passed"`
	Reason         string `json:"reason,omitempty"`
	Message        string `json:"message,omitempty"`
	Observed       string `json:"observed,omitempty"`
	DurationMillis int64  `json:"duration_ms,omitempty"`
}

// newResponse converts report into a safe HTTP response DTO.
func newResponse(report health.Report, passed bool, policy health.TargetPolicy, detail DetailLevel) Response {
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
func selectChecks(report health.Report, policy health.TargetPolicy, detail DetailLevel) []health.Result {
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
func formatTimestamp(ts time.Time) string {
	if ts.IsZero() {
		return ""
	}

	return ts.UTC().Format(time.RFC3339Nano)
}

// durationMillis converts duration into whole milliseconds for public responses.
func durationMillis(duration time.Duration) int64 {
	if duration <= 0 {
		return 0
	}

	return duration.Milliseconds()
}

// formatReason converts reason into the public reason representation.
func formatReason(reason health.Reason) string {
	if reason == health.ReasonNone {
		return ""
	}
	if !reason.IsValid() {
		return reason.String()
	}

	return string(reason)
}
