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

// Uint16Descriptor builds uint16 descriptors.
//
// Uint16Descriptor records portable fixed-width uint16 constraints. It deliberately
// avoids Go platform-sized uint semantics, so descriptors remain stable
// across architectures and code generators.
type Uint16Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorUint16 constraints under construction.
	payload uint16Payload
}

// Uint16 returns a uint16 descriptor builder.
//
// Typical reusable declaration:
//
//	portType := Uint16().Range(1, 65535)
func Uint16() Uint16Descriptor {
	return Uint16Descriptor{header: newHeader(DescriptorUint16)}
}

// Nullable returns a uint16 descriptor that admits null values.
func (desc Uint16Descriptor) Nullable() Uint16Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive uint16 lower bound.
func (desc Uint16Descriptor) Min(n uint16) Uint16Descriptor {
	desc.payload.min = limit[uint16]{value: n, set: true}

	return desc
}

// Max sets the inclusive uint16 upper bound.
func (desc Uint16Descriptor) Max(n uint16) Uint16Descriptor {
	desc.payload.max = limit[uint16]{value: n, set: true}

	return desc
}

// Range sets the inclusive uint16 lower and upper bounds.
func (desc Uint16Descriptor) Range(min, max uint16) Uint16Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted uint16 literals in declaration order.
func (desc Uint16Descriptor) Enum(values ...uint16) Uint16Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Uint16Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.uint16 = cloneUint16Payload(desc.payload)

	return out
}

// descriptorExpr marks Uint16Descriptor as a sealed DescriptorExpr implementation.
func (desc Uint16Descriptor) descriptorExpr() {}
