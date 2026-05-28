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

// Uint64Type builds uint64 descriptors.
//
// Uint64Type records portable fixed-width uint64 constraints. Codecs and
// generators may need special handling for targets that cannot precisely
// represent the full uint64 range.
type Uint64Type struct {
	// header stores the descriptor family and cross-family flags under construction.
	header typeHeader
	// payload stores TypeUint64 constraints under construction.
	payload uint64Payload
}

// Uint64 returns a uint64 descriptor builder.
//
// Typical reusable declaration:
//
//	sizeType := Uint64().Min(0)
func Uint64() Uint64Type {
	return Uint64Type{header: newHeader(TypeUint64)}
}

// Nullable returns a uint64 descriptor that admits null values.
func (t Uint64Type) Nullable() Uint64Type {
	t.header = t.header.withNullable()
	return t
}

// Min sets the inclusive uint64 lower bound.
func (t Uint64Type) Min(n uint64) Uint64Type {
	t.payload.min = uint64Limit{value: n, set: true}
	return t
}

// Max sets the inclusive uint64 upper bound.
func (t Uint64Type) Max(n uint64) Uint64Type {
	t.payload.max = uint64Limit{value: n, set: true}
	return t
}

// Range sets the inclusive uint64 lower and upper bounds.
func (t Uint64Type) Range(min, max uint64) Uint64Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted uint64 literals in declaration order.
func (t Uint64Type) Enum(values ...uint64) Uint64Type {
	t.payload.enum = cloneUint64s(values)
	return t
}

// Type returns a detached Type descriptor.
func (t Uint64Type) Type() Type {
	out := typeFromHeader(t.header)
	out.uint64 = cloneUint64Payload(t.payload)
	return out
}

// typeExpr marks Uint64Type as a sealed TypeExpr implementation.
func (t Uint64Type) typeExpr() {}
