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

// Float32 starts a field descriptor whose value descriptor is a fixed-width float32.
//
// The descriptor records structural bounds and enum literals. ValidateResolved
// rejects NaN and infinities because descriptors must remain portable across
// codecs and schema exporters.
//
// Field builder flow:
//
//	Field("ratio").Float32().
//		Optional().
//		Nullable().
//		Range(0, 1).
//		Enum(0, 0.5, 1).
//		Description("Ratio from zero to one.")
func (b FieldBuilder) Float32() Float32Field {
	return Float32Field{field: b.state(), descriptor: Float32()}
}
