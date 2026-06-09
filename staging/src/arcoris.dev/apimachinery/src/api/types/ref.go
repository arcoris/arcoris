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

// RefDescriptor builds references to named structural Definition values.
//
// DescriptorRef is the descriptor reuse mechanism. It never represents arbitrary Go
// implementations, package globals, reflection types, or runtime object types.
// Recursive Definition graphs are not supported; recursive schemas need a
// future explicit design pass before DescriptorRef can carry those semantics.
type RefDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores the exact reference target under construction.
	payload refPayload
}

// Ref returns a reference descriptor builder for name.
//
// The name parameter accepts string-like values so descriptor declarations can
// use string literals while the stored payload remains a TypeName.
//
// Typical reusable declaration:
//
//	nameRef := Ref("meta.arcoris.dev.Name").
//		Nullable()
func Ref[N ~string](name N) RefDescriptor {
	return RefDescriptor{
		header:  newHeader(DescriptorRef),
		payload: refPayload{name: TypeName(name)},
	}
}

// Nullable returns a reference descriptor that admits null values.
func (desc RefDescriptor) Nullable() RefDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc RefDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.ref = cloneRefPayload(desc.payload)

	return out
}

// descriptorExpr marks RefDescriptor as a sealed DescriptorExpr implementation.
func (desc RefDescriptor) descriptorExpr() {}
