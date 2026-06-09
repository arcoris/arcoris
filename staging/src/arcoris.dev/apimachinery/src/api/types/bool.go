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

// BoolDescriptor builds boolean descriptors.
//
// BoolDescriptor records the structural contract for boolean API values. It has no
// value constraints in this design pass and exists as a closed descriptor
// builder rather than a Go bool wrapper.
type BoolDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
}

// Bool returns a descriptor builder for boolean values.
//
// Typical reusable declaration:
//
//	enabledType := Bool().Nullable()
func Bool() BoolDescriptor {
	return BoolDescriptor{header: newHeader(DescriptorBool)}
}

// Nullable returns a boolean descriptor that admits null values.
func (desc BoolDescriptor) Nullable() BoolDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc BoolDescriptor) Descriptor() Descriptor {
	return descriptorFromHeader(desc.header)
}

// descriptorExpr marks BoolDescriptor as a sealed DescriptorExpr implementation.
func (desc BoolDescriptor) descriptorExpr() {}
