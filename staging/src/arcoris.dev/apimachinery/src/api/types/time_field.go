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

// TimeField builds object fields whose value type is a time of day.
//
// The wrapper keeps object-field metadata beside the time builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type TimeField struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ TimeType
}

// Required marks the field key as required.
func (f TimeField) Required() TimeField { f.field = f.field.withRequired(); return f }

// Optional marks the field key as optional.
func (f TimeField) Optional() TimeField { f.field = f.field.withOptional(); return f }

// Description attaches human-facing descriptor text.
func (f TimeField) Description(text string) TimeField {
	f.field = f.field.withDescription(text)
	return f
}

// Nullable admits null in addition to time values.
func (f TimeField) Nullable() TimeField { f.typ = f.typ.Nullable(); return f }

// Field returns a detached finalized field descriptor.
func (f TimeField) Field() FieldDescriptor { return f.field.fieldWithType(f.typ.Type()) }

// fieldExpr marks TimeField as a sealed FieldExpr implementation.
func (f TimeField) fieldExpr() {}
