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
	headerAllow = "Allow"
)

var (
	// allowedMethodsHeader is the stable value for the HTTP Allow header emitted
	// by health HTTP endpoints.
	allowedMethodsHeader = http.MethodGet + ", " + http.MethodHead
)

// methodAllowed reports whether method is accepted by health HTTP handlers.
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
func writeMethodNotAllowed(w http.ResponseWriter) {
	w.Header().Set(headerAllow, allowedMethodsHeader)
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
