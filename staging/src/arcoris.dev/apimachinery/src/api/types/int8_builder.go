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

// Int8 starts a field descriptor whose value type is a fixed-width int8.
//
// The descriptor records portable int8 constraints. It never uses Go's
// platform-sized int contract, so generated APIs and remote descriptors remain
// architecture-independent.
//
// Field builder flow:
//
//	Field("priority").Int8().
//		Optional().
//		Nullable().
//		Range(0, 10).
//		Enum(0, 5, 10).
//		Description("Scheduling priority.")
func (b FieldBuilder) Int8() Int8Field {
	return Int8Field{field: b.state(), typ: Int8()}
}
