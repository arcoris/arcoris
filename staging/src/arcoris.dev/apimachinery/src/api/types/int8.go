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

// Int8Type builds int8 descriptors.
//
// Int8Type records portable fixed-width int8 constraints. It does not use Go
// platform-sized int semantics, so descriptors remain stable across
// architectures and code generators.
type Int8Type struct {
	// header stores the descriptor family and cross-family flags under construction.
	header typeHeader
	// payload stores TypeInt8 constraints under construction.
	payload int8Payload
}

// Int8 returns an int8 descriptor builder.
//
// Typical reusable declaration:
//
//	priorityType := Int8().Range(0, 10)
func Int8() Int8Type {
	return Int8Type{header: newHeader(TypeInt8)}
}

// Nullable returns an int8 descriptor that admits null values.
func (t Int8Type) Nullable() Int8Type {
	t.header = t.header.withNullable()
	return t
}

// Min sets the inclusive int8 lower bound.
func (t Int8Type) Min(n int8) Int8Type {
	t.payload.min = int8Limit{value: n, set: true}
	return t
}

// Max sets the inclusive int8 upper bound.
func (t Int8Type) Max(n int8) Int8Type {
	t.payload.max = int8Limit{value: n, set: true}
	return t
}

// Range sets the inclusive int8 lower and upper bounds.
func (t Int8Type) Range(min, max int8) Int8Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted int8 literals in declaration order.
func (t Int8Type) Enum(values ...int8) Int8Type {
	t.payload.enum = cloneInt8s(values)
	return t
}

// Type returns a detached Type descriptor.
func (t Int8Type) Type() Type {
	out := typeFromHeader(t.header)
	out.int8 = cloneInt8Payload(t.payload)
	return out
}

// typeExpr marks Int8Type as a sealed TypeExpr implementation.
func (t Int8Type) typeExpr() {}
