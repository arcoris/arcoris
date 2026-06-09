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

// Bool starts a field descriptor whose value descriptor is boolean.
//
// The returned wrapper combines the captured field name with boolean descriptor
// metadata. Presence remains unset until Required or Optional is called, so
// ValidateResolved can reject incomplete field declarations with a precise path.
//
// Field builder flow:
//
//	Field("enabled").Bool().
//		Required().
//		Nullable().
//		Description("Whether the feature is enabled.")
func (b FieldBuilder) Bool() BoolField {
	return BoolField{field: b.state(), descriptor: Bool()}
}
