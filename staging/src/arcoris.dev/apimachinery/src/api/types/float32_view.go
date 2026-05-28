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

// Float32View exposes read-only TypeFloat32 payload data.
type Float32View struct {
	// payload is a detached copy of the TypeFloat32 payload.
	payload float32Payload
}

// Float32 returns a read-only TypeFloat32 view when t is a float32 descriptor.
func (t Type) Float32() (Float32View, bool) {
	return Float32View{payload: cloneFloat32Payload(t.float32)}, t.code == TypeFloat32
}

// Min returns the inclusive lower bound when one is set.
func (v Float32View) Min() (float32, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Float32View) Max() (float32, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted float32 literals detached from descriptor storage.
func (v Float32View) Enum() []float32 {
	return cloneFloat32s(v.payload.enum)
}
