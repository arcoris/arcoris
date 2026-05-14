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
	"net/http"
	"testing"
)

func TestDefaultStatusCodeConstants(t *testing.T) {
	t.Parallel()

	if DefaultPassedStatus != http.StatusOK {
		t.Fatalf("DefaultPassedStatus = %d, want %d", DefaultPassedStatus, http.StatusOK)
	}
	if DefaultFailedStatus != http.StatusServiceUnavailable {
		t.Fatalf("DefaultFailedStatus = %d, want %d", DefaultFailedStatus, http.StatusServiceUnavailable)
	}
	if DefaultErrorStatus != http.StatusInternalServerError {
		t.Fatalf("DefaultErrorStatus = %d, want %d", DefaultErrorStatus, http.StatusInternalServerError)
	}
}

func TestDefaultStatusCodes(t *testing.T) {
	t.Parallel()

	codes := DefaultStatusCodes()

	if codes.Passed != DefaultPassedStatus {
		t.Fatalf("Passed = %d, want %d", codes.Passed, DefaultPassedStatus)
	}
	if codes.Failed != DefaultFailedStatus {
		t.Fatalf("Failed = %d, want %d", codes.Failed, DefaultFailedStatus)
	}
	if codes.Error != DefaultErrorStatus {
		t.Fatalf("Error = %d, want %d", codes.Error, DefaultErrorStatus)
	}
}

func TestHTTPStatusCodesNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		codes HTTPStatusCodes
		want  HTTPStatusCodes
	}{
		{
			name:  "zero",
			codes: HTTPStatusCodes{},
			want:  DefaultStatusCodes(),
		},
		{
			name: "partial passed",
			codes: HTTPStatusCodes{
				Passed: http.StatusAccepted,
			},
			want: HTTPStatusCodes{
				Passed: http.StatusAccepted,
				Failed: DefaultFailedStatus,
				Error:  DefaultErrorStatus,
			},
		},
		{
			name: "partial failed",
			codes: HTTPStatusCodes{
				Failed: http.StatusBadGateway,
			},
			want: HTTPStatusCodes{
				Passed: DefaultPassedStatus,
				Failed: http.StatusBadGateway,
				Error:  DefaultErrorStatus,
			},
		},
		{
			name: "partial error",
			codes: HTTPStatusCodes{
				Error: http.StatusBadGateway,
			},
			want: HTTPStatusCodes{
				Passed: DefaultPassedStatus,
				Failed: DefaultFailedStatus,
				Error:  http.StatusBadGateway,
			},
		},
		{
			name: "complete",
			codes: HTTPStatusCodes{
				Passed: http.StatusNoContent,
				Failed: http.StatusTooManyRequests,
				Error:  http.StatusBadGateway,
			},
			want: HTTPStatusCodes{
				Passed: http.StatusNoContent,
				Failed: http.StatusTooManyRequests,
				Error:  http.StatusBadGateway,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.codes.Normalize(); got != tc.want {
				t.Fatalf("Normalize() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestHTTPStatusCodesValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		codes HTTPStatusCodes
		want  bool
	}{
		{name: "defaults", codes: DefaultStatusCodes(), want: true},
		{name: "zero normalizes", codes: HTTPStatusCodes{}, want: true},
		{
			name: "custom valid",
			codes: HTTPStatusCodes{
				Passed: http.StatusNoContent,
				Failed: http.StatusTooManyRequests,
				Error:  http.StatusBadGateway,
			},
			want: true,
		},
		{
			name: "passed below range",
			codes: HTTPStatusCodes{
				Passed: 99,
				Failed: DefaultFailedStatus,
				Error:  DefaultErrorStatus,
			},
			want: false,
		},
		{
			name: "passed not 2xx",
			codes: HTTPStatusCodes{
				Passed: http.StatusServiceUnavailable,
				Failed: DefaultFailedStatus,
				Error:  DefaultErrorStatus,
			},
			want: false,
		},
		{
			name: "failed not error class",
			codes: HTTPStatusCodes{
				Passed: DefaultPassedStatus,
				Failed: http.StatusOK,
				Error:  DefaultErrorStatus,
			},
			want: false,
		},
		{
			name: "error not 5xx",
			codes: HTTPStatusCodes{
				Passed: DefaultPassedStatus,
				Failed: DefaultFailedStatus,
				Error:  http.StatusBadRequest,
			},
			want: false,
		},
		{
			name: "error above range",
			codes: HTTPStatusCodes{
				Passed: DefaultPassedStatus,
				Failed: DefaultFailedStatus,
				Error:  600,
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.codes.Validate()
			if got := err == nil; got != tc.want {
				t.Fatalf("Validate() ok = %v, want %v; err=%v", got, tc.want, err)
			}
		})
	}
}

func TestHTTPStatusCodesStatusForReport(t *testing.T) {
	t.Parallel()

	codes := HTTPStatusCodes{
		Passed: http.StatusNoContent,
		Failed: http.StatusTooManyRequests,
		Error:  http.StatusBadGateway,
	}

	if got := codes.statusForReport(true); got != http.StatusNoContent {
		t.Fatalf("statusForReport(true) = %d, want %d", got, http.StatusNoContent)
	}
	if got := codes.statusForReport(false); got != http.StatusTooManyRequests {
		t.Fatalf("statusForReport(false) = %d, want %d", got, http.StatusTooManyRequests)
	}
}

func TestHTTPStatusCodesStatusForError(t *testing.T) {
	t.Parallel()

	codes := HTTPStatusCodes{
		Passed: http.StatusNoContent,
		Failed: http.StatusTooManyRequests,
		Error:  http.StatusBadGateway,
	}

	if got := codes.statusForError(); got != http.StatusBadGateway {
		t.Fatalf("statusForError() = %d, want %d", got, http.StatusBadGateway)
	}
}

func TestValidHTTPStatusCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code int
		want bool
	}{
		{name: "below", code: 99, want: false},
		{name: "min", code: 100, want: true},
		{name: "ok", code: http.StatusOK, want: true},
		{name: "max", code: 599, want: true},
		{name: "above", code: 600, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := validHTTPStatusCode(tc.code); got != tc.want {
				t.Fatalf("validHTTPStatusCode(%d) = %v, want %v", tc.code, got, tc.want)
			}
		})
	}
}

func TestStatusClass(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		code  int
		class int
		want  bool
	}{
		{name: "2xx ok", code: http.StatusOK, class: 2, want: true},
		{name: "2xx no content", code: http.StatusNoContent, class: 2, want: true},
		{name: "5xx service unavailable", code: http.StatusServiceUnavailable, class: 5, want: true},
		{name: "4xx too many requests", code: http.StatusTooManyRequests, class: 4, want: true},
		{name: "2xx false for 503", code: http.StatusServiceUnavailable, class: 2, want: false},
		{name: "5xx false for 404", code: http.StatusNotFound, class: 5, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := statusClass(tc.code, tc.class); got != tc.want {
				t.Fatalf("statusClass(%d, %d) = %v, want %v", tc.code, tc.class, got, tc.want)
			}
		})
	}
}
