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

// Float32Field builds object fields whose value descriptor is float32.
//
// The wrapper keeps object-field metadata beside the float32 builder,
// allowing fluent field declarations without making fields reusable unnamed
// types.
type Float32Field struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// descriptor stores the typed descriptor builder for this field value.
	descriptor Float32Descriptor
}

// Required marks the field key as required.
func (f Float32Field) Required() Float32Field {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f Float32Field) Optional() Float32Field {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f Float32Field) Description(text string) Float32Field {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to float32 values.
func (f Float32Field) Nullable() Float32Field {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// Min sets the inclusive float32 lower bound.
func (f Float32Field) Min(n float32) Float32Field {
	f.descriptor = f.descriptor.Min(n)

	return f
}

// Max sets the inclusive float32 upper bound.
func (f Float32Field) Max(n float32) Float32Field {
	f.descriptor = f.descriptor.Max(n)

	return f
}

// Range sets the inclusive float32 lower and upper bounds.
func (f Float32Field) Range(min, max float32) Float32Field {
	f.descriptor = f.descriptor.Range(min, max)

	return f
}

// Enum stores accepted float32 literals for the field.
func (f Float32Field) Enum(values ...float32) Float32Field {
	f.descriptor = f.descriptor.Enum(values...)

	return f
}

// Field returns a detached finalized field descriptor.
func (f Float32Field) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks Float32Field as a sealed FieldExpr implementation.
func (f Float32Field) fieldExpr() {}
