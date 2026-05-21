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

// Response is the safe HTTP representation of a health report.
//
// Response intentionally does not embed health.Report. The adapter owns its own
// response contract so future additions to the core health package cannot widen
// HTTP exposure accidentally.
type Response struct {
	Target         string          `json:"target"`
	Status         string          `json:"status"`
	Passed         bool            `json:"passed"`
	Observed       string          `json:"observed,omitempty"`
	DurationMillis int64           `json:"duration_ms,omitempty"`
	Checks         []CheckResponse `json:"checks,omitempty"`
}

// CheckResponse is the safe HTTP representation of one health check result.
//
// The DTO includes only fields that are explicitly safe for HTTP exposure. It
// never contains Result.Cause, panic stacks, raw errors, or context causes.
type CheckResponse struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	Passed         bool   `json:"passed"`
	Reason         string `json:"reason,omitempty"`
	Message        string `json:"message,omitempty"`
	Observed       string `json:"observed,omitempty"`
	DurationMillis int64  `json:"duration_ms,omitempty"`
}
