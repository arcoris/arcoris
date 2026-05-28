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

// Int32View exposes read-only TypeInt32 payload data.
type Int32View struct {
	// payload is a detached copy of the TypeInt32 payload.
	payload int32Payload
}

// Int32 returns a read-only TypeInt32 view when t is an int32 descriptor.
func (t Type) Int32() (Int32View, bool) {
	if t.code != TypeInt32 {
		return Int32View{}, false
	}
	return Int32View{payload: cloneInt32Payload(t.int32)}, true
}

// Min returns the inclusive lower bound when one is set.
func (v Int32View) Min() (int32, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Int32View) Max() (int32, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted int32 literals detached from descriptor storage.
func (v Int32View) Enum() []int32 {
	return cloneSlice(v.payload.enum)
}
