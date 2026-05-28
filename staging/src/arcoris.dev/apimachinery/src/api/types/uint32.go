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

// Uint32Type builds uint32 descriptors.
//
// Uint32Type records portable fixed-width uint32 constraints. It deliberately
// avoids Go platform-sized uint semantics, so descriptors remain stable
// across architectures and code generators.
type Uint32Type struct {
	// header stores the descriptor family and cross-family flags under construction.
	header typeHeader
	// payload stores TypeUint32 constraints under construction.
	payload uint32Payload
}

// Uint32 returns a uint32 descriptor builder.
//
// Typical reusable declaration:
//
//	generationType := Uint32().Min(1)
func Uint32() Uint32Type {
	return Uint32Type{header: newHeader(TypeUint32)}
}

// Nullable returns a uint32 descriptor that admits null values.
func (t Uint32Type) Nullable() Uint32Type {
	t.header = t.header.withNullable()
	return t
}

// Min sets the inclusive uint32 lower bound.
func (t Uint32Type) Min(n uint32) Uint32Type {
	t.payload.min = uint32Limit{value: n, set: true}
	return t
}

// Max sets the inclusive uint32 upper bound.
func (t Uint32Type) Max(n uint32) Uint32Type {
	t.payload.max = uint32Limit{value: n, set: true}
	return t
}

// Range sets the inclusive uint32 lower and upper bounds.
func (t Uint32Type) Range(min, max uint32) Uint32Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted uint32 literals in declaration order.
func (t Uint32Type) Enum(values ...uint32) Uint32Type {
	t.payload.enum = cloneUint32s(values)
	return t
}

// Type returns a detached Type descriptor.
func (t Uint32Type) Type() Type {
	out := typeFromHeader(t.header)
	out.uint32 = cloneUint32Payload(t.payload)
	return out
}

// typeExpr marks Uint32Type as a sealed TypeExpr implementation.
func (t Uint32Type) typeExpr() {}
