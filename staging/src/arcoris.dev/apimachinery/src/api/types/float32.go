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

// Float32Descriptor builds float32 descriptors.
//
// Float32Descriptor records portable binary32 floating-point constraints.
// Validation rejects NaN and infinities in descriptor rules because those
// values are not stable common API literals.
//
// Float32Descriptor describes portable IEEE-754 binary32 API values. It stores
// structural bounds and enum literals as descriptor data only; it does not
// implement value parsing, coercion, defaulting, validation, codec mapping, or
// OpenAPI/JSON Schema export.
type Float32Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorFloat32 constraints under construction.
	payload float32Payload
}

// Float32 returns a float32 descriptor builder.
//
// Typical reusable declaration:
//
//	ratioType := Float32().Range(0, 1)
func Float32() Float32Descriptor {
	return Float32Descriptor{header: newHeader(DescriptorFloat32)}
}

// Nullable returns a float32 descriptor that admits null values.
func (desc Float32Descriptor) Nullable() Float32Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive float32 lower bound.
func (desc Float32Descriptor) Min(n float32) Float32Descriptor {
	desc.payload.min = limit[float32]{value: n, set: true}

	return desc
}

// Max sets the inclusive float32 upper bound.
func (desc Float32Descriptor) Max(n float32) Float32Descriptor {
	desc.payload.max = limit[float32]{value: n, set: true}

	return desc
}

// Range sets the inclusive float32 lower and upper bounds.
func (desc Float32Descriptor) Range(min, max float32) Float32Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted float32 literals in declaration order.
func (desc Float32Descriptor) Enum(values ...float32) Float32Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Float32Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.float32 = cloneFloat32Payload(desc.payload)

	return out
}

// descriptorExpr marks Float32Descriptor as a sealed DescriptorExpr implementation.
func (desc Float32Descriptor) descriptorExpr() {}
