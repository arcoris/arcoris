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

import "arcoris.dev/health"

// newResponse converts a core health report into the adapter's safe DTO.
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

// newCheckResponses converts selected results into safe HTTP DTOs.
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

// newCheckResponse converts one core health result into a safe DTO.
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
