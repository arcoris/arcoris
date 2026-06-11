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

package valuefieldset

import "testing"

func TestRecordMemberPathAcceptsStructurallyValidName(t *testing.T) {
	base := rootField("spec")

	got, err := recordMemberPath(base, "x-y")
	requireNoError(t, err)

	if got.String() != `$.spec."x-y"` {
		t.Fatalf("path = %s", got)
	}
}

func TestRecordMemberPathRejectsInvalidFieldName(t *testing.T) {
	_, err := recordMemberPath(rootField("spec"), "")

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorReason(t, err, ErrorReasonInvalidFieldName)
	requireErrorPath(t, err, "$.spec")
}

func TestMapMemberPathAcceptsStructurallyValidKey(t *testing.T) {
	base := rootField("labels")

	got, err := mapMemberPath(base, "x-y")
	requireNoError(t, err)

	if got.String() != `$.labels["x-y"]` {
		t.Fatalf("path = %s", got)
	}
}

func TestMapMemberPathRejectsInvalidMapKey(t *testing.T) {
	_, err := mapMemberPath(rootField("labels"), "")

	requireErrorIs(t, err, ErrInvalidValue)
	requireErrorReason(t, err, ErrorReasonInvalidMapKey)
	requireErrorPath(t, err, "$.labels")
}
