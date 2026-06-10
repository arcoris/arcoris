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

import "testing"

func TestRecordViewAccessors(t *testing.T) {
	value := mustRecord(t,
		MustRecordMember("name", StringValue("worker")),
		MustRecordMember("payload", BytesValue([]byte{1, 2})),
	)
	view, ok := value.AsRecord()
	requireEqual(t, ok, true)

	requireEqual(t, view.Len(), 2)
	requireEqual(t, view.IsEmpty(), false)
	requireEqual(t, view.Has(MustMemberName("name")), true)
	requireEqual(t, view.Has(MustMemberName("missing")), false)
	requireMemberNamesEqual(t, view.Names(), []MemberName{"name", "payload"})

	got, ok := view.Get(MustMemberName("name"))
	requireEqual(t, ok, true)
	name, ok := got.AsString()
	requireEqual(t, ok, true)
	requireEqual(t, name, "worker")

	member, ok := view.Member(1)
	requireEqual(t, ok, true)
	requireEqual(t, member.Name, MemberName("payload"))

	_, ok = view.Get(MustMemberName("missing"))
	requireEqual(t, ok, false)
}

func TestRecordViewEmptySlicesAreNonNil(t *testing.T) {
	value := mustRecord(t)
	view, ok := value.AsRecord()
	requireEqual(t, ok, true)

	members := view.Members()
	if members == nil {
		t.Fatal("Members() returned nil")
	}
	requireEqual(t, len(members), 0)

	names := view.Names()
	if names == nil {
		t.Fatal("Names() returned nil")
	}
	requireEqual(t, len(names), 0)
}

func TestRecordViewMembersReturnsDetachedSliceWithoutDeepCloningValues(t *testing.T) {
	value := mustRecord(t, MustRecordMember("payload", BytesValue([]byte{1, 2})))
	view, ok := value.AsRecord()
	requireEqual(t, ok, true)

	members := view.Members()
	members[0].Name = "changed"
	members[0].Value.bytesValue[0] = 9

	requireEqual(t, view.Has(MustMemberName("payload")), true)
	got, ok := view.Get(MustMemberName("payload"))
	requireEqual(t, ok, true)

	bytes, ok := got.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{9, 2})
}

func TestRecordViewCloneMembersDeepClonesValues(t *testing.T) {
	value := mustRecord(t, MustRecordMember("payload", BytesValue([]byte{1, 2})))
	view, ok := value.AsRecord()
	requireEqual(t, ok, true)

	members := view.CloneMembers()
	members[0].Value.bytesValue[0] = 9

	got, ok := view.Get(MustMemberName("payload"))
	requireEqual(t, ok, true)
	bytes, ok := got.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{1, 2})
}

func TestRecordViewCloneGetDeepClonesValue(t *testing.T) {
	value := mustRecord(t, MustRecordMember("payload", BytesValue([]byte{1, 2})))
	view, ok := value.AsRecord()
	requireEqual(t, ok, true)

	got, ok := view.CloneGet(MustMemberName("payload"))
	requireEqual(t, ok, true)
	got.bytesValue[0] = 9

	again, ok := view.Get(MustMemberName("payload"))
	requireEqual(t, ok, true)
	bytes, ok := again.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{1, 2})
}

func TestRecordViewForEachVisitsInOrderAndStopsEarly(t *testing.T) {
	value := mustRecord(t,
		MustRecordMember("first", StringValue("one")),
		MustRecordMember("second", StringValue("two")),
	)
	view, ok := value.AsRecord()
	requireEqual(t, ok, true)

	var names []MemberName
	view.ForEach(func(index int, member RecordMember) bool {
		names = append(names, member.Name)
		return false
	})

	requireMemberNamesEqual(t, names, []MemberName{"first"})
}

func TestRecordWrongKindAccessorReturnsFalse(t *testing.T) {
	_, ok := NullValue().AsRecord()
	requireEqual(t, ok, false)
}
