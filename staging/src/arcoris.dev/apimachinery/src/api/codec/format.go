// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codec

// Format identifies a codec's wire/document format family.
//
// Format is open-world. The constants in this package are common names, but
// concrete codecs may define additional canonical tokens. Format is not a
// media type and must not contain HTTP content-type syntax such as "/" or
// parameters.
type Format string

const (
	// FormatJSON identifies JSON document codecs.
	FormatJSON Format = "json"

	// FormatYAML identifies YAML document codecs.
	FormatYAML Format = "yaml"

	// FormatCBOR identifies CBOR document codecs.
	FormatCBOR Format = "cbor"
)

// String returns the raw format text.
//
// Call Validate when canonical lowercase token form is required.
func (f Format) String() string {
	return string(f)
}

// IsZero reports whether f is absent.
//
// A zero Format is never valid codec metadata, but IsZero is intentionally only
// an absence check.
func (f Format) IsZero() bool {
	return f == ""
}
