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

// Int32 starts a field descriptor whose value descriptor is a fixed-width int32.
//
// The descriptor records portable int32 constraints. It never uses Go's
// platform-sized int contract, so generated APIs and remote descriptors remain
// architecture-independent.
//
// Field builder flow:
//
//	Field("replicas").Int32().
//		Optional().
//		Nullable().
//		Range(1, 1000).
//		Enum(1, 3, 5).
//		Description("Desired replica count.")
func (b FieldBuilder) Int32() Int32Field {
	return Int32Field{field: b.state(), descriptor: Int32()}
}
