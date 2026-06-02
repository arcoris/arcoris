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

// StringField builds object fields whose value type is string.
//
// The wrapper keeps object-field metadata beside the string builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type StringField struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ StringType
}

// Required marks the field key as required.
func (f StringField) Required() StringField {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f StringField) Optional() StringField {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f StringField) Description(text string) StringField {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to string values.
func (f StringField) Nullable() StringField {
	f.typ = f.typ.Nullable()

	return f
}

// MinLen sets the inclusive minimum string length.
func (f StringField) MinLen(n int) StringField {
	f.typ = f.typ.MinLen(n)

	return f
}

// MaxLen sets the inclusive maximum string length.
func (f StringField) MaxLen(n int) StringField {
	f.typ = f.typ.MaxLen(n)

	return f
}

// Pattern stores a textual regular expression for the string field.
func (f StringField) Pattern(pattern string) StringField {
	f.typ = f.typ.Pattern(pattern)

	return f
}

// Enum stores accepted string literals for the field.
func (f StringField) Enum(values ...string) StringField {
	f.typ = f.typ.Enum(values...)

	return f
}

// Field returns a detached finalized field descriptor.
func (f StringField) Field() FieldDescriptor {
	return f.field.fieldWithType(f.typ.Type())
}

// fieldExpr marks StringField as a sealed FieldExpr implementation.
func (f StringField) fieldExpr() {}
