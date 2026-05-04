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
	// headerContentType is the HTTP header used to declare response body type.
	//
	// The Go standard library does not provide a constant for this header name.
	headerContentType = "Content-Type"

	// headerCacheControl is the HTTP header used to prevent health response
	// caching.
	//
	// Health endpoints expose current runtime state. Their responses should not
	// be cached by clients, proxies, or intermediate infrastructure.
	headerCacheControl = "Cache-Control"

	// headerXContentTypeOptions is the HTTP header used to prevent content type
	// sniffing by clients.
	headerXContentTypeOptions = "X-Content-Type-Options"
)

const (
	// headerValueNoStore disables storing health responses.
	headerValueNoStore = "no-store"

	// headerValueNoSniff disables response content sniffing.
	headerValueNoSniff = "nosniff"
)

const (
	// textOK is the compact successful text response body.
	textOK = "ok\n"

	// textUnhealthy is the compact failed-health text response body.
	textUnhealthy = "unhealthy\n"

	// textHandlerError is the compact adapter/evaluator error response body.
	//
	// The message is intentionally generic. Raw evaluator errors, causes, panic
	// details, and internal diagnostic values must not be rendered by default.
	textHandlerError = "health handler error\n"
)

// errorResponse is the safe JSON representation of an adapter/evaluator boundary
// error.
//
// It intentionally does not contain the raw error string. Handler construction,
// evaluator, and adapter failures may include internal implementation details
// that should stay in logs or owner-controlled diagnostics.
type errorResponse struct {
	Error string `json:"error"`
}

// renderReport writes an HTTP response for a successfully evaluated health
// report.
//
// passed must be computed with the same policy stored in config. renderReport
// uses passed both for HTTP status selection and response DTO construction.
//
// renderReport never exposes Result.Cause. For detailed output it delegates to
// response DTO construction, which copies only safe fields.
func renderReport(w http.ResponseWriter, r *http.Request, config config, report health.Report, passed bool) {
	switch config.format {
	case FormatJSON:
		renderJSONReport(w, r, config, report, passed)
	default:
		renderTextReport(w, r, config, report, passed)
	}
}

// renderHandlerError writes an HTTP response for adapter/evaluator boundary
// errors.
//
// The concrete error is intentionally not accepted as an argument. This keeps
// the rendering boundary safe by construction: callers cannot accidentally pass
// raw internal errors into public HTTP responses.
func renderHandlerError(w http.ResponseWriter, r *http.Request, config config) {
	format := config.format
	if !format.IsValid() {
		format = FormatText
	}

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
//
// DetailNone returns a compact probe-oriented body. DetailFailed and DetailAll
// include safe check-level diagnostics selected by DetailLevel. Result.Cause is
// never rendered.
func textReportBody(
	report health.Report,
	passed bool,
	policy health.TargetPolicy,
	detail DetailLevel,
) string {
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
//
// The line includes check name, status, optional reason, and optional safe
// message. It intentionally ignores Result.Cause.
func writeTextCheck(builder *strings.Builder, result health.Result, policy health.TargetPolicy) {
	marker := "[-]"
	if policy.Passes(result.Status) {
		marker = "[+]"
	}

	fmt.Fprintf(
		builder,
		"%s %s %s",
		marker,
		result.Name,
		result.Status.String(),
	)

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
//
// HEAD responses should include the same status and headers as GET responses but
// no body. A nil request is treated as a normal body-allowed call so tests and
// defensive internal usage remain simple.
func suppressBody(r *http.Request) bool {
	return r != nil && r.Method == http.MethodHead
}
