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

func TestValidateObjectMemberRejectsEmptyName(t *testing.T) {
	err := validateObjectMember(0, Member{Name: "", Value: NullValue()}, nil)
	valueErr := requireValueError(
		t,
		err,
		ErrEmptyName,
		objectMemberNamePath(0),
		ErrorReasonEmptyName,
	)

	requireErrorIs(t, err, ErrInvalidObject)
	requireEqual(t, valueErr.Cause == nil, true)
}

func TestValidateObjectMemberRejectsInvalidValue(t *testing.T) {
	err := validateObjectMember(0, Member{Name: "name"}, nil)

	requireValueError(
		t,
		err,
		ErrInvalidMember,
		objectMemberValuePath(0),
		ErrorReasonInvalidValue,
	)
	requireErrorIs(t, err, ErrInvalidObject)
}

func TestValidateObjectMemberRejectsDuplicateName(t *testing.T) {
	err := validateObjectMember(
		1,
		ObjectMember("name", NullValue()),
		[]Member{ObjectMember("name", NullValue())},
	)

	requireValueError(
		t,
		err,
		ErrDuplicateName,
		objectMemberNamePath(1),
		ErrorReasonDuplicateName,
	)
	requireErrorIs(t, err, ErrInvalidObject)
}
