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

// Int32Type builds int32 descriptors.
//
// Int32Type records portable fixed-width int32 constraints. It does not use
// Go platform-sized int semantics, so descriptors remain stable across
// architectures and code generators.
type Int32Type struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores TypeInt32 constraints under construction.
	payload int32Payload
}

// Int32 returns an int32 descriptor builder.
//
// Typical reusable declaration:
//
//	replicasType := Int32().Range(1, 1000)
func Int32() Int32Type {
	return Int32Type{header: newHeader(TypeInt32)}
}

// Nullable returns an int32 descriptor that admits null values.
func (t Int32Type) Nullable() Int32Type {
	t.header = t.header.withNullable()

	return t
}

// Min sets the inclusive int32 lower bound.
func (t Int32Type) Min(n int32) Int32Type {
	t.payload.min = limit[int32]{value: n, set: true}

	return t
}

// Max sets the inclusive int32 upper bound.
func (t Int32Type) Max(n int32) Int32Type {
	t.payload.max = limit[int32]{value: n, set: true}

	return t
}

// Range sets the inclusive int32 lower and upper bounds.
func (t Int32Type) Range(min, max int32) Int32Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted int32 literals in declaration order.
func (t Int32Type) Enum(values ...int32) Int32Type {
	t.payload.enum = cloneSlice(values)

	return t
}

// Type returns a detached Type descriptor.
func (t Int32Type) Type() Type {
	out := typeFromHeader(t.header)
	out.int32 = cloneInt32Payload(t.payload)

	return out
}

// typeExpr marks Int32Type as a sealed TypeExpr implementation.
func (t Int32Type) typeExpr() {}
