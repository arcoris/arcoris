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

package types

import "testing"

func TestObjectValidateRejectsDuplicateInvalidFieldAndUnknownPolicy(t *testing.T) {
	requireErrorIs(t, ValidateType(Object(Field("name").String().Required(), Field("name").String().Optional()).Type(), nil), ErrDuplicateField)
	requireErrorIs(t, ValidateType(Object(Field("bad-name").String().Required()).Type(), nil), ErrInvalidField)
	requireErrorIs(t, ValidateType(Object(Field("name").String()).Type(), nil), ErrInvalidField)

	invalidPolicy := Object().Type()
	invalidPolicy.object.unknown = UnknownFieldPolicy(99)
	requireErrorIs(t, ValidateType(invalidPolicy, nil), ErrInvalidType)
}
