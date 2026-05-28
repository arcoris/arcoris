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

// MapView exposes read-only map payload data.
type MapView struct {
	// payload is a detached copy of the map descriptor payload.
	payload mapPayload
}

// Map returns a map view when t is TypeMap.
func (t Type) Map() (MapView, bool) {
	return MapView{payload: cloneType(t).mapType}, t.code == TypeMap
}

// Key returns the map key type.
func (v MapView) Key() MapKeyType {
	return v.payload.key
}

// Value returns a detached map value descriptor.
func (v MapView) Value() Type {
	if v.payload.value == nil {
		return Type{}
	}
	return cloneType(*v.payload.value)
}

// MinLen returns the map minimum size rule.
func (v MapView) MinLen() (int, bool) {
	return v.payload.minLen.value, v.payload.minLen.set
}

// MaxLen returns the map maximum size rule.
func (v MapView) MaxLen() (int, bool) {
	return v.payload.maxLen.value, v.payload.maxLen.set
}
