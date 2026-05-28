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

// BytesType builds byte-sequence descriptors with length constraints.
//
// BytesType records binary payload structure and length constraints. It
// deliberately does not choose a wire encoding such as base64 or raw bytes;
// codecs own that mapping.
type BytesType struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores the exact byte-sequence constraints under construction.
	payload bytesPayload
}

// Bytes returns a descriptor builder for byte sequences.
//
// Typical reusable declaration:
//
//	payloadType := Bytes()
//	payloadType = payloadType.MinLen(1)
//	payloadType = payloadType.MaxLen(4096)
func Bytes() BytesType {
	return BytesType{header: newHeader(TypeBytes)}
}

// Nullable returns a bytes descriptor that admits null values.
func (t BytesType) Nullable() BytesType {
	t.header = t.header.withNullable()
	return t
}

// MinLen sets the inclusive minimum byte length.
func (t BytesType) MinLen(n int) BytesType {
	t.payload.minLen = limit[int]{value: n, set: true}
	return t
}

// MaxLen sets the inclusive maximum byte length.
func (t BytesType) MaxLen(n int) BytesType {
	t.payload.maxLen = limit[int]{value: n, set: true}
	return t
}

// Type returns a detached Type descriptor.
func (t BytesType) Type() Type {
	out := typeFromHeader(t.header)
	out.bytes = cloneBytesPayload(t.payload)
	return out
}

// typeExpr marks BytesType as a sealed TypeExpr implementation.
func (t BytesType) typeExpr() {}
