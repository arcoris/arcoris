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

// Uint8Descriptor builds uint8 descriptors.
//
// Uint8Descriptor records portable fixed-width uint8 constraints. It deliberately
// avoids Go platform-sized uint semantics, so descriptors remain stable
// across architectures and code generators.
type Uint8Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorUint8 constraints under construction.
	payload uint8Payload
}

// Uint8 returns a uint8 descriptor builder.
//
// Typical reusable declaration:
//
//	percentageType := Uint8().Max(100)
func Uint8() Uint8Descriptor {
	return Uint8Descriptor{header: newHeader(DescriptorUint8)}
}

// Nullable returns a uint8 descriptor that admits null values.
func (desc Uint8Descriptor) Nullable() Uint8Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive uint8 lower bound.
func (desc Uint8Descriptor) Min(n uint8) Uint8Descriptor {
	desc.payload.min = limit[uint8]{value: n, set: true}

	return desc
}

// Max sets the inclusive uint8 upper bound.
func (desc Uint8Descriptor) Max(n uint8) Uint8Descriptor {
	desc.payload.max = limit[uint8]{value: n, set: true}

	return desc
}

// Range sets the inclusive uint8 lower and upper bounds.
func (desc Uint8Descriptor) Range(min, max uint8) Uint8Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted uint8 literals in declaration order.
func (desc Uint8Descriptor) Enum(values ...uint8) Uint8Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Uint8Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.uint8 = cloneUint8Payload(desc.payload)

	return out
}

// descriptorExpr marks Uint8Descriptor as a sealed DescriptorExpr implementation.
func (desc Uint8Descriptor) descriptorExpr() {}
