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

// Int8Field builds object fields whose value type is int8.
//
// The wrapper keeps object-field metadata beside the int8 builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type Int8Field struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ Int8Type
}

// Required marks the field key as required.
func (f Int8Field) Required() Int8Field { f.field = f.field.withRequired(); return f }

// Optional marks the field key as optional.
func (f Int8Field) Optional() Int8Field { f.field = f.field.withOptional(); return f }

// Description attaches human-facing descriptor text.
func (f Int8Field) Description(text string) Int8Field {
	f.field = f.field.withDescription(text)
	return f
}

// Nullable admits null in addition to int8 values.
func (f Int8Field) Nullable() Int8Field { f.typ = f.typ.Nullable(); return f }

// Min sets the inclusive int8 lower bound.
func (f Int8Field) Min(n int8) Int8Field { f.typ = f.typ.Min(n); return f }

// Max sets the inclusive int8 upper bound.
func (f Int8Field) Max(n int8) Int8Field { f.typ = f.typ.Max(n); return f }

// Range sets the inclusive int8 lower and upper bounds.
func (f Int8Field) Range(min, max int8) Int8Field { f.typ = f.typ.Range(min, max); return f }

// Enum stores accepted int8 literals for the field.
func (f Int8Field) Enum(values ...int8) Int8Field { f.typ = f.typ.Enum(values...); return f }

// Field returns a detached finalized field descriptor.
func (f Int8Field) Field() FieldDescriptor { return f.field.fieldWithType(f.typ.Type()) }

// fieldExpr marks Int8Field as a sealed FieldExpr implementation.
func (f Int8Field) fieldExpr() {}
