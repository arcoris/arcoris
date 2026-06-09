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

// Int8Descriptor builds int8 descriptors.
//
// Int8Descriptor records portable fixed-width int8 constraints. It does not use Go
// platform-sized int semantics, so descriptors remain stable across
// architectures and code generators.
type Int8Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorInt8 constraints under construction.
	payload int8Payload
}

// Int8 returns an int8 descriptor builder.
//
// Typical reusable declaration:
//
//	priorityType := Int8().Range(0, 10)
func Int8() Int8Descriptor {
	return Int8Descriptor{header: newHeader(DescriptorInt8)}
}

// Nullable returns an int8 descriptor that admits null values.
func (desc Int8Descriptor) Nullable() Int8Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive int8 lower bound.
func (desc Int8Descriptor) Min(n int8) Int8Descriptor {
	desc.payload.min = limit[int8]{value: n, set: true}

	return desc
}

// Max sets the inclusive int8 upper bound.
func (desc Int8Descriptor) Max(n int8) Int8Descriptor {
	desc.payload.max = limit[int8]{value: n, set: true}

	return desc
}

// Range sets the inclusive int8 lower and upper bounds.
func (desc Int8Descriptor) Range(min, max int8) Int8Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted int8 literals in declaration order.
func (desc Int8Descriptor) Enum(values ...int8) Int8Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Int8Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.int8 = cloneInt8Payload(desc.payload)

	return out
}

// descriptorExpr marks Int8Descriptor as a sealed DescriptorExpr implementation.
func (desc Int8Descriptor) descriptorExpr() {}
