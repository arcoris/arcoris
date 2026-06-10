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

func TestValidateRecordMemberRejectsEmptyName(t *testing.T) {
	err := validateRecordMember(0, RecordMember{Name: "", Value: NullValue()})
	valueErr := requireValueError(
		t,
		err,
		ErrEmptyMemberName,
		recordMemberNamePath(0),
		ErrorReasonEmptyMemberName,
	)

	requireErrorIs(t, err, ErrInvalidRecord)
	requireEqual(t, valueErr.Cause == nil, true)
}

func TestValidateRecordMemberRejectsInvalidValue(t *testing.T) {
	err := validateRecordMember(0, RecordMember{Name: MustMemberName("name")})

	requireValueError(
		t,
		err,
		ErrInvalidRecordMember,
		recordMemberValuePath(0),
		ErrorReasonInvalidValue,
	)
	requireErrorIs(t, err, ErrInvalidRecord)
}
