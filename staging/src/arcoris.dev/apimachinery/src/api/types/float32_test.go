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

func TestFloat32TypeDescriptor(t *testing.T) {
	desc := Float32().Range(0.5, 2.5).Enum(1.5, 2.5).Nullable().Descriptor()
	requireExactNumericDescriptor(t, exactNumericDescriptorCase[float32]{
		descriptor:  desc,
		code:        DescriptorFloat32,
		min:         func(desc Descriptor) (float32, bool) { return requireFloat32View(t, desc).Min() },
		max:         func(desc Descriptor) (float32, bool) { return requireFloat32View(t, desc).Max() },
		enum:        func(desc Descriptor) []float32 { return requireFloat32View(t, desc).Enum() },
		wrong:       func(desc Descriptor) bool { _, ok := desc.AsFloat64(); return ok },
		wantMin:     0.5,
		wantMax:     2.5,
		wantFirst:   1.5,
		replaceWith: 9,
	})
}

func TestFloat32DescriptorExprMarker(t *testing.T) {
	Float32().descriptorExpr()
}
