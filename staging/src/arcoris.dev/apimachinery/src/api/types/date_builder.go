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

// Date starts a field descriptor for a calendar date.
//
// The descriptor records the temporal kind only. Calendar text formats,
// timezone interpretation, and codec-specific wire forms are outside this
// package.
//
// Field builder flow:
//
//	Field("effectiveDate").Date().
//		Optional().
//		Nullable().
//		Description("Effective calendar date.")
func (b FieldBuilder) Date() DateField {
	return DateField{field: b.state(), descriptor: Date()}
}
