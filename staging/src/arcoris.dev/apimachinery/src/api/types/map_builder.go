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

// MapOf starts a field descriptor whose value type is a string-keyed map.
//
// Only string keys are part of this first descriptor design. Value type
// validation is deferred to ValidateType so map.value diagnostics include the
// surrounding descriptor path.
//
// Field builder flow:
//
//	Field("labels").MapOf(
//		String().
//			MinLen(1),
//	).
//		Optional().
//		Nullable().
//		MinLen(1).
//		MaxLen(64).
//		Description("User-provided labels.")
func (b FieldBuilder) MapOf(value TypeExpr) MapField {
	return MapField{field: b.state(), typ: MapOf(value)}
}
