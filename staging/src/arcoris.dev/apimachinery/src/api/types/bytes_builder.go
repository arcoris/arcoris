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

// Bytes starts a field descriptor whose value descriptor is an opaque byte sequence.
//
// Bytes fields support structural length constraints only. Encoding choices
// such as base64, hex, or binary framing belong to codecs and exporters, not to
// this descriptor package.
//
// Field builder flow:
//
//	Field("payload").Bytes().
//		Optional().
//		Nullable().
//		MinBytes(1).
//		MaxBytes(4096).
//		Description("Opaque payload bytes.")
func (b FieldBuilder) Bytes() BytesField {
	return BytesField{field: b.state(), descriptor: Bytes()}
}
