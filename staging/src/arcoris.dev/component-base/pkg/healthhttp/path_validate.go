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
	"strings"
	"unicode"
)

// ValidatePath validates path as a local HTTP route path for a health endpoint.
//
// A valid path:
//
//   - is not empty;
//   - starts with "/";
//   - is not the root path "/";
//   - does not contain a query component;
//   - does not contain a fragment component;
//   - is not an absolute URL;
//   - is not a protocol-relative URL;
//   - does not contain whitespace;
//   - does not contain ASCII control characters;
//   - does not contain backslashes.
//
// ValidatePath intentionally does not enforce router-specific path pattern
// syntax. Different HTTP muxes support different matching rules, wildcards, and
// path normalization behavior. healthhttp only requires the path to be a stable
// local route path that is safe to pass to a mux.
//
// Examples of valid paths:
//
//   - "/startupz";
//   - "/livez";
//   - "/readyz";
//   - "/healthz";
//   - "/health";
//   - "/internal/health/ready".
//
// Examples of invalid paths:
//
//   - "";
//   - "readyz";
//   - "/";
//   - "/readyz?verbose";
//   - "/readyz#fragment";
//   - "http://localhost/readyz";
//   - "//localhost/readyz";
//   - "/readyz debug";
//   - "/readyz\\debug".
func ValidatePath(path string) error {
	if !validPath(path) {
		return InvalidPathError{Path: path}
	}

	return nil
}

// validPath reports whether path is a safe local route path.
//
// The helper is private so the package can keep ValidatePath as the public
// error-returning boundary while still allowing table-driven tests to focus on
// boolean path semantics if needed.
func validPath(path string) bool {
	if path == "" || path == "/" {
		return false
	}
	if !strings.HasPrefix(path, "/") {
		return false
	}
	if strings.HasPrefix(path, "//") {
		return false
	}
	if strings.ContainsAny(path, "?#\\") {
		return false
	}
	if strings.Contains(path, "://") {
		return false
	}
	if strings.IndexFunc(path, invalidPathRune) >= 0 {
		return false
	}

	return true
}

// invalidPathRune reports whether r is unsafe inside a health HTTP route path.
//
// The validation is intentionally conservative. Health endpoints should be
// stable infrastructure routes, not user-facing URLs or router-specific
// patterns. Whitespace and control characters make diagnostics, tests, logs, and
// mux registration harder to reason about and are rejected.
func invalidPathRune(r rune) bool {
	return unicode.IsSpace(r) || r < 0x20 || r == 0x7f
}
