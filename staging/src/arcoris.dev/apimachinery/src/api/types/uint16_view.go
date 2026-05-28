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

// Uint16View exposes read-only TypeUint16 payload data.
type Uint16View struct {
	// payload is a detached copy of the TypeUint16 payload.
	payload uint16Payload
}

// Uint16 returns a read-only TypeUint16 view when t is a uint16 descriptor.
func (t Type) Uint16() (Uint16View, bool) {
	if t.code != TypeUint16 {
		return Uint16View{}, false
	}
	return Uint16View{payload: cloneUint16Payload(t.uint16)}, true
}

// Min returns the inclusive lower bound when one is set.
func (v Uint16View) Min() (uint16, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Uint16View) Max() (uint16, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted uint16 literals detached from descriptor storage.
func (v Uint16View) Enum() []uint16 {
	return cloneSlice(v.payload.enum)
}
