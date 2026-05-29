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

package types

// DecimalType builds exact decimal descriptors.
//
// DecimalType records exact base-10 numeric shape. Exact decimal comparison
// and literal representation need a dedicated value model, so this pass only
// records precision and scale.
//
// Decimal min/max rules are intentionally absent in this package because exact
// decimal values need a dedicated value representation design before structural
// comparison rules can be portable.
type DecimalType struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores the exact decimal constraints under construction.
	payload decimalPayload
}

// Decimal returns an exact decimal descriptor builder.
//
// Typical reusable declaration:
//
//	priceType := Decimal()
//	priceType = priceType.Precision(12)
//	priceType = priceType.Scale(2)
func Decimal() DecimalType { return DecimalType{header: newHeader(TypeDecimal)} }

// Nullable returns a decimal descriptor that admits null values.
func (t DecimalType) Nullable() DecimalType { t.header = t.header.withNullable(); return t }

// Precision sets the maximum number of significant decimal digits.
func (t DecimalType) Precision(n int) DecimalType {
	t.payload.precision = limit[int]{n, true}
	return t
}

// Scale sets the number of fractional decimal digits.
//
// Scale may be set without Precision. In that form it records fractional
// shape without bounding total significant digits. When both rules are set,
// validation requires Scale to be less than or equal to Precision.
func (t DecimalType) Scale(n int) DecimalType { t.payload.scale = limit[int]{n, true}; return t }

// Type returns a detached Type descriptor.
func (t DecimalType) Type() Type {
	out := typeFromHeader(t.header)
	out.decimal = cloneDecimalPayload(t.payload)
	return out
}

// typeExpr marks DecimalType as a sealed TypeExpr implementation.
func (t DecimalType) typeExpr() {}
