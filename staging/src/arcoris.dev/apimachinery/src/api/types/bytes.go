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

// BytesDescriptor builds byte-sequence descriptors with length constraints.
//
// BytesDescriptor records binary payload structure and length constraints. It
// deliberately does not choose a wire encoding such as base64 or raw bytes;
// codecs own that mapping.
type BytesDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores the exact byte-sequence constraints under construction.
	payload bytesPayload
}

// Bytes returns a descriptor builder for byte sequences.
//
// Typical reusable declaration:
//
//	payloadType := Bytes().
//		MinBytes(1).
//		MaxBytes(4096)
func Bytes() BytesDescriptor {
	return BytesDescriptor{header: newHeader(DescriptorBytes)}
}

// Nullable returns a bytes descriptor that admits null values.
func (desc BytesDescriptor) Nullable() BytesDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// MinBytes sets the inclusive minimum byte length.
func (desc BytesDescriptor) MinBytes(n int) BytesDescriptor {
	desc.payload.minBytes = limit[int]{value: n, set: true}

	return desc
}

// MaxBytes sets the inclusive maximum byte length.
func (desc BytesDescriptor) MaxBytes(n int) BytesDescriptor {
	desc.payload.maxBytes = limit[int]{value: n, set: true}

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc BytesDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.bytes = cloneBytesPayload(desc.payload)

	return out
}

// descriptorExpr marks BytesDescriptor as a sealed DescriptorExpr implementation.
func (desc BytesDescriptor) descriptorExpr() {}
