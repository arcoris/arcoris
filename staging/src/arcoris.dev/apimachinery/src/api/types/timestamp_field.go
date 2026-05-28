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

// TimestampField builds object fields whose value type is timestamp.
//
// The wrapper keeps object-field metadata beside the timestamp builder,
// allowing fluent field declarations without making fields reusable unnamed
// types.
type TimestampField struct {
	// field stores name, presence, and description shared by all field families.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ TimestampType
}

// Required marks the field key as required.
func (f TimestampField) Required() TimestampField { f.field = f.field.withRequired(); return f }

// Optional marks the field key as optional.
func (f TimestampField) Optional() TimestampField { f.field = f.field.withOptional(); return f }

// Description attaches human-facing descriptor text.
func (f TimestampField) Description(text string) TimestampField {
	f.field = f.field.withDescription(text)
	return f
}

// Nullable admits null in addition to timestamp values.
func (f TimestampField) Nullable() TimestampField { f.typ = f.typ.Nullable(); return f }

// Field returns a detached finalized field descriptor.
func (f TimestampField) Field() FieldDescriptor { return f.field.fieldWithType(f.typ.Type()) }

// fieldExpr marks TimestampField as a sealed FieldExpr implementation.
func (f TimestampField) fieldExpr() {}
