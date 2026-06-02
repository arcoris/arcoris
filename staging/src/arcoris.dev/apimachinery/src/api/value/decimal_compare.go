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

// Compare orders decimal numeric values exactly.
//
// Compare ignores representational scale differences that do not change the
// numeric value, so 1.0 and 1.00 compare equal even though Equal keeps payload
// representation equality strict. The comparison never converts through binary
// floating point.
func (d Decimal) Compare(other Decimal) int {
	left := decimalComparisonParts(d)
	right := decimalComparisonParts(other)

	switch {
	case left.negative && !right.negative:
		return -1
	case !left.negative && right.negative:
		return 1
	}

	result := compareDecimalMagnitude(left, right)
	if left.negative {
		return -result
	}

	return result
}

// decimalCompareParts stores a Decimal in numeric-comparison form.
type decimalCompareParts struct {
	negative    bool
	coefficient string
	scale       int
}

// decimalComparisonParts normalizes zero and insignificant fractional zeros.
func decimalComparisonParts(decimal Decimal) decimalCompareParts {
	coefficient := decimal.Coefficient()
	scale := int(decimal.Scale())

	if decimal.IsZero() {
		return decimalCompareParts{coefficient: "0"}
	}

	for scale > 0 && strings.HasSuffix(coefficient, "0") {
		coefficient = strings.TrimSuffix(coefficient, "0")
		scale--
	}

	return decimalCompareParts{
		negative:    decimal.IsNegative(),
		coefficient: coefficient,
		scale:       scale,
	}
}

// compareDecimalMagnitude compares absolute decimal values.
func compareDecimalMagnitude(left, right decimalCompareParts) int {
	leftIntegerDigits := len(left.coefficient) - left.scale
	rightIntegerDigits := len(right.coefficient) - right.scale

	if leftIntegerDigits < rightIntegerDigits {
		return -1
	}
	if leftIntegerDigits > rightIntegerDigits {
		return 1
	}

	commonScale := max(left.scale, right.scale)
	leftDigits := left.coefficient + strings.Repeat("0", commonScale-left.scale)
	rightDigits := right.coefficient + strings.Repeat("0", commonScale-right.scale)

	return strings.Compare(leftDigits, rightDigits)
}
