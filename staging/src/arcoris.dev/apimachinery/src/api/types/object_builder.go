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

// Object starts a field descriptor whose value type is a structural object.
//
// The nested field order is preserved exactly as declared. Duplicate names,
// missing presence, invalid field names, and invalid nested types are reported
// by ValidateType rather than by this constructor.
//
// Field builder flow:
//
//	Field("spec").Object(
//		Field("image").String().
//			Required().
//			MinLen(1),
//	).
//		Required().
//		Nullable().
//		UnknownFields(UnknownReject).
//		Description("Desired object state.")
func (b FieldBuilder) Object(fields ...FieldExpr) ObjectField {
	return ObjectField{field: b.state(), typ: Object(fields...)}
}
