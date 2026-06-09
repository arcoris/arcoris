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

// Null starts a field descriptor whose value descriptor is the null literal.
//
// Use Null for fields that are structurally constrained to the literal null
// value. Do not use it as a nullable marker for another family; nullable
// string, object, list, map, number, temporal, and reference fields use their
// own Nullable methods.
//
// Field builder flow:
//
//	Field("deletedAt").Null().
//		Optional().
//		Description("Explicit null marker.")
func (b FieldBuilder) Null() NullField {
	return NullField{field: b.state(), descriptor: Null()}
}
