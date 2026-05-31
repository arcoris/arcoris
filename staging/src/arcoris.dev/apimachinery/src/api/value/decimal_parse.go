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

package value

import "strings"

// NewDecimal parses first-pass decimal syntax without exponent notation.
//
// Supported input is a plain optional negative sign, an integer part, and an
// optional fractional part. Leading integer zeros are removed, fractional scale
// is preserved for non-zero values, and negative zero normalizes to "0".
func NewDecimal(text string) (Decimal, error) {
	negative, unsigned := splitDecimalSign(text)
	if unsigned == "" {
		return Decimal{}, invalidDecimal("decimal text is empty or sign-only")
	}
	if strings.HasPrefix(unsigned, "+") {
		return Decimal{}, invalidDecimal("plus sign is not supported")
	}

	integerPart, fractionPart, hasFraction, err := splitDecimalParts(unsigned)
	if err != nil {
		return Decimal{}, err
	}
	if err := validateDecimalDigits(integerPart, fractionPart, hasFraction); err != nil {
		return Decimal{}, err
	}

	coefficient, scale := canonicalDecimalParts(integerPart, fractionPart)
	if coefficient == "0" {
		return Decimal{coefficient: "0"}, nil
	}

	return Decimal{negative: negative, coefficient: coefficient, scale: scale}, nil
}

// MustDecimal parses a Decimal or panics when text is malformed.
//
// It is intended for fixtures and static definitions where invalid decimal text
// is a programmer error. Runtime parsing paths should use NewDecimal.
func MustDecimal(text string) Decimal {
	decimal, err := NewDecimal(text)
	if err != nil {
		panic(err)
	}

	return decimal
}

// splitDecimalSign separates the optional negative sign.
//
// A plus sign is intentionally not accepted by this first-pass grammar and is
// left in the unsigned text so NewDecimal can reject it explicitly.
func splitDecimalSign(text string) (bool, string) {
	if strings.HasPrefix(text, "-") {
		return true, strings.TrimPrefix(text, "-")
	}

	return false, text
}

// splitDecimalParts separates integer and fractional digits.
//
// The function performs only structural dot splitting. Digit validation stays
// in validateDecimalDigits so error reasons remain focused. hasFraction lets
// callers distinguish "no dot" from "dot with an empty fractional part".
func splitDecimalParts(text string) (string, string, bool, error) {
	parts := strings.Split(text, ".")
	if len(parts) > 2 {
		return "", "", false, invalidDecimal("decimal text contains multiple decimal points")
	}

	integerPart := parts[0]
	fractionPart := ""
	hasFraction := len(parts) == 2
	if len(parts) == 2 {
		fractionPart = parts[1]
	}

	return integerPart, fractionPart, hasFraction, nil
}
