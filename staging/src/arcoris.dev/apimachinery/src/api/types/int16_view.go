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

// Int16View exposes read-only TypeInt16 payload data.
type Int16View struct {
	// payload is a detached copy of the TypeInt16 payload.
	payload int16Payload
}

// Int16 returns a read-only TypeInt16 view when t is an int16 descriptor.
func (t Type) Int16() (Int16View, bool) {
	return Int16View{payload: cloneInt16Payload(t.int16)}, t.code == TypeInt16
}

// Min returns the inclusive lower bound when one is set.
func (v Int16View) Min() (int16, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Int16View) Max() (int16, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted int16 literals detached from descriptor storage.
func (v Int16View) Enum() []int16 {
	return cloneInt16s(v.payload.enum)
}
