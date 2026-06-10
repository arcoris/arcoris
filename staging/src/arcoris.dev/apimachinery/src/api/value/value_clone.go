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
// inside record/list payloads are copied recursively. Clone never validates
// payloads because Values can only be built through constructors that already
// enforce construction invariants.
func (v Value) Clone() Value {
	clone := v

	switch v.kind {
	case KindBytes:
		clone.bytesValue = cloneBytes(v.bytesValue)
	case KindRecord:
		clone.recordValue = cloneRecordPayload(v.recordValue)
	case KindList:
		clone.listValue = cloneListPayload(v.listValue)
	}

	return clone
}

// cloneRecordPayload deep-copies record members.
//
// No member-name index is rebuilt because record payloads intentionally use
// linear lookup.
func cloneRecordPayload(payload recordPayload) recordPayload {
	return recordPayload{members: cloneMembers(payload.members)}
}

// cloneListPayload deep-copies list items.
//
// Nil slices remain nil so empty payload compaction is preserved across clones.
func cloneListPayload(payload listPayload) listPayload {
	return listPayload{items: cloneValues(payload.items)}
}

// cloneMembers deep-copies record members and nested values.
func cloneMembers(members []RecordMember) []RecordMember {
	if members == nil {
		return nil
	}

	cloned := make([]RecordMember, len(members))
	for i, member := range members {
		cloned[i] = cloneRecordMember(member)
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
