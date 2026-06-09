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

// RefView exposes read-only reference payload data.
type RefView struct {
	// payload is a detached copy of the reference descriptor payload.
	payload refPayload
}

// AsRef returns a reference view when desc is DescriptorRef.
func (desc Descriptor) AsRef() (RefView, bool) {
	if desc.code != DescriptorRef {
		return RefView{}, false
	}

	return RefView{payload: cloneRefPayload(desc.ref)}, true
}

// Name returns the referenced type name.
func (v RefView) Name() TypeName {
	return v.payload.name
}
