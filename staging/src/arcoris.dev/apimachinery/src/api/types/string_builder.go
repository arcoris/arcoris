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

// String starts a field descriptor whose value descriptor is UTF-8 text.
//
// String fields can record portable length, pattern, and enum constraints.
// Pattern text is kept as descriptor data rather than a compiled regexp so
// future codecs and schema exporters can choose their own representation.
//
// Field builder flow:
//
//	Field("image").String().
//		Required().
//		Nullable().
//		MinBytes(1).
//		MaxBytes(253).
//		Pattern("^[a-z][a-z0-9-]*$").
//		Enum("app", "sidecar").
//		Description("Container image reference.")
func (b FieldBuilder) String() StringField {
	return StringField{field: b.state(), descriptor: String()}
}
