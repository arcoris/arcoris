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

// Int64Field builds object fields whose value descriptor is int64.
//
// The wrapper keeps object-field metadata beside the int64 builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type Int64Field struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// descriptor stores the typed descriptor builder for this field value.
	descriptor Int64Descriptor
}

// Required marks the field key as required.
func (f Int64Field) Required() Int64Field {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f Int64Field) Optional() Int64Field {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f Int64Field) Description(text string) Int64Field {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to int64 values.
func (f Int64Field) Nullable() Int64Field {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// Min sets the inclusive int64 lower bound.
func (f Int64Field) Min(n int64) Int64Field {
	f.descriptor = f.descriptor.Min(n)

	return f
}

// Max sets the inclusive int64 upper bound.
func (f Int64Field) Max(n int64) Int64Field {
	f.descriptor = f.descriptor.Max(n)

	return f
}

// Range sets the inclusive int64 lower and upper bounds.
func (f Int64Field) Range(min, max int64) Int64Field {
	f.descriptor = f.descriptor.Range(min, max)

	return f
}

// Enum stores accepted int64 literals for the field.
func (f Int64Field) Enum(values ...int64) Int64Field {
	f.descriptor = f.descriptor.Enum(values...)

	return f
}

// Field returns a detached finalized field descriptor.
func (f Int64Field) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks Int64Field as a sealed FieldExpr implementation.
func (f Int64Field) fieldExpr() {}
