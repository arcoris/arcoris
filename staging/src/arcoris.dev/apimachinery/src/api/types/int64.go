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

// Int64Descriptor builds int64 descriptors.
//
// Int64Descriptor records portable fixed-width int64 constraints. It does not use
// Go platform-sized int semantics, so descriptors remain stable across
// architectures and code generators.
//
// The builder stores descriptor rules, not runtime values. Min, Max, Range, and
// Enum record portable int64 constraints that ValidateResolved later checks for
// consistency. The builder is immutable-by-value: each method returns a copy so
// fluent declarations can be reused without sharing mutable slices or maps.
type Int64Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorInt64 constraints under construction.
	payload int64Payload
}

// Int64 returns an int64 descriptor builder.
//
// Typical reusable declaration:
//
//	maxConcurrencyType := Int64().Min(1)
func Int64() Int64Descriptor {
	return Int64Descriptor{header: newHeader(DescriptorInt64)}
}

// Nullable returns an int64 descriptor that admits null values.
func (desc Int64Descriptor) Nullable() Int64Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive int64 lower bound.
//
// The bound is structural metadata. This package does not compare concrete API
// values against it; that belongs to future value-validation layers.
func (desc Int64Descriptor) Min(n int64) Int64Descriptor {
	desc.payload.min = limit[int64]{n, true}

	return desc
}

// Max sets the inclusive int64 upper bound.
//
// The bound is retained as a limit[int64] so an explicit zero can be
// distinguished from an unset maximum without allocating a pointer.
func (desc Int64Descriptor) Max(n int64) Int64Descriptor {
	desc.payload.max = limit[int64]{n, true}

	return desc
}

// Range sets the inclusive int64 lower and upper bounds.
func (desc Int64Descriptor) Range(min, max int64) Int64Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted int64 literals in declaration order.
//
// The input slice is cloned. Later caller mutation of the variadic backing
// array cannot rewrite the descriptor returned by Descriptor.
func (desc Int64Descriptor) Enum(values ...int64) Int64Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Int64Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.int64 = cloneInt64Payload(desc.payload)

	return out
}

// descriptorExpr marks Int64Descriptor as a sealed DescriptorExpr implementation.
func (desc Int64Descriptor) descriptorExpr() {}
