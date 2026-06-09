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

// Uint64Descriptor builds uint64 descriptors.
//
// Uint64Descriptor records portable fixed-width uint64 constraints. Codecs and
// generators may need special handling for targets that cannot precisely
// represent the full uint64 range.
type Uint64Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorUint64 constraints under construction.
	payload uint64Payload
}

// Uint64 returns a uint64 descriptor builder.
//
// Typical reusable declaration:
//
//	sizeType := Uint64().Min(0)
func Uint64() Uint64Descriptor {
	return Uint64Descriptor{header: newHeader(DescriptorUint64)}
}

// Nullable returns a uint64 descriptor that admits null values.
func (desc Uint64Descriptor) Nullable() Uint64Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive uint64 lower bound.
func (desc Uint64Descriptor) Min(n uint64) Uint64Descriptor {
	desc.payload.min = limit[uint64]{value: n, set: true}

	return desc
}

// Max sets the inclusive uint64 upper bound.
func (desc Uint64Descriptor) Max(n uint64) Uint64Descriptor {
	desc.payload.max = limit[uint64]{value: n, set: true}

	return desc
}

// Range sets the inclusive uint64 lower and upper bounds.
func (desc Uint64Descriptor) Range(min, max uint64) Uint64Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted uint64 literals in declaration order.
func (desc Uint64Descriptor) Enum(values ...uint64) Uint64Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Uint64Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.uint64 = cloneUint64Payload(desc.payload)

	return out
}

// descriptorExpr marks Uint64Descriptor as a sealed DescriptorExpr implementation.
func (desc Uint64Descriptor) descriptorExpr() {}
