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

// Float32Type builds float32 descriptors.
//
// Float32Type records portable binary32 floating-point constraints.
// Validation rejects NaN and infinities in descriptor rules because those
// values are not stable common API literals.
//
// Float32Type describes portable IEEE-754 binary32 API values. It stores
// structural bounds and enum literals as descriptor data only; it does not
// implement value parsing, coercion, defaulting, validation, codec mapping, or
// OpenAPI/JSON Schema export.
type Float32Type struct {
	// header stores the descriptor family and cross-family flags under construction.
	header typeHeader
	// payload stores TypeFloat32 constraints under construction.
	payload float32Payload
}

// Float32 returns a float32 descriptor builder.
//
// Typical reusable declaration:
//
//	ratioType := Float32().Range(0, 1)
func Float32() Float32Type {
	return Float32Type{header: newHeader(TypeFloat32)}
}

// Nullable returns a float32 descriptor that admits null values.
func (t Float32Type) Nullable() Float32Type {
	t.header = t.header.withNullable()
	return t
}

// Min sets the inclusive float32 lower bound.
func (t Float32Type) Min(n float32) Float32Type {
	t.payload.min = float32Limit{value: n, set: true}
	return t
}

// Max sets the inclusive float32 upper bound.
func (t Float32Type) Max(n float32) Float32Type {
	t.payload.max = float32Limit{value: n, set: true}
	return t
}

// Range sets the inclusive float32 lower and upper bounds.
func (t Float32Type) Range(min, max float32) Float32Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted float32 literals in declaration order.
func (t Float32Type) Enum(values ...float32) Float32Type {
	t.payload.enum = cloneFloat32s(values)
	return t
}

// Type returns a detached Type descriptor.
func (t Float32Type) Type() Type {
	out := typeFromHeader(t.header)
	out.float32 = cloneFloat32Payload(t.payload)
	return out
}

// typeExpr marks Float32Type as a sealed TypeExpr implementation.
func (t Float32Type) typeExpr() {}
