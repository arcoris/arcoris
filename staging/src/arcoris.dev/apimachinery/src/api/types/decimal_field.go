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

// DecimalField builds object fields whose value descriptor is exact decimal.
//
// The wrapper keeps object-field metadata beside the decimal builder,
// allowing fluent field declarations without making fields reusable unnamed
// types.
type DecimalField struct {
	// field stores the name, presence, and description shared by field wrappers.
	field fieldState
	// descriptor stores the decimal descriptor builder for this field value.
	descriptor DecimalDescriptor
}

// Required marks the field key as required.
func (f DecimalField) Required() DecimalField {
	f.field = f.field.withRequired()

	return f
}

// Optional marks the field key as optional.
func (f DecimalField) Optional() DecimalField {
	f.field = f.field.withOptional()

	return f
}

// Description attaches human-facing descriptor text.
func (f DecimalField) Description(text string) DecimalField {
	f.field = f.field.withDescription(text)

	return f
}

// Nullable admits null in addition to decimal values.
func (f DecimalField) Nullable() DecimalField {
	f.descriptor = f.descriptor.Nullable()

	return f
}

// Precision sets the maximum number of significant decimal digits.
func (f DecimalField) Precision(n int) DecimalField {
	f.descriptor = f.descriptor.Precision(n)

	return f
}

// Scale sets the number of fractional decimal digits.
func (f DecimalField) Scale(n int) DecimalField {
	f.descriptor = f.descriptor.Scale(n)

	return f
}

// Field returns a detached finalized field descriptor.
func (f DecimalField) Field() FieldDescriptor {
	return f.field.fieldWithType(f.descriptor.Descriptor())
}

// fieldExpr marks DecimalField as a sealed FieldExpr implementation.
func (f DecimalField) fieldExpr() {}
