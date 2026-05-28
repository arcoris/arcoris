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

// Int16Field builds object fields whose value type is int16.
//
// The wrapper keeps object-field metadata beside the int16 builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type Int16Field struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ Int16Type
}

// Required marks the field key as required.
func (f Int16Field) Required() Int16Field { f.field = f.field.withRequired(); return f }

// Optional marks the field key as optional.
func (f Int16Field) Optional() Int16Field { f.field = f.field.withOptional(); return f }

// Description attaches human-facing descriptor text.
func (f Int16Field) Description(text string) Int16Field {
	f.field = f.field.withDescription(text)
	return f
}

// Nullable admits null in addition to int16 values.
func (f Int16Field) Nullable() Int16Field { f.typ = f.typ.Nullable(); return f }

// Min sets the inclusive int16 lower bound.
func (f Int16Field) Min(n int16) Int16Field { f.typ = f.typ.Min(n); return f }

// Max sets the inclusive int16 upper bound.
func (f Int16Field) Max(n int16) Int16Field { f.typ = f.typ.Max(n); return f }

// Range sets the inclusive int16 lower and upper bounds.
func (f Int16Field) Range(min, max int16) Int16Field { f.typ = f.typ.Range(min, max); return f }

// Enum stores accepted int16 literals for the field.
func (f Int16Field) Enum(values ...int16) Int16Field { f.typ = f.typ.Enum(values...); return f }

// Field returns a detached finalized field descriptor.
func (f Int16Field) Field() FieldDescriptor { return f.field.fieldWithType(f.typ.Type()) }

// fieldExpr marks Int16Field as a sealed FieldExpr implementation.
func (f Int16Field) fieldExpr() {}
