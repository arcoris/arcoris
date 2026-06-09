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

// TimestampDescriptor builds descriptors for absolute points in time.
//
// TimestampDescriptor records an absolute point-in-time descriptor. It does not
// choose RFC3339, Unix time, or another encoding; codecs own concrete
// representations.
type TimestampDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores the exact timestamp constraints under construction.
	payload timestampPayload
}

// Timestamp returns a timestamp descriptor builder.
//
// Typical reusable declaration:
//
//	observedAtType := Timestamp().Nullable()
func Timestamp() TimestampDescriptor {
	return TimestampDescriptor{header: newHeader(DescriptorTimestamp)}
}

// Nullable returns a timestamp descriptor that admits null values.
func (desc TimestampDescriptor) Nullable() TimestampDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc TimestampDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.timestamp = cloneTimestampPayload(desc.payload)

	return out
}

// descriptorExpr marks TimestampDescriptor as a sealed DescriptorExpr implementation.
func (desc TimestampDescriptor) descriptorExpr() {}
