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

func TestRecordValueConstructsRecordValue(t *testing.T) {
	value, err := RecordValue(
		MustRecordMember("name", StringValue("worker")),
		MustRecordMember("replicas", Int64Value(3)),
	)
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindRecord)
	requireEqual(t, len(value.recordValue.members), 2)
	requireEqual(t, value.recordValue.members[0].Name, MemberName("name"))
	requireEqual(t, value.recordValue.members[1].Name, MemberName("replicas"))
}

func TestRecordValueAcceptsEmptyRecord(t *testing.T) {
	value, err := RecordValue()
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindRecord)
	requireEqual(t, len(value.recordValue.members), 0)
}

func TestRecordValueAcceptsExplicitNullMember(t *testing.T) {
	value, err := RecordValue(MustRecordMember("optional", NullValue()))
	requireNoError(t, err)

	view, ok := value.AsRecord()
	requireEqual(t, ok, true)

	memberValue, ok := view.Get(MustMemberName("optional"))
	requireEqual(t, ok, true)
	requireEqual(t, memberValue.IsNull(), true)
}

func TestRecordValueRejectsMalformedMembers(t *testing.T) {
	tests := []struct {
		name   string
		input  []RecordMember
		target error
		path   string
		reason ErrorReason
	}{
		{
			name:   "empty name",
			input:  []RecordMember{{Name: "", Value: NullValue()}},
			target: ErrEmptyMemberName,
			path:   recordMemberNamePath(0),
			reason: ErrorReasonEmptyMemberName,
		},
		{
			name:   "invalid zero value",
			input:  []RecordMember{{Name: MustMemberName("name")}},
			target: ErrInvalidRecordMember,
			path:   recordMemberValuePath(0),
			reason: ErrorReasonInvalidValue,
		},
		{
			name: "duplicate name",
			input: []RecordMember{
				MustRecordMember("name", NullValue()),
				MustRecordMember("name", StringValue("worker")),
			},
			target: ErrDuplicateMemberName,
			path:   recordMemberNamePath(1),
			reason: ErrorReasonDuplicateMemberName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RecordValue(tt.input...)

			requireValueError(t, err, tt.target, tt.path, tt.reason)
			requireErrorIs(t, err, ErrInvalidRecord)
		})
	}
}

func TestRecordValueRejectsLargeDuplicateMemberName(t *testing.T) {
	members := make([]RecordMember, 0, recordDuplicateMapThreshold+2)
	for i := 0; i < recordDuplicateMapThreshold+1; i++ {
		members = append(members, MustRecordMember("name-"+string(rune('a'+i)), NullValue()))
	}
	members = append(members, MustRecordMember("name-a", StringValue("duplicate")))

	_, err := RecordValue(members...)

	requireValueError(
		t,
		err,
		ErrDuplicateMemberName,
		recordMemberNamePath(recordDuplicateMapThreshold+1),
		ErrorReasonDuplicateMemberName,
	)
}

func TestRecordValueClonesMemberValuesAtConstructionBoundary(t *testing.T) {
	source := BytesValue([]byte{1, 2})
	record := mustRecord(t, NewRecordMember(MustMemberName("payload"), source))

	source.bytesValue[0] = 9

	view, ok := record.AsRecord()
	requireEqual(t, ok, true)
	got, ok := view.Get(MustMemberName("payload"))
	requireEqual(t, ok, true)
	bytes, ok := got.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{1, 2})
}

func TestMustRecordValuePanicsOnMalformedMembers(t *testing.T) {
	requirePanic(t, func() {
		MustRecordValue(RecordMember{Name: "", Value: NullValue()})
	})
}
