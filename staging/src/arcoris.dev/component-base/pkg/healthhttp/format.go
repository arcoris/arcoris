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

// Format identifies the HTTP response representation used by health handlers.
//
// Format controls only the representation of the response body. It does not
// affect health evaluation, target policy, HTTP status mapping, endpoint paths,
// request methods, detail selection, logging, metrics, or diagnostics policy.
//
// The zero value is FormatText. This is intentional: plain text is the safest
// default for infrastructure probes, load balancers, simple curl checks, and
// environments where the response body is usually ignored.
type Format uint8

const (
	// FormatText renders a compact text response.
	//
	// Text is the default format. It is intended for health probes and simple
	// operational checks where the status code is the primary signal and the
	// response body should remain small and stable.
	FormatText Format = iota

	// FormatJSON renders a structured JSON response.
	//
	// JSON is intended for diagnostics and integrations that need a structured
	// view of the evaluated health report. JSON rendering must still respect the
	// package exposure model: Result.Cause, panic stacks, raw errors, and other
	// internal diagnostics must not be exposed.
	FormatJSON
)

const (
	// contentTypeText is the Content-Type value for FormatText responses.
	//
	// The value is package-private because content negotiation and response
	// writing are owned by render.go. Format only owns the mapping from format
	// identity to the media type used by the renderer.
	contentTypeText = "text/plain; charset=utf-8"

	// contentTypeJSON is the Content-Type value for FormatJSON responses.
	//
	// The value intentionally omits a charset parameter. JSON is Unicode text
	// encoded as UTF-8 by convention, and "application/json" is the most common
	// interoperable content type for HTTP JSON responses.
	contentTypeJSON = "application/json"
)

// String returns the stable diagnostic name for format.
//
// String is intended for diagnostics, tests, and error messages. It is not a
// wire-format negotiation mechanism. Unknown values return "invalid".
func (f Format) String() string {
	switch f {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	default:
		return "invalid"
	}
}

// IsValid reports whether f is a supported health HTTP response format.
func (f Format) IsValid() bool {
	switch f {
	case FormatText, FormatJSON:
		return true
	default:
		return false
	}
}

// contentType returns the HTTP Content-Type value for f.
//
// Invalid formats return an empty string. Callers that accept user-provided or
// option-provided formats should validate them before rendering.
func (f Format) contentType() string {
	switch f {
	case FormatText:
		return contentTypeText
	case FormatJSON:
		return contentTypeJSON
	default:
		return ""
	}
}

// validateFormat returns an error if format is not supported.
//
// The helper is package-private because callers can use Format.IsValid for
// boolean checks, while option parsing and constructor code need an
// error-returning boundary.
func validateFormat(format Format) error {
	if !format.IsValid() {
		return InvalidFormatError{Format: format}
	}

	return nil
}
