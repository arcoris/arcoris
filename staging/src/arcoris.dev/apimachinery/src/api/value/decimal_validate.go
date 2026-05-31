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

// validateDecimalDigits checks first-pass decimal digit grammar.
//
// Decimal text must have a non-empty integer part and may contain only ASCII
// digits after dot splitting. Unicode digits and exponent notation are not part
// of the first-pass portable value grammar.
func validateDecimalDigits(integerPart string, fractionPart string, hasFraction bool) error {
	if integerPart == "" {
		return invalidDecimal("decimal text must include an integer part")
	}
	if hasFraction && fractionPart == "" {
		return invalidDecimal("decimal text must include fractional digits after decimal point")
	}
	if !allASCIIDigits(integerPart) || !allASCIIDigits(fractionPart) {
		return invalidDecimal("decimal text must contain only ASCII digits and one decimal point")
	}

	return nil
}

// canonicalDecimalParts trims leading integer zeros and preserves fractional scale.
//
// Zero is collapsed to coefficient "0" and scale 0. Non-zero values preserve
// the number of fractional digits from input so diagnostics can round-trip
// intentional scale.
func canonicalDecimalParts(integerPart, fractionPart string) (string, uint32) {
	integerPart = strings.TrimLeft(integerPart, "0")

	digits := strings.TrimLeft(integerPart+fractionPart, "0")
	if digits == "" {
		return "0", 0
	}

	return digits, uint32(len(fractionPart))
}

// allASCIIDigits reports whether s contains only ASCII decimal digits.
//
// The helper intentionally avoids unicode digit classes because API payload
// number syntax is protocol text, not display text.
func allASCIIDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	return true
}

// isZeroDigits reports whether every byte in s is the ASCII digit zero.
//
// It is used to normalize manually assembled Decimal values that carry zero in
// a non-canonical coefficient.
func isZeroDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '0' {
			return false
		}
	}

	return true
}

// invalidDecimal returns a structured decimal construction error.
//
// Decimal parsing has no nested causes, so the detail text is the primary
// human-readable diagnostic while the sentinel/reason remain stable.
func invalidDecimal(detail string) error {
	return &Error{
		Path:   pathDecimal,
		Err:    ErrInvalidDecimal,
		Reason: ErrorReasonInvalidDecimal,
		Detail: detail,
	}
}
