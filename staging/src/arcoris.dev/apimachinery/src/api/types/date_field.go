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

// DateField builds object fields whose value type is a calendar date.
//
// The wrapper keeps object-field metadata beside the date builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type DateField struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ DateType
}

// Required marks the field key as required.
func (f DateField) Required() DateField {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f DateField) Optional() DateField {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f DateField) Description(text string) DateField {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to date values.
func (f DateField) Nullable() DateField {
	f.typ = f.typ.Nullable()

	return f
}

// Field returns a detached finalized field descriptor.
func (f DateField) Field() FieldDescriptor {
	return f.field.fieldWithType(f.typ.Type())
}

// fieldExpr marks DateField as a sealed FieldExpr implementation.
func (f DateField) fieldExpr() {}
