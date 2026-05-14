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

	"arcoris.dev/health"
)

// renderReport writes a response for a successfully evaluated health report.
func renderReport(w http.ResponseWriter, r *http.Request, cfg config, report health.Report, passed bool) {
	switch cfg.format {
	case FormatJSON:
		renderJSONReport(w, r, cfg, report, passed)
	default:
		renderTextReport(w, r, cfg, report, passed)
	}
}

// renderHandlerError writes a generic adapter-boundary failure response.
//
// The function never accepts a raw error value and therefore cannot leak one by
// accident. Callers choose only the configured format and status-code mapping.
func renderHandlerError(w http.ResponseWriter, r *http.Request, cfg config) {
	format := cfg.format
	if !format.IsValid() {
		format = FormatText
	}
	cfg.format = format

	switch format {
	case FormatJSON:
		renderJSONError(w, r, cfg)
	default:
		renderTextError(w, r, cfg)
	}
}
