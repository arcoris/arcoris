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

// Int64 starts a field descriptor whose value type is a fixed-width int64.
//
// The descriptor records portable int64 constraints. It never uses Go's
// platform-sized int contract, so generated APIs and remote descriptors remain
// architecture-independent.
//
// Field builder flow:
//
//	Field("maxConcurrency").Int64().
//		Required().
//		Nullable().
//		Min(1).
//		Max(1000).
//		Enum(1, 10, 100).
//		Description("Maximum concurrent work.")
func (b FieldBuilder) Int64() Int64Field {
	return Int64Field{field: b.state(), typ: Int64()}
}
