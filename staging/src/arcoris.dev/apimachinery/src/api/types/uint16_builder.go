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

// Uint16 starts a field descriptor whose value descriptor is a fixed-width uint16.
//
// The descriptor records portable uint16 constraints. It never uses Go's
// platform-sized uint contract, so generated APIs and remote descriptors remain
// architecture-independent.
//
// Field builder flow:
//
//	Field("port").Uint16().
//		Required().
//		Nullable().
//		Range(1, 65535).
//		Enum(80, 443).
//		Description("Network port.")
func (b FieldBuilder) Uint16() Uint16Field {
	return Uint16Field{field: b.state(), descriptor: Uint16()}
}
