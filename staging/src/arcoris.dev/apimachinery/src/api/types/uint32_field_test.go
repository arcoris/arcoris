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

import "testing"

func TestUint32FieldWrapper(t *testing.T) {
	field := Field("value").Uint32().Required().Nullable().Range(0, 2).Enum(0, 1, 2).Description("value").Field()

	requireEqual(t, field.Name(), FieldName("value"))
	requireEqual(t, field.IsRequired(), true)
	requireEqual(t, field.Description(), "value")
	requireEqual(t, field.Descriptor().Code(), DescriptorUint32)
	requireEqual(t, field.Descriptor().Nullable(), true)
	requireNoError(t, ValidateLocal(objectTypeForField(field)))
}

func TestUint32FieldExprMarker(t *testing.T) {
	Field("value").Uint32().fieldExpr()
}

func TestUint32FieldOptionalPath(t *testing.T) {
	field := Field("value").Uint32().Optional().Min(1).Max(2).Field()

	requireEqual(t, field.IsOptional(), true)
}
