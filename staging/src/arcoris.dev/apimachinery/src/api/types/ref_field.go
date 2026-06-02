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

// RefField builds object fields whose value type is a named type reference.
type RefField struct {
	field fieldState
	typ   RefType
}

// Required marks the field key as required.
func (f RefField) Required() RefField {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f RefField) Optional() RefField {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f RefField) Description(text string) RefField {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to referenced values.
func (f RefField) Nullable() RefField {
	f.typ = f.typ.Nullable()

	return f
}

// Field returns a detached finalized field descriptor.
func (f RefField) Field() FieldDescriptor {
	return f.field.fieldWithType(f.typ.Type())
}

// fieldExpr marks RefField as a sealed FieldExpr implementation.
func (f RefField) fieldExpr() {}
