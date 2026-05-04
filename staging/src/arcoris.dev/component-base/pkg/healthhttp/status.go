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
	// DefaultPassedStatus is returned when a health report passes policy.
	DefaultPassedStatus = http.StatusOK

	// DefaultFailedStatus is returned when evaluation succeeds but report fails.
	DefaultFailedStatus = http.StatusServiceUnavailable

	// DefaultErrorStatus is returned for adapter/evaluator boundary errors.
	DefaultErrorStatus = http.StatusInternalServerError
)

// HTTPStatusCodes defines the HTTP status code mapping used by health handlers.
type HTTPStatusCodes struct {
	Passed int
	Failed int
	Error  int
}

// DefaultStatusCodes returns the default HTTP status code mapping.
func DefaultStatusCodes() HTTPStatusCodes {
	return HTTPStatusCodes{
		Passed: DefaultPassedStatus,
		Failed: DefaultFailedStatus,
		Error:  DefaultErrorStatus,
	}
}

// Normalize returns a copy of codes with zero fields replaced by defaults.
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
func (codes HTTPStatusCodes) Validate() error {
	codes = codes.Normalize()

	if !validHTTPStatusCode(codes.Passed) || !statusClass(codes.Passed, 2) {
		return InvalidHTTPStatusCodeError{Field: "passed", Code: codes.Passed}
	}
	if !validHTTPStatusCode(codes.Failed) || !(statusClass(codes.Failed, 4) || statusClass(codes.Failed, 5)) {
		return InvalidHTTPStatusCodeError{Field: "failed", Code: codes.Failed}
	}
	if !validHTTPStatusCode(codes.Error) || !statusClass(codes.Error, 5) {
		return InvalidHTTPStatusCodeError{Field: "error", Code: codes.Error}
	}

	return nil
}

// statusForReport returns the HTTP status code for a successfully evaluated
// health report.
func (codes HTTPStatusCodes) statusForReport(passed bool) int {
	if passed {
		return codes.Passed
	}

	return codes.Failed
}

// statusForError returns the HTTP status code for an adapter/evaluator error.
func (codes HTTPStatusCodes) statusForError() int {
	return codes.Error
}

// validHTTPStatusCode reports whether code is inside the standard HTTP status
// code range.
func validHTTPStatusCode(code int) bool {
	return code >= 100 && code <= 599
}

// statusClass reports whether code belongs to the given HTTP status class.
func statusClass(code int, class int) bool {
	return code >= class*100 && code < (class+1)*100
}
