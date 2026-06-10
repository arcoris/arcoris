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

// RecordView exposes read-only record payload data.
//
// A view references private payload data owned by an immutable Value. Read
// methods return nested Values by value; explicit Clone methods deep-copy nested
// payloads when a caller needs a detached ownership boundary.
type RecordView struct {
	// payload is the private record payload the view reads from.
	payload recordPayload
}

// AsRecord returns a read-only record view when v is KindRecord.
//
// For every other kind, AsRecord returns ok=false. The returned view preserves
// record member order without eagerly deep-cloning the whole record payload.
func (v Value) AsRecord() (RecordView, bool) {
	if v.kind != KindRecord {
		return RecordView{}, false
	}

	return RecordView{payload: v.recordValue}, true
}

// Len returns the number of record members.
func (o RecordView) Len() int {
	return len(o.payload.members)
}

// IsEmpty reports whether the record has no members.
func (o RecordView) IsEmpty() bool {
	return len(o.payload.members) == 0
}

// Has reports whether the record contains name.
//
// Lookup is linear by design. The value model favors cheap construction and
// cloning over storing per-record indexes for small payloads.
func (o RecordView) Has(name MemberName) bool {
	return findRecordMember(o.payload.members, name) >= 0
}

// Get returns a member value by name without deep-cloning it.
//
// Value has no external mutators, so returning the nested Value by value is
// safe. Missing names return the zero Value and ok=false.
func (o RecordView) Get(name MemberName) (Value, bool) {
	index := findRecordMember(o.payload.members, name)
	if index < 0 {
		return Value{}, false
	}

	return o.payload.members[index].Value, true
}

// CloneGet returns a deep-cloned member value by name.
func (o RecordView) CloneGet(name MemberName) (Value, bool) {
	nested, ok := o.Get(name)
	if !ok {
		return Value{}, false
	}

	return nested.Clone(), true
}

// Member returns one member by index without deep-cloning its nested Value.
func (o RecordView) Member(index int) (RecordMember, bool) {
	if index < 0 || index >= len(o.payload.members) {
		return RecordMember{}, false
	}

	return o.payload.members[index], true
}

// Members returns a detached member slice in original order.
//
// The slice is caller-owned. Nested Values are returned by value without deep
// cloning; use CloneMembers when a recursive copy is required.
func (o RecordView) Members() []RecordMember {
	if len(o.payload.members) == 0 {
		return []RecordMember{}
	}

	members := make([]RecordMember, len(o.payload.members))
	copy(members, o.payload.members)
	return members
}

// CloneMembers returns detached members with deep-cloned nested Values.
func (o RecordView) CloneMembers() []RecordMember {
	if len(o.payload.members) == 0 {
		return []RecordMember{}
	}

	return cloneMembers(o.payload.members)
}

// Names returns detached member names in original order.
//
// Names mirrors Members order and returns a caller-owned slice.
func (o RecordView) Names() []MemberName {
	names := make([]MemberName, 0, len(o.payload.members))
	for _, member := range o.payload.members {
		names = append(names, member.Name)
	}

	return names
}

// ForEach visits members in original order until fn returns false.
func (o RecordView) ForEach(fn func(index int, member RecordMember) bool) {
	for i, member := range o.payload.members {
		if !fn(i, member) {
			return
		}
	}
}
