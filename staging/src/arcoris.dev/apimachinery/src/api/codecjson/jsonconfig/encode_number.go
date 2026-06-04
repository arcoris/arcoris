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

// DefaultEncodeMaxNumberDigits bounds JSON number output by default.
const DefaultEncodeMaxNumberDigits = 4096

// EncodeNumberConfig controls JSON number rendering.
type EncodeNumberConfig struct {
	// MaxDigits bounds generated JSON number text.
	MaxDigits int

	// DecimalScale controls decimal scale preservation.
	DecimalScale DecimalScaleMode

	// FloatFormat controls float value rendering.
	FloatFormat FloatFormatMode

	// NegativeZero controls how negative floating zero is handled.
	NegativeZero NegativeZeroMode
}

// defaultEncodeNumberConfig returns safe number output policy.
func defaultEncodeNumberConfig() EncodeNumberConfig {
	return EncodeNumberConfig{
		MaxDigits:    DefaultEncodeMaxNumberDigits,
		DecimalScale: DecimalScalePreserve,
		FloatFormat:  FloatFormatShortestRoundTrip,
		NegativeZero: NegativeZeroNormalize,
	}
}

// resolveEncodeNumberConfig applies number encode defaults in place.
func resolveEncodeNumberConfig(config *EncodeNumberConfig) {
	if config.MaxDigits == 0 {
		config.MaxDigits = DefaultEncodeMaxNumberDigits
	}
	if config.DecimalScale == DecimalScaleDefault {
		config.DecimalScale = DecimalScalePreserve
	}
	if config.FloatFormat == FloatFormatDefault {
		config.FloatFormat = FloatFormatShortestRoundTrip
	}
	if config.NegativeZero == NegativeZeroDefault {
		config.NegativeZero = NegativeZeroNormalize
	}
}

// validateEncodeNumberConfig checks JSON number output policy.
func validateEncodeNumberConfig(config EncodeNumberConfig) error {
	switch {
	case config.MaxDigits <= 0:
		return invalidConfig("encode.numbers.max_digits", "must be greater than zero")
	case !isKnownDecimalScaleMode(config.DecimalScale):
		return invalidConfig("encode.numbers.decimal_scale", "unknown decimal scale mode %d", config.DecimalScale)
	case config.DecimalScale != DecimalScalePreserve:
		return unsupportedConfig("encode.numbers.decimal_scale", "decimal scale mode %d is not implemented", config.DecimalScale)
	case !isKnownFloatFormatMode(config.FloatFormat):
		return invalidConfig("encode.numbers.float_format", "unknown float format mode %d", config.FloatFormat)
	case !isKnownNegativeZeroMode(config.NegativeZero):
		return invalidConfig("encode.numbers.negative_zero", "unknown negative zero mode %d", config.NegativeZero)
	default:
		return nil
	}
}
