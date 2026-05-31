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

func TestObjectValueConstructsObjectValue(t *testing.T) {
	value, err := ObjectValue(
		ObjectMember("name", StringValue("worker")),
		ObjectMember("replicas", Int64Value(3)),
	)
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindObject)
	requireEqual(t, len(value.objectValue.members), 2)
	requireEqual(t, value.objectValue.members[0].Name, "name")
	requireEqual(t, value.objectValue.members[1].Name, "replicas")
}

func TestObjectValueAcceptsEmptyObject(t *testing.T) {
	value, err := ObjectValue()
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindObject)
	requireEqual(t, len(value.objectValue.members), 0)
}

func TestObjectValueAcceptsExplicitNullMember(t *testing.T) {
	value, err := ObjectValue(ObjectMember("optional", NullValue()))
	requireNoError(t, err)

	view, ok := value.Object()
	requireEqual(t, ok, true)

	memberValue, ok := view.Get("optional")
	requireEqual(t, ok, true)
	requireEqual(t, memberValue.IsNull(), true)
}

func TestObjectValueRejectsMalformedMembers(t *testing.T) {
	tests := []struct {
		name   string
		input  []Member
		target error
		path   string
		reason ErrorReason
	}{
		{
			name:   "empty name",
			input:  []Member{{Name: "", Value: NullValue()}},
			target: ErrEmptyName,
			path:   objectMemberNamePath(0),
			reason: ErrorReasonEmptyName,
		},
		{
			name:   "invalid zero value",
			input:  []Member{{Name: "name"}},
			target: ErrInvalidMember,
			path:   objectMemberValuePath(0),
			reason: ErrorReasonInvalidValue,
		},
		{
			name: "duplicate name",
			input: []Member{
				ObjectMember("name", NullValue()),
				ObjectMember("name", StringValue("worker")),
			},
			target: ErrDuplicateName,
			path:   objectMemberNamePath(1),
			reason: ErrorReasonDuplicateName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ObjectValue(tt.input...)

			requireValueError(t, err, tt.target, tt.path, tt.reason)
			requireErrorIs(t, err, ErrInvalidObject)
		})
	}
}

func TestMustObjectValuePanicsOnMalformedMembers(t *testing.T) {
	requirePanic(t, func() {
		MustObjectValue(Member{Name: "", Value: NullValue()})
	})
}
