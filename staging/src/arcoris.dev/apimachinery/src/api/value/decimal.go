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

// Decimal stores one exact base-10 decimal value.
//
// The coefficient stores significant decimal digits without a sign or decimal
// point. Scale is the number of fractional digits. Precision and scale limits
// are descriptor validation concerns; this type only preserves an exact
// canonical payload.
type Decimal struct {
	// negative records the sign for non-zero coefficients.
	negative bool
	// coefficient stores canonical ASCII digits, with "0" for zero.
	coefficient string
	// scale stores the number of fractional decimal digits.
	scale uint32
}

// IsNegative reports whether d is strictly less than zero.
//
// Decimal zero is never negative, even when a malformed private value has
// negative=true.
func (d Decimal) IsNegative() bool {
	return d.negative && !d.IsZero()
}

// Coefficient returns canonical decimal digits without sign or decimal point.
//
// It returns "0" for malformed private zero values with an empty coefficient so
// external behavior stays canonical.
func (d Decimal) Coefficient() string {
	if d.coefficient == "" {
		return "0"
	}

	return d.coefficient
}

// Scale returns the number of fractional decimal digits.
//
// Scale is preserved from parsed input for non-zero decimals. It is not a
// validation constraint by itself.
func (d Decimal) Scale() uint32 {
	return d.scale
}

// IsZero reports whether d represents decimal zero.
//
// Empty coefficient is treated as zero to keep externally visible behavior
// robust for the private zero Decimal value.
func (d Decimal) IsZero() bool {
	return d.coefficient == "" || d.coefficient == "0"
}

// String returns the canonical diagnostic decimal text for d.
//
// The text preserves fractional scale for non-zero values. It is diagnostic
// text, not a package-level codec contract.
func (d Decimal) String() string {
	if d.IsZero() {
		return "0"
	}

	return d.formatNonZero()
}

// Equal reports whether d and other represent the same canonical decimal.
//
// Scale participates in equality, so 1.20 and 1.2 remain distinct payloads in
// this first-pass value model.
func (d Decimal) Equal(other Decimal) bool {
	return d.IsNegative() == other.IsNegative() &&
		d.Coefficient() == other.Coefficient() &&
		d.Scale() == other.Scale()
}

// formatNonZero formats d after zero normalization has been handled.
//
// The formatter reconstructs the decimal point from coefficient length and
// scale without using binary floating-point conversion.
func (d Decimal) formatNonZero() string {
	coefficient := d.Coefficient()
	scale := int(d.scale)
	intDigits := len(coefficient) - scale

	var builder strings.Builder
	if d.IsNegative() {
		builder.WriteByte('-')
	}

	switch {
	case scale == 0:
		builder.WriteString(coefficient)
	case intDigits > 0:
		builder.WriteString(coefficient[:intDigits])
		builder.WriteByte('.')
		builder.WriteString(coefficient[intDigits:])
	default:
		builder.WriteString("0.")
		builder.WriteString(strings.Repeat("0", -intDigits))
		builder.WriteString(coefficient)
	}

	return builder.String()
}
