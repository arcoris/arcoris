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

// ListField builds object fields whose value descriptor is a list.
type ListField struct {
	field      fieldState
	descriptor ListDescriptor
}

// Required marks the field key as required.
func (f ListField) Required() ListField {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f ListField) Optional() ListField {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f ListField) Description(text string) ListField {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to list values.
func (f ListField) Nullable() ListField {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// MinItems sets the inclusive minimum list item count.
func (f ListField) MinItems(n int) ListField {
	f.descriptor = f.descriptor.MinItems(n)

	return f
}

// MaxItems sets the inclusive maximum list item count.
func (f ListField) MaxItems(n int) ListField {
	f.descriptor = f.descriptor.MaxItems(n)

	return f
}

// Atomic records atomic list semantics.
func (f ListField) Atomic() ListField {
	f.descriptor = f.descriptor.Atomic()

	return f
}

// Ordered records index-addressable list semantics.
func (f ListField) Ordered() ListField {
	f.descriptor = f.descriptor.Ordered()

	return f
}

// Set records set-like list semantics.
func (f ListField) Set() ListField {
	f.descriptor = f.descriptor.Set()

	return f
}

// Map records map-like list semantics keyed by object field names.
func (f ListField) Map(keys ...string) ListField {
	f.descriptor = f.descriptor.Map(keys...)

	return f
}

// Field returns a detached finalized field descriptor.
func (f ListField) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks ListField as a sealed FieldExpr implementation.
func (f ListField) fieldExpr() {}
