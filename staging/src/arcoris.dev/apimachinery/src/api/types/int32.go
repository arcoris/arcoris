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

// Int32Descriptor builds int32 descriptors.
//
// Int32Descriptor records portable fixed-width int32 constraints. It does not use
// Go platform-sized int semantics, so descriptors remain stable across
// architectures and code generators.
type Int32Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorInt32 constraints under construction.
	payload int32Payload
}

// Int32 returns an int32 descriptor builder.
//
// Typical reusable declaration:
//
//	replicasType := Int32().Range(1, 1000)
func Int32() Int32Descriptor {
	return Int32Descriptor{header: newHeader(DescriptorInt32)}
}

// Nullable returns an int32 descriptor that admits null values.
func (desc Int32Descriptor) Nullable() Int32Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive int32 lower bound.
func (desc Int32Descriptor) Min(n int32) Int32Descriptor {
	desc.payload.min = limit[int32]{value: n, set: true}

	return desc
}

// Max sets the inclusive int32 upper bound.
func (desc Int32Descriptor) Max(n int32) Int32Descriptor {
	desc.payload.max = limit[int32]{value: n, set: true}

	return desc
}

// Range sets the inclusive int32 lower and upper bounds.
func (desc Int32Descriptor) Range(min, max int32) Int32Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted int32 literals in declaration order.
func (desc Int32Descriptor) Enum(values ...int32) Int32Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Int32Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.int32 = cloneInt32Payload(desc.payload)

	return out
}

// descriptorExpr marks Int32Descriptor as a sealed DescriptorExpr implementation.
func (desc Int32Descriptor) descriptorExpr() {}
