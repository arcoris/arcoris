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

// Uint64View exposes read-only DescriptorUint64 payload data.
type Uint64View struct {
	// payload is a detached copy of the DescriptorUint64 payload.
	payload uint64Payload
}

// AsUint64 returns a read-only DescriptorUint64 view when desc is a uint64 descriptor.
func (desc Descriptor) AsUint64() (Uint64View, bool) {
	if desc.code != DescriptorUint64 {
		return Uint64View{}, false
	}

	return Uint64View{payload: cloneUint64Payload(desc.uint64)}, true
}

// Min returns the inclusive lower bound when one is set.
func (v Uint64View) Min() (uint64, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Uint64View) Max() (uint64, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted uint64 literals detached from descriptor storage.
func (v Uint64View) Enum() []uint64 {
	return cloneSlice(v.payload.enum)
}
