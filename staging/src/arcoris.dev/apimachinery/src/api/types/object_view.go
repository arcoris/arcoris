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

// ObjectView exposes read-only object payload data.
type ObjectView struct {
	// payload is a detached copy of the object descriptor payload.
	payload objectPayload
}

// Object returns an object view when t is TypeObject.
func (t Type) Object() (ObjectView, bool) {
	if t.code != TypeObject {
		return ObjectView{}, false
	}
	return ObjectView{payload: cloneObjectPayload(t.object)}, true
}

// Fields returns detached object fields in declaration order.
func (v ObjectView) Fields() []FieldDescriptor {
	return cloneFields(v.payload.fields)
}

// UnknownFields returns the object unknown-field policy.
func (v ObjectView) UnknownFields() UnknownFieldPolicy {
	return v.payload.unknown
}
