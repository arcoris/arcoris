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

// Float64 starts a field descriptor whose value type is a fixed-width float64.
//
// The descriptor records structural bounds and enum literals. ValidateType
// rejects NaN and infinities because descriptors must remain portable across
// codecs and schema exporters.
//
// Field builder flow:
//
//	Field("weight").Float64().
//		Optional().
//		Nullable().
//		Min(0).
//		Max(1000).
//		Enum(1, 10, 100).
//		Description("Relative scheduling weight.")
func (b FieldBuilder) Float64() Float64Field {
	return Float64Field{field: b.state(), typ: Float64()}
}
