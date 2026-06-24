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

package jsonconfig

// EncodeStringConfig controls JSON string output.
type EncodeStringConfig struct {
	// EscapeHTML enables encoding/json-compatible HTML escaping.
	EscapeHTML bool

	// InvalidUTF8 controls invalid string payloads before quoting.
	InvalidUTF8 InvalidUTF8Mode
}

// defaultEncodeStringConfig returns safe string output policy.
func defaultEncodeStringConfig() EncodeStringConfig {
	return EncodeStringConfig{InvalidUTF8: InvalidUTF8Reject}
}

// resolveEncodeStringConfig applies string encode defaults in place.
func resolveEncodeStringConfig(config *EncodeStringConfig) {
	if config.InvalidUTF8 == InvalidUTF8Default {
		config.InvalidUTF8 = InvalidUTF8Reject
	}
}

// validateEncodeStringConfig checks JSON string output policy.
func validateEncodeStringConfig(config EncodeStringConfig) error {
	if !isKnownInvalidUTF8Mode(config.InvalidUTF8) {
		return invalidConfig("encode.strings.invalid_utf8", "unknown invalid UTF-8 mode %d", config.InvalidUTF8)
	}
	if config.InvalidUTF8 != InvalidUTF8Reject {
		return unsupportedConfig("encode.strings.invalid_utf8", "invalid UTF-8 mode %d is not implemented", config.InvalidUTF8)
	}

	return nil
}
