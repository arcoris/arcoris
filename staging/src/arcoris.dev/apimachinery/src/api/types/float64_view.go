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

// Float64View exposes read-only TypeFloat64 payload data.
type Float64View struct {
	// payload is a detached copy of the TypeFloat64 payload.
	payload float64Payload
}

// Float64 returns a read-only TypeFloat64 view when t is a float64 descriptor.
func (t Type) Float64() (Float64View, bool) {
	if t.code != TypeFloat64 {
		return Float64View{}, false
	}

	return Float64View{payload: cloneFloat64Payload(t.float64)}, true
}

// Min returns the inclusive lower bound when one is set.
func (v Float64View) Min() (float64, bool) {
	return v.payload.min.value, v.payload.min.set
}

// Max returns the inclusive upper bound when one is set.
func (v Float64View) Max() (float64, bool) {
	return v.payload.max.value, v.payload.max.set
}

// Enum returns accepted float64 literals detached from descriptor storage.
func (v Float64View) Enum() []float64 {
	return cloneSlice(v.payload.enum)
}
