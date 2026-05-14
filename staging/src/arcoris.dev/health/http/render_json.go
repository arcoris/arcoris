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
	"encoding/json"
	"net/http"
	"strings"

	"arcoris.dev/health"
)

// errorResponse is the safe JSON representation of an adapter-boundary error.
type errorResponse struct {
	Error string `json:"error"`
}

// renderJSONReport writes a JSON response for a health report.
func renderJSONReport(w http.ResponseWriter, r *http.Request, config config, report health.Report, passed bool) {
	writeCommonHeaders(w, FormatJSON)
	w.WriteHeader(config.statusCodes.statusForReport(passed))

	if suppressBody(r) {
		return
	}

	response := newResponse(report, passed, config.policy, config.detailLevel)

	// Encoder failures can only arise from the response writer boundary after
	// headers have been committed. There is no meaningful recovery path here.
	_ = json.NewEncoder(w).Encode(response)
}

// renderJSONError writes a generic JSON response for adapter-boundary failures.
func renderJSONError(w http.ResponseWriter, r *http.Request, config config) {
	writeCommonHeaders(w, FormatJSON)
	w.WriteHeader(config.statusCodes.statusForError())

	if suppressBody(r) {
		return
	}

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: strings.TrimSpace(textHandlerError),
	})
}
