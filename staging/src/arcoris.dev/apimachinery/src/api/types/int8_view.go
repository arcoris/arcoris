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

// Int8View exposes read-only DescriptorInt8 payload data.
type Int8View struct {
	// payload is a detached copy of the DescriptorInt8 payload.
	payload int8Payload
}

// AsInt8 returns a read-only DescriptorInt8 view when desc is an int8 descriptor.
func (desc Descriptor) AsInt8() (Int8View, bool) {
	if desc.code != DescriptorInt8 {
		return Int8View{}, false
	}

	return Int8View{payload: cloneInt8Payload(desc.int8)}, true
}

// Min returns the inclusive lower bound when one is set.
func (v Int8View) Min() (int8, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Int8View) Max() (int8, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted int8 literals detached from descriptor storage.
func (v Int8View) Enum() []int8 {
	return cloneSlice(v.payload.enum)
}
