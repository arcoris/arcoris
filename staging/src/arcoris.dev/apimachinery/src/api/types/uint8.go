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

// Uint8Type builds uint8 descriptors.
//
// Uint8Type records portable fixed-width uint8 constraints. It deliberately
// avoids Go platform-sized uint semantics, so descriptors remain stable
// across architectures and code generators.
type Uint8Type struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores TypeUint8 constraints under construction.
	payload uint8Payload
}

// Uint8 returns a uint8 descriptor builder.
//
// Typical reusable declaration:
//
//	percentageType := Uint8().Max(100)
func Uint8() Uint8Type {
	return Uint8Type{header: newHeader(TypeUint8)}
}

// Nullable returns a uint8 descriptor that admits null values.
func (t Uint8Type) Nullable() Uint8Type {
	t.header = t.header.withNullable()
	return t
}

// Min sets the inclusive uint8 lower bound.
func (t Uint8Type) Min(n uint8) Uint8Type {
	t.payload.min = limit[uint8]{value: n, set: true}
	return t
}

// Max sets the inclusive uint8 upper bound.
func (t Uint8Type) Max(n uint8) Uint8Type {
	t.payload.max = limit[uint8]{value: n, set: true}
	return t
}

// Range sets the inclusive uint8 lower and upper bounds.
func (t Uint8Type) Range(min, max uint8) Uint8Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted uint8 literals in declaration order.
func (t Uint8Type) Enum(values ...uint8) Uint8Type {
	t.payload.enum = cloneSlice(values)
	return t
}

// Type returns a detached Type descriptor.
func (t Uint8Type) Type() Type {
	out := typeFromHeader(t.header)
	out.uint8 = cloneUint8Payload(t.payload)
	return out
}

// typeExpr marks Uint8Type as a sealed TypeExpr implementation.
func (t Uint8Type) typeExpr() {}
