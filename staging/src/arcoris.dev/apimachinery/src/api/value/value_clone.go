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

// Clone returns a deep copy of v.
//
// Scalar payloads are copied by value. Mutable payloads and every nested Value
// inside object/list payloads are copied recursively. Clone never validates
// payloads because Values can only be built through constructors that already
// enforce construction invariants.
func (v Value) Clone() Value {
	clone := v

	switch v.kind {
	case KindBytes:
		clone.bytesValue = cloneBytes(v.bytesValue)
	case KindObject:
		clone.objectValue = cloneObjectPayload(v.objectValue)
	case KindList:
		clone.listValue = cloneListPayload(v.listValue)
	}

	return clone
}

// cloneObjectPayload deep-copies object members.
//
// No member-name index is rebuilt because object payloads intentionally use
// linear lookup.
func cloneObjectPayload(payload objectPayload) objectPayload {
	return objectPayload{members: cloneMembers(payload.members)}
}

// cloneListPayload deep-copies list items.
//
// Nil slices remain nil so empty payload compaction is preserved across clones.
func cloneListPayload(payload listPayload) listPayload {
	return listPayload{items: cloneValues(payload.items)}
}

// cloneMembers deep-copies object members and nested values.
//
// ObjectMember performs the nested Value clone, keeping clone semantics aligned
// with object construction semantics.
func cloneMembers(members []Member) []Member {
	if members == nil {
		return nil
	}

	cloned := make([]Member, len(members))
	for i, member := range members {
		cloned[i] = ObjectMember(member.Name, member.Value)
	}

	return cloned
}

// cloneValues deep-copies nested Value slices.
//
// The returned slice is caller-owned. Nil input remains nil to preserve compact
// empty payload shape.
func cloneValues(values []Value) []Value {
	if values == nil {
		return nil
	}

	cloned := make([]Value, len(values))
	for i, item := range values {
		cloned[i] = item.Clone()
	}

	return cloned
}
