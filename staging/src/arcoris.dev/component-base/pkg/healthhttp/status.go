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

import "net/http"

const (
	// DefaultPassedStatus is the HTTP status code returned when a health report
	// passes the handler target policy.
	//
	// Health probe clients, load balancers, and orchestrators primarily consume
	// the status code. A 200 response means the evaluated target is acceptable
	// according to the handler policy.
	DefaultPassedStatus = http.StatusOK

	// DefaultFailedStatus is the HTTP status code returned when health evaluation
	// succeeds but the report fails the handler target policy.
	//
	// 503 is used instead of 500 because the handler itself worked correctly. The
	// component is unavailable for the evaluated health contract, most commonly
	// readiness or liveness.
	DefaultFailedStatus = http.StatusServiceUnavailable

	// DefaultErrorStatus is the HTTP status code returned when the HTTP adapter
	// cannot produce a reliable health response because of a handler, evaluator,
	// or configuration boundary error.
	//
	// This status is reserved for adapter-level failures, not for ordinary failed
	// health reports. A target that evaluates successfully but fails policy should
	// use DefaultFailedStatus.
	DefaultErrorStatus = http.StatusInternalServerError
)

// HTTPStatusCodes defines the HTTP status code mapping used by health handlers.
//
// The mapping is intentionally small:
//
//   - Passed is used when the health report passes the configured target policy.
//   - Failed is used when health evaluation succeeds but the report fails policy.
//   - Error is used when the HTTP adapter cannot produce a reliable health
//     response because evaluation or handler configuration failed.
//
// HTTPStatusCodes does not include method errors. Method validation is owned by
// method.go and uses http.StatusMethodNotAllowed.
type HTTPStatusCodes struct {
	// Passed is returned when the evaluated target passes policy.
	Passed int

	// Failed is returned when the evaluated target fails policy.
	Failed int

	// Error is returned when the handler cannot produce a reliable health
	// response because of an adapter or evaluator boundary error.
	Error int
}

// DefaultStatusCodes returns the default HTTP status code mapping for health
// handlers.
//
// The returned value is independent and may be copied freely.
func DefaultStatusCodes() HTTPStatusCodes {
	return HTTPStatusCodes{
		Passed: DefaultPassedStatus,
		Failed: DefaultFailedStatus,
		Error:  DefaultErrorStatus,
	}
}

// Normalize returns a copy of codes with zero fields replaced by defaults.
//
// Normalize lets callers override only part of the mapping without having to
// repeat every default. It does not validate non-zero values. Call Validate after
// Normalize at configuration boundaries.
func (codes HTTPStatusCodes) Normalize() HTTPStatusCodes {
	defaults := DefaultStatusCodes()

	if codes.Passed == 0 {
		codes.Passed = defaults.Passed
	}
	if codes.Failed == 0 {
		codes.Failed = defaults.Failed
	}
	if codes.Error == 0 {
		codes.Error = defaults.Error
	}

	return codes
}

// Validate reports an error if codes contains unsupported HTTP status codes.
//
// Validation accepts only status codes in the inclusive range [100, 599] and
// requires class-appropriate mappings:
//
//   - Passed must be a 2xx status.
//   - Failed must be a 4xx or 5xx status.
//   - Error must be a 5xx status.
//
// These restrictions keep health endpoints predictable for infrastructure
// clients and prevent misconfiguration such as returning 200 for a failed target
// or 404 for an adapter error.
func (codes HTTPStatusCodes) Validate() error {
	codes = codes.Normalize()

	if !validHTTPStatusCode(codes.Passed) || !statusClass(codes.Passed, 2) {
		return InvalidHTTPStatusCodeError{
			Field: "passed",
			Code:  codes.Passed,
		}
	}
	if !validHTTPStatusCode(codes.Failed) || !(statusClass(codes.Failed, 4) || statusClass(codes.Failed, 5)) {
		return InvalidHTTPStatusCodeError{
			Field: "failed",
			Code:  codes.Failed,
		}
	}
	if !validHTTPStatusCode(codes.Error) || !statusClass(codes.Error, 5) {
		return InvalidHTTPStatusCodeError{
			Field: "error",
			Code:  codes.Error,
		}
	}

	return nil
}

// statusForReport returns the HTTP status code for a successfully evaluated
// health report.
//
// The function assumes codes has already been normalized and validated by the
// handler configuration boundary.
func (codes HTTPStatusCodes) statusForReport(passed bool) int {
	if passed {
		return codes.Passed
	}

	return codes.Failed
}

// statusForError returns the HTTP status code for an adapter/evaluator boundary
// error.
//
// The function assumes codes has already been normalized and validated by the
// handler configuration boundary.
func (codes HTTPStatusCodes) statusForError() int {
	return codes.Error
}

// validHTTPStatusCode reports whether code is inside the standard HTTP status
// code range.
//
// The helper intentionally validates numeric range only. Class-specific
// constraints are checked by HTTPStatusCodes.Validate.
func validHTTPStatusCode(code int) bool {
	return code >= 100 && code <= 599
}

// statusClass reports whether code belongs to the given HTTP status class.
//
// For example, statusClass(http.StatusOK, 2) is true and
// statusClass(http.StatusServiceUnavailable, 5) is true.
func statusClass(code int, class int) bool {
	return code >= class*100 && code < (class+1)*100
}
