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

// DecimalValue constructs a decimal Value from an already canonical Decimal.
//
// The function still normalizes zero defensively because Decimal has private
// fields and tests in this package may construct edge values directly.
func DecimalValue(v Decimal) Value {
	if v.coefficient == "" || isZeroDigits(v.coefficient) {
		v = Decimal{coefficient: "0"}
	}

	return Value{kind: KindDecimal, decimalValue: v}
}
