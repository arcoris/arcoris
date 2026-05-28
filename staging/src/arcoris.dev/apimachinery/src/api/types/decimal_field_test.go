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

func TestDecimalFieldWrapper(t *testing.T) {
	field := Field("value").Decimal().Required().Nullable().Precision(10).Scale(2).Description("value").Field()

	requireEqual(t, field.Name(), FieldName("value"))
	requireEqual(t, field.IsRequired(), true)
	requireEqual(t, field.Description(), "value")
	requireEqual(t, field.Type().Code(), TypeDecimal)
	requireEqual(t, field.Type().Nullable(), true)
	requireNoError(t, ValidateType(objectTypeForField(field), nil))
}

func TestDecimalFieldExprMarker(t *testing.T) {
	Field("value").Decimal().fieldExpr()
}

func TestDecimalFieldOptionalPath(t *testing.T) {
	field := Field("value").Decimal().Optional().Precision(2).Scale(1).Field()

	requireEqual(t, field.IsOptional(), true)
}
