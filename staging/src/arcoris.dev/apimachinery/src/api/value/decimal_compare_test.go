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

import "testing"

func TestDecimalCompareEqualDifferentScale(t *testing.T) {
	requireEqual(t, MustDecimal("1.0").Compare(MustDecimal("1.00")), 0)
	requireEqual(t, MustDecimal("-1.0").Compare(MustDecimal("-1.00")), 0)
}

func TestDecimalCompareLessGreater(t *testing.T) {
	requireEqual(t, MustDecimal("1.23").Compare(MustDecimal("1.24")), -1)
	requireEqual(t, MustDecimal("10").Compare(MustDecimal("2")), 1)
	requireEqual(t, MustDecimal("0.001").Compare(MustDecimal("0.01")), -1)
}

func TestDecimalCompareNegativeValues(t *testing.T) {
	requireEqual(t, MustDecimal("-2").Compare(MustDecimal("-1")), -1)
	requireEqual(t, MustDecimal("-0.01").Compare(MustDecimal("0")), -1)
	requireEqual(t, MustDecimal("1").Compare(MustDecimal("-1")), 1)
}

func TestDecimalCompareLargeCoefficient(t *testing.T) {
	left := MustDecimal("123456789012345678901234567890.01")
	right := MustDecimal("123456789012345678901234567890.02")

	requireEqual(t, left.Compare(right), -1)
	requireEqual(t, right.Compare(left), 1)
}

func TestDecimalCompareZero(t *testing.T) {
	requireEqual(t, MustDecimal("0").Compare(MustDecimal("-0")), 0)
	requireEqual(t, Decimal{}.Compare(MustDecimal("0")), 0)
}
