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
	"fmt"
	"net/http"
	"strings"

	"arcoris.dev/component-base/pkg/health"
)

const (
	headerContentType         = "Content-Type"
	headerCacheControl        = "Cache-Control"
	headerXContentTypeOptions = "X-Content-Type-Options"
)

const (
	headerValueNoStore = "no-store"
	headerValueNoSniff = "nosniff"
)

const (
	textOK           = "ok\n"
	textUnhealthy    = "unhealthy\n"
	textHandlerError = "health handler error\n"
)

// errorResponse is the safe JSON representation of an adapter/evaluator error.
type errorResponse struct {
	Error string `json:"error"`
}

// renderReport writes an HTTP response for a successfully evaluated health
// report.
func renderReport(w http.ResponseWriter, r *http.Request, config config, report health.Report, passed bool) {
	switch config.format {
	case FormatJSON:
		renderJSONReport(w, r, config, report, passed)
	default:
		renderTextReport(w, r, config, report, passed)
	}
}

// renderHandlerError writes an HTTP response for adapter/evaluator errors.
func renderHandlerError(w http.ResponseWriter, r *http.Request, config config) {
	format := config.format
	if !format.IsValid() {
		format = FormatText
	}
	config.format = format

	switch format {
	case FormatJSON:
		renderJSONError(w, r, config)
	default:
		renderTextError(w, r, config)
	}
}

// renderTextReport writes a text response for a health report.
func renderTextReport(w http.ResponseWriter, r *http.Request, config config, report health.Report, passed bool) {
	writeCommonHeaders(w, FormatText)
	w.WriteHeader(config.statusCodes.statusForReport(passed))

	if suppressBody(r) {
		return
	}

	_, _ = w.Write([]byte(textReportBody(report, passed, config.policy, config.detailLevel)))
}

// renderJSONReport writes a JSON response for a health report.
func renderJSONReport(w http.ResponseWriter, r *http.Request, config config, report health.Report, passed bool) {
	writeCommonHeaders(w, FormatJSON)
	w.WriteHeader(config.statusCodes.statusForReport(passed))

	if suppressBody(r) {
		return
	}

	response := newResponse(report, passed, config.policy, config.detailLevel)
	_ = json.NewEncoder(w).Encode(response)
}

// renderTextError writes a safe text response for adapter/evaluator errors.
func renderTextError(w http.ResponseWriter, r *http.Request, config config) {
	writeCommonHeaders(w, FormatText)
	w.WriteHeader(config.statusCodes.statusForError())

	if suppressBody(r) {
		return
	}

	_, _ = w.Write([]byte(textHandlerError))
}

// renderJSONError writes a safe JSON response for adapter/evaluator errors.
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

// textReportBody builds a safe text response body for report.
func textReportBody(report health.Report, passed bool, policy health.TargetPolicy, detail DetailLevel) string {
	if !detail.IncludesChecks() {
		if passed {
			return textOK
		}

		return textUnhealthy
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "%s: %s\n", report.Target.String(), report.Status.String())
	fmt.Fprintf(&builder, "passed: %t\n", passed)

	for _, check := range selectChecks(report, policy, detail) {
		writeTextCheck(&builder, check, policy)
	}

	return builder.String()
}

// writeTextCheck appends one safe text check line to builder.
func writeTextCheck(builder *strings.Builder, result health.Result, policy health.TargetPolicy) {
	marker := "[-]"
	if policy.Passes(result.Status) {
		marker = "[+]"
	}

	fmt.Fprintf(builder, "%s %s %s", marker, result.Name, result.Status.String())

	if reason := formatReason(result.Reason); reason != "" {
		fmt.Fprintf(builder, " %s", reason)
	}

	if result.Message != "" {
		fmt.Fprintf(builder, ": %s", result.Message)
	}

	builder.WriteByte('\n')
}

// writeCommonHeaders writes headers shared by all health HTTP responses.
func writeCommonHeaders(w http.ResponseWriter, format Format) {
	w.Header().Set(headerContentType, format.contentType())
	w.Header().Set(headerCacheControl, headerValueNoStore)
	w.Header().Set(headerXContentTypeOptions, headerValueNoSniff)
}

// suppressBody reports whether a response body must be suppressed.
func suppressBody(r *http.Request) bool {
	return r != nil && r.Method == http.MethodHead
}
