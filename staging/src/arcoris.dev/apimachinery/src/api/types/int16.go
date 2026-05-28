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

// Int16Type builds int16 descriptors.
//
// Int16Type records portable fixed-width int16 constraints. It does not use
// Go platform-sized int semantics, so descriptors remain stable across
// architectures and code generators.
type Int16Type struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores TypeInt16 constraints under construction.
	payload int16Payload
}

// Int16 returns an int16 descriptor builder.
//
// Typical reusable declaration:
//
//	shardType := Int16().Min(0)
func Int16() Int16Type {
	return Int16Type{header: newHeader(TypeInt16)}
}

// Nullable returns an int16 descriptor that admits null values.
func (t Int16Type) Nullable() Int16Type {
	t.header = t.header.withNullable()
	return t
}

// Min sets the inclusive int16 lower bound.
func (t Int16Type) Min(n int16) Int16Type {
	t.payload.min = limit[int16]{value: n, set: true}
	return t
}

// Max sets the inclusive int16 upper bound.
func (t Int16Type) Max(n int16) Int16Type {
	t.payload.max = limit[int16]{value: n, set: true}
	return t
}

// Range sets the inclusive int16 lower and upper bounds.
func (t Int16Type) Range(min, max int16) Int16Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted int16 literals in declaration order.
func (t Int16Type) Enum(values ...int16) Int16Type {
	t.payload.enum = cloneSlice(values)
	return t
}

// Type returns a detached Type descriptor.
func (t Int16Type) Type() Type {
	out := typeFromHeader(t.header)
	out.int16 = cloneInt16Payload(t.payload)
	return out
}

// typeExpr marks Int16Type as a sealed TypeExpr implementation.
func (t Int16Type) typeExpr() {}
