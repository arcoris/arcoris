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

func TestListFieldWrapper(t *testing.T) {
	field := Field("value").ListOf(Object(Field("type").String().Required())).Required().Nullable().MinLen(1).MaxLen(3).Map("type").Description("value").Field()

	requireEqual(t, field.Type().Code(), TypeList)
	requireEqual(t, field.Type().Nullable(), true)
	requireEqual(t, field.Description(), "value")
	requireNoError(t, ValidateType(objectTypeForField(field), nil))
}

func TestListFieldExprMarker(t *testing.T) {
	Field("value").ListOf(String()).fieldExpr()
}

func TestListFieldOptionalPath(t *testing.T) {
	field := Field("value").ListOf(String()).Optional().Atomic().Set().Field()

	requireEqual(t, field.IsOptional(), true)
}
