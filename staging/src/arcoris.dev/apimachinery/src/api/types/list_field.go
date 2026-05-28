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

// ListField builds object fields whose value type is a list.
type ListField struct {
	field fieldState
	typ   ListType
}

// Required marks the field key as required.
func (f ListField) Required() ListField { f.field = f.field.withRequired(); return f }

// Optional marks the field key as optional.
func (f ListField) Optional() ListField { f.field = f.field.withOptional(); return f }

// Description attaches human-facing descriptor text.
func (f ListField) Description(text string) ListField {
	f.field = f.field.withDescription(text)
	return f
}

// Nullable admits null in addition to list values.
func (f ListField) Nullable() ListField { f.typ = f.typ.Nullable(); return f }

// MinLen sets the inclusive minimum list length.
func (f ListField) MinLen(n int) ListField { f.typ = f.typ.MinLen(n); return f }

// MaxLen sets the inclusive maximum list length.
func (f ListField) MaxLen(n int) ListField { f.typ = f.typ.MaxLen(n); return f }

// Atomic records atomic list semantics.
func (f ListField) Atomic() ListField { f.typ = f.typ.Atomic(); return f }

// Set records set-like list semantics.
func (f ListField) Set() ListField { f.typ = f.typ.Set(); return f }

// Map records map-like list semantics keyed by object field names.
func (f ListField) Map(keys ...string) ListField { f.typ = f.typ.Map(keys...); return f }

// Field returns a detached finalized field descriptor.
func (f ListField) Field() FieldDescriptor { return f.field.fieldWithType(f.typ.Type()) }

// fieldExpr marks ListField as a sealed FieldExpr implementation.
func (f ListField) fieldExpr() {}
