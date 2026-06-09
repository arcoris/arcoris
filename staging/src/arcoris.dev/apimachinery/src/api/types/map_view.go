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

// AsMap returns a map view when desc is DescriptorMap.
func (desc Descriptor) AsMap() (MapView, bool) {
	if desc.code != DescriptorMap {
		return MapView{}, false
	}

	return MapView{payload: cloneMapPayload(desc.mapType)}, true
}

// Key returns a detached map key descriptor.
func (v MapView) Key() Descriptor {
	if v.payload.key == nil {
		return Descriptor{}
	}

	return cloneDescriptor(*v.payload.key)
}

// Value returns a detached map value descriptor.
func (v MapView) Value() Descriptor {
	if v.payload.value == nil {
		return Descriptor{}
	}

	return cloneDescriptor(*v.payload.value)
}

// MinEntries returns the map minimum entry count rule.
func (v MapView) MinEntries() (int, bool) {
	return v.payload.minLen.value, v.payload.minLen.set
}

// MaxEntries returns the map maximum entry count rule.
func (v MapView) MaxEntries() (int, bool) {
	return v.payload.maxLen.value, v.payload.maxLen.set
}
