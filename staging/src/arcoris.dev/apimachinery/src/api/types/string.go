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

// StringDescriptor builds UTF-8 text descriptors with portable string constraints.
//
// StringDescriptor records portable text constraints such as length, pattern text,
// and enum literals. Patterns are stored as strings so future exporters and
// codecs can choose their own compiled representation.
type StringDescriptor struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header descriptorHeader
	// payload stores the exact string constraints under construction.
	payload stringPayload
}

// String returns a descriptor builder for UTF-8 text values.
//
// Typical reusable declaration:
//
//	nameType := String().
//		MinBytes(1).
//		MaxBytes(253).
//		Pattern(namePattern)
func String() StringDescriptor {
	return StringDescriptor{header: newHeader(DescriptorString)}
}

// Nullable returns a string descriptor that admits null values.
func (desc StringDescriptor) Nullable() StringDescriptor {
	desc.header = desc.header.withNullable()

	return desc
}

// MinBytes sets the inclusive minimum UTF-8 byte length.
func (desc StringDescriptor) MinBytes(n int) StringDescriptor {
	desc.payload.minBytes = limit[int]{value: n, set: true}

	return desc
}

// MaxBytes sets the inclusive maximum UTF-8 byte length.
func (desc StringDescriptor) MaxBytes(n int) StringDescriptor {
	desc.payload.maxBytes = limit[int]{value: n, set: true}

	return desc
}

// MinRunes sets the inclusive minimum Unicode code point count.
//
// The descriptor counts Go runes. It deliberately does not count grapheme
// clusters, because grapheme segmentation is locale-sensitive and belongs to a
// higher text-processing layer if ARCORIS ever needs it.
func (desc StringDescriptor) MinRunes(n int) StringDescriptor {
	desc.payload.minRunes = limit[int]{value: n, set: true}

	return desc
}

// MaxRunes sets the inclusive maximum Unicode code point count.
func (desc StringDescriptor) MaxRunes(n int) StringDescriptor {
	desc.payload.maxRunes = limit[int]{value: n, set: true}

	return desc
}

// Pattern stores a portable textual regular expression for string values.
func (desc StringDescriptor) Pattern(pattern string) StringDescriptor {
	desc.payload.pattern = pattern
	desc.payload.hasPattern = true

	return desc
}

// Enum stores accepted string literals in declaration order.
func (desc StringDescriptor) Enum(values ...string) StringDescriptor {
	desc.payload.enum = cloneSlice(values)

	return desc
}

// Descriptor returns a detached Descriptor descriptor.
func (desc StringDescriptor) Descriptor() Descriptor {
	out := descriptorFromHeader(desc.header)
	out.string = cloneStringPayload(desc.payload)

	return out
}

// descriptorExpr marks StringDescriptor as a sealed DescriptorExpr implementation.
func (desc StringDescriptor) descriptorExpr() {}
