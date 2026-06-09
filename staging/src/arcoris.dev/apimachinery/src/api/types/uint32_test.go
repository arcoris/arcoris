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
	desc := Uint32().Range(0, 2).Enum(0, 1, 2).Nullable().Descriptor()
	requireExactNumericDescriptor(t, exactNumericDescriptorCase[uint32]{
		descriptor:  desc,
		code:        DescriptorUint32,
		min:         func(desc Descriptor) (uint32, bool) { return requireUint32View(t, desc).Min() },
		max:         func(desc Descriptor) (uint32, bool) { return requireUint32View(t, desc).Max() },
		enum:        func(desc Descriptor) []uint32 { return requireUint32View(t, desc).Enum() },
		wrong:       func(desc Descriptor) bool { _, ok := desc.AsUint8(); return ok },
		wantMin:     0,
		wantMax:     2,
		wantFirst:   0,
		replaceWith: 9,
	})
}

func TestUint32DescriptorExprMarker(t *testing.T) {
	Uint32().descriptorExpr()
}
