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

// Uint64 starts a field descriptor whose value type is a fixed-width uint64.
//
// The descriptor records portable uint64 constraints. It never uses Go's
// platform-sized uint contract, so generated APIs and remote descriptors remain
// architecture-independent.
//
// Field builder flow:
//
//	Field("observedBytes").Uint64().
//		Optional().
//		Nullable().
//		Min(0).
//		Max(1099511627776).
//		Enum(0, 1024).
//		Description("Observed byte count.")
func (b FieldBuilder) Uint64() Uint64Field {
	return Uint64Field{field: b.state(), typ: Uint64()}
}
