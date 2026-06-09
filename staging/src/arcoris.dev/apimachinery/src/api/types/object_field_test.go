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

func TestObjectFieldWrapper(t *testing.T) {
	field := Field("value").
		Object(
			Field("name").String().Required(),
		).
		Required().
		Nullable().
		UnknownFields(UnknownPrune).
		Description("value").
		Field()

	requireEqual(t, field.Descriptor().Code(), DescriptorObject)
	requireEqual(t, field.Descriptor().Nullable(), true)
	requireEqual(t, field.Description(), "value")
	requireNoError(t, ValidateLocal(objectTypeForField(field)))
}

func TestObjectFieldExprMarker(t *testing.T) {
	Field("value").Object().fieldExpr()
}

func TestObjectFieldOptionalPath(t *testing.T) {
	field := Field("value").Object().Optional().Field()

	requireEqual(t, field.IsOptional(), true)
}
