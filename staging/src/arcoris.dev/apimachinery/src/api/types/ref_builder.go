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

// Ref starts a field descriptor whose value descriptor is a named structural reference.
//
// The name is recorded as descriptor data. ValidateResolved checks name syntax and,
// when supplied a Resolver, verifies that the reference resolves and does not
// introduce a cycle.
//
// Field builder flow:
//
//	Field("name").Ref("arcoris.meta.Name").
//		Required().
//		Nullable().
//		Description("Reference to a reusable name type.")
func (b FieldBuilder) Ref(name string) RefField {
	return RefField{field: b.state(), descriptor: Ref(name)}
}
