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

// Uint32View exposes read-only TypeUint32 payload data.
type Uint32View struct {
	// payload is a detached copy of the TypeUint32 payload.
	payload uint32Payload
}

// Uint32 returns a read-only TypeUint32 view when t is a uint32 descriptor.
func (t Type) Uint32() (Uint32View, bool) {
	if t.code != TypeUint32 {
		return Uint32View{}, false
	}
	return Uint32View{payload: cloneUint32Payload(t.uint32)}, true
}

// Min returns the inclusive lower bound when one is set.
func (v Uint32View) Min() (uint32, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Uint32View) Max() (uint32, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted uint32 literals detached from descriptor storage.
func (v Uint32View) Enum() []uint32 {
	return cloneSlice(v.payload.enum)
}
