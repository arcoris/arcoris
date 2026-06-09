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

// Int16Descriptor builds int16 descriptors.
//
// Int16Descriptor records portable fixed-width int16 constraints. It does not use
// Go platform-sized int semantics, so descriptors remain stable across
// architectures and code generators.
type Int16Descriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores DescriptorInt16 constraints under construction.
	payload int16Payload
}

// Int16 returns an int16 descriptor builder.
//
// Typical reusable declaration:
//
//	shardType := Int16().Min(0)
func Int16() Int16Descriptor {
	return Int16Descriptor{header: newHeader(DescriptorInt16)}
}

// Nullable returns an int16 descriptor that admits null values.
func (desc Int16Descriptor) Nullable() Int16Descriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Min sets the inclusive int16 lower bound.
func (desc Int16Descriptor) Min(n int16) Int16Descriptor {
	desc.payload.min = limit[int16]{value: n, set: true}

	return desc
}

// Max sets the inclusive int16 upper bound.
func (desc Int16Descriptor) Max(n int16) Int16Descriptor {
	desc.payload.max = limit[int16]{value: n, set: true}

	return desc
}

// Range sets the inclusive int16 lower and upper bounds.
func (desc Int16Descriptor) Range(min, max int16) Int16Descriptor {
	return desc.Min(min).Max(max)
}

// Enum stores accepted int16 literals in declaration order.
func (desc Int16Descriptor) Enum(values ...int16) Int16Descriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc Int16Descriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.int16 = cloneInt16Payload(desc.payload)

	return out
}

// descriptorExpr marks Int16Descriptor as a sealed DescriptorExpr implementation.
func (desc Int16Descriptor) descriptorExpr() {}
