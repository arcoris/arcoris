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

// MapOf starts a field descriptor whose value descriptor is a string-keyed map.
//
// Concrete keys are strings. Keys can add structural constraints to those
// string tokens, while value descriptor validation is deferred to ValidateResolved so
// map.value diagnostics include the surrounding descriptor path.
//
// Field builder flow:
//
//	Field("labels").MapOf(
//		String().
//			MinBytes(1),
//	).Optional().
//		Nullable().
//		MinEntries(1).
//		MaxEntries(64).
//		Description("User-provided labels.")
func (b FieldBuilder) MapOf(value DescriptorExpr) MapField {
	return MapField{field: b.state(), descriptor: MapOf(value)}
}
