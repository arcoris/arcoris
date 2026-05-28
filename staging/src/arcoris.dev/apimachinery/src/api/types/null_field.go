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

// NullField builds object fields whose value type is the null literal.
//
// The wrapper carries field metadata separately from the null type descriptor
// so presence, description, and value shape remain distinct.
type NullField struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// typ stores the typed descriptor builder for this field value.
	typ NullType
}

// Required marks the field key as required.
func (f NullField) Required() NullField { f.field = f.field.withRequired(); return f }

// Optional marks the field key as optional.
func (f NullField) Optional() NullField { f.field = f.field.withOptional(); return f }

// Description attaches human-facing descriptor text.
func (f NullField) Description(text string) NullField {
	f.field = f.field.withDescription(text)
	return f
}

// Field returns a detached finalized field descriptor.
func (f NullField) Field() FieldDescriptor { return f.field.fieldWithType(f.typ.Type()) }

// fieldExpr marks NullField as a sealed FieldExpr implementation.
func (f NullField) fieldExpr() {}
