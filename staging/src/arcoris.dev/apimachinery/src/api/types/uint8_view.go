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

// Uint8View exposes read-only DescriptorUint8 payload data.
type Uint8View struct {
	// payload is a detached copy of the DescriptorUint8 payload.
	payload uint8Payload
}

// AsUint8 returns a read-only DescriptorUint8 view when desc is a uint8 descriptor.
func (desc Descriptor) AsUint8() (Uint8View, bool) {
	if desc.code != DescriptorUint8 {
		return Uint8View{}, false
	}

	return Uint8View{payload: cloneUint8Payload(desc.uint8)}, true
}

// Min returns the inclusive lower bound when one is set.
func (v Uint8View) Min() (uint8, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Uint8View) Max() (uint8, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted uint8 literals detached from descriptor storage.
func (v Uint8View) Enum() []uint8 {
	return cloneSlice(v.payload.enum)
}
