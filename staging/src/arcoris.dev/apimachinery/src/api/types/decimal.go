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

// DecimalDescriptor builds exact decimal descriptors.
//
// DecimalDescriptor records exact base-10 numeric shape. Exact decimal comparison
// and literal representation need a dedicated value model, so this pass only
// records precision and scale.
//
// Decimal min/max rules are intentionally absent in this package because exact
// decimal values need a dedicated value representation design before structural
// comparison rules can be portable.
type DecimalDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
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
func Decimal() DecimalDescriptor {
	return DecimalDescriptor{header: newHeader(DescriptorDecimal)}
}

// Nullable returns a decimal descriptor that admits null values.
func (desc DecimalDescriptor) Nullable() DecimalDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Precision sets the maximum number of significant decimal digits.
func (desc DecimalDescriptor) Precision(n int) DecimalDescriptor {
	desc.payload.precision = limit[int]{n, true}

	return desc
}

// Scale sets the number of fractional decimal digits.
//
// Scale may be set without Precision. In that form it records fractional
// shape without bounding total significant digits. When both rules are set,
// validation requires Scale to be less than or equal to Precision.
func (desc DecimalDescriptor) Scale(n int) DecimalDescriptor {
	desc.payload.scale = limit[int]{n, true}

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc DecimalDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.decimal = cloneDecimalPayload(desc.payload)

	return out
}

// descriptorExpr marks DecimalDescriptor as a sealed DescriptorExpr implementation.
func (desc DecimalDescriptor) descriptorExpr() {}
