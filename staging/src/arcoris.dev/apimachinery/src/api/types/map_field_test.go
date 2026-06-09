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

func TestMapFieldWrapper(t *testing.T) {
	field := Field("value").MapOf(String()).Required().Nullable().MinEntries(1).MaxEntries(3).Description("value").Field()

	requireEqual(t, field.Descriptor().Code(), DescriptorMap)
	requireEqual(t, field.Descriptor().Nullable(), true)
	requireEqual(t, field.Description(), "value")
	requireNoError(t, ValidateLocal(objectTypeForField(field)))
}

func TestMapFieldExprMarker(t *testing.T) {
	Field("value").MapOf(String()).fieldExpr()
}

func TestMapFieldOptionalPath(t *testing.T) {
	field := Field("value").MapOf(String()).Optional().Field()

	requireEqual(t, field.IsOptional(), true)
}
