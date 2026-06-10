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

func TestNewMemberNameValidatesNonEmptyNames(t *testing.T) {
	name, err := NewMemberName("payload")
	requireNoError(t, err)

	requireEqual(t, name.String(), "payload")
	requireEqual(t, name.IsZero(), false)
}

func TestNewMemberNameRejectsEmptyName(t *testing.T) {
	err := MemberName("").ValidateLexical()

	requireValueError(
		t,
		err,
		ErrEmptyMemberName,
		pathMemberName,
		ErrorReasonEmptyMemberName,
	)
	requireErrorIs(t, err, ErrInvalidRecord)
}
