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
	"fmt"
	"net/http"
	"strings"

	"arcoris.dev/health"
)

const (
	textOK           = "ok\n"
	textUnhealthy    = "unhealthy\n"
	textHandlerError = "health handler error\n"
)

// renderTextReport writes a text response for a health report.
func renderTextReport(w http.ResponseWriter, r *http.Request, cfg config, report health.Report, passed bool) {
	writeCommonHeaders(w, FormatText)
	w.WriteHeader(cfg.statusCodes.statusForReport(passed))

	if suppressBody(r) {
		return
	}

	_, _ = w.Write([]byte(textReportBody(report, passed, cfg.policy, cfg.detailLevel)))
}

// renderTextError writes a generic text response for adapter-boundary failures.
func renderTextError(w http.ResponseWriter, r *http.Request, cfg config) {
	writeCommonHeaders(w, FormatText)
	w.WriteHeader(cfg.statusCodes.statusForError())

	if suppressBody(r) {
		return
	}

	_, _ = w.Write([]byte(textHandlerError))
}

// textReportBody builds a safe text body for report.
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

// writeTextCheck appends one safe check line to builder.
func writeTextCheck(builder *strings.Builder, res health.Result, policy health.TargetPolicy) {
	marker := "[-]"
	if policy.Passes(res.Status) {
		marker = "[+]"
	}

	fmt.Fprintf(builder, "%s %s %s", marker, res.Name, res.Status.String())

	if reason := formatReason(res.Reason); reason != "" {
		fmt.Fprintf(builder, " %s", reason)
	}

	if res.Message != "" {
		fmt.Fprintf(builder, ": %s", res.Message)
	}

	builder.WriteByte('\n')
}
