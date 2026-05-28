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

func TestInt8FieldWrapper(t *testing.T) {
	field := Field("value").Int8().Required().Nullable().Range(-1, 1).Enum(-1, 0, 1).Description("value").Field()

	requireEqual(t, field.Name(), FieldName("value"))
	requireEqual(t, field.IsRequired(), true)
	requireEqual(t, field.Description(), "value")
	requireEqual(t, field.Type().Code(), TypeInt8)
	requireEqual(t, field.Type().Nullable(), true)
	requireNoError(t, ValidateType(objectTypeForField(field), nil))
}

func TestInt8FieldExprMarker(t *testing.T) {
	Field("value").Int8().fieldExpr()
}

func TestInt8FieldOptionalPath(t *testing.T) {
	field := Field("value").Int8().Optional().Min(-1).Max(1).Field()

	requireEqual(t, field.IsOptional(), true)
}
