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

// DecodeStringConfig controls raw and decoded JSON string handling.
type DecodeStringConfig struct {
	// InvalidUTF8 controls invalid raw UTF-8 input.
	InvalidUTF8 InvalidUTF8Mode
}

// defaultDecodeStringConfig returns safe string decode policy.
func defaultDecodeStringConfig() DecodeStringConfig {
	return DecodeStringConfig{InvalidUTF8: InvalidUTF8Reject}
}

// resolveDecodeStringConfig applies string decode defaults in place.
func resolveDecodeStringConfig(config *DecodeStringConfig) {
	if config.InvalidUTF8 == InvalidUTF8Default {
		config.InvalidUTF8 = InvalidUTF8Reject
	}
}

// validateDecodeStringConfig checks raw JSON string decode policy.
func validateDecodeStringConfig(config DecodeStringConfig) error {
	if !isKnownInvalidUTF8Mode(config.InvalidUTF8) {
		return invalidConfig("decode.strings.invalid_utf8", "unknown invalid UTF-8 mode %d", config.InvalidUTF8)
	}
	if config.InvalidUTF8 != InvalidUTF8Reject {
		return unsupportedConfig("decode.strings.invalid_utf8", "invalid UTF-8 mode %d is not implemented", config.InvalidUTF8)
	}

	return nil
}
