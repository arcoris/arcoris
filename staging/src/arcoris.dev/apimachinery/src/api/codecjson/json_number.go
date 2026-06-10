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

package codecjson

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecjson/jsonconfig"
	"arcoris.dev/apimachinery/api/value"
)

// parseJSONNumber converts one JSON number literal without float64-first parsing.
//
// The token text comes from encoding/json's json.Number path, so JSON lexical
// grammar has already been checked. This helper still classifies the literal by
// source form so integer-looking input remains an exact integer and any input
// containing a fraction or exponent remains a decimal, even when its
// mathematical value is integral.
func parseJSONNumber(path jsonPath, text string, config resolvedDecodeConfig) (value.Value, error) {
	if numberDigitCost(text) > config.maxNumberDigits {
		return value.Value{}, invalidNumber(path, "JSON number exceeds maximum digit budget")
	}

	if !strings.ContainsAny(text, ".eE") {
		return parseJSONInteger(path, text)
	}

	plain, err := expandJSONDecimal(text, config.maxNumberDigits)
	if err != nil {
		return value.Value{}, wrapAt(path, ErrInvalidNumber, codec.ErrInvalidNumber, ErrorReasonInvalidNumber, "invalid JSON decimal", err)
	}
	decimal, err := value.ParseDecimal(plain)
	if err != nil {
		return value.Value{}, wrapAt(path, ErrInvalidNumber, codec.ErrInvalidNumber, ErrorReasonInvalidNumber, "invalid JSON decimal", err)
	}

	return value.DecimalValue(decimal), nil
}

// parseJSONInteger converts one integer literal into the exact integer union.
//
// Negative JSON integers must fit int64. Non-negative integers use uint64 so
// the generic JSON codec can represent the full unsigned side of the API value
// model without passing through a decimal or float fallback.
func parseJSONInteger(path jsonPath, text string) (value.Value, error) {
	if strings.HasPrefix(text, "-") {
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return value.Value{}, invalidNumber(path, "JSON integer is outside int64 range")
		}
		return value.Int64Value(v), nil
	}

	v, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return value.Value{}, invalidNumber(path, "JSON integer is outside uint64 range")
	}

	return value.Uint64Value(v), nil
}

// expandJSONDecimal expands exponent notation into plain decimal text.
//
// Decimal expansion is deliberately string-based. It preserves decimal
// precision, avoids binary floating-point rounding, and bounds allocation before
// exponent expansion can create an unreasonably large string.
func expandJSONDecimal(text string, maxDigits int) (string, error) {
	negative := strings.HasPrefix(text, "-")
	unsigned := strings.TrimPrefix(text, "-")

	mantissa, exponentText, hasExponent := splitExponent(unsigned)
	exponent := 0
	if hasExponent {
		parsed, err := strconv.ParseInt(exponentText, 10, 32)
		if err != nil {
			return "", err
		}
		exponent = int(parsed)
	}

	integerPart, fractionPart, _ := strings.Cut(mantissa, ".")
	digits := integerPart + fractionPart
	if len(digits)+absInt(exponent) > maxDigits {
		return "", errors.New("expanded decimal exceeds maximum digit budget")
	}

	scale := len(fractionPart) - exponent
	var plain string
	switch {
	case scale <= 0:
		plain = digits + strings.Repeat("0", -scale)
	case scale >= len(digits):
		plain = "0." + strings.Repeat("0", scale-len(digits)) + digits
	default:
		point := len(digits) - scale
		plain = digits[:point] + "." + digits[point:]
	}
	if negative {
		plain = "-" + plain
	}

	return plain, nil
}

// splitExponent separates decimal mantissa and exponent text.
//
// The bool result reports whether an exponent marker was present so callers can
// distinguish an empty exponent from a plain mantissa. The JSON tokenizer has
// already rejected syntactically empty exponents.
func splitExponent(text string) (string, string, bool) {
	if index := strings.IndexAny(text, "eE"); index >= 0 {
		return text[:index], text[index+1:], true
	}

	return text, "", false
}

// numberDigitCost counts digit characters in the original JSON number.
//
// This is a cheap pre-expansion guard. It catches very large integer/fraction
// payloads before deeper parsing and complements the exponent-aware budget in
// expandJSONDecimal.
func numberDigitCost(text string) int {
	count := 0
	for _, r := range text {
		if r >= '0' && r <= '9' {
			count++
		}
	}

	return count
}

// invalidNumber returns a dual codecjson/api codec number diagnostic.
//
// The local sentinel keeps codecjson-specific matching available, while the
// root codec sentinel preserves the format-independent ErrInvalidNumber
// contract for callers working through api/codec.
func invalidNumber(path jsonPath, detail string) error {
	return errorAt(path, ErrInvalidNumber, codec.ErrInvalidNumber, ErrorReasonInvalidNumber, detail)
}

// finiteFloatText returns a round-trip JSON number for finite floats.
//
// value.FloatValue rejects non-finite values, but this defensive check keeps
// the JSON writer correct even if a future value construction path changes.
func finiteFloatText(path jsonPath, f float64, config resolvedEncodeConfig) (string, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "", invalidNumber(path, "non-finite float cannot be encoded as JSON")
	}
	if f == 0 && math.Signbit(f) {
		switch config.negativeZero {
		case jsonconfig.NegativeZeroReject:
			return "", invalidNumber(path, "negative zero cannot be encoded as JSON")
		case jsonconfig.NegativeZeroNormalize:
			return "0", nil
		}
	}
	if config.floatFormat == jsonconfig.FloatFormatReject {
		return "", errorAt(path, ErrUnsupportedValue, errors.Join(codec.ErrEncodeFailed, codec.ErrUnsupportedFeature), ErrorReasonUnsupportedValue, "float values are disabled by JSON codec config")
	}

	return strconv.FormatFloat(f, 'g', -1, 64), nil
}

// absInt returns |v| without importing math for integers.
//
// The helper is intentionally tiny because math.Abs works on float64 and would
// be the wrong tool for exponent-budget arithmetic.
func absInt(v int) int {
	if v < 0 {
		return -v
	}

	return v
}
