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

// Uint32Descriptor builds uint32 descriptors.
//
// Uint32Descriptor records portable fixed-width uint32 constraints. It deliberately
// avoids Go platform-sized uint semantics, so descriptors remain stable
// across architectures and code generators.
type Uint32Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorUint32 constraints under construction.
	payload uint32Payload
}

// Uint32 returns a uint32 descriptor builder.
//
// Typical reusable declaration:
//
//	generationType := Uint32().Min(1)
func Uint32() Uint32Descriptor {
	return Uint32Descriptor{header: newHeader(DescriptorUint32)}
}

// Nullable returns a uint32 descriptor that admits null values.
func (desc Uint32Descriptor) Nullable() Uint32Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive uint32 lower bound.
func (desc Uint32Descriptor) Min(n uint32) Uint32Descriptor {
	desc.payload.min = limit[uint32]{value: n, set: true}

	return desc
}

// Max sets the inclusive uint32 upper bound.
func (desc Uint32Descriptor) Max(n uint32) Uint32Descriptor {
	desc.payload.max = limit[uint32]{value: n, set: true}

	return desc
}

// Range sets the inclusive uint32 lower and upper bounds.
func (desc Uint32Descriptor) Range(min, max uint32) Uint32Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted uint32 literals in declaration order.
func (desc Uint32Descriptor) Enum(values ...uint32) Uint32Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Uint32Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.uint32 = cloneUint32Payload(desc.payload)

	return out
}

// descriptorExpr marks Uint32Descriptor as a sealed DescriptorExpr implementation.
func (desc Uint32Descriptor) descriptorExpr() {}
