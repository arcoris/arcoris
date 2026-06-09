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

// Int32Field builds object fields whose value descriptor is int32.
//
// The wrapper keeps object-field metadata beside the int32 builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type Int32Field struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// descriptor stores the typed descriptor builder for this field value.
	descriptor Int32Descriptor
}

// Required marks the field key as required.
func (f Int32Field) Required() Int32Field {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f Int32Field) Optional() Int32Field {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f Int32Field) Description(text string) Int32Field {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to int32 values.
func (f Int32Field) Nullable() Int32Field {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// Min sets the inclusive int32 lower bound.
func (f Int32Field) Min(n int32) Int32Field {
	f.descriptor = f.descriptor.Min(n)

	return f
}

// Max sets the inclusive int32 upper bound.
func (f Int32Field) Max(n int32) Int32Field {
	f.descriptor = f.descriptor.Max(n)

	return f
}

// Range sets the inclusive int32 lower and upper bounds.
func (f Int32Field) Range(min, max int32) Int32Field {
	f.descriptor = f.descriptor.Range(min, max)

	return f
}

// Enum stores accepted int32 literals for the field.
func (f Int32Field) Enum(values ...int32) Int32Field {
	f.descriptor = f.descriptor.Enum(values...)

	return f
}

// Field returns a detached finalized field descriptor.
func (f Int32Field) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks Int32Field as a sealed FieldExpr implementation.
func (f Int32Field) fieldExpr() {}
