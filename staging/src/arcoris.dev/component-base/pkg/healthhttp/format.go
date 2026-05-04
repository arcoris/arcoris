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
type Format uint8

const (
	// FormatText renders a compact text response.
	FormatText Format = iota

	// FormatJSON renders a structured JSON response.
	FormatJSON
)

const (
	contentTypeText = "text/plain; charset=utf-8"
	contentTypeJSON = "application/json"
)

// String returns the stable diagnostic name for format.
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
func validateFormat(format Format) error {
	if !format.IsValid() {
		return InvalidFormatError{Format: format}
	}

	return nil
}
