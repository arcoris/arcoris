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

func TestUint32TypeDescriptor(t *testing.T) {
	typ := Uint32().Range(0, 2).Enum(0, 1, 2).Nullable().Type()

	requireEqual(t, typ.Code(), TypeUint32)
	requireEqual(t, typ.Nullable(), true)
	view, ok := typ.Uint32()
	requireEqual(t, ok, true)
	min, ok := view.Min()
	requireEqual(t, ok, true)
	requireEqual(t, min, uint32(0))
	max, ok := view.Max()
	requireEqual(t, ok, true)
	requireEqual(t, max, uint32(2))
	enum := view.Enum()
	enum[0] = 9
	requireEqual(t, view.Enum()[0], uint32(0))
	_, ok = typ.Uint8()
	requireEqual(t, ok, false)
	requireNoError(t, ValidateType(typ, nil))
}

func TestUint32TypeExprMarker(t *testing.T) {
	Uint32().typeExpr()
}
