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

// Float64Type builds float64 descriptors.
//
// Float64Type records portable binary64 floating-point constraints.
// Validation rejects NaN and infinities in descriptor rules because those
// values are not stable common API literals.
//
// Float64Type describes portable IEEE-754 binary64 API values. Descriptor
// validation rejects NaN and infinities in bounds and enums because those
// values are not stable common API literals across formats.
type Float64Type struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores TypeFloat64 constraints under construction.
	payload float64Payload
}

// Float64 returns a float64 descriptor builder.
//
// Typical reusable declaration:
//
//	weightType := Float64().Min(0)
func Float64() Float64Type {
	return Float64Type{header: newHeader(TypeFloat64)}
}

// Nullable returns a float64 descriptor that admits null values.
func (t Float64Type) Nullable() Float64Type {
	t.header = t.header.withNullable()
	return t
}

// Min sets the inclusive float64 lower bound.
func (t Float64Type) Min(n float64) Float64Type {
	t.payload.min = limit[float64]{value: n, set: true}
	return t
}

// Max sets the inclusive float64 upper bound.
func (t Float64Type) Max(n float64) Float64Type {
	t.payload.max = limit[float64]{value: n, set: true}
	return t
}

// Range sets the inclusive float64 lower and upper bounds.
func (t Float64Type) Range(min, max float64) Float64Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted float64 literals in declaration order.
func (t Float64Type) Enum(values ...float64) Float64Type {
	t.payload.enum = cloneSlice(values)
	return t
}

// Type returns a detached Type descriptor.
func (t Float64Type) Type() Type {
	out := typeFromHeader(t.header)
	out.float64 = cloneFloat64Payload(t.payload)
	return out
}

// typeExpr marks Float64Type as a sealed TypeExpr implementation.
func (t Float64Type) typeExpr() {}
