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
	// headerAllow is the HTTP response header used to advertise methods accepted
	// by a resource after a 405 Method Not Allowed response.
	//
	// The Go standard library provides constants for HTTP methods and status
	// codes, but it does not provide a dedicated constant for the Allow header.
	headerAllow = "Allow"
)

var (
	// allowedMethodsHeader is the stable value for the HTTP Allow header emitted
	// by health HTTP endpoints.
	//
	// Method names come from net/http constants so the package does not duplicate
	// standard HTTP method strings in control-flow logic. The comma-separated
	// formatting is owned by this package because it is response-header
	// presentation, not method identity.
	allowedMethodsHeader = http.MethodGet + ", " + http.MethodHead
)

// methodAllowed reports whether method is accepted by health HTTP handlers.
//
// Health endpoints are read-only infrastructure endpoints. They accept GET for
// normal probe and diagnostic reads, and HEAD for clients that only need response
// headers and status code. Mutating methods and preflight/admin-style methods
// are intentionally rejected by default.
//
// The check is case-sensitive because HTTP method tokens are case-sensitive and
// the net/http package exposes canonical method constants such as http.MethodGet
// and http.MethodHead.
func methodAllowed(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead:
		return true
	default:
		return false
	}
}

// writeMethodNotAllowed writes a 405 Method Not Allowed response for unsupported
// health endpoint methods.
//
// The response includes an Allow header listing the methods supported by this
// package. The body is intentionally generic and does not mention handler
// internals, health targets, evaluator state, checks, reports, or reasons.
func writeMethodNotAllowed(w http.ResponseWriter) {
	w.Header().Set(headerAllow, allowedMethodsHeader)
	http.Error(
		w,
		http.StatusText(http.StatusMethodNotAllowed),
		http.StatusMethodNotAllowed,
	)
}
