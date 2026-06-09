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

// Float64Field builds object fields whose value descriptor is float64.
//
// The wrapper keeps object-field metadata beside the float64 builder,
// allowing fluent field declarations without making fields reusable unnamed
// types.
type Float64Field struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// descriptor stores the typed descriptor builder for this field value.
	descriptor Float64Descriptor
}

// Required marks the field key as required.
func (f Float64Field) Required() Float64Field {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f Float64Field) Optional() Float64Field {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f Float64Field) Description(text string) Float64Field {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to float64 values.
func (f Float64Field) Nullable() Float64Field {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// Min sets the inclusive float64 lower bound.
func (f Float64Field) Min(n float64) Float64Field {
	f.descriptor = f.descriptor.Min(n)

	return f
}

// Max sets the inclusive float64 upper bound.
func (f Float64Field) Max(n float64) Float64Field {
	f.descriptor = f.descriptor.Max(n)

	return f
}

// Range sets the inclusive float64 lower and upper bounds.
func (f Float64Field) Range(min, max float64) Float64Field {
	f.descriptor = f.descriptor.Range(min, max)

	return f
}

// Enum stores accepted float64 literals for the field.
func (f Float64Field) Enum(values ...float64) Float64Field {
	f.descriptor = f.descriptor.Enum(values...)

	return f
}

// Field returns a detached finalized field descriptor.
func (f Float64Field) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks Float64Field as a sealed FieldExpr implementation.
func (f Float64Field) fieldExpr() {}
