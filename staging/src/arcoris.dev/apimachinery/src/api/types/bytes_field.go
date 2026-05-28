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

// BytesField builds object fields whose value type is bytes.
//
// The wrapper keeps object-field metadata beside the bytes builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type BytesField struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ BytesType
}

// Required marks the field key as required.
func (f BytesField) Required() BytesField { f.field = f.field.withRequired(); return f }

// Optional marks the field key as optional.
func (f BytesField) Optional() BytesField { f.field = f.field.withOptional(); return f }

// Description attaches human-facing descriptor text.
func (f BytesField) Description(text string) BytesField {
	f.field = f.field.withDescription(text)
	return f
}

// Nullable admits null in addition to byte-sequence values.
func (f BytesField) Nullable() BytesField { f.typ = f.typ.Nullable(); return f }

// MinLen sets the inclusive minimum byte length.
func (f BytesField) MinLen(n int) BytesField { f.typ = f.typ.MinLen(n); return f }

// MaxLen sets the inclusive maximum byte length.
func (f BytesField) MaxLen(n int) BytesField { f.typ = f.typ.MaxLen(n); return f }

// Field returns a detached finalized field descriptor.
func (f BytesField) Field() FieldDescriptor { return f.field.fieldWithType(f.typ.Type()) }

// fieldExpr marks BytesField as a sealed FieldExpr implementation.
func (f BytesField) fieldExpr() {}
