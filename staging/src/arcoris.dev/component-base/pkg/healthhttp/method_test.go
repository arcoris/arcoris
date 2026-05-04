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
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMethodAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		want   bool
	}{
		{name: "get", method: http.MethodGet, want: true},
		{name: "head", method: http.MethodHead, want: true},
		{name: "post", method: http.MethodPost, want: false},
		{name: "put", method: http.MethodPut, want: false},
		{name: "patch", method: http.MethodPatch, want: false},
		{name: "delete", method: http.MethodDelete, want: false},
		{name: "options", method: http.MethodOptions, want: false},
		{name: "lowercase_get", method: strings.ToLower(http.MethodGet), want: false},
		{name: "empty", method: "", want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := methodAllowed(test.method); got != test.want {
				t.Fatalf("methodAllowed(%q) = %v, want %v", test.method, got, test.want)
			}
		})
	}
}

func TestAllowedMethodsHeader(t *testing.T) {
	t.Parallel()

	want := http.MethodGet + ", " + http.MethodHead
	if allowedMethodsHeader != want {
		t.Fatalf("allowedMethodsHeader = %q, want %q", allowedMethodsHeader, want)
	}
}

func TestWriteMethodNotAllowed(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, DefaultReadyPath, nil)

	writeMethodNotAllowed(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusMethodNotAllowed)
	}

	if got := response.Header.Get(headerAllow); got != allowedMethodsHeader {
		t.Fatalf("Allow header = %q, want %q", got, allowedMethodsHeader)
	}

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		t.Fatalf("Content-Type = %q, want text/plain", contentType)
	}

	if got := response.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q, want nosniff", got)
	}
	if got := response.Header.Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
}
