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

import "testing"

func TestDefaultEncodeNumberConfig(t *testing.T) {
	t.Parallel()

	config := defaultEncodeNumberConfig()

	if config.MaxDigits != DefaultEncodeMaxNumberDigits {
		t.Fatalf("max digits = %d; want %d", config.MaxDigits, DefaultEncodeMaxNumberDigits)
	}
	if config.DecimalScale != DecimalScalePreserve {
		t.Fatalf("decimal scale = %d; want preserve", config.DecimalScale)
	}
	if config.FloatFormat != FloatFormatShortestRoundTrip {
		t.Fatalf("float format = %d; want shortest round trip", config.FloatFormat)
	}
	if config.NegativeZero != NegativeZeroNormalize {
		t.Fatalf("negative zero = %d; want normalize", config.NegativeZero)
	}
}

func TestResolveEncodeNumberConfig(t *testing.T) {
	t.Parallel()

	config := EncodeNumberConfig{}
	resolveEncodeNumberConfig(&config)

	if config.MaxDigits != DefaultEncodeMaxNumberDigits {
		t.Fatalf("max digits = %d; want %d", config.MaxDigits, DefaultEncodeMaxNumberDigits)
	}
	if config.DecimalScale == DecimalScaleDefault {
		t.Fatalf("decimal scale still default")
	}
	if config.FloatFormat == FloatFormatDefault {
		t.Fatalf("float format still default")
	}
	if config.NegativeZero == NegativeZeroDefault {
		t.Fatalf("negative zero still default")
	}
}

func TestValidateEncodeNumberConfigRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config EncodeNumberConfig
		target error
		path   string
	}{
		"max digits": {
			config: EncodeNumberConfig{
				MaxDigits:    -1,
				DecimalScale: DecimalScalePreserve,
				FloatFormat:  FloatFormatShortestRoundTrip,
				NegativeZero: NegativeZeroNormalize,
			},
			target: ErrInvalidConfig,
			path:   "encode.numbers.max_digits",
		},
		"decimal scale unknown": {
			config: EncodeNumberConfig{
				MaxDigits:    DefaultEncodeMaxNumberDigits,
				DecimalScale: DecimalScaleMode(99),
				FloatFormat:  FloatFormatShortestRoundTrip,
				NegativeZero: NegativeZeroNormalize,
			},
			target: ErrInvalidConfig,
			path:   "encode.numbers.decimal_scale",
		},
		"decimal scale unsupported": {
			config: EncodeNumberConfig{
				MaxDigits:    DefaultEncodeMaxNumberDigits,
				DecimalScale: DecimalScaleTrimTrailingZeros,
				FloatFormat:  FloatFormatShortestRoundTrip,
				NegativeZero: NegativeZeroNormalize,
			},
			target: ErrUnsupportedConfig,
			path:   "encode.numbers.decimal_scale",
		},
		"float format unknown": {
			config: EncodeNumberConfig{
				MaxDigits:    DefaultEncodeMaxNumberDigits,
				DecimalScale: DecimalScalePreserve,
				FloatFormat:  FloatFormatMode(99),
				NegativeZero: NegativeZeroNormalize,
			},
			target: ErrInvalidConfig,
			path:   "encode.numbers.float_format",
		},
		"negative zero unknown": {
			config: EncodeNumberConfig{
				MaxDigits:    DefaultEncodeMaxNumberDigits,
				DecimalScale: DecimalScalePreserve,
				FloatFormat:  FloatFormatShortestRoundTrip,
				NegativeZero: NegativeZeroMode(99),
			},
			target: ErrInvalidConfig,
			path:   "encode.numbers.negative_zero",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			err := validateEncodeNumberConfig(testCase.config)
			requireConfigErrorIs(t, err, testCase.target)
			requireErrorTextContains(t, err, testCase.path)
		})
	}
}
