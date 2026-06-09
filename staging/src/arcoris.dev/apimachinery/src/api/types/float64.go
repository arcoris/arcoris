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

// Float64Descriptor builds float64 descriptors.
//
// Float64Descriptor records portable binary64 floating-point constraints.
// Validation rejects NaN and infinities in descriptor rules because those
// values are not stable common API literals.
//
// Float64Descriptor describes portable IEEE-754 binary64 API values. Descriptor
// validation rejects NaN and infinities in bounds and enums because those
// values are not stable common API literals across formats.
type Float64Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorFloat64 constraints under construction.
	payload float64Payload
}

// Float64 returns a float64 descriptor builder.
//
// Typical reusable declaration:
//
//	weightType := Float64().Min(0)
func Float64() Float64Descriptor {
	return Float64Descriptor{header: newHeader(DescriptorFloat64)}
}

// Nullable returns a float64 descriptor that admits null values.
func (desc Float64Descriptor) Nullable() Float64Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive float64 lower bound.
func (desc Float64Descriptor) Min(n float64) Float64Descriptor {
	desc.payload.min = limit[float64]{value: n, set: true}

	return desc
}

// Max sets the inclusive float64 upper bound.
func (desc Float64Descriptor) Max(n float64) Float64Descriptor {
	desc.payload.max = limit[float64]{value: n, set: true}

	return desc
}

// Range sets the inclusive float64 lower and upper bounds.
func (desc Float64Descriptor) Range(min, max float64) Float64Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted float64 literals in declaration order.
func (desc Float64Descriptor) Enum(values ...float64) Float64Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Float64Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.float64 = cloneFloat64Payload(desc.payload)

	return out
}

// descriptorExpr marks Float64Descriptor as a sealed DescriptorExpr implementation.
func (desc Float64Descriptor) descriptorExpr() {}
