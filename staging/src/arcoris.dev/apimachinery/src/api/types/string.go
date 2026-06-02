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

// StringType builds UTF-8 text descriptors with portable string constraints.
//
// StringType records portable text constraints such as length, pattern text,
// and enum literals. Patterns are stored as strings so future exporters and
// codecs can choose their own compiled representation.
type StringType struct {
	// header stores the descriptor kind and descriptor-wide flags under construction.
	header typeHeader
	// payload stores the exact string constraints under construction.
	payload stringPayload
}

// String returns a descriptor builder for UTF-8 text values.
//
// Typical reusable declaration:
//
//	nameType := String()
//	nameType = nameType.MinLen(1)
//	nameType = nameType.MaxLen(253)
//	nameType = nameType.Pattern(namePattern)
func String() StringType {
	return StringType{header: newHeader(TypeString)}
}

// Nullable returns a string descriptor that admits null values.
func (t StringType) Nullable() StringType {
	t.header = t.header.withNullable()

	return t
}

// MinLen sets the inclusive minimum string length.
func (t StringType) MinLen(n int) StringType {
	t.payload.minLen = limit[int]{value: n, set: true}

	return t
}

// MaxLen sets the inclusive maximum string length.
func (t StringType) MaxLen(n int) StringType {
	t.payload.maxLen = limit[int]{value: n, set: true}

	return t
}

// Pattern stores a portable textual regular expression for string values.
func (t StringType) Pattern(pattern string) StringType {
	t.payload.pattern = pattern
	t.payload.hasPattern = true

	return t
}

// Enum stores accepted string literals in declaration order.
func (t StringType) Enum(values ...string) StringType {
	t.payload.enum = cloneSlice(values)

	return t
}

// Type returns a detached Type descriptor.
func (t StringType) Type() Type {
	out := typeFromHeader(t.header)
	out.string = cloneStringPayload(t.payload)

	return out
}

// typeExpr marks StringType as a sealed TypeExpr implementation.
func (t StringType) typeExpr() {}
