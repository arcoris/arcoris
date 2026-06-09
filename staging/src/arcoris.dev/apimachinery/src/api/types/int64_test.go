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

func TestInt64TypeDescriptor(t *testing.T) {
	desc := Int64().Range(-1, 1).Enum(-1, 0, 1).Nullable().Descriptor()
	requireExactNumericDescriptor(t, exactNumericDescriptorCase[int64]{
		descriptor: desc,
		code:       DescriptorInt64,
		min:        func(desc Descriptor) (int64, bool) { return requireInt64View(t, desc).Min() },
		max:        func(desc Descriptor) (int64, bool) { return requireInt64View(t, desc).Max() },
		enum:       func(desc Descriptor) []int64 { return requireInt64View(t, desc).Enum() },
		wrong:      func(desc Descriptor) bool { _, ok := desc.AsInt8(); return ok },
		wantMin:    -1,
		wantMax:    1,
		wantFirst:  -1,
	})
}

func TestInt64DescriptorExprMarker(t *testing.T) {
	Int64().descriptorExpr()
}
