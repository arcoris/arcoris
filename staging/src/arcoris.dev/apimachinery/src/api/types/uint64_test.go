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

func TestUint64TypeDescriptor(t *testing.T) {
	typ := Uint64().Range(0, 2).Enum(0, 1, 2).Nullable().Type()
	requireExactNumericDescriptor(t, exactNumericDescriptorCase[uint64]{
		typ:         typ,
		code:        TypeUint64,
		min:         func(typ Type) (uint64, bool) { return requireUint64View(t, typ).Min() },
		max:         func(typ Type) (uint64, bool) { return requireUint64View(t, typ).Max() },
		enum:        func(typ Type) []uint64 { return requireUint64View(t, typ).Enum() },
		wrong:       func(typ Type) bool { _, ok := typ.Uint32(); return ok },
		wantMin:     0,
		wantMax:     2,
		wantFirst:   0,
		replaceWith: 9,
	})
}

func TestUint64TypeExprMarker(t *testing.T) {
	Uint64().typeExpr()
}
