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

package valueapply

import "testing"

func TestApplyMapKeyNoConflict(t *testing.T) {
	result, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("app", str("old"))),
		Applied:    obj(member("app", str("new"))),
		Descriptor: mapDescriptor(),
		Ownership:  state(),
	}, Options{})
	requireNoError(t, err)

	requireStringMember(t, result.Value, "app", "new")
	requireSet(t, result.AppliedFields, `$["app"]`)
}

func TestApplyMapKeyConflict(t *testing.T) {
	_, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("app", str("old"))),
		Applied:    obj(member("app", str("new"))),
		Descriptor: mapDescriptor(),
		Ownership:  state(entry("other", labelPath())),
	}, Options{})

	requireErrorIs(t, err, ErrConflict)
}

func TestApplyMapKeyForce(t *testing.T) {
	result, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       obj(member("app", str("old"))),
		Applied:    obj(member("app", str("new"))),
		Descriptor: mapDescriptor(),
		Ownership:  state(entry("other", labelPath())),
	}, Options{Force: true})
	requireNoError(t, err)

	requireOwnersOf(t, result.Ownership, labelPath(), "user")
}

func TestApplyAtomicListConflictAtListPath(t *testing.T) {
	_, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(str("a")),
		Applied:    list(str("b")),
		Descriptor: atomicStringListDescriptor(),
		Ownership:  state(entry("other", root())),
	}, Options{})

	requireErrorIs(t, err, ErrConflict)
}

func TestApplyAtomicListForceReplacesWholeList(t *testing.T) {
	result, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(str("a")),
		Applied:    list(str("b")),
		Descriptor: atomicStringListDescriptor(),
		Ownership:  state(entry("other", root())),
	}, Options{Force: true})
	requireNoError(t, err)

	requireListStrings(t, result.Value, "b")
	requireOwnersOf(t, result.Ownership, root(), "user")
}

func TestApplyOrderedListIndexConflict(t *testing.T) {
	_, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(str("a"), str("b")),
		Applied:    list(str("a"), str("B")),
		Descriptor: orderedStringListDescriptor(),
		Ownership:  state(entry("other", root().Index(1))),
	}, Options{})

	requireErrorIs(t, err, ErrConflict)
}

func TestApplyOrderedListIndexForce(t *testing.T) {
	result, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(str("a"), str("b")),
		Applied:    list(str("a"), str("B")),
		Descriptor: orderedStringListDescriptor(),
		Ownership:  state(entry("other", root().Index(1))),
	}, Options{Force: true})
	requireNoError(t, err)

	requireListStrings(t, result.Value, "a", "B")
	requireOwnersOf(t, result.Ownership, root().Index(1), "user")
}

func TestApplyListMapConditionStatusNoConflict(t *testing.T) {
	result, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(readyCondition("False")),
		Applied:    list(readyCondition("True")),
		Descriptor: conditionsDescriptor(),
		Ownership:  state(),
	}, Options{})
	requireNoError(t, err)

	item := requireListItem(t, result.Value, 0)
	requireStringMember(t, item, "status", "True")
	requireSet(t, result.AppliedFields, `$[{"type":"Ready"}].status`, `$[{"type":"Ready"}].type`)
}

func TestApplyListMapConditionStatusConflict(t *testing.T) {
	_, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(readyCondition("False")),
		Applied:    list(readyCondition("True")),
		Descriptor: conditionsDescriptor(),
		Ownership:  state(entry("other", readyStatusPath())),
	}, Options{})

	requireErrorIs(t, err, ErrConflict)
}

func TestApplyListMapConditionStatusForce(t *testing.T) {
	result, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(readyCondition("False")),
		Applied:    list(readyCondition("True")),
		Descriptor: conditionsDescriptor(),
		Ownership:  state(entry("other", readyStatusPath())),
	}, Options{Force: true})
	requireNoError(t, err)

	requireOwnersOf(t, result.Ownership, readyStatusPath(), "user")
}

func TestApplyListMapSameValueSharedOwnership(t *testing.T) {
	result, err := Apply(Request{
		Path:       root(),
		Owner:      owner("user"),
		Live:       list(readyCondition("True")),
		Applied:    list(readyCondition("True")),
		Descriptor: conditionsDescriptor(),
		Ownership:  state(entry("other", readyStatusPath())),
	}, Options{})
	requireNoError(t, err)

	requireOwnersOf(t, result.Ownership, readyStatusPath(), "other", "user")
}
