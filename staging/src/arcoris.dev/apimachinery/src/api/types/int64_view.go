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

// Int64View exposes read-only TypeInt64 payload data.
type Int64View struct {
	// payload is a detached copy of the TypeInt64 payload.
	payload int64Payload
}

// Int64 returns a read-only TypeInt64 view when t is an int64 descriptor.
func (t Type) Int64() (Int64View, bool) {
	if t.code != TypeInt64 {
		return Int64View{}, false
	}

	return Int64View{payload: cloneInt64Payload(t.int64)}, true
}

// Min returns the inclusive lower bound when one is set.
func (v Int64View) Min() (int64, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Int64View) Max() (int64, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted int64 literals detached from descriptor storage.
func (v Int64View) Enum() []int64 {
	return cloneSlice(v.payload.enum)
}
