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

// BoolField builds object fields whose value descriptor is boolean.
//
// The wrapper carries field metadata separately from the boolean descriptor
// descriptor so presence, description, and value shape remain distinct.
type BoolField struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// descriptor stores the typed descriptor builder for this field value.
	descriptor BoolDescriptor
}

// Required marks the field key as required.
func (f BoolField) Required() BoolField {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f BoolField) Optional() BoolField {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f BoolField) Description(text string) BoolField {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to boolean values.
func (f BoolField) Nullable() BoolField {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// Field returns a detached finalized field descriptor.
func (f BoolField) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks BoolField as a sealed FieldExpr implementation.
func (f BoolField) fieldExpr() {}
