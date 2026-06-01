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

package value

// ObjectView exposes read-only object payload data.
//
// A view references private immutable-by-convention payload data. Methods that
// return values or slices clone those results, so callers cannot mutate the
// original Value or the view through returned data.
type ObjectView struct {
	// payload is the private object payload the view reads from.
	payload objectPayload
}

// Object returns a read-only object view when v is KindObject.
//
// For every other kind, Object returns ok=false. The returned view preserves
// object member order without eagerly deep-cloning the whole object payload.
func (v Value) Object() (ObjectView, bool) {
	if v.kind != KindObject {
		return ObjectView{}, false
	}

	return ObjectView{payload: v.objectValue}, true
}

// Len returns the number of object members.
func (o ObjectView) Len() int {
	return len(o.payload.members)
}

// IsEmpty reports whether the object has no members.
func (o ObjectView) IsEmpty() bool {
	return len(o.payload.members) == 0
}

// Has reports whether the object contains name.
//
// Lookup is linear by design. The value model favors cheap construction and
// cloning over storing per-object indexes for small payloads.
func (o ObjectView) Has(name string) bool {
	return findObjectMember(o.payload.members, name) >= 0
}

// Get returns a cloned member value by name.
//
// The returned Value is detached from the view. Missing names return the zero
// Value and ok=false.
func (o ObjectView) Get(name string) (Value, bool) {
	index := findObjectMember(o.payload.members, name)
	if index < 0 {
		return Value{}, false
	}

	return o.payload.members[index].Value.Clone(), true
}

// Members returns detached members in original order.
//
// Both the slice and every nested Value are cloned. Mutating the result cannot
// affect the view or source Value.
func (o ObjectView) Members() []Member {
	if len(o.payload.members) == 0 {
		return []Member{}
	}

	return cloneMembers(o.payload.members)
}

// Names returns detached member names in original order.
//
// Names mirrors Members order and returns a caller-owned slice.
func (o ObjectView) Names() []string {
	names := make([]string, 0, len(o.payload.members))
	for _, member := range o.payload.members {
		names = append(names, member.Name)
	}

	return names
}
