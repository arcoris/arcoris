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

// NullDescriptor builds the DescriptorNull literal descriptor.
//
// NullDescriptor models the null literal as its own structural descriptor. It is not a
// nullable marker for other descriptors; builders for non-null families carry
// nullability through Descriptor flags instead.
//
// NullDescriptor deliberately has no Nullable method. DescriptorNull already describes the
// null literal itself and must not also carry nullable semantics.
type NullDescriptor struct {
	// header stores the descriptor kind under construction.
	header descriptorHeader
}

// Null returns a descriptor for the null literal type.
//
// Typical reusable declaration:
//
//	nullLiteral := Null()
func Null() NullDescriptor {
	return NullDescriptor{header: newHeader(DescriptorNull)}
}

// Descriptor returns a detached Descriptor descriptor.
func (desc NullDescriptor) Descriptor() Descriptor {
	return descriptorFromHeader(desc.header)
}

// descriptorExpr marks NullDescriptor as a sealed DescriptorExpr implementation.
func (desc NullDescriptor) descriptorExpr() {}
