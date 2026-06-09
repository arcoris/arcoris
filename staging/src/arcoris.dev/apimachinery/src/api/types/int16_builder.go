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

// Int16 starts a field descriptor whose value descriptor is a fixed-width int16.
//
// The descriptor records portable int16 constraints. It never uses Go's
// platform-sized int contract, so generated APIs and remote descriptors remain
// architecture-independent.
//
// Field builder flow:
//
//	Field("shard").Int16().
//		Required().
//		Nullable().
//		Min(0).
//		Max(1024).
//		Enum(0, 1, 2).
//		Description("Shard identifier.")
func (b FieldBuilder) Int16() Int16Field {
	return Int16Field{field: b.state(), descriptor: Int16()}
}
