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

// ObjectField builds object fields whose value descriptor is another object.
type ObjectField struct {
	field      fieldState
	descriptor ObjectDescriptor
}

// Required marks the field key as required.
func (f ObjectField) Required() ObjectField {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f ObjectField) Optional() ObjectField {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f ObjectField) Description(text string) ObjectField {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to object values.
func (f ObjectField) Nullable() ObjectField {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// UnknownFields records the structural policy for undeclared nested fields.
func (f ObjectField) UnknownFields(policy UnknownFieldPolicy) ObjectField {
	f.descriptor = f.descriptor.UnknownFields(policy)

	return f
}

// Field returns a detached finalized field descriptor.
func (f ObjectField) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks ObjectField as a sealed FieldExpr implementation.
func (f ObjectField) fieldExpr() {}
