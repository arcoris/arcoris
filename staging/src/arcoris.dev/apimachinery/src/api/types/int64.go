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

// Int64Type builds int64 descriptors.
//
// Int64Type records portable fixed-width int64 constraints. It does not use
// Go platform-sized int semantics, so descriptors remain stable across
// architectures and code generators.
//
// The builder stores descriptor rules, not runtime values. Min, Max, Range, and
// Enum record portable int64 constraints that ValidateType later checks for
// consistency. The builder is immutable-by-value: each method returns a copy so
// fluent declarations can be reused without sharing mutable slices or maps.
type Int64Type struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores TypeInt64 constraints under construction.
	payload int64Payload
}

// Int64 returns an int64 descriptor builder.
//
// Typical reusable declaration:
//
//	maxConcurrencyType := Int64().Min(1)
func Int64() Int64Type {
	return Int64Type{header: newHeader(TypeInt64)}
}

// Nullable returns an int64 descriptor that admits null values.
func (t Int64Type) Nullable() Int64Type {
	t.header = t.header.withNullable()
	return t
}

// Min sets the inclusive int64 lower bound.
//
// The bound is structural metadata. This package does not compare concrete API
// values against it; that belongs to future value-validation layers.
func (t Int64Type) Min(n int64) Int64Type {
	t.payload.min = limit[int64]{n, true}
	return t
}

// Max sets the inclusive int64 upper bound.
//
// The bound is retained as a limit[int64] so an explicit zero can be
// distinguished from an unset maximum without allocating a pointer.
func (t Int64Type) Max(n int64) Int64Type {
	t.payload.max = limit[int64]{n, true}
	return t
}

// Range sets the inclusive int64 lower and upper bounds.
func (t Int64Type) Range(min, max int64) Int64Type {
	return t.Min(min).Max(max)
}

// Enum stores accepted int64 literals in declaration order.
//
// The input slice is cloned. Later caller mutation of the variadic backing
// array cannot rewrite the descriptor returned by Type.
func (t Int64Type) Enum(values ...int64) Int64Type {
	t.payload.enum = cloneSlice(values)
	return t
}

// Type returns a detached Type descriptor.
func (t Int64Type) Type() Type {
	out := typeFromHeader(t.header)
	out.int64 = cloneInt64Payload(t.payload)
	return out
}

// typeExpr marks Int64Type as a sealed TypeExpr implementation.
func (t Int64Type) typeExpr() {}
