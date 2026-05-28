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

func TestInt16TypeDescriptor(t *testing.T) {
	typ := Int16().Range(-1, 1).Enum(-1, 0, 1).Nullable().Type()
	requireExactNumericDescriptor(t, exactNumericDescriptorCase[int16]{
		typ:       typ,
		code:      TypeInt16,
		min:       func(typ Type) (int16, bool) { return requireInt16View(t, typ).Min() },
		max:       func(typ Type) (int16, bool) { return requireInt16View(t, typ).Max() },
		enum:      func(typ Type) []int16 { return requireInt16View(t, typ).Enum() },
		wrong:     func(typ Type) bool { _, ok := typ.Int8(); return ok },
		wantMin:   -1,
		wantMax:   1,
		wantFirst: -1,
	})
}

func TestInt16TypeExprMarker(t *testing.T) {
	Int16().typeExpr()
}
