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

package valuecompare

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestCompareUnknownObjectMembersPruneIgnoresUnknowns(t *testing.T) {
	oldObject, _ := valueObject("extra", "old").Object()
	newObject, _ := valueObject("extra", "new").Object()

	got, err := newComparer(Options{}).compareUnknownObjectMembers(rootField("spec"), oldObject, newObject, nil, types.UnknownPrune)
	requireNoError(t, err)

	requireResult(t, got, nil, nil, nil)
}

func TestCompareUnknownObjectMembersInvalidPolicy(t *testing.T) {
	oldObject, _ := valueObject().Object()
	newObject, _ := valueObject().Object()

	_, err := newComparer(Options{}).compareUnknownObjectMembers(rootField("spec"), oldObject, newObject, nil, types.UnknownFieldPolicy(99))

	requireErrorIs(t, err, ErrInvalidDescriptor)
}
