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

// Uint8Field builds object fields whose value descriptor is uint8.
//
// The wrapper keeps object-field metadata beside the uint8 builder, allowing
// fluent field declarations without making fields reusable unnamed types.
type Uint8Field struct {
	// field stores name, presence, and description shared by all field wrappers.
	field fieldState
	// descriptor stores the typed descriptor builder for this field value.
	descriptor Uint8Descriptor
}

// Required marks the field key as required.
func (f Uint8Field) Required() Uint8Field {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f Uint8Field) Optional() Uint8Field {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f Uint8Field) Description(text string) Uint8Field {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to uint8 values.
func (f Uint8Field) Nullable() Uint8Field {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// Min sets the inclusive uint8 lower bound.
func (f Uint8Field) Min(n uint8) Uint8Field {
	f.descriptor = f.descriptor.Min(n)

	return f
}

// Max sets the inclusive uint8 upper bound.
func (f Uint8Field) Max(n uint8) Uint8Field {
	f.descriptor = f.descriptor.Max(n)

	return f
}

// Range sets the inclusive uint8 lower and upper bounds.
func (f Uint8Field) Range(min, max uint8) Uint8Field {
	f.descriptor = f.descriptor.Range(min, max)

	return f
}

// Enum stores accepted uint8 literals for the field.
func (f Uint8Field) Enum(values ...uint8) Uint8Field {
	f.descriptor = f.descriptor.Enum(values...)

	return f
}

// Field returns a detached finalized field descriptor.
func (f Uint8Field) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks Uint8Field as a sealed FieldExpr implementation.
func (f Uint8Field) fieldExpr() {}
