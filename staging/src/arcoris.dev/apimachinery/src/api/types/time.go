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

// TimeDescriptor builds descriptors for times of day without calendar dates.
//
// TimeDescriptor records a time-of-day descriptor without a calendar date. It does
// not choose an external textual or numeric encoding.
type TimeDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores the exact time-of-day constraints under construction.
	payload timePayload
}

// Time returns a time-of-day descriptor builder.
//
// Typical reusable declaration:
//
//	startTimeType := Time().Nullable()
func Time() TimeDescriptor {
	return TimeDescriptor{header: newHeader(DescriptorTime)}
}

// Nullable returns a time descriptor that admits null values.
func (desc TimeDescriptor) Nullable() TimeDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc TimeDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.timeOfDay = cloneTimePayload(desc.payload)

	return out
}

// descriptorExpr marks TimeDescriptor as a sealed DescriptorExpr implementation.
func (desc TimeDescriptor) descriptorExpr() {}
