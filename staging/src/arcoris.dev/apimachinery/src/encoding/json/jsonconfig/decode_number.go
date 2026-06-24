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

// DecodeNumberConfig controls JSON number interpretation.
type DecodeNumberConfig struct {
	// Mode controls generic JSON number preservation.
	Mode NumberMode
}

// defaultDecodeNumberConfig returns exact number decoding.
func defaultDecodeNumberConfig() DecodeNumberConfig {
	return DecodeNumberConfig{Mode: NumberModeExact}
}

// resolveDecodeNumberConfig applies number defaults in place.
func resolveDecodeNumberConfig(config *DecodeNumberConfig) {
	if config.Mode == NumberModeDefault {
		config.Mode = NumberModeExact
	}
}

// validateDecodeNumberConfig checks the configured number parser policy.
func validateDecodeNumberConfig(config DecodeNumberConfig) error {
	if !isKnownNumberMode(config.Mode) {
		return invalidConfig("decode.numbers.mode", "unknown number mode %d", config.Mode)
	}
	if config.Mode != NumberModeExact {
		return unsupportedConfig("decode.numbers.mode", "number mode %d is not implemented", config.Mode)
	}

	return nil
}
