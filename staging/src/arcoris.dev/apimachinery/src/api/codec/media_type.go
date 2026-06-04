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

// MediaType identifies one concrete content media type supported by a codec.
//
// MediaType is open-world. The constants in this package are common standard
// values, but concrete codecs may expose vendor-specific media types. Codec
// media types are exact type/subtype identifiers; HTTP parameters and content
// negotiation rules are outside this package.
type MediaType string

const (
	// MediaTypeJSON identifies JSON API documents.
	MediaTypeJSON MediaType = "application/json"

	// MediaTypeYAML identifies YAML API documents.
	MediaTypeYAML MediaType = "application/yaml"

	// MediaTypeCBOR identifies CBOR API documents.
	MediaTypeCBOR MediaType = "application/cbor"
)

// String returns the raw media type text.
//
// Call Validate when canonical lowercase type/subtype form is required.
func (m MediaType) String() string {
	return string(m)
}

// IsZero reports whether m is absent.
//
// A zero MediaType is never valid codec metadata, but IsZero is intentionally
// only an absence check.
func (m MediaType) IsZero() bool {
	return m == ""
}
