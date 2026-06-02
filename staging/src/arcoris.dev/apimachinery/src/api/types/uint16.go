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

// Uint16Type builds uint16 descriptors.
//
// Uint16Type records portable fixed-width uint16 constraints. It deliberately
// avoids Go platform-sized uint semantics, so descriptors remain stable
// across architectures and code generators.
type Uint16Type struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores TypeUint16 constraints under construction.
	payload uint16Payload
}

// Uint16 returns a uint16 descriptor builder.
//
// Typical reusable declaration:
//
//	portType := Uint16().Range(1, 65535)
func Uint16() Uint16Type {
	return Uint16Type{header: newHeader(TypeUint16)}
}

// Nullable returns a uint16 descriptor that admits null values.
func (t Uint16Type) Nullable() Uint16Type {
	t.header = t.header.withNullable()

	return t
}

// Min sets the inclusive uint16 lower bound.
func (t Uint16Type) Min(n uint16) Uint16Type {
	t.payload.min = limit[uint16]{value: n, set: true}

	return t
}

// Max sets the inclusive uint16 upper bound.
func (t Uint16Type) Max(n uint16) Uint16Type {
	t.payload.max = limit[uint16]{value: n, set: true}

	return t
}

// Range sets the inclusive uint16 lower and upper bounds.
func (t Uint16Type) Range(min, max uint16) Uint16Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted uint16 literals in declaration order.
func (t Uint16Type) Enum(values ...uint16) Uint16Type {
	t.payload.enum = cloneSlice(values)

	return t
}

// Type returns a detached Type descriptor.
func (t Uint16Type) Type() Type {
	out := typeFromHeader(t.header)
	out.uint16 = cloneUint16Payload(t.payload)

	return out
}

// typeExpr marks Uint16Type as a sealed TypeExpr implementation.
func (t Uint16Type) typeExpr() {}
