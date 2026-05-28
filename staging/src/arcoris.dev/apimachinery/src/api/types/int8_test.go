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

func TestInt8TypeDescriptor(t *testing.T) {
	typ := Int8().Range(-1, 1).Enum(-1, 0, 1).Nullable().Type()

	requireEqual(t, typ.Code(), TypeInt8)
	requireEqual(t, typ.Nullable(), true)
	view, ok := typ.Int8()
	requireEqual(t, ok, true)
	min, ok := view.Min()
	requireEqual(t, ok, true)
	requireEqual(t, min, int8(-1))
	max, ok := view.Max()
	requireEqual(t, ok, true)
	requireEqual(t, max, int8(1))
	enum := view.Enum()
	enum[0] = 9
	requireEqual(t, view.Enum()[0], int8(-1))
	_, ok = typ.Int16()
	requireEqual(t, ok, false)
	requireNoError(t, ValidateType(typ, nil))
}

func TestInt8TypeExprMarker(t *testing.T) {
	Int8().typeExpr()
}
