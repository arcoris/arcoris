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

func TestValidateMapEntryRejectsEmptyKey(t *testing.T) {
	err := validateMapEntry(0, Entry{Key: "", Value: Null()}, nil)

	requireValueError(
		t,
		err,
		ErrEmptyKey,
		mapEntryKeyPath(0),
		ErrorReasonEmptyKey,
	)

	requireErrorIs(t, err, ErrInvalidMap)
}

func TestValidateMapEntryRejectsInvalidValue(t *testing.T) {
	err := validateMapEntry(0, Entry{Key: "name"}, nil)

	requireValueError(
		t,
		err,
		ErrInvalidEntry,
		mapEntryValuePath(0),
		ErrorReasonInvalidValue,
	)

	requireErrorIs(t, err, ErrInvalidMap)
}

func TestValidateMapEntryRejectsDuplicateKey(t *testing.T) {
	err := validateMapEntry(
		1,
		MapEntry("name", Null()),
		[]Entry{MapEntry("name", Null())},
	)

	requireValueError(
		t,
		err,
		ErrDuplicateKey,
		mapEntryKeyPath(1),
		ErrorReasonDuplicateKey,
	)

	requireErrorIs(t, err, ErrInvalidMap)
}
