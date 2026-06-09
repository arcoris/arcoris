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

// DateDescriptor builds descriptors for calendar dates without a time of day.
//
// DateDescriptor records a calendar-date descriptor without time-of-day semantics.
// It does not choose an external textual or numeric encoding.
type DateDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores the exact date constraints under construction.
	payload datePayload
}

// Date returns a date descriptor builder.
//
// Typical reusable declaration:
//
//	effectiveDateType := Date().Nullable()
func Date() DateDescriptor {
	return DateDescriptor{header: newHeader(DescriptorDate)}
}

// Nullable returns a date descriptor that admits null values.
func (desc DateDescriptor) Nullable() DateDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc DateDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.date = cloneDatePayload(desc.payload)

	return out
}

// descriptorExpr marks DateDescriptor as a sealed DescriptorExpr implementation.
func (desc DateDescriptor) descriptorExpr() {}
