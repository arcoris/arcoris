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

// ListView exposes read-only list payload data.
type ListView struct {
	// payload is a detached copy of the list descriptor payload.
	payload listPayload
}

// AsList returns a list view when desc is DescriptorList.
func (desc Descriptor) AsList() (ListView, bool) {
	if desc.code != DescriptorList {
		return ListView{}, false
	}

	return ListView{payload: cloneListPayload(desc.list)}, true
}

// Element returns a detached list element descriptor.
func (v ListView) Element() Descriptor {
	if v.payload.elem == nil {
		return Descriptor{}
	}

	return cloneDescriptor(*v.payload.elem)
}

// MinItems returns the list minimum item count rule.
func (v ListView) MinItems() (int, bool) {
	return v.payload.minLen.value, v.payload.minLen.set
}

// MaxItems returns the list maximum item count rule.
func (v ListView) MaxItems() (int, bool) {
	return v.payload.maxLen.value, v.payload.maxLen.set
}

// Semantics returns the list semantic policy.
func (v ListView) Semantics() ListSemantics {
	return v.payload.semantics
}

// MapKeys returns detached list map keys.
func (v ListView) MapKeys() []FieldName {
	return cloneSlice(v.payload.mapKeys)
}
