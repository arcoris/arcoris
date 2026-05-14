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
// The zero value is FormatText. Plain text is the safest default for probes and
// load balancers.
//
// Format controls representation only. It does not change health evaluation,
// target-policy pass/fail decisions, diagnostics safety rules, or content
// negotiation behavior.
type Format uint8

const (
	// FormatText renders a compact text response.
	//
	// Text is the default because it is easy for probes, load balancers, shell
	// tooling, and humans to consume without depending on structured parsers.
	FormatText Format = iota

	// FormatJSON renders a structured JSON response.
	//
	// JSON uses adapter-owned DTOs rather than exposing package health structs
	// directly.
	FormatJSON
)

const (
	contentTypeText = "text/plain; charset=utf-8"
	contentTypeJSON = "application/json"
)

// String returns the stable diagnostic name for format.
//
// Invalid values return "invalid" so diagnostics stay explicit without
// inventing a fallback format silently.
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
//
// Unsupported values are rejected during handler construction instead of being
// normalized implicitly.
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
// The helper is package-private because callers configure representation
// through Format rather than by managing raw content-type strings.
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
// Validation remains package-local so configuration failures preserve the
// adapter's typed error surface.
func validateFormat(format Format) error {
	if !format.IsValid() {
		return InvalidFormatError{Format: format}
	}

	return nil
}
