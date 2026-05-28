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

// Uint32Field builds object fields whose value type is uint32.
//
// The wrapper keeps object-field metadata beside the uint32 builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type Uint32Field struct {
	// field stores name, presence, and description shared by all field families.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ Uint32Type
}

// Required marks the field key as required.
func (f Uint32Field) Required() Uint32Field { f.field = f.field.withRequired(); return f }

// Optional marks the field key as optional.
func (f Uint32Field) Optional() Uint32Field { f.field = f.field.withOptional(); return f }

// Description attaches human-facing descriptor text.
func (f Uint32Field) Description(text string) Uint32Field {
	f.field = f.field.withDescription(text)
	return f
}

// Nullable admits null in addition to uint32 values.
func (f Uint32Field) Nullable() Uint32Field { f.typ = f.typ.Nullable(); return f }

// Min sets the inclusive uint32 lower bound.
func (f Uint32Field) Min(n uint32) Uint32Field { f.typ = f.typ.Min(n); return f }

// Max sets the inclusive uint32 upper bound.
func (f Uint32Field) Max(n uint32) Uint32Field { f.typ = f.typ.Max(n); return f }

// Range sets the inclusive uint32 lower and upper bounds.
func (f Uint32Field) Range(min, max uint32) Uint32Field { f.typ = f.typ.Range(min, max); return f }

// Enum stores accepted uint32 literals for the field.
func (f Uint32Field) Enum(values ...uint32) Uint32Field { f.typ = f.typ.Enum(values...); return f }

// Field returns a detached finalized field descriptor.
func (f Uint32Field) Field() FieldDescriptor { return f.field.fieldWithType(f.typ.Type()) }

// fieldExpr marks Uint32Field as a sealed FieldExpr implementation.
func (f Uint32Field) fieldExpr() {}
